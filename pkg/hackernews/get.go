package hackernews

import (
	"bitbucket.org/xivart/hacker-news-api/pkg/models"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"net/http"
	"net/http/httptrace"
	"time"
)

type types interface {
	models.Item | models.User | models.TopStories
}

// get generic func making get calls
func get[T types](url string) (T, error) {
	var t T
	req, _ := http.NewRequest("GET", url, nil)
	var start time.Time
	var ttfb time.Duration
	trace := &httptrace.ClientTrace{
		GotFirstResponseByte: func() {
			ttfb = time.Since(start)
		},
	}
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
	start = time.Now()
	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		return t, err
	}
	log.Info().Dur("TotalTime", time.Since(start)).Dur("TTFB", ttfb).Str("URL", url).Send()
	if err := json.NewDecoder(resp.Body).Decode(&t); err != nil {
		return t, err
	}
	return t, nil
}

func (s *service) GetItem(id int) (models.Item, error) {
	url := fmt.Sprintf("https://hacker-news.firebaseio.com/v0/item/%d.json", id)
	return get[models.Item](url)
}

func (s *service) GetUser(username string) (models.User, error) {
	url := fmt.Sprintf("https://hacker-news.firebaseio.com/v0/user/%s.json", username)
	return get[models.User](url)
}

func (s *service) GetTopStories() ([]models.Item, error) {
	topStories, err := get[models.TopStories]("https://hacker-news.firebaseio.com/v0/topstories.json")
	if err != nil {
		return nil, err
	}
	items, errs := s.parallelizeGetItem(topStories)
	if len(errs) > 0 {
		for _, err1 := range errs {
			log.Err(err1).Send()
		}
		return nil, errors.New("multiple errors in getting stories")
	}
	return items, nil
}

// parallelizeGetItem uses worker pool to fetch stories using story id
func (s *service) parallelizeGetItem(topStories models.TopStories) ([]models.Item, []error) {
	type Result struct {
		item models.Item
		err  error
	}
	worker := func(jobs <-chan int, results chan<- Result) {
		for j := range jobs {
			item, err := s.GetItem(j)
			results <- Result{
				item,
				err,
			}
		}
	}
	numJobs := len(topStories)
	jobs := make(chan int, numJobs)
	results := make(chan Result, numJobs)
	for i := 0; i < s.maxWorkers; i++ {
		go worker(jobs, results)
	}
	for _, ts := range topStories {
		jobs <- ts
	}
	close(jobs)
	items := make([]models.Item, 0)
	errs := make([]error, 0)
	for i := 0; i < numJobs; i++ {
		result := <-results
		if result.err == error(nil) {
			items = append(items, result.item)
		} else {
			errs = append(errs, result.err)
		}
	}
	return items, errs
}

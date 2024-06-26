package hackernews

import (
	"errors"
	"fmt"
	"github.com/patrickmn/go-cache"
	"math"
	"sort"

	"github.com/rs/zerolog/log"

	"github.com/snigdhasambitak/hackernews-api/pkg/models"
)

func (s *service) Curated50(minKarma int) ([]models.Story, error) {
	return curated50(s.cache, minKarma, s.maxWorkers, s.GetTopStories, s.GetUser)
}

// curated50 func extracted to use func instead of pointer receiver for better testing
func curated50(cache *cache.Cache, minKarma int, maxWorkers int, getTopStories func() ([]models.Item, error), getUser func(username string) (models.User, error)) ([]models.Story, error) {
	cacheKey := fmt.Sprintf("curated50-%d", minKarma)

	// check the cache and return result if found
	if item, found := cache.Get(cacheKey); found {
		log.Info().Str("cache-key", cacheKey).Msg("returning cached result for curated50")
		return item.([]models.Story), nil
	}

	topStories, err := getTopStories()
	if err != nil {
		return nil, err
	}

	// filter the topStories based on author karma
	stories, errs := parallelizeFilterStories(topStories, minKarma, maxWorkers, getUser)
	if len(errs) > 0 {
		for _, err := range errs {
			log.Error().Msg(err.Error())
		}
		return nil, errors.New("multiple errors in getting author")
	}

	// sort and truncate top 50 stories based on comments
	sorted := append(make([]models.Story, 0, len(stories)), stories...)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Comments > sorted[j].Comments
	})
	sorted = sorted[:(int)(math.Min(float64(len(sorted)), 50))]
	// if sorted slice is expected iterate over sorted to add position value
	// add to cache and return sorted slice

	// this iteration keeps the order of topStories so resulted slice is not sorted
	toReturn := make([]models.Story, 0)
	for _, story := range stories {
		if len(toReturn) == 50 {
			break
		}
		if position, found := contains(sorted, story); found {
			story.Position = position + 1
			toReturn = append(toReturn, story)
		}
	}

	// add to cache
	cache.Set(cacheKey, toReturn, 0)
	return toReturn, nil
}

// parallelizeFilterStories uses worker pool to fetch authors of stories and filters them based on minKarma
func parallelizeFilterStories(topStories []models.Item, minKarma int, maxWorkers int, getUser func(username string) (models.User, error)) ([]models.Story, []error) {
	type Result struct {
		story models.Item
		user  models.User
		err   error
	}
	worker := func(jobs <-chan models.Item, results chan<- Result) {
		for j := range jobs {
			user, err := getUser(j.By)
			results <- Result{
				j,
				user,
				err,
			}
		}
	}
	numJobs := len(topStories)
	jobs := make(chan models.Item, numJobs)
	results := make(chan Result, 0)
	for i := 0; i < maxWorkers; i++ {
		go worker(jobs, results)
	}
	for _, ts := range topStories {
		jobs <- ts
	}
	close(jobs)
	stories := make([]models.Story, 0)
	errs := make([]error, 0)
	for i := 0; i < numJobs; i++ {
		result := <-results
		if result.err == error(nil) {
			if result.user.Karma > minKarma {
				stories = append(stories, models.Story{
					Author:   result.user.ID,
					Karma:    result.user.Karma,
					Comments: result.story.Descendants,
					Title:    result.story.Title,
				})
			}
		} else {
			errs = append(errs, result.err)
		}
	}
	return stories, errs
}

// contains checks if given slice of stories contains given story
func contains(stories []models.Story, story models.Story) (int, bool) {
	for i, s := range stories {
		if s.Title == story.Title && s.Author == story.Author {
			return i, true
		}
	}
	return -1, false
}

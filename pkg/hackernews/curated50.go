package hackernews

import (
	"bitbucket.org/xivart/hacker-news-api/pkg/models"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"sort"
)

func (s *service) Curated50(minKarma int) ([]models.Story, error) {
	cacheKey := fmt.Sprintf("curated50-%d", minKarma)

	// check the cache and return result if found
	if item, found := s.cache.Get(cacheKey); found {
		log.Info().Str("cache-key", cacheKey).Msg("returning cached result for curated50")
		return item.([]models.Story), nil
	}

	topStories, err := s.GetTopStories()
	if err != nil {
		return nil, err
	}

	// filter the topStories based on author karma
	stories, errs := s.parallelizeFilterStories(topStories, minKarma)
	if len(errs) > 0 {
		for _, err := range errs {
			log.Error().Msg(err.Error())
		}
		return nil, errors.New("multiple errors in getting author")
	}

	// sort and truncate top 50 stories based on comments
	sorted := append(make([]models.Story, 0, len(stories)), stories...)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Comments < sorted[j].Comments
	})
	sorted = sorted[0:50]
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
	s.cache.Set(cacheKey, toReturn, 0)
	return toReturn, nil
}

// parallelizeFilterStories uses worker pool to fetch authors of stories and filters them based on minKarma
func (s *service) parallelizeFilterStories(topStories []models.Item, minKarma int) ([]models.Story, []error) {
	type Result struct {
		story models.Item
		user  models.User
		err   error
	}
	worker := func(jobs chan models.Item, results chan Result) {
		for j := range jobs {
			item, err := s.GetUser(j.By)
			results <- Result{
				j,
				item,
				err,
			}
		}
	}
	numJobs := len(topStories)
	jobs := make(chan models.Item, numJobs)
	results := make(chan Result, 0)
	for i := 0; i < s.maxWorkers; i++ {
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

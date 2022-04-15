package hackernews

import (
	"github.com/patrickmn/go-cache"
	"github.com/snigdhasambitak/hackernews-api/pkg/models"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_curated50(t *testing.T) {
	c := cache.New(5*time.Minute, 10*time.Minute)
	minKarma := 15
	maxWorkers := 50
	var getStoriesCalled, getUserCalled int
	// mock getTopStories
	getTopStories := func() ([]models.Item, error) {
		getStoriesCalled++
		return []models.Item{
			{By: "yl", Descendants: 10, Title: "First"},
			{By: "na1", Descendants: 11, Title: "Ignored1"},
			{By: "sa", Descendants: 18, Title: "Second"},
			{By: "abd", Descendants: 25, Title: "third"},
			{By: "na2", Descendants: 2, Title: "Ignored2"},
			{By: "xyz", Descendants: 9, Title: "fourth"},
			{By: "dd", Descendants: 7, Title: "fifth"},
			{By: "na3", Descendants: 42, Title: "Ignored3"},
		}, nil
	}
	// mocked users
	users := map[string]models.User{
		"yl":  {ID: "yl", Karma: 20},
		"na1": {ID: "na1", Karma: 9},
		"sa":  {ID: "sa", Karma: 30},
		"abd": {ID: "abd", Karma: 25},
		"na2": {ID: "na2", Karma: 2},
		"xyz": {ID: "xyz", Karma: 77},
		"dd":  {ID: "dd", Karma: 90},
		"na3": {ID: "na3", Karma: 7},
	}
	// mock getUser
	getUser := func(username string) (models.User, error) {
		getUserCalled++
		return users[username], nil
	}
	actual, err := curated50(c, minKarma, maxWorkers, getTopStories, getUser)
	assertT := assert.New(t)
	assertT.NoError(err)
	assertT.Equal([]models.Story{
		{Author: "yl", Karma: 20, Comments: 10, Title: "First", Position: 3},
		{Author: "sa", Karma: 30, Comments: 18, Title: "Second", Position: 2},
		{Author: "abd", Karma: 25, Comments: 25, Title: "third", Position: 1},
		{Author: "xyz", Karma: 77, Comments: 9, Title: "fourth", Position: 4},
		{Author: "dd", Karma: 90, Comments: 7, Title: "fifth", Position: 5},
	}, actual)
	assertT.Equal(1, getStoriesCalled)
	assertT.Equal(8, getUserCalled)
}

func Test_curated50_WithCache(t *testing.T) {
	// TODO: test cache is hit instead of API calls
	assert.True(t, true)
}

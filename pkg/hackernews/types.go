package hackernews

import (
	"time"

	"github.com/patrickmn/go-cache"

	"github.com/snigdhasambitak/hackernews-api/pkg/models"
)

type Service interface {
	// GetItem returns item from hackernews API for given id
	GetItem(id int) (models.Item, error)
	// GetUser returns user from hackernews API for given username
	GetUser(username string) (models.User, error)
	// GetTopStories returns top 500 stories from hackernews API
	GetTopStories() ([]models.Item, error)
	// Curated50 returns top 50 of the latest 500 stories where the author has karma above 2413 with most comments
	Curated50(minKarma int) ([]models.Story, error)
}

type service struct {
	cache      *cache.Cache
	maxWorkers int
}

func NewService() Service {
	return &service{
		cache:      cache.New(5*time.Minute, 10*time.Minute),
		maxWorkers: 50,
	}
}

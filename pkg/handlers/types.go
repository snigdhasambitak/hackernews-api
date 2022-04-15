package handlers

import (
	"bitbucket.org/xivart/hacker-news-api/pkg/hackernews"
	"github.com/gin-gonic/gin"
)

type Handlers interface {
	GetStories(c *gin.Context)
}

type handlers struct {
	hn hackernews.Service
}

func NewHandlers() Handlers {
	return &handlers{
		hn: hackernews.NewService(),
	}
}

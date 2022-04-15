package handlers

import (
	"github.com/gin-gonic/gin"

	"github.com/snigdhasambitak/hackernews-api/pkg/hackernews"
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

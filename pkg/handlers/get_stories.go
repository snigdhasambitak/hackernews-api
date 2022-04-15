package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
)

func (s *handlers) GetStories(c *gin.Context) {
	start := time.Now()
	stories, err := s.hn.Curated50(2413)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		c.JSON(200, stories)
		log.Info().Dur("TotalTime", time.Since(start)).Msg("sent /stories")
	}
}

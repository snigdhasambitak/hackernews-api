package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func (s *handlers) GetStories(c *gin.Context) {
	start := time.Now()
	stories, err := s.hn.Curated50(2413)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, stories)
		log.Info().Dur("TotalTime", time.Since(start)).Msg("sent /stories")
	}
}

package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
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
		c.JSON(http.StatusOK, stories)
		log.Info().Dur("TotalTime", time.Since(start)).Msg("sent /stories")
	}
	TotalTime := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "TotalTime",
		Help: "Total request latency time taken by each request",
	})
	prometheus.Register(TotalTime)

	//r := gin.Default()
	//// Middleware to set Total request latency time for all requests
	//r.Use(func(context *gin.Context) {
	//	TotalTime.Set(float64(time.Since(start)))
	//})

}

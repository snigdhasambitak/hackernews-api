package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/snigdhasambitak/hackernews-api/pkg/handlers"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	r := gin.Default()
	h := handlers.NewHandlers()
	r.GET("/stories", h.GetStories)
	r.GET("/metrics", func(c *gin.Context) {
		handler := promhttp.Handler()
		handler.ServeHTTP(c.Writer, c.Request)
	})
	panic(r.Run())
}

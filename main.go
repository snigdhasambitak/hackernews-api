package main

import (
	"bitbucket.org/xivart/hacker-news-api/pkg/handlers"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	r := gin.Default()
	h := handlers.NewHandlers()
	r.GET("/stories", h.GetStories)
	panic(r.Run())
}

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// Configure zerolog
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Info().Msg("Starting OctaCart E-Commerce Platform...")

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		log.Debug().Msg("Ping endpoint called")
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	if err := r.Run(); err != nil {
		log.Fatal().Err(err).Msg("Server failed to start")
	}
}

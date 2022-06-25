package main

import (
	api "github.com/bellamy-labs/auth-api/api/v1"
	"github.com/bellamy-labs/auth-api/config"
	"github.com/bellamy-labs/auth-api/store"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// configure logging library
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	err := config.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err)
	}

	err = store.InitDB()
	if err != nil {
		log.Fatal().Err(err)
	}

	server := api.Server{}
	server.Run(":8080")
}

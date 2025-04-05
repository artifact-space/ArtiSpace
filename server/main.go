package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/artifact-space/ArtiSpace/config"
	"github.com/artifact-space/ArtiSpace/db"
	"github.com/artifact-space/ArtiSpace/log"
	"github.com/artifact-space/ArtiSpace/routes"
	"github.com/artifact-space/ArtiSpace/storage"
)

func main() {

	config, err := config.Config()
	if err != nil {
		log.Logger().Fatal().Msg("error in loading configuration")
	}

	err = storage.InitStorage(&config.Storage)
	if err != nil {
		log.Logger().Fatal().Msg("error in loading storage")
	}

	err = db.InitDB(&config.Database)
	if err != nil {
		log.Logger().Fatal().Msg("error in connecting to database")
	}
	defer db.Release()

	router := routes.InitRouter()
	dockerV2Router := routes.InitDockerV2Router()

	address := fmt.Sprintf(":%d", config.Server.Port)
	dockerV2Address := fmt.Sprintf(":%d", config.Server.DockerV2Port)

	server := &http.Server{
		Addr:    address,
		Handler: router,
	}

	dockerV2Server := &http.Server{
		Addr:    dockerV2Address,
		Handler: dockerV2Router,
	}

	go func(server *http.Server) {

		log.Logger().Info().Msgf("Listening on: %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Logger().Error().Err(err).Msgf("Unable to start server on: %s", server.Addr)
		}

	}(server)

	go func(server *http.Server) {

		log.Logger().Info().Msgf("Listening on: %s to serve Docker V2 APIs", server.Addr)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Logger().Error().Err(err).Msgf("Enable to start server on: %s", server.Addr)
		}
	}(dockerV2Server)

	//graceful shutdown
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	<-shutdown
	log.Logger().Info().Msg("Shutting down servers .....")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Logger().Error().Err(err).Msgf("Error in shutting down server: %s", server.Addr)
	}

	if err := dockerV2Server.Shutdown(ctx); err != nil {
		log.Logger().Error().Err(err).Msgf("Error in shutting down server: %s", dockerV2Server.Addr)
	}

	log.Logger().Info().Msg("Servers shut down gracefully.")
}

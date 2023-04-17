package main

import (
	"chassit-on-repeat/internal"
	"chassit-on-repeat/internal/service"
	"chassit-on-repeat/internal/utils"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Setup logger
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	level, err := zerolog.ParseLevel(utils.GetStringEnv("LOG_LEVEL", "info"))
	if err != nil {
		log.Fatal().Str("tag", "main").Err(err).Msg("Invalid log level")
	}
	log.Logger = log.Logger.Level(level)

	// Setup required environment variables
	if err := godotenv.Load(); err != nil {
		log.Info().Str("tag", "main").Msg("No .env file found")
	}
	utils.ValidateEnvs()

	// Create file handler and watcher
	handler, stop := internal.NewFileHandler()
	defer stop()

	// Create overrides handler
	overrides, oStop := internal.NewOverridesHandler()
	defer oStop()

	// Create a new service
	svc := service.NewService(handler, overrides)
	defer svc.Shutdown()

	// Create a new goroutine to listen for shutdown
	go func() {
		signChan := make(chan os.Signal, 1)

		signal.Notify(signChan, os.Interrupt, syscall.SIGTERM|syscall.SIGINT)
		sig := <-signChan
		log.Info().Str("tag", "main").Str("sig", sig.String()).Msg("Shutting down")

		svc.Shutdown()
	}()

	// Start service
	if err = svc.Start(); err != nil {
		log.Fatal().Str("tag", "main").Err(err).Send()
	}
	log.Info().Str("tag", "main").Msg("Goodbye!")
}

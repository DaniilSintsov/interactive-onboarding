package main

import (
	"log"

	"github.com/DaniilSintsov/interactive-onboarding/backend/internal/app"
	"github.com/DaniilSintsov/interactive-onboarding/backend/internal/config"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()

	if err != nil {
		log.Fatalf("can not initialize logger: %s", err)
	}

	cfg, err := config.New()

	if err != nil {
		logger.Fatal("can not initialize config", zap.Error(err))
	}

	if err = app.Run(logger, cfg); err != nil {
		logger.Fatal("application stopped with error", zap.Error(err))
	}
}

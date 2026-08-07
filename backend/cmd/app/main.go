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
		log.Fatalf("can not initialize config: %s", err)
	}

	app.Run(logger, cfg)
}

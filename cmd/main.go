package main

import (
	"github.com/blazee5/hezzl-test/internal/app"
	"github.com/blazee5/hezzl-test/internal/config"
	"github.com/blazee5/hezzl-test/lib/logger"
)

func main() {
	cfg := config.Load()

	log := logger.NewLogger(cfg)

	app.Run(log, cfg)
}

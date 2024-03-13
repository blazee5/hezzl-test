package logger

import (
	"github.com/blazee5/hezzl-test/internal/config"
	"go.uber.org/zap"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func NewLogger(cfg *config.Config) *zap.SugaredLogger {
	var log *zap.Logger

	switch cfg.Env {
	case envLocal:
		log, _ = zap.NewDevelopment()
	case envDev:
		log, _ = zap.NewDevelopment()
	case envProd:
		log, _ = zap.NewProduction()
	}

	defer log.Sync()

	return log.Sugar()
}

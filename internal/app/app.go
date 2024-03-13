package app

import (
	"context"
	"github.com/blazee5/hezzl-test/internal/config"
	"github.com/blazee5/hezzl-test/internal/handler"
	"github.com/blazee5/hezzl-test/internal/repository"
	"github.com/blazee5/hezzl-test/internal/service"
	"github.com/blazee5/hezzl-test/lib/db/postgres"
	"github.com/blazee5/hezzl-test/lib/db/redis"
	"github.com/blazee5/hezzl-test/lib/nats"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

func Run(log *zap.SugaredLogger, cfg *config.Config) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db := postgres.New(ctx, cfg)
	rdb := redis.New(ctx, cfg)
	nc := nats.New(cfg)
	repositories := repository.NewRepository(db, rdb)
	services := service.NewService(cfg, log, repositories)
	srv := handler.NewServer(log, cfg, db, rdb, nc, services)

	go func() {
		if err := srv.Run(srv.InitRoutes()); err != nil {
			log.Fatalf("Error while start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Infof("Error occured on server shutting down: %v", err)
	}
}

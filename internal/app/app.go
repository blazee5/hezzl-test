package app

import (
	"context"
	"github.com/blazee5/hezzl-test/internal/config"
	"github.com/blazee5/hezzl-test/internal/handler"
	nats2 "github.com/blazee5/hezzl-test/internal/nats"
	"github.com/blazee5/hezzl-test/internal/repository"
	"github.com/blazee5/hezzl-test/internal/service"
	"github.com/blazee5/hezzl-test/lib/db/clickhouse"
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
	nc, js, stream := nats.New(ctx, cfg)
	producer := nats2.NewProducer(js)

	clickConn := clickhouse.New(ctx, cfg)
	repositories := repository.NewRepository(db, rdb, clickConn)
	services := service.NewService(cfg, log, repositories, producer)
	srv := handler.NewServer(log, cfg, services, nc)
	consumer := nats2.NewConsumer(log, cfg, stream, repositories)

	go func() {
		if err := srv.Run(srv.InitRoutes()); err != nil {
			log.Fatalf("Error while start server: %v", err)
		}
	}()

	go func() {
		if err := consumer.Run(ctx); err != nil {
			log.Fatalf("error while run consumer: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	if err := srv.Shutdown(ctx); err != nil {
		log.Infof("Error occured on server shutting down: %v", err)
	}

	db.Close()
	if err := rdb.Close(); err != nil {
		log.Infof("error while close redis conn: %v", err)
	}

	if err := nc.Drain(); err != nil {
		log.Infof("error while close redis conn: %v", err)
	}

	nc.Close()
}

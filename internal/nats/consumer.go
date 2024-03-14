package nats

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/blazee5/hezzl-test/internal/config"
	"github.com/blazee5/hezzl-test/internal/domain"
	"github.com/blazee5/hezzl-test/internal/repository"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"time"
)

type Consumer struct {
	log    *zap.SugaredLogger
	cfg    *config.Config
	stream jetstream.Stream
	repo   *repository.Repository
}

func NewConsumer(log *zap.SugaredLogger, cfg *config.Config, stream jetstream.Stream, repo *repository.Repository) *Consumer {
	return &Consumer{log: log, cfg: cfg, stream: stream, repo: repo}
}

func (c *Consumer) Run(ctx context.Context) error {
	consumer, err := c.stream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Durable: "processor",
	})

	if err != nil {
		c.log.Errorf("error while create consumer: %v", err)
		return err
	}

	eg, ctx := errgroup.WithContext(ctx)
	for i := 0; i < c.cfg.Nats.WorkersCount; i++ {
		eg.Go(func() error {
			return c.RunConsumer(ctx, consumer)
		})
	}

	if err := eg.Wait(); err != nil {
		c.log.Errorf("error while running consumers: %v", err)
	}

	return nil
}

func (c *Consumer) RunConsumer(ctx context.Context, consumer jetstream.Consumer) error {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			msgs, err := consumer.FetchNoWait(c.cfg.Nats.BatchSize)
			if err != nil {
				if errors.Is(err, nats.ErrTimeout) {
					continue
				}
				c.log.Errorf("error fetching messages: %v", err)
				return err
			}

			var goods []domain.Good
			for msg := range msgs.Messages() {
				var good domain.Good
				if err := json.Unmarshal(msg.Data(), &good); err != nil {
					c.log.Errorf("error unmarshalling message: %v", err)
					msg.Nak()
					continue
				}
				goods = append(goods, good)
				msg.Ack()
			}

			if len(goods) > 0 {
				if err := c.repo.GoodClickhouse.InsertGoods(ctx, goods); err != nil {
					c.log.Errorf("error while inserting goods: %v", err)
				}
			}
		}
	}
}

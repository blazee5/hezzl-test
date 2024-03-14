package nats

import (
	"context"
	"fmt"
	"github.com/blazee5/hezzl-test/internal/config"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"log"
)

func New(ctx context.Context, cfg *config.Config) (*nats.Conn, jetstream.JetStream, jetstream.Stream) {
	nc, err := nats.Connect(fmt.Sprintf("nats://%s", cfg.Nats.Host))

	if err != nil {
		log.Fatalf("error while connect to nats: %v", err)
	}

	js, err := jetstream.New(nc)

	if err != nil {
		log.Fatalf("error while create jetstream: %v", err)
	}

	stream, err := js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:     "goods",
		Subjects: []string{"goods.>"},
	})

	if err != nil {
		log.Fatalf("error while create stream: %v", err)
	}

	return nc, js, stream
}

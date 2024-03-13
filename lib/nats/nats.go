package nats

import (
	"fmt"
	"github.com/blazee5/hezzl-test/internal/config"
	"github.com/nats-io/nats.go"
	"log"
)

func New(cfg *config.Config) *nats.Conn {
	nc, err := nats.Connect(fmt.Sprintf("nats://%s", cfg.Nats.Host))

	if err != nil {
		log.Fatalf("error while connect to nats: %v", err)
	}

	return nc
}

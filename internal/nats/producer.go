package nats

import (
	"context"
	"encoding/json"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/zap"
)

type Producer struct {
	log *zap.SugaredLogger
	js  jetstream.JetStream
}

func NewProducer(js jetstream.JetStream) *Producer {
	return &Producer{js: js}
}

func (p *Producer) Publish(ctx context.Context, good any) error {
	bytes, err := json.Marshal(&good)

	if err != nil {
		p.log.Infof("error while encode good: %v", err)
		return err
	}

	_, err = p.js.Publish(ctx, "goods.test", bytes)

	if err != nil {
		p.log.Infof("error while send good in nats: %v", err)
	}

	return nil
}

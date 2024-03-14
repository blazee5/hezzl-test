package clickhouse

import (
	"context"
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/blazee5/hezzl-test/internal/config"
	"log"
	"time"
)

func New(ctx context.Context, cfg *config.Config) driver.Conn {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{cfg.Clickhouse.Host},
		Auth: clickhouse.Auth{
			Database: cfg.Clickhouse.Name,
			Username: cfg.Clickhouse.User,
			Password: cfg.Clickhouse.Password,
		},
		Settings: clickhouse.Settings{
			"max_execution_time": 60,
		},
		DialTimeout: 5 * time.Second,
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		},
	})

	if err != nil {
		log.Fatalf("error while connect to clickhouse: %v", err)
	}

	if err := conn.Ping(ctx); err != nil {
		log.Fatalf("error while ping clickhouse: %v", err)
	}

	return conn
}

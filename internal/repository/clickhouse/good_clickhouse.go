package clickhouse

import (
	"context"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/blazee5/hezzl-test/internal/domain"
	"time"
)

type GoodRepository struct {
	conn driver.Conn
}

func NewGoodRepository(conn driver.Conn) *GoodRepository {
	return &GoodRepository{conn: conn}
}

func (repo *GoodRepository) InsertGoods(ctx context.Context, rows []domain.Good) error {
	batch, err := repo.conn.PrepareBatch(ctx, "INSERT INTO goods (Id, ProjectId, Name, Description, Priority, Removed, EventTime) VALUES ($1, $2, $3, $4. $5, $6, $7)")

	if err != nil {
		return err
	}

	for _, v := range rows {
		err := batch.Append(v.ID, v.ProjectID, v.Name, v.Description, v.Priority, v.Removed, time.Now())

		if err != nil {
			return err
		}
	}

	return batch.Send()
}

package repository

import (
	"context"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/blazee5/hezzl-test/internal/domain"
	"github.com/blazee5/hezzl-test/internal/models"
	"github.com/blazee5/hezzl-test/internal/repository/clickhouse"
	"github.com/blazee5/hezzl-test/internal/repository/postgres"
	redisRepo "github.com/blazee5/hezzl-test/internal/repository/redis"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Repository struct {
	Good
	GoodRedis
	GoodClickhouse
}

type Good interface {
	Get(ctx context.Context, limit, offset int) (domain.GoodList, error)
	GetByID(ctx context.Context, projectId, id int) (models.Good, error)
	Create(ctx context.Context, projectId int, input domain.CreateGoodRequest) (models.Good, error)
	Update(ctx context.Context, projectId, id int, input domain.UpdateGoodRequest) (models.Good, error)
	Reprioritize(ctx context.Context, projectId, id int, input domain.ReprioritizeRequest) (models.GoodPriorities, error)
	Delete(ctx context.Context, projectId, id int) (models.DeletedGood, error)
}

type GoodRedis interface {
	GetByIDCtx(ctx context.Context, key string) (*models.Good, error)
	SetGoodCtx(ctx context.Context, key string, seconds int, good *models.Good) error
	DeleteGoodCtx(ctx context.Context, key string) error
}

type GoodClickhouse interface {
	InsertGoods(ctx context.Context, rows []domain.Good) error
}

func NewRepository(db *pgxpool.Pool, rdb *redis.Client, clickConn driver.Conn) *Repository {
	return &Repository{
		Good:           postgres.NewGoodRepository(db),
		GoodRedis:      redisRepo.NewGoodRepository(rdb),
		GoodClickhouse: clickhouse.NewGoodRepository(clickConn),
	}
}

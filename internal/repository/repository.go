package repository

import (
	"context"
	"github.com/blazee5/hezzl-test/internal/domain"
	"github.com/blazee5/hezzl-test/internal/models"
	"github.com/blazee5/hezzl-test/internal/repository/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Repository struct {
	Good
}

type Good interface {
	Get(ctx context.Context, limit, offset int) (domain.GoodList, error)
	Create(ctx context.Context, projectId int, input domain.CreateGoodRequest) (models.Good, error)
	Update(ctx context.Context, projectId, id int, input domain.UpdateGoodRequest) (models.Good, error)
	Reprioritize(ctx context.Context, projectId, id int, input domain.ReprioritizeRequest) (models.GoodPriorities, error)
	Delete(ctx context.Context, projectId, id int) (models.DeletedGood, error)
}

func NewRepository(db *pgxpool.Pool, rdb *redis.Client) *Repository {
	return &Repository{
		Good: postgres.NewGoodRepository(db),
	}
}

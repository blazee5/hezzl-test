package redis

import (
	"context"
	"encoding/json"
	"github.com/blazee5/hezzl-test/internal/models"
	"github.com/redis/go-redis/v9"
	"time"
)

type GoodRepository struct {
	rdb *redis.Client
}

func NewGoodRepository(rdb *redis.Client) *GoodRepository {
	return &GoodRepository{rdb: rdb}
}

func (repo *GoodRepository) GetByIDCtx(ctx context.Context, key string) (*models.Good, error) {
	goodBytes, err := repo.rdb.Get(ctx, "good:"+key).Bytes()

	if err != nil {
		return nil, err
	}

	var good *models.Good

	if err = json.Unmarshal(goodBytes, &good); err != nil {
		return nil, err
	}

	return good, nil
}

func (repo *GoodRepository) SetGoodCtx(ctx context.Context, key string, seconds int, good *models.Good) error {
	goodBytes, err := json.Marshal(good)

	if err != nil {
		return err
	}

	if err := repo.rdb.Set(ctx, "good:"+key, goodBytes, time.Second*time.Duration(seconds)).Err(); err != nil {
		return err
	}

	return nil
}

func (repo *GoodRepository) DeleteGoodCtx(ctx context.Context, key string) error {
	if err := repo.rdb.Del(ctx, "good:"+key).Err(); err != nil {
		return err
	}

	return nil
}

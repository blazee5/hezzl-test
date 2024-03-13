package redis

import (
	"context"
	"encoding/json"
	"github.com/blazee5/hezzl-test/internal/models"
	"github.com/redis/go-redis/v9"
	"time"
)

type GoodRedisRepository struct {
	rdb *redis.Client
}

func NewGoodRedisRepository(rdb *redis.Client) *GoodRedisRepository {
	return &GoodRedisRepository{rdb: rdb}
}

func (repo *GoodRedisRepository) GetByIDCtx(ctx context.Context, key string) (*models.Good, error) {
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

func (repo *GoodRedisRepository) SetGoodCtx(ctx context.Context, key string, seconds int, good *models.Good) error {
	goodBytes, err := json.Marshal(good)

	if err != nil {
		return err
	}

	if err := repo.rdb.Set(ctx, "good:"+key, goodBytes, time.Second*time.Duration(seconds)).Err(); err != nil {
		return err
	}

	return nil
}

func (repo *GoodRedisRepository) DeleteGoodCtx(ctx context.Context, key string) error {
	if err := repo.rdb.Del(ctx, "good:"+key).Err(); err != nil {
		return err
	}

	return nil
}

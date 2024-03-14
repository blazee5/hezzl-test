package service

import (
	"context"
	"github.com/blazee5/hezzl-test/internal/config"
	"github.com/blazee5/hezzl-test/internal/domain"
	"github.com/blazee5/hezzl-test/internal/models"
	"github.com/blazee5/hezzl-test/internal/nats"
	"github.com/blazee5/hezzl-test/internal/repository"
	"go.uber.org/zap"
	"strconv"
)

type GoodService struct {
	cfg      *config.Config
	log      *zap.SugaredLogger
	repo     *repository.Repository
	producer *nats.Producer
}

func NewGoodService(cfg *config.Config, log *zap.SugaredLogger, repo *repository.Repository, producer *nats.Producer) *GoodService {
	return &GoodService{cfg: cfg, log: log, repo: repo, producer: producer}
}

func (s *GoodService) GetGoods(ctx context.Context, limit, offset int) (domain.GoodList, error) {
	return s.repo.Good.Get(ctx, limit, offset)
}

func (s *GoodService) GetGoodByID(ctx context.Context, projectID, id int) (models.Good, error) {
	cachedGood, err := s.repo.GoodRedis.GetByIDCtx(ctx, strconv.Itoa(id))

	if err != nil {
		s.log.Infof("cannot get good by id in redis: %v", err)
	}

	if cachedGood != nil {
		return *cachedGood, nil
	}

	good, err := s.repo.Good.GetByID(ctx, projectID, id)

	if err != nil {
		return models.Good{}, err
	}

	if err := s.repo.GoodRedis.SetGoodCtx(ctx, strconv.Itoa(id), 60, &good); err != nil {
		s.log.Infof("error while save good to redis: %v", err)
	}

	return good, nil
}

func (s *GoodService) CreateGood(ctx context.Context, projectID int, input domain.CreateGoodRequest) (models.Good, error) {
	good, err := s.repo.Good.Create(ctx, projectID, input)

	err = s.producer.Publish(ctx, good)

	if err != nil {
		s.log.Infof("error while send good in nats: %v", err)
	}

	return good, nil
}

func (s *GoodService) UpdateGood(ctx context.Context, projectID, id int, input domain.UpdateGoodRequest) (models.Good, error) {
	good, err := s.repo.Good.Update(ctx, projectID, id, input)

	if err != nil {
		return models.Good{}, err
	}

	if err = s.repo.GoodRedis.DeleteGoodCtx(ctx, strconv.Itoa(id)); err != nil {
		s.log.Infof("error while delete good from redis: %v", err)
	}

	err = s.producer.Publish(ctx, good)

	if err != nil {
		s.log.Infof("error while send good in nats: %v", err)
	}

	return good, nil
}

func (s *GoodService) ReprioritizeGood(ctx context.Context, projectID, id int, input domain.ReprioritizeRequest) (models.GoodPriorities, error) {
	goods, err := s.repo.Good.Reprioritize(ctx, projectID, id, input)

	if err != nil {
		return models.GoodPriorities{}, err
	}

	for _, good := range goods.Priorities {
		if err = s.repo.GoodRedis.DeleteGoodCtx(ctx, strconv.Itoa(good.ID)); err != nil {
			s.log.Infof("error while delete good from redis: %v", err)
		}

		err = s.producer.Publish(ctx, good)

		if err != nil {
			s.log.Infof("error while send good in nats: %v", err)
		}
	}

	return goods, nil
}

func (s *GoodService) DeleteGood(ctx context.Context, projectID, id int) (models.DeletedGood, error) {
	good, err := s.repo.Good.Delete(ctx, projectID, id)

	if err != nil {
		return models.DeletedGood{}, err
	}

	if err = s.repo.GoodRedis.DeleteGoodCtx(ctx, strconv.Itoa(id)); err != nil {
		s.log.Infof("error while delete good from redis: %v", err)
	}

	err = s.producer.Publish(ctx, good)

	if err != nil {
		s.log.Infof("error while send good in nats: %v", err)
	}

	return good, nil
}

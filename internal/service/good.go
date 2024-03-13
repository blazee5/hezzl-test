package service

import (
	"context"
	"github.com/blazee5/hezzl-test/internal/config"
	"github.com/blazee5/hezzl-test/internal/domain"
	"github.com/blazee5/hezzl-test/internal/models"
	"github.com/blazee5/hezzl-test/internal/repository"
	"go.uber.org/zap"
)

type GoodService struct {
	cfg  *config.Config
	log  *zap.SugaredLogger
	repo *repository.Repository
}

func NewGoodService(cfg *config.Config, log *zap.SugaredLogger, repo *repository.Repository) *GoodService {
	return &GoodService{cfg: cfg, log: log, repo: repo}
}

func (s *GoodService) GetGoods(ctx context.Context, limit, offset int) (domain.GoodList, error) {
	return s.repo.Good.Get(ctx, limit, offset)
}

func (s *GoodService) CreateGood(ctx context.Context, projectID int, input domain.CreateGoodRequest) (models.Good, error) {
	return s.repo.Good.Create(ctx, projectID, input)
}

func (s *GoodService) UpdateGood(ctx context.Context, projectID, id int, input domain.UpdateGoodRequest) (models.Good, error) {
	return s.repo.Good.Update(ctx, projectID, id, input)
}

func (s *GoodService) ReprioritizeGood(ctx context.Context, projectID, id int, input domain.ReprioritizeRequest) (models.GoodPriorities, error) {
	return s.repo.Good.Reprioritize(ctx, projectID, id, input)
}

func (s *GoodService) DeleteGood(ctx context.Context, projectID, id int) (models.DeletedGood, error) {
	return s.repo.Good.Delete(ctx, projectID, id)
}

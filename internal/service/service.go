package service

import (
	"context"
	"github.com/blazee5/hezzl-test/internal/config"
	"github.com/blazee5/hezzl-test/internal/domain"
	"github.com/blazee5/hezzl-test/internal/models"
	"github.com/blazee5/hezzl-test/internal/repository"
	"go.uber.org/zap"
)

type Service struct {
	Good
}

type Good interface {
	GetGoods(ctx context.Context, limit, offset int) (domain.GoodList, error)
	CreateGood(ctx context.Context, projectID int, input domain.CreateGoodRequest) (models.Good, error)
	UpdateGood(ctx context.Context, projectID, id int, input domain.UpdateGoodRequest) (models.Good, error)
	ReprioritizeGood(ctx context.Context, projectID, id int, input domain.ReprioritizeRequest) (models.GoodPriorities, error)
	DeleteGood(ctx context.Context, projectID, id int) (models.DeletedGood, error)
}

func NewService(cfg *config.Config, log *zap.SugaredLogger, repo *repository.Repository) *Service {
	return &Service{Good: NewGoodService(cfg, log, repo)}
}
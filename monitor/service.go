package monitor

import (
	"context"
	"dogker/lintang/monitor-service/domain"

	"go.uber.org/zap"
)

// type ContainerRepository interface {
// 	GetById(ctx context.Context) string
// 	GetAllUserContainer(ctx context.Context, userId string) ([]domain.Container, error)
	
// }

type Service struct {
	containerRepo ContainerRepository
}

func NewService(c ContainerRepository) *Service {
	return &Service{
		containerRepo: c,
	}
}


func (m *Service) GetAllUserContainerService(ctx context.Context, userId string) (*[]domain.Container, error) {
	res, err := m.containerRepo.GetAllUserContainer(ctx, userId)
	if err != nil {
		zap.L().Debug("m.containerRepo.GetAllUserContainer", zap.String("userID", userId))
		return nil, nil
	}
	return res, nil
} 

func (m *Service) TesDoang(ctx context.Context) (string, error) {

	zap.L().Debug("Hello sadoakdaas", zap.String("user", "lintang"),
		zap.Int("age", 20))
	zap.L().Error("Hello error", zap.String("user", "lintang"),
		zap.Int("age", 20))

	zap.L().Info("Hello error", zap.String("user", "lintang"),
		zap.Int("age", 20))

	zap.L().Warn("Hello error", zap.String("user", "lintang"),
		zap.Int("age", 20))
	return "tess", nil
}

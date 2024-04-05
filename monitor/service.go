package monitor

import (
	"context"

	"go.uber.org/zap"
)

type ContainerRepository interface {
	GetById(ctx context.Context) string
}



type Service struct {
	containerRepo ContainerRepository
}

func NewService(c ContainerRepository) *Service {
	return &Service{
		containerRepo: c,
	}
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

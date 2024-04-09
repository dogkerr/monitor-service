package monitor

import (
	"context"
	"dogker/lintang/monitor-service/domain"
	"errors"

	"go.uber.org/zap"
)

type DashboardRepository interface {
	CreateDashboard(ctx context.Context, dashboard *domain.Dashboard) error
	GetByUserIDAndType(ctx context.Context, userID, dbType string) (*domain.Dashboard, error)
}

type GrafanaAPI interface {
	CreateDashboard(ctx context.Context, userID string) (*domain.Dashboard, error)
}

type Service struct {
	containerRepo ContainerRepository
	grafanaClient GrafanaAPI
	dashboardRepo DashboardRepository
}

func NewService(c ContainerRepository, grf GrafanaAPI, db DashboardRepository) *Service {
	return &Service{
		containerRepo: c,
		grafanaClient: grf,
		dashboardRepo: db,
	}
}

// get dashboard user di tabel dashboard kalo ada, kalo gak ada buatin dashboard grafana baru dan simpen uid dashboard nya ke tabel dashboard
func (m *Service) GetUserMonitorDashboard(ctx context.Context, userID string) (*domain.Dashboard, error) {
	res, err := m.dashboardRepo.GetByUserIDAndType(ctx, userID, "monitor")
	if err != nil && errors.Is(err, domain.ErrNotFound) {
		newDashboard, err := m.grafanaClient.CreateDashboard(ctx, userID)
		if err != nil {
			zap.L().Error("cant create grafana dashboard", zap.String("userID", userID))
			return nil, err
		}
		err = m.dashboardRepo.CreateDashboard(ctx, newDashboard)
		if err != nil {
			zap.L().Error("cant create grafana dashboard", zap.String("userID", userID))
			return nil, err
		}
		return newDashboard, nil
	}
	return res, err
}

// buat testing doang
func (m *Service) GetAllUserContainerService(ctx context.Context, userID string) (*[]domain.Container, error) {
	res, err := m.containerRepo.GetAllUserContainer(ctx, userID)
	if err != nil {
		zap.L().Debug("m.containerRepo.GetAllUserContainer", zap.String("userID", userID))
		return nil, nil
	}
	return res, nil
}

// buat testing daong
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

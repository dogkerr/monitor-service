package monitor

import (
	"context"
	"dogker/lintang/monitor-service/domain"
	"errors"
	"time"

	"go.uber.org/zap"
)

type DashboardRepository interface {
	CreateDashboard(ctx context.Context, dashboard *domain.Dashboard) error
	GetByUserIDAndType(ctx context.Context, userID, dbType string) (*domain.Dashboard, error)
}

type GrafanaAPI interface {
	CreateMonitorDashboard(ctx context.Context, userID string) (*domain.Dashboard, error)
	CreateLogsDashboard(ctx context.Context, userID string) (*domain.Dashboard, error)
}

type UserRepository interface {
	GetAllUsers(ctx context.Context) (*[]domain.User, error)
}

type MonitorMQ interface {
	SendAllUserMetrics(ctx context.Context, usersAllMetrics []domain.UserMetricsMessage) error
}

type Service struct {
	containerRepo ContainerRepository
	grafanaClient GrafanaAPI
	dashboardRepo DashboardRepository
	userRepo      UserRepository
	promeAPI      PrometheusAPI
	monitorMQ     MonitorMQ
}

func NewService(c ContainerRepository, grf GrafanaAPI, db DashboardRepository, userDb UserRepository, prome PrometheusAPI,
	mtqMq MonitorMQ) *Service {
	return &Service{
		containerRepo: c,
		grafanaClient: grf,
		dashboardRepo: db,
		userRepo:      userDb,
		promeAPI:      prome,
		monitorMQ:     mtqMq,
	}
}

// get dashboard prometheus user di tabel dashboard kalo ada, kalo gak ada buatin dashboard grafana baru dan simpen uid dashboard nya ke tabel dashboard
func (m *Service) GetUserMonitorDashboard(ctx context.Context, userID string) (*domain.Dashboard, error) {
	res, err := m.dashboardRepo.GetByUserIDAndType(ctx, userID, "monitor")
	if err != nil && errors.Is(err, domain.ErrNotFound) {
		newDashboard, err := m.grafanaClient.CreateMonitorDashboard(ctx, userID)
		if err != nil {
			zap.L().Error("cant create grafana prometheus dashboard", zap.String("userID", userID))
			return nil, err
		}
		err = m.dashboardRepo.CreateDashboard(ctx, newDashboard)
		if err != nil {
			zap.L().Error("cant insert prometheus dashboard to db", zap.String("userID", userID))
			return nil, err
		}
		return newDashboard, nil
	}
	return res, err
}

// get dashboard loki user di tabel dashboard kalo ada, kalo gak ada buatin dashboard grafana baru dan simpen uid dashboard nya ke tabel dashboard
func (m *Service) GetLogsDashboard(ctx context.Context, userID string) (*domain.Dashboard, error) {
	res, err := m.dashboardRepo.GetByUserIDAndType(ctx, userID, "log")
	if err != nil && errors.Is(err, domain.ErrNotFound) {
		newLogsDashboard, err := m.grafanaClient.CreateLogsDashboard(ctx, userID)
		if err != nil {
			zap.L().Error("cant create grafana loki dashboard", zap.String("userID", userID))
			return nil, err
		}
		err = m.dashboardRepo.CreateDashboard(ctx, newLogsDashboard)
		if err != nil {
			zap.L().Error("cant insert loki dashboard to db", zap.String("userID", userID))
			return nil, err
		}
		return newLogsDashboard, nil
	}
	return res, nil
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

func (m *Service) SendAllUsersMetricsToRMQ(ctx context.Context) error {
	users, err := m.userRepo.GetAllUsers(ctx) // mendapatkan semua users
	if err != nil {
		zap.L().Error("m.userRepo.GetAllUsers", zap.Error(err))
		return err
	}

	// var allUsersMetrics domain.AllUsersMetricsMessage
	var allUsersMetrics []domain.UserMetricsMessage
	for _, user := range *users {
		// iterate all users
		userContainers, err := m.containerRepo.GetAllUserContainer(ctx, user.ID.String())
		if err != nil {
			// users belum punya container
			continue
		}

		for i := 0; i < len(*userContainers); i++ {
			// iterate semua container  milik user
			// set metrics container for all usersemua
			ctr := *userContainers

			ctrMetrics, err := m.promeAPI.GetMetricsByServiceIDNotGRPC(ctx, ctr[i].ServiceID, time.Now().Add(-30*time.Minute))
			if err != nil {
				zap.L().Error("server.prome.GetMetricsByServiceId", zap.Error(err))
				return err
			}
			allUsersMetrics = append(allUsersMetrics, domain.UserMetricsMessage{
				ContainerID:         ctr[i].ID.String(),
				UserID:              user.ID.String(),
				CpuUsage:            ctrMetrics.CpuUsage,
				MemoryUsage:         ctrMetrics.MemoryUsage,
				NetworkIngressUsage: ctrMetrics.NetworkIngressUsage,
				NetworkEgressUsage:  ctrMetrics.NetworkEgressUsage,
			})
			// append usersMetricsMessag
			// allUsersMetrics.AllUsersMetrics = append(allUsersMetrics.AllUsersMetrics,
			// domain.UserMetricsMessage{
			// 	ContainerID:         ctr[i].ID.String(),
			// 	UserID:              user.ID.String(),
			// 	CpuUsage:            ctrMetrics.CpuUsage,
			// 	MemoryUsage:         ctrMetrics.MemoryUsage,
			// 	NetworkIngressUsage: ctrMetrics.NetworkIngressUsage,
			// 	NetworkEgressUsage:  ctrMetrics.NetworkEgressUsage,
			// },
			// )

		}
	}
	err = m.monitorMQ.SendAllUserMetrics(ctx, allUsersMetrics)
	if err != nil {
		zap.L().Error("error pas SendAllUserMetrics ke rabbittmq: ", zap.Error(err))
		return err
	}
	return nil
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

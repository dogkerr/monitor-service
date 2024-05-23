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
	GetDashboardOwner(ctx context.Context, dashboardUID string, userID string) error
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

type ContainerServiceClient interface {
	SendTerminatedContainerToCtrService(ctx context.Context, terminatedInstances []string) error
}

type Service struct {
	containerRepo ContainerRepository
	grafanaClient GrafanaAPI
	dashboardRepo DashboardRepository
	userRepo      UserRepository
	promeAPI      PrometheusAPI
	monitorMQ     MonitorMQ
	ctrClient     ContainerServiceClient
}

func NewService(c ContainerRepository, grf GrafanaAPI, db DashboardRepository, userDb UserRepository, prome PrometheusAPI,
	mtqMq MonitorMQ, ctrClient ContainerServiceClient) *Service {
	return &Service{
		containerRepo: c,
		grafanaClient: grf,
		dashboardRepo: db,
		userRepo:      userDb,
		promeAPI:      prome,
		monitorMQ:     mtqMq,
		ctrClient:     ctrClient,
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

			var ctrMetrics *domain.Metric
			// get metrics dari prometehus
			ctrMetrics, err = m.promeAPI.GetMetricsByServiceIDNotGRPC(ctx, ctr[i].ServiceID, time.Now().Add(-30*time.Minute))

			if err != nil {
				zap.L().Error("server.prome.GetMetricsByServiceId", zap.Error(err))
				return err
			}

			if ctrMetrics.CpuUsage == 0 {
				// kalo metrics cpu prometheus gak ada berarti containernya udah pernah diterminate, dan harus ambil metrics dari postgres
				// tapi ini swarm servicenya harus dihapus lewat endpoint kalo gak lewat endpoint nanti error karena di tabel container-metrics belum ada rownya

				ctrMetrics, err = m.containerRepo.GetSpecificConatainerMetrics(ctx, ctr[i].ID.String())
				if err != nil {
					zap.L().Error("m.containerRepo.GetSpecificConatainerMetrics", zap.Error(err))
					return err
				}
			}

			allUsersMetrics = append(allUsersMetrics, domain.UserMetricsMessage{
				ContainerID:         ctr[i].ID.String(),
				UserID:              user.ID.String(),
				CpuUsage:            ctrMetrics.CpuUsage,
				MemoryUsage:         ctrMetrics.MemoryUsage,
				NetworkIngressUsage: ctrMetrics.NetworkIngressUsage,
				NetworkEgressUsage:  ctrMetrics.NetworkEgressUsage,
			})

		}
	}
	err = m.monitorMQ.SendAllUserMetrics(ctx, allUsersMetrics)
	if err != nil {
		zap.L().Error("error pas SendAllUserMetrics ke rabbittmq: ", zap.Error(err))
		return err
	}
	return nil
}

func (m *Service) SendTerminatedInstanceToContainerService(ctx context.Context) error {
	deadSvc, err := m.promeAPI.GetTerminatedContainers(ctx)
	if err != nil {
		zap.L().Error(" m.promeAPI.GetTerminatedContainers(ctx) (SendTerminatedInstanceToContainerService) (MonoitorService)", zap.Error(err))
		return err
	}

	err = m.ctrClient.SendTerminatedContainerToCtrService(ctx, deadSvc)
	if err != nil {
		return err
	}
	return nil
}

func (m *Service) AuthorizeGrafanaDashboardAccess(ctx context.Context, ctrID string, userID string) error {
	err := m.dashboardRepo.GetDashboardOwner(ctx, ctrID, userID)
	if err != nil {
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

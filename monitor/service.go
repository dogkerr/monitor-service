package monitor

import (
	"context"
	"dogker/lintang/monitor-service/domain"
	"errors"
	"time"

	"github.com/google/uuid"
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
	SendDownContainerToCtrService(ctx context.Context, terminatedInstances []string) error
	GetContainerStatus(ctx context.Context, serviceID string) (bool, error)
}

type MailingWebAPI interface {
	SendDownSwarmServiceToMailingService(ctx context.Context, ml domain.CommonLabelsMailing) error
}

type Service struct {
	containerRepo ContainerRepository
	grafanaClient GrafanaAPI
	dashboardRepo DashboardRepository
	userRepo      UserRepository
	promeAPI      PrometheusAPI
	monitorMQ     MonitorMQ
	ctrClient     ContainerServiceClient
	mailingClient MailingWebAPI
}

func NewService(c ContainerRepository, grf GrafanaAPI, db DashboardRepository, userDb UserRepository, prome PrometheusAPI,
	mtqMq MonitorMQ, ctrClient ContainerServiceClient, mailingClient MailingWebAPI) *Service {
	return &Service{
		containerRepo: c,
		grafanaClient: grf,
		dashboardRepo: db,
		userRepo:      userDb,
		promeAPI:      prome,
		monitorMQ:     mtqMq,
		ctrClient:     ctrClient,
		mailingClient: mailingClient,
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
			// get metrics dari prometehus setiap container
			ctrMetrics, err = m.promeAPI.GetMetricsByServiceIDNotGRPC(ctx, ctr[i].ServiceID, time.Now().Add(-30*time.Minute))

			if err != nil {
				//
				zap.L().Error("server.prome.GetMetricsByServiceId", zap.Error(err))
				return err
			}

			if ctrMetrics.CpuUsage == 0 {
				// ketika di promethus gak ada
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

// SendTerminatedInstanceToContainerService
// @Desc Service ini dijalankan setiap 4 detik oleh cron job
// dia listen prometheus buat tau swarm service mana yang mati (bukan dihapus) bukan karena dimatiiinn user lewat api container-servicee
// terus accidentally stoped swarm service dikirim ke container-service
// container-service masukin metrics terakhir dari ctr tersebut ke tabel metrics
// karna stop swarm service itu bisa berkali kali , misal user stop terus start terus stop swarm service lagi
// func (m *Service) SendDownInstanceToContainerServiceAndMailingService(ctx context.Context) error {
// 	deadSvc, currentDownReplicas, err := m.promeAPI.GetStoppedContainers(ctx) // get serviceIDs container yang stopped > 10s, tapi gak swarm service/semua container swarm service tsb yang mati , jadi ada 1 ctr dari swarm service mati padahal total replica ada 5 pun tetep direturn method ini
// 	if err != nil {
// 		zap.L().Error(" m.promeAPI.GetTerminatedContainers(ctx) (SendTerminatedInstanceToContainerService) (MonoitorService)", zap.Error(err))
// 		return err
// 	}

// 	notProcessedDownServiceIDs, notProcessedDownCtrIDs, err := m.containerRepo.GetProcessedContainers(ctx, deadSvc, currentDownReplicas) // cek apakah container sebelumnya pernah diproses service ini , returnnya ctrIDs dan serviceIDs yang belum pernah diproses service ini
// 	if err != nil {
// 		zap.L().Error("m.containerRepo.GetProcessedContainers(ctx, deadSvc) (SendTerminatedInstanceToContainerService) (ContainerService)", zap.Error(err))
// 		return err
// 	}
// 	notProcessedDownCtrIDs = m.removeDuplicateUUIDArray(notProcessedDownCtrIDs, len(notProcessedDownCtrIDs))
// 	notProcessedDownServiceIDs = m.removeDuplicateStringArray(notProcessedDownServiceIDs, len(notProcessedDownServiceIDs))
// 	for i, notProcessedSvcID := range notProcessedDownServiceIDs {
// 		if notProcessedSvcID == "" {
// 			notProcessedDownServiceIDs = notProcessedDownServiceIDs[:i]
// 			notProcessedDownCtrIDs = notProcessedDownCtrIDs[:i]
// 		}
// 	}
// 	// cek apakah emang container down bukan karena di stop oleh user pemilik container . Caranya dg cek last status container
// 	for i, _ := range notProcessedDownCtrIDs { // ini bikin error
// 		latestCtrILife, err := m.containerRepo.GetLatestContainerLifecycleByCtrID(ctx, notProcessedDownCtrIDs[i].String())
// 		if err != nil {
// 			zap.L().Error("m.containerRepo.GetLatestContainerLifecycleByCtrID (SendDownInstanceToContainerServiceAndMailingService) (ContainerService)", zap.Error(err))
// 			// gak usah return?
// 			return err
// 		}

// 		if latestCtrILife.Status == domain.STOP {
// 			// jika status latest containerLifeCycle adalah stopped , berarti emang sengaja distop user (di handler endpoint stopContainer, update latest ctrlIfe jadi stopped )
// 			// delete inplace elemen array ini

// 			notProcessedDownCtrIDs[i] = notProcessedDownCtrIDs[len(notProcessedDownCtrIDs)-1]
// 			notProcessedDownCtrIDs = notProcessedDownCtrIDs[:len(notProcessedDownCtrIDs)-1]

// 			// itk notProcessedDownServiceIDs
// 			notProcessedDownServiceIDs[i] = notProcessedDownServiceIDs[len(notProcessedDownServiceIDs)-1]
// 			notProcessedDownServiceIDs = notProcessedDownServiceIDs[:len(notProcessedDownServiceIDs)-1]
// 		}

// 	}

// 		/// yang ini grpc bikin error...

// 	// for i, svcID := range notProcessedDownServiceIDs {

// 	// 	ctrStatusIsRun, err := m.ctrClient.GetContainerStatus(ctx, svcID)
// 	// 	if err != nil {
// 	// 		zap.L().Error(" m.ctrClient.GetContainerStatus (SendDownInstanceToContainerServiceAndMailingService)", zap.Error(err))
// 	// 		return err
// 	// 	}
// 	// 	if ctrStatusIsRun == true {
// 	// 		// kalau masih run (sattus tabel container & status di dockerAPI) hapus dari array

// 	// 		notProcessedDownServiceIDs[i] = notProcessedDownServiceIDs[len(notProcessedDownServiceIDs)-1]
// 	// 		notProcessedDownServiceIDs = notProcessedDownServiceIDs[:len(notProcessedDownServiceIDs)-1]

// 	// 		notProcessedDownCtrIDs[i] = notProcessedDownCtrIDs[len(notProcessedDownCtrIDs)-1]
// 	// 		notProcessedDownCtrIDs = notProcessedDownCtrIDs[:len(notProcessedDownCtrIDs)-1]

// 	// 	}
// 	// }

// 	// cek lagi di docker api dia masih stopped atau gak...

// 	// kalau latestCtrLife.Sttaus == RUN dan sebelumnya distop ->  tetep dikirim ke ctr service

// 	// ini gaperlu karena nanti bakal kestart lagi servicenya (bikin error jg), dan gaperlu update status container jg 
// 	// kalau container stopped bukan karena aksi sengaja user , kirim ke container svc & email svc
// 	// err = m.ctrClient.SendDownContainerToCtrService(ctx, notProcessedDownServiceIDs) // kiirm terminated container ke container-service buat
// 	// if err != nil {
// 	// 	return err
// 	// }

// 	downSwarmServicesDetail, err := m.containerRepo.GetSwarmServicesDetail(ctx, notProcessedDownServiceIDs) // get swarm service detail buat diikiirm ke mailing service
// 	if err != nil {
// 		return err
// 	}

// 	for i, _ := range notProcessedDownCtrIDs {
		
// 		// insert terminated conatiner yang baru aja diprocess ke tabel terminated container
// 		newProcessedCtr := notProcessedDownCtrIDs[i]
// 		err := m.containerRepo.InsertTerminatedContainer(ctx, newProcessedCtr.String())
// 		if err != nil {
// 			zap.L().Error("m.containerRepo.InsertTerminatedContainer(ctx, newProcessedCtr) (SendTerminatedInstanceToContainerService) (ContainerService)", zap.Error(err))
// 			return err
// 		}
// 	}

// 	if len(downSwarmServicesDetail) != 0 {
// 		// send down service message to mailing (udah pasti 1 service down aja tiap 1 menit yang dikirim)
// 		for i := range downSwarmServicesDetail {
// 			// send down swarm service detail to mailing service
// 			err := m.mailingClient.SendDownSwarmServiceToMailingService(ctx, downSwarmServicesDetail[i])
// 			if err != nil {
// 				zap.L().Error(" m.mailingClient.SendDownSwarmServiceToMailingService (SendDownInstanceToContainerServiceAndMailingService) (MonitorService)", zap.Error(err))
// 				return err
// 			}
// 		}
// 	}

// 	return nil
// }



func (m *Service) AuthorizeGrafanaDashboardAccess(ctx context.Context, ctrID string, userID string) error {
	err := m.dashboardRepo.GetDashboardOwner(ctx, ctrID, userID)
	if err != nil {
		return err
	}
	return nil
}

func (m *Service) removeDuplicateStringArray(arr []string, n int) []string {
	if n == 0 || n == 1 {
		return arr
	}
	var temp []string = make([]string, n)
	var j int = 0

	for i := 0; i < n-1; i++ {
		if arr[i] != arr[i+1] {
			temp[j] = arr[i]
			j += 1
		}
	}
	temp[j] = arr[n-1]
	j += 1

	for i := 0; i < j; i++ {
		arr[i] = temp[i]

	}
	arr = arr[:j]
	return arr
}

func (m *Service) removeDuplicateUUIDArray(arr []uuid.UUID, n int) []uuid.UUID {
	if n == 0 || n == 1 {
		return arr
	}
	var temp []uuid.UUID = make([]uuid.UUID, n)
	var j int = 0

	for i := 0; i < n-1; i++ {
		if arr[i] != arr[i+1] {
			temp[j] = arr[i]
			j += 1
		}
	}
	temp[j] = arr[n-1]
	j += 1

	for i := 0; i < j; i++ {
		arr[i] = temp[i]

	}
	arr = arr[:j]
	return arr
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

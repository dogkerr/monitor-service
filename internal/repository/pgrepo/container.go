package pgrepo

import (
	"context"
	"database/sql"
	"dogker/lintang/monitor-service/domain"
	"dogker/lintang/monitor-service/internal/repository/pgrepo/queries"
	"dogker/lintang/monitor-service/pkg/postgres"
	"fmt"
	"time"

	gofrsuuid "github.com/gofrs/uuid"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ContainerRepository struct {
	db *postgres.Postgres
}

func NewContainerRepo(db *postgres.Postgres) *ContainerRepository {
	return &ContainerRepository{db}
}

func (r *ContainerRepository) GetAllUserContainer(ctx context.Context, userID string) (*[]domain.Container, error) {
	q := queries.New(r.db.Pool)
	userUUID, err := gofrsuuid.FromString(userID)
	if err != nil {
		zap.L().Error("uuid fromString", zap.Error(err), zap.String("userUUID", userID))
		return nil, err
	}

	ctrs, err := q.GetAllUserContainer(ctx, uuid.UUID(userUUID))
	if err != nil {
		zap.L().Error("GetAllUserContainer", zap.Error(err), zap.String("userID", userID))
		return nil, err
	}

	var res []domain.Container
	for _, ctr := range ctrs {
		var cLife domain.ContainerLifecycle
		var clifeStatus domain.ContainerStatus
		clifeStatus = domain.STOP
		if ctr.Lifecyclestatus.ContainerStatus == queries.ContainerStatusRUN {
			clifeStatus = domain.RUN
		}
		cLife = domain.ContainerLifecycle{
			ID:        ctr.Lifecycleid.UUID,
			StartTime: ctr.Lifecyclestarttime.Time,
			StopTime:  ctr.Lifecyclestoptime.Time,
			Replica:   uint64(ctr.Lifecyclereplica.Int32),
			Status:    clifeStatus,
		}

		if (len(res) > 0 && res[len(res)-1].ID != ctr.ID) || len(res) == 0 {
			var newCl []domain.ContainerLifecycle

			var terminatedtime time.Time
			var publicPort int

			if ctr.TerminatedTime.Valid {
				terminatedtime = ctr.TerminatedTime.Time
			}
			if ctr.PublicPort.Valid {
				publicPort = int(ctr.PublicPort.Int32)
			}

			var ctrStatus domain.ServiceStatus
			ctrStatus = domain.ServiceStopped
			if ctr.Status == queries.ServiceStatusRUN {
				ctrStatus = domain.ServiceRun
			}

			res = append(res, domain.Container{
				ID:                  ctr.ID,
				UserID:              ctr.UserID,
				Image:               ctr.Image,
				Status:              ctrStatus,
				Name:                ctr.Name,
				ContainerPort:       int(ctr.ContainerPort),
				PublicPort:          int(publicPort),
				CreatedTime:         ctr.CreatedTime,
				ServiceID:           ctr.ServiceID,
				TerminatedTime:      terminatedtime,
				ContainerLifecycles: append(newCl, cLife),
			})
		} else {
			res[len(res)-1].ContainerLifecycles = append(res[len(res)-1].ContainerLifecycles,
				cLife,
			)
		}
	}

	return &res, nil
}

func (r *ContainerRepository) Get(ctx context.Context, serviceID string) (*domain.Container, error) {
	q := queries.New(r.db.Pool)

	ctrs, err := q.GetContainer(ctx, serviceID)
	if err != nil {
		if err == sql.ErrNoRows {
			zap.L().Debug("GetContainer (containerRepository)", zap.Error(err), zap.String("serviceID", serviceID))

			return nil, domain.WrapErrorf(err, domain.ErrNotFound, "container dengan id: "+serviceID+" tidak ada di database")
		}
		zap.L().Error("GetContainer (containerRepository)", zap.Error(err), zap.String("serviceID", serviceID))
		return nil, domain.WrapErrorf(err, domain.ErrInternalServerError, "internal server error")
	}
	var res domain.Container
	for _, ctr := range ctrs {
		var clifeStatus domain.ContainerStatus
		clifeStatus = domain.STOP
		if ctr.Lifecyclestatus.ContainerStatus == queries.ContainerStatusRUN {
			clifeStatus = domain.RUN
		}

		cLife := domain.ContainerLifecycle{
			ID:          ctr.Lifeid.UUID,
			ContainerID: ctr.ID,
			StartTime:   ctr.Lifecyclestarttime.Time,
			StopTime:    ctr.Lifecyclestoptime.Time,
			Replica:     uint64(ctr.Lifecyclereplica.Int32),
			Status:      clifeStatus,
		}

		if res.Name == "" {
			var newCl []domain.ContainerLifecycle
			var publicPort int
			var terminatedtime time.Time
			if ctr.PublicPort.Valid {
				publicPort = int(ctr.PublicPort.Int32)
			}
			if ctr.TerminatedTime.Valid {
				terminatedtime = ctr.TerminatedTime.Time
			}

			var ctrStatus domain.ServiceStatus
			ctrStatus = domain.ServiceStopped
			if ctr.Status == queries.ServiceStatusRUN {
				ctrStatus = domain.ServiceRun
			}
			res = domain.Container{
				ID:                  ctr.ID,
				UserID:              ctr.UserID,
				Image:               ctr.Image,
				Status:              ctrStatus,
				Name:                ctr.Name,
				ContainerPort:       int(ctr.ContainerPort),
				PublicPort:          publicPort,
				CreatedTime:         ctr.CreatedTime,
				ServiceID:           serviceID,
				TerminatedTime:      terminatedtime,
				ContainerLifecycles: append(newCl, cLife),
			}
		} else {
			res.ContainerLifecycles = append(res.ContainerLifecycles,
				cLife,
			)
		}
	}
	return &res, nil
}

// GetSpecificConatainerMetrics
func (r *ContainerRepository) GetSpecificConatainerMetrics(ctx context.Context, ctrID string) (*domain.Metric, error) {
	q := queries.New(r.db.Pool)
	ctrUUID, err := gofrsuuid.FromString(ctrID)
	if err != nil {
		zap.L().Error("uuid fromString", zap.Error(err), zap.String("ctrID", ctrID))
		return nil, err
	}

	metrics, err := q.GetSpecificContainerMetrics(ctx, uuid.UUID(ctrUUID))
	if err != nil {
		zap.L().Error("GetSpecificContainerMetrics", zap.Error(err), zap.String("ctrID", ctrID))
		return nil, err
	}

	// var metr domain.Metric
	var metr queries.GetSpecificContainerMetricsRow
	if len(metrics) != 0 {
		metr = qSortWaktu(metrics) // get last all user container metrics (berdasarkan last date) (kalau enag row di tabel metrics)
		// kalau gak ada metrics di tabel row berarti value dari metrics metr 0 semua
	}

	return &domain.Metric{
		CpuUsage:            float32(metr.Cpus),
		MemoryUsage:         float32(metr.Memory),
		NetworkIngressUsage: float32(metr.NetworkIngress),
		NetworkEgressUsage:  float32(metr.NetworkEgress),
	}, nil
}

type containerWithItsReplica struct {
	CtrID     uuid.UUID
	Replica   int
	ServiceID string
}

// get containerID yang sebelumnya pernah diproses sama service SendDownInstanceToContainerServiceAndMailingService
func (r *ContainerRepository) GetProcessedContainers(ctx context.Context, serviceIDs []string, downContainersReplica map[string]int) ([]string, []uuid.UUID, error) {
	q := queries.New(r.db.Pool)

	rows, err := q.GetContainerByServiceIDs(ctx, serviceIDs) // dapetin replica sama serviceId sama ctrID
	if err != nil {
		if err != nil {
			zap.L().Error(" q.GetContainerByServiceIDs(ctx, serviceIDs) (GetProcessedContainers) (ContainerRepositort)", zap.Error(err))
			return nil, nil, err
		}
	}

	var containerIDs []uuid.UUID
	for i, _ := range rows {
		containerIDs = append(containerIDs, rows[i].ID)
	}

	processedTerminatedCtrIDs, err := q.GetProcessedContainers(ctx, containerIDs)
	if err != nil {
		zap.L().Debug("q.GetProcessedContainers(ctx, serviceUUIDs) (GetProcessedContainers)", zap.Error(err))
		return nil, nil, err
	}

	var processedCtrSet = make(map[uuid.UUID]bool) // bikin set isinya processed containerIDs
	for i, _ := range processedTerminatedCtrIDs {
		if time.Now().Sub(processedTerminatedCtrIDs[i].DownTime) > 1*time.Minute {
			// ini mastiin kalau swarm service down , dalam 1 menit cuma bisa kirim 1 email doang
			// jadi menit kedua kalau swarm service down lagi ya kirim email lagi , begitu seterus nya utk menit menit berikutnya (selisih 1 menit)
			processedCtrSet[processedTerminatedCtrIDs[i].ContainerID] = true
		}
	}

	// mastiin lastestDownServiceWithItsReplica hanya berisi swarm service yang down (replica yang down >= jumlah replica swarm service)
	var lastestDownServiceWithItsReplica []containerWithItsReplica //

	var ctrSet = make(map[uuid.UUID]bool) // container set
	increasingSortedRowsByStartTime := qSortLastReplica(rows)
	latestRowsByStartTime := reverseInplaceCtrArray(increasingSortedRowsByStartTime, 0, len(increasingSortedRowsByStartTime)-1) // latest container replicas

	// mastiin lastestDownServiceWithItsReplica hanya berisi swarm service yang down (replica yg down >= jumlah replica swarm service)
	for i, _ := range latestRowsByStartTime {
		if !ctrSet[latestRowsByStartTime[i].ID] {
			// buatmastiin  latestCtrReplica adalah yang replica ctr yang paling terbaru
			latestCtrReplica := int(latestRowsByStartTime[i].Lifecyclereplica.Int32)
			downReplica := downContainersReplica[latestRowsByStartTime[i].ServiceID]
			if downReplica >= latestCtrReplica {
				// kalau replica yang down  >= jumlah replica servicenya append ke array lastestDownServiceWithItsReplica
				lastestDownServiceWithItsReplica = append(lastestDownServiceWithItsReplica, containerWithItsReplica{CtrID: latestRowsByStartTime[i].ID,
					Replica: int(latestRowsByStartTime[i].Lifecyclereplica.Int32), ServiceID: latestRowsByStartTime[i].ServiceID})
			}

		} else {
			ctrSet[latestRowsByStartTime[i].ID] = true
		}
	}

	// filter down servic dg  service yang sebelumnya belum pernah diproses
	for i, _ := range lastestDownServiceWithItsReplica {
		if processedCtrSet[lastestDownServiceWithItsReplica[i].CtrID] {
			// kalau row ctrID ada di processedCTRID -> delete inplace element ke i di array
			lastestDownServiceWithItsReplica[i] = lastestDownServiceWithItsReplica[len(lastestDownServiceWithItsReplica)-1]
			lastestDownServiceWithItsReplica = lastestDownServiceWithItsReplica[:len(lastestDownServiceWithItsReplica)-1]
		}
	}

	var notProcessedServiceIDs []string // berisi serviceID yang sebelumnya belum pernah diproses
	var notProcessedCtrIDs []uuid.UUID  // berisi ctrID yang sebelumnya belum pernah diproses
	for i, _ := range lastestDownServiceWithItsReplica {
		notProcessedServiceIDs = append(notProcessedServiceIDs, lastestDownServiceWithItsReplica[i].ServiceID)
		notProcessedCtrIDs = append(notProcessedCtrIDs, lastestDownServiceWithItsReplica[i].CtrID)
	}

	return notProcessedServiceIDs, notProcessedCtrIDs, nil
}

func (r *ContainerRepository) GetSwarmServicesDetail(ctx context.Context, serviceIDs []string) ([]domain.CommonLabelsMailing, error) {
	q := queries.New(r.db.Pool)

	swarmServicesDetail, err := q.GetSwarmServiceDetailByServiceIDs(ctx, serviceIDs)
	if err != nil {
		if err == sql.ErrNoRows {
			return []domain.CommonLabelsMailing{}, nil
		}
		zap.L().Error("q.GetSwarmServiceDetailByServiceIDs (GetSwarmServicesDetail) (ContainerRepository)", zap.Error(err))
		return nil, domain.WrapErrorf(err, domain.ErrInternalServerError, domain.MessageInternalServerError)
	}
	var commonLabels []domain.CommonLabelsMailing
	for i, _ := range swarmServicesDetail {
		commonLabels = append(commonLabels, domain.CommonLabelsMailing{
			Alertname: fmt.Sprintf("swarm service %s down", swarmServicesDetail[i].ServiceID ),
			ContainerSwarmServiceID: swarmServicesDetail[i].ServiceID,
			ContainerDockerSwarmServiceName: swarmServicesDetail[i].Name,
			ContainerLabelUserID: swarmServicesDetail[i].UserID.String(),
		})
	}
	return commonLabels, nil 
}

// insert accidentally terminated container  (container yg mati bukan karena dimattiin user lewat api container-service)
func (r *ContainerRepository) InsertTerminatedContainer(ctx context.Context, containerID string) error {
	q := queries.New(r.db.Pool)

	containerUUID, err := gofrsuuid.FromString(containerID)
	if err != nil {
		zap.L().Error("gofrsuuid.FromString(containerID) (InsertTerminatedContainer) (ContainerRepository)", zap.Error(err))
		return err
	}
	err = q.InsertTerminatedContainer(ctx, queries.InsertTerminatedContainerParams{ContainerID: uuid.UUID(containerUUID), DownTime: time.Now()})
	if err != nil {
		zap.L().Error("q.InsertTerminatedContainer(ctx, uuid.UUID(containerUUID)) (InsertTerminatedContainer) (ContainerRepository)", zap.Error(err),
			zap.String("containerID", containerID))
		return err
	}
	return nil
}

func reverseInplaceCtrArray(arr []queries.GetContainerByServiceIDsRow, start int, end int) []queries.GetContainerByServiceIDsRow {
	for start < end {
		temp := arr[start]
		arr[start] = arr[end]
		arr[end] = temp
		start++
		end--
	}
	return arr
}

// / sorting containerlifecycle berdassarkan waktu distartnya dari terdahulu ke terkini
func qSortLastReplica(arr []queries.GetContainerByServiceIDsRow) []queries.GetContainerByServiceIDsRow {
	var recurse func(left int, right int)
	var partition func(left int, right int, pivot int) int

	partition = func(left int, right int, pivot int) int {
		v := arr[pivot]
		right--
		arr[pivot], arr[right] = arr[right], arr[pivot]

		for i := left; i < right; i++ {
			if arr[i].StartTime.Time.Unix() <= v.StartTime.Time.Unix() {
				arr[i], arr[left] = arr[left], arr[i]
				left++
			}
		}

		arr[left], arr[right] = arr[right], arr[left]
		return left
	}

	recurse = func(left int, right int) {
		if left < right {
			pivot := (right + left) / 2
			pivot = partition(left, right, pivot)
			recurse(left, pivot)
			recurse(pivot+1, right)
		}
	}

	recurse(0, len(arr))
	return arr
}

func qSortWaktu(arr []queries.GetSpecificContainerMetricsRow) queries.GetSpecificContainerMetricsRow {

	var recurse func(left int, right int)
	var partition func(left int, right int, pivot int) int

	partition = func(left int, right int, pivot int) int {
		v := arr[pivot]
		right--
		arr[pivot], arr[right] = arr[right], arr[pivot]

		for i := left; i < right; i++ {
			if arr[i].CreatedTime.Unix() <= v.CreatedTime.Unix() {
				arr[i], arr[left] = arr[left], arr[i]
				left++
			}
		}

		arr[left], arr[right] = arr[right], arr[left]
		return left
	}

	recurse = func(left int, right int) {
		if left < right {
			pivot := (right + left) / 2
			pivot = partition(left, right, pivot)
			recurse(left, pivot)
			recurse(pivot+1, right)
		}
	}

	recurse(0, len(arr))
	return arr[len(arr)-1]
}

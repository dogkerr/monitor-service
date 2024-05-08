package pgrepo

import (
	"context"
	"database/sql"
	"dogker/lintang/monitor-service/domain"
	"dogker/lintang/monitor-service/internal/repository/pgrepo/queries"
	"dogker/lintang/monitor-service/pkg/postgres"
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

			var ctrStatus domain.ContainerStatus
			ctrStatus = domain.STOP
			if ctr.Status == queries.ContainerStatusRUN {
				ctrStatus = domain.RUN
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

			var ctrStatus domain.ContainerStatus
			ctrStatus = domain.STOP
			if ctr.Status == queries.ContainerStatusRUN {
				ctrStatus = domain.RUN
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

	metr, err := q.GetSpecificContainerMetrics(ctx, uuid.UUID(ctrUUID))
	if err != nil {
		zap.L().Error("GetSpecificContainerMetrics", zap.Error(err), zap.String("ctrID", ctrID))
		return nil, err
	}
	return &domain.Metric{
		CpuUsage:            float32(metr.Cpus),
		MemoryUsage:         float32(metr.Memory),
		NetworkIngressUsage: float32(metr.NetworkIngress),
		NetworkEgressUsage:  float32(metr.NetworkEgress),
	}, nil
}

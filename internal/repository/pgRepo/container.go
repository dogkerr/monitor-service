package pgRepo

import (
	"context"
	"database/sql"
	"dogker/lintang/monitor-service/domain"
	"dogker/lintang/monitor-service/pkg/postgres"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ContainerRepository struct {
	db *postgres.Postgres
}

func NewContainerRepo(db *postgres.Postgres) *ContainerRepository {
	return &ContainerRepository{db}
}


// query ke postgres buat dapetin semua container yang dimiliki user
func (r *ContainerRepository) GetAllUserContainer(ctx context.Context, userId string) (*[]domain.Container, error) {
	rows, err := r.db.Pool.QueryContext(ctx, `SELECT c.id, c.user_id, c.image_url, c.status, c.name, c.container_port, c.public_port, c.created_time,c.serviceId, c.terminated_time,
			cl.id as lifecycleId, cl.start_time as lifecycleStartTime, cl.stop_time as lifecycleStopTime, 
			cl.replica as lifecycleReplica, cl.status FROM containers c  LEFT JOIN container_lifecycles cl ON cl.container_id=c.id
			WHERE c.user_id=$1`, userId)
	if err != nil {
		zap.L().Debug(fmt.Sprintf("r.db.Pool.QueryContext  user %s", userId))
		return nil, err
	}

	defer rows.Close()

	var res []domain.Container

	for rows.Next() {
		var containerId uuid.UUID
		var userId uuid.UUID
		var imageUrl string
		var status domain.Status
		var name string
		var containerPort int
		var publicPort int
		var createdTime time.Time
		var terminatedTimeNull sql.NullTime
		var terminatedTime time.Time
		var serviceIdNull sql.NullString
		var serviceId string

		var ctrStatus string

		var cLifeId uuid.UUID
		var lStartTime time.Time
		var lStopTime time.Time
		var lReplica uint64
		var clifeStatus domain.Status
		var clStatus string

		var cLife domain.ContainerLifecycle

		if err := rows.Scan(&containerId, &userId, &imageUrl, &ctrStatus, &name, &containerPort, &publicPort, &createdTime, &serviceIdNull, &terminatedTimeNull,
			&cLifeId, &lStartTime, &lStopTime, &lReplica, &clStatus); err != nil {

			zap.L().Error("rows.Scan", zap.Error(err), zap.String("userId", userId.String()))
			return nil, err
		}

		if serviceIdNull.Valid {
			serviceId = serviceIdNull.String
		}
		if terminatedTimeNull.Valid {
			terminatedTime = terminatedTimeNull.Time
		}
		if ctrStatus == "RUN" {
			status = domain.RUN
		} else {
			status = domain.STOPPED
		}

		if clStatus == "RUN" {
			clifeStatus = domain.RUN
		} else {
			clifeStatus = domain.STOPPED
		}

		cLife = domain.ContainerLifecycle{
			ID:        cLifeId,
			StartTime: lStartTime,
			StopTime:  lStopTime,
			Replica:   lReplica,
			Status:    clifeStatus,
		}

		if (len(res) > 0 && res[len(res)-1].ID != containerId) || len(res) == 0 {

			var newCl []domain.ContainerLifecycle
			res = append(res, domain.Container{
				ID:                  containerId,
				UserId:              userId,
				ImageUrl:            imageUrl,
				Status:              status,
				Name:                name,
				ContainerPort:       containerPort,
				PublicPort:          publicPort,
				CreatedTime:         createdTime,
				ServiceId:           serviceId,
				TerminatedTime:      terminatedTime,
				ContainerLifecycles: append(newCl, cLife),
			})
		} else {

			res[len(res)-1].ContainerLifecycles = append(res[len(res)-1].ContainerLifecycles,
				cLife,
			)

		}

	}

	if len(res) == 0 {

		return nil, domain.ErrNotFound
	}
	return &res, nil
}

// Get mendapatkan container bedasarkan idnya
func (r *ContainerRepository) Get(ctx context.Context, cId string) (*domain.Container, error) {
	rows, err := r.db.Pool.QueryContext(ctx, `SELECT c.id, c.user_id, c.image_url, c.status, c.name, c.container_port, c.public_port,
			 c.terminated_time,
			c.created_time, c.serviceId, cl.id as lifeId, cl.start_time, cl.stop_time, cl.status, cl.replica 
			FROM containers c LEFT JOIN container_lifecycles cl ON cl.container_id=c.id
			WHERE c.id=$1  `, cId)

	if err != nil {
		zap.L().Error(fmt.Sprintf("r.db.Pool.QueryContext "), zap.String("containerId", cId))
	}

	var res domain.Container // container dg id cId

	defer rows.Close()

	for rows.Next() {
		var containerId uuid.UUID
		var userId uuid.UUID
		var imageUrl string
		var status domain.Status
		var name string
		var containerPort int
		var publicPort int
		var createdTime time.Time
		var terminatedTimeNull sql.NullTime
		var terminatedTime time.Time
		var serviceIdNull sql.NullString
		var serviceId string

		var ctrStatus string

		var cLifeId uuid.UUID
		var lStartTime time.Time
		var lStopTime time.Time
		var lReplica uint64
		var clifeStatus domain.Status
		var clStatus string

		var cLife domain.ContainerLifecycle

		if err := rows.Scan(&containerId, &userId, &imageUrl, &ctrStatus, &name, &containerPort, &publicPort, &createdTime, &serviceIdNull, &terminatedTimeNull,
			&cLifeId, &lStartTime, &lStopTime, &lReplica, &clStatus); err != nil {

			zap.L().Error("rows.Scan", zap.Error(err), zap.String("userId", userId.String()))
			return nil, err
		}

		if serviceIdNull.Valid {
			serviceId = serviceIdNull.String
		}
		if terminatedTimeNull.Valid {
			terminatedTime = terminatedTimeNull.Time
		}
		if ctrStatus == "RUN" {
			status = domain.RUN
		} else {
			status = domain.STOPPED
		}

		if clStatus == "RUN" {
			clifeStatus = domain.RUN
		} else {
			clifeStatus = domain.STOPPED
		}

		cLife = domain.ContainerLifecycle{
			ID:        cLifeId,
			StartTime: lStartTime,
			StopTime:  lStopTime,
			Replica:   lReplica,
			Status:    clifeStatus,
		}

		if res.Name == "" {
			var newCl []domain.ContainerLifecycle
			res = domain.Container{
				ID:                  containerId,
				UserId:              userId,
				ImageUrl:            imageUrl,
				Status:              status,
				Name:                name,
				ContainerPort:       containerPort,
				PublicPort:          publicPort,
				CreatedTime:         createdTime,
				ServiceId:           serviceId,
				TerminatedTime:      terminatedTime,
				ContainerLifecycles: append(newCl, cLife),
			}
		}else {
			res.ContainerLifecycles = append(res.ContainerLifecycles,
				cLife,
			)
		}
	}
	if  res.Name == "" {
		zap.L().Debug("container not found", zap.String("containerId", cId))
		return nil, domain.ErrNotFound
	}

	return &res, nil 

}

// func (r *ContainerRepository)

/*
Fetch(ctx context.Context, created_time string, limit int, page int) (res []domain.Container, pagination repository.Pagination, err error)
	GetById(ctx context.Context, id uuid.UUID) (domain.Container, error)
	SearchContainer(ctx context.Context, query string) (domain.Container, error)
	Update(ctx context.Context, c *domain.Container) error
	Store(ctx context.Context, c *domain.Container) error
	Delete(ctx context.Context, id uuid.UUID) error


*/

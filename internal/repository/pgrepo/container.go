package pgrepo

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
func (r *ContainerRepository) GetAllUserContainer(ctx context.Context, userID string) (*[]domain.Container, error) {
	rows, err := r.db.Pool.QueryContext(ctx, `SELECT c.id, c.user_id, c.image_url, c.status, c.name, c.container_port, c.public_port, c.created_time,c.service_id, c.terminated_time,
			cl.id as lifecycleId, cl.start_time as lifecycleStartTime, cl.stop_time as lifecycleStopTime, 
			cl.replica as lifecycleReplica, cl.status FROM containers c  LEFT JOIN container_lifecycles cl ON cl.container_id=c.id
			WHERE c.user_id=$1`, userID)
	if err != nil {
		zap.L().Debug(fmt.Sprintf("r.db.Pool.QueryContext  user %s", userID))
		return nil, err
	}

	defer rows.Close()

	var res []domain.Container

	for rows.Next() {
		var containerID uuid.UUID
		var userID uuid.UUID
		var imageURL string
		var status domain.Status
		var name string
		var containerPort int
		var publicPort int
		var createdTime time.Time
		var terminatedTimeNull sql.NullTime
		var terminatedTime time.Time
		var serviceIDNull sql.NullString
		var serviceID string

		var ctrStatus string

		var cLifeID uuid.UUID

		var lStartTimeNull sql.NullTime
		var lStopTimeNull sql.NullTime

		var lStartTime time.Time
		var lStopTime time.Time
		var lReplica uint64
		var clifeStatus domain.Status
		var clStatus string

		var cLife domain.ContainerLifecycle

		if err := rows.Scan(&containerID, &userID, &imageURL, &ctrStatus, &name, &containerPort, &publicPort, &createdTime, &serviceIDNull, &terminatedTimeNull,
			&cLifeID, &lStartTimeNull, &lStopTimeNull, &lReplica, &clStatus); err != nil {
			zap.L().Error("rows.Scan", zap.Error(err), zap.String("userID", userID.String()))
			return nil, err
		}
		if lStartTimeNull.Valid {
			lStartTime = lStartTimeNull.Time
		}
		if lStopTimeNull.Valid {
			lStopTime = lStopTimeNull.Time
		}

		if serviceIDNull.Valid {
			serviceID = serviceIDNull.String
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
			ID:        cLifeID,
			StartTime: lStartTime,
			StopTime:  lStopTime,
			Replica:   lReplica,
			Status:    clifeStatus,
		}

		if (len(res) > 0 && res[len(res)-1].ID != containerID) || len(res) == 0 {
			var newCl []domain.ContainerLifecycle
			res = append(res, domain.Container{
				ID:                  containerID,
				UserID:              userID,
				ImageURL:            imageURL,
				Status:              status,
				Name:                name,
				ContainerPort:       containerPort,
				PublicPort:          publicPort,
				CreatedTime:         createdTime,
				ServiceID:           serviceID,
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
func (r *ContainerRepository) Get(ctx context.Context, serviceID string) (*domain.Container, error) {
	rows, err := r.db.Pool.QueryContext(ctx, `SELECT c.id, c.user_id, c.image_url, c.status, c.name, c.container_port, c.public_port,c.created_time,
	c.service_id,c.terminated_time, cl.id as lifeId, cl.start_time, cl.stop_time, cl.replica , cl.status
	FROM containers c LEFT JOIN container_lifecycles cl ON cl.container_id=c.id
	WHERE c.service_id=$1`, serviceID)

	if err != nil {
		zap.L().Error("r.db.Pool.QueryContext ", zap.String("containerID", serviceID))
		return nil, err
	}

	var res domain.Container // container dg id serviceID

	defer rows.Close()

	for rows.Next() {
		var containerID uuid.UUID
		var userID uuid.UUID
		var imageURL string
		var status domain.Status
		var name string
		var containerPort int
		var publicPort int
		var createdTime time.Time
		var terminatedTimeNull sql.NullTime
		var terminatedTime time.Time
		var serviceIDNull sql.NullString
		var serviceID string

		var ctrStatus string

		var cLifeID uuid.UUID
		var lStartTime time.Time
		var lStopTime time.Time
		var lReplica uint64
		var clifeStatus domain.Status
		var clStatus string

		var cLife domain.ContainerLifecycle

		if err := rows.Scan(&containerID, &userID, &imageURL, &ctrStatus, &name, &containerPort, &publicPort, &createdTime, &serviceIDNull, &terminatedTimeNull,
			&cLifeID, &lStartTime, &lStopTime, &lReplica, &clStatus); err != nil {
			zap.L().Error("rows.Scan", zap.Error(err), zap.String("userID", userID.String()))
			return nil, err
		}

		if serviceIDNull.Valid {
			serviceID = serviceIDNull.String
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
			ID:        cLifeID,
			StartTime: lStartTime,
			StopTime:  lStopTime,
			Replica:   lReplica,
			Status:    clifeStatus,
		}

		if res.Name == "" {
			var newCl []domain.ContainerLifecycle
			res = domain.Container{
				ID:                  containerID,
				UserID:              userID,
				ImageURL:            imageURL,
				Status:              status,
				Name:                name,
				ContainerPort:       containerPort,
				PublicPort:          publicPort,
				CreatedTime:         createdTime,
				ServiceID:           serviceID,
				TerminatedTime:      terminatedTime,
				ContainerLifecycles: append(newCl, cLife),
			}
		} else {
			res.ContainerLifecycles = append(res.ContainerLifecycles,
				cLife,
			)
		}
	}
	if res.Name == "" {
		zap.L().Debug("container not found", zap.String("containerID", serviceID))
		return nil, domain.ErrNotFound
	}

	return &res, nil
}

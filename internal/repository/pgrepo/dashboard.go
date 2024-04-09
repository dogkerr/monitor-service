package pgrepo

import (
	"context"
	"database/sql"
	"dogker/lintang/monitor-service/domain"
	"dogker/lintang/monitor-service/pkg/postgres"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type DashboardRepository struct {
	db *postgres.Postgres
}

func NewDashboardRepo(db *postgres.Postgres) *DashboardRepository {
	return &DashboardRepository{db}
}

// dapetin dashboard by user_id dan db_type
func (r *DashboardRepository) GetByUserIDAndType(ctx context.Context, userID, dbType string) (*domain.Dashboard, error) {
	var id, owner uuid.UUID
	var uid, dbTypePrs string
	userIDDB, _ := uuid.Parse(userID)
	row := r.db.Pool.QueryRowContext(ctx, "SELECT * FROM dashboards WHERE owner=$1 and  db_type=$2", userIDDB, dbType)
	if err := row.Scan(&id, &uid, &owner, &dbTypePrs); err != nil {
		if err == sql.ErrNoRows {
			zap.L().Debug("dashboard not found", zap.String("userID", userID))
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &domain.Dashboard{
		Id:    id.String(),
		Owner: owner.String(),
		Type:  dbTypePrs,
		Uid:   uid,
	}, nil
}

// create dashboard dg type ('log'/'monitoring')
func (r *DashboardRepository) CreateDashboard(ctx context.Context, dashboard *domain.Dashboard) error {

	_, err := r.db.Pool.ExecContext(ctx, "INSERT INTO dashboards(uid, owner, db_type)  VALUES($1, $2, $3)",
		dashboard.Uid, dashboard.Owner, dashboard.Type)
	if err != nil {
		zap.L().Error("r.db.Pool.ExecContext(ctx, 'INSERT INTO dashboards(uid, owner, db_type)  VALUES(?, ?, ?)'")
	}
	return err
}

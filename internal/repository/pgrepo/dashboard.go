package pgrepo

import (
	"context"
	"database/sql"
	"dogker/lintang/monitor-service/domain"
	"dogker/lintang/monitor-service/internal/repository/pgrepo/queries"
	"dogker/lintang/monitor-service/pkg/postgres"
	"errors"
	"fmt"

	gofrsuuid "github.com/gofrs/uuid"
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

func (r *DashboardRepository) GetDashboardOwner(ctx context.Context, dashboardUID string, userID string) error {
	q := queries.New(r.db.Pool)
	dashboardUUID, err := gofrsuuid.FromString(dashboardUID)
	if err != nil {
		zap.L().Error("uuid fromString", zap.Error(err), zap.String("dashboardUID", dashboardUID))
		return err
	}
	dbOwner, err := q.GetContainerOwnerByID(ctx, dashboardUUID.String())
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.WrapErrorf(err, domain.ErrNotFound, fmt.Sprintf(`dashboard %s not found`, dashboardUID))
		}
	}
	if dbOwner.Owner.String() != userID {
		return domain.WrapErrorf(errors.New("you are not authorized to access this grafana dashboard"), domain.ErrUnauthorized, "you are not authorized to access this grafana dashboard")
	}

	return nil
}

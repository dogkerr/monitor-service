package pgrepo

import (
	"context"
	"dogker/lintang/monitor-service/domain"
	"dogker/lintang/monitor-service/pkg/postgres"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type UserRepository struct {
	db *postgres.Postgres
}

func NewUserRepo(db *postgres.Postgres) *UserRepository {
	return &UserRepository{db}
}

// buat dapetin semua users di postgres
func (r *UserRepository) GetAllUsers(ctx context.Context) (*[]domain.User, error) {
	rows, err := r.db.Pool.QueryContext(ctx, `SELECT u.id, u.username, u.email FROM users u`)
	if err != nil {
		zap.L().Error("r.db.Pool.QueryContext  select all users: ", zap.Error(err))
		return nil, err
	}

	defer rows.Close()

	var res []domain.User

	for rows.Next() {
		var userID uuid.UUID
		var username string
		var email string

		if err := rows.Scan(&userID, &username, &email); err != nil {
			zap.L().Error("rows.Scan", zap.Error(err), zap.String("userID", userID.String()))
			return nil, err
		}

		res = append(res, domain.User{
			ID:       userID,
			Username: username,
			Email:    email,
		})

	}

	if len(res) == 0 {
		return nil, domain.ErrNotFound
	}

	return &res, nil
}

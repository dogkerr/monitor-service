package postgres

import (
	"context"
	"database/sql"
	"dogker/lintang/monitor-service/config"
	"net/url"

	"time"

	"go.uber.org/zap"
)

type Postgres struct {
	Pool *sql.DB
}

func NewPostgres(cfg *config.Config) *Postgres {
	dsn := url.URL{
		Scheme: cfg.Postgres.PGScheme,
		Host:   cfg.Postgres.PGURL, // "localhost:5432"
		User:   url.UserPassword(cfg.Postgres.Username, cfg.Postgres.Password),
		Path:   cfg.Postgres.PGDB,
	}

	q := dsn.Query()
	q.Add("sslmode", "disable")

	dsn.RawQuery = q.Encode()

	db, err := sql.Open("pgx", dsn.String())
	if err != nil {
		zap.L().Fatal("sql.Open", zap.Error(err))
	}

	db.SetMaxIdleConns(20)
	db.SetMaxOpenConns(250)
	db.SetConnMaxIdleTime(5 * time.Minute)
	db.SetConnMaxLifetime(60 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		zap.L().Fatal("db.PingContext", zap.Error(err))
	}
	return &Postgres{db}
}

func ClosePostgres(pg *sql.DB) error {
	err := pg.Close()
	return err
}

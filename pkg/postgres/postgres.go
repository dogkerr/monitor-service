package postgres

import (
	"context"
	"database/sql"
	"dogker/lintang/monitor-service/config"
	"net/url"

	"go.uber.org/zap"
)


type Postgres struct {
	Pool *sql.DB
}



func NewPostgres(cfg *config.Config) *Postgres {
	dsn := url.URL{
		Scheme: "postgres",
		Host: "localhost:5432",
		User: url.UserPassword("postgres", "pass"),
		Path: "dogker",
	}

	q := dsn.Query()
	q.Add("sslmode", "disable")

	dsn.RawQuery = q.Encode()

	db, err := sql.Open("pgx", dsn.String())
	if err != nil {
		zap.L().Fatal("sql.Open", zap.Error(err))
	}

	// defer func() {
	// 	_ = db.Close()
	// 	zap.L().Error("postgres closed!")
	// }()

	if err := db.PingContext(context.Background()); err != nil {
		zap.L().Fatal("db.PingContext", zap.Error(err))
	}
	

	return &Postgres{db}
}
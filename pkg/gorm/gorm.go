package gorm

import (
	"dogker/lintang/monitor-service/config"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Gorm struct {
	Pool *gorm.DB
}

func NewGorm(cfg *config.Config) *Gorm {
	dsn := "host=localhost user=" + cfg.Postgres.Username + " password=" + cfg.Postgres.Password + " dbname=dogker port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		zap.L().Fatal("Error NewGorm", zap.Error(err))
	}
	gorm := &Gorm{
		Pool: db,
	}

	return gorm
}

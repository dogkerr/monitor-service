//go:build wireinject
// +build wireinject

package di

import (
	"dogker/lintang/monitor-service/app/start"
	"dogker/lintang/monitor-service/config"
	"dogker/lintang/monitor-service/pkg/postgres"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)


func InitApp(cfg *config.Config, handler *gin.Engine) *postgres.Postgres {
	wire.Build(
		start.InitHTTPandGRPC,
	)

	return nil
}

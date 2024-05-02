//go:build wireinject
// +build wireinject

package di

import (
	"dogker/lintang/monitor-service/app/start"
	"dogker/lintang/monitor-service/config"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func InitApp(cfg *config.Config, handler *gin.Engine) *start.InitWireApp {
	wire.Build(
		start.InitHTTPandGRPC,
	)

	return nil
}

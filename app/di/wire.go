//go:build wireinject
// +build wireinject

package di

import (
	"dogker/lintang/monitor-service/internal/repository/postgres"
	"dogker/lintang/monitor-service/internal/rest"
	"dogker/lintang/monitor-service/monitor"
	"dogker/lintang/monitor-service/pkg/gorm"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

// var ProviderSet = wire.NewSet(gorm.NewGorm)
var monitorSet = wire.NewSet(
	postgres.NewContainerRepo,
	monitor.NewService,
	wire.Bind(new(monitor.ContainerRepository), new(*postgres.ContainerRepository)),
	wire.Bind(new(rest.MonitorService), new(*monitor.Service)),
)

func InitRouterApi(*gorm.Gorm, *gin.Engine) *gin.RouterGroup {
	wire.Build(
		monitorSet,
		rest.NewRouter,
	)
	return nil
}

//go:build wireinject
// +build wireinject

package di

import (
	"dogker/lintang/monitor-service/config"
	"dogker/lintang/monitor-service/internal/grpc"
	"dogker/lintang/monitor-service/internal/repository/postgres"
	"dogker/lintang/monitor-service/internal/rest"
	"dogker/lintang/monitor-service/internal/webapi"
	"dogker/lintang/monitor-service/monitor"
	"dogker/lintang/monitor-service/pb"
	"dogker/lintang/monitor-service/pkg/gorm"
	"net"

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

func InitRouterApi(*config.Config, *gin.Engine) *gin.RouterGroup {
	wire.Build(
		gorm.NewGorm,
		monitorSet,
		rest.NewRouter,
	)
	return nil
}

var monitorGrpcSet = wire.NewSet(
	webapi.NewPrometheusAPI,
	monitor.NewMonitorServer,
	wire.Bind(new(monitor.PrometheusApi), new(*webapi.PrometheusAPIImpl)),
	wire.Bind(new(pb.MonitorServiceServer), new(*monitor.MonitorServerImpl)),
)

func InitGrpcMonitorApi(promeAddress string, listener net.Listener) error {
	wire.Build(
		monitorGrpcSet,
		grpc.RunGRPCServer,
	)
	return nil
}

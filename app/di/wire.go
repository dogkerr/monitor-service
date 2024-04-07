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
	pgPool "dogker/lintang/monitor-service/pkg/postgres"
	"net"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

// var ProviderSet = wire.NewSet(gorm.NewGorm)
var monitorSet = wire.NewSet(
	postgres.NewContainerRepo,
	monitor.NewService,
	wire.Bind(new(monitor.ContainerRepository), new(*postgres.ContainerRepositoryI)),
	wire.Bind(new(rest.MonitorService), new(*monitor.Service)),
)

func InitRouterApi(*config.Config, *gin.Engine) *gin.RouterGroup {
	wire.Build(
		pgPool.NewPostgres,
		monitorSet,
		rest.NewRouter,
	)
	return nil
}

var monitorGrpcSet = wire.NewSet(
	postgres.NewContainerRepo,
	wire.Bind(new(monitor.ContainerRepository), new(*postgres.ContainerRepositoryI)),
	webapi.NewPrometheusAPI,
	monitor.NewMonitorServer,
	wire.Bind(new(monitor.PrometheusApi), new(*webapi.PrometheusAPIImpl)),
	wire.Bind(new(pb.MonitorServiceServer), new(*monitor.MonitorServerImpl)),
)

func InitGrpcMonitorApi(promeAddress string, listener net.Listener, cfg *config.Config, ) error {
	wire.Build(
		pgPool.NewPostgres,
		monitorGrpcSet,
		grpc.RunGRPCServer,
	)
	return nil
}

// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package di

import (
	"dogker/lintang/monitor-service/config"
	"dogker/lintang/monitor-service/internal/grpc"
	postgres2 "dogker/lintang/monitor-service/internal/repository/postgres"
	"dogker/lintang/monitor-service/internal/rest"
	"dogker/lintang/monitor-service/internal/webapi"
	"dogker/lintang/monitor-service/monitor"
	"dogker/lintang/monitor-service/pb"
	"dogker/lintang/monitor-service/pkg/postgres"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"net"
)

// Injectors from wire.go:

func InitRouterApi(configConfig *config.Config, engine *gin.Engine) *gin.RouterGroup {
	postgresPostgres := postgres.NewPostgres(configConfig)
	containerRepositoryI := postgres2.NewContainerRepo(postgresPostgres)
	service := monitor.NewService(containerRepositoryI)
	routerGroup := rest.NewRouter(engine, service)
	return routerGroup
}

func InitGrpcMonitorApi(promeAddress string, listener net.Listener, cfg *config.Config) error {
	prometheusAPIImpl := webapi.NewPrometheusAPI(promeAddress)
	postgresPostgres := postgres.NewPostgres(cfg)
	containerRepositoryI := postgres2.NewContainerRepo(postgresPostgres)
	monitorServerImpl := monitor.NewMonitorServer(prometheusAPIImpl, containerRepositoryI)
	error2 := grpc.RunGRPCServer(monitorServerImpl, listener)
	return error2
}

// wire.go:

// var ProviderSet = wire.NewSet(gorm.NewGorm)
var monitorSet = wire.NewSet(postgres2.NewContainerRepo, monitor.NewService, wire.Bind(new(monitor.ContainerRepository), new(*postgres2.ContainerRepositoryI)), wire.Bind(new(rest.MonitorService), new(*monitor.Service)))

var monitorGrpcSet = wire.NewSet(postgres2.NewContainerRepo, wire.Bind(new(monitor.ContainerRepository), new(*postgres2.ContainerRepositoryI)), webapi.NewPrometheusAPI, monitor.NewMonitorServer, wire.Bind(new(monitor.PrometheusApi), new(*webapi.PrometheusAPIImpl)), wire.Bind(new(pb.MonitorServiceServer), new(*monitor.MonitorServerImpl)))

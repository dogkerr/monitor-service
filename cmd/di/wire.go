// go:build wireinject
//go:build wireinject
// +build wireinject

package di

import (
	"dogker/lintang/monitor-service/config"
	"dogker/lintang/monitor-service/pkg/rabbitmq"

	"dogker/lintang/monitor-service/internal/grpc"
	"dogker/lintang/monitor-service/internal/repository/pgrepo"
	"dogker/lintang/monitor-service/internal/repository/rabbitmqrepo"
	"dogker/lintang/monitor-service/internal/rest"
	"dogker/lintang/monitor-service/internal/webapi"
	"dogker/lintang/monitor-service/monitor"
	"dogker/lintang/monitor-service/pb"
	"dogker/lintang/monitor-service/pkg/postgres"

	"github.com/google/wire"
)

var ProviderSet wire.ProviderSet = wire.NewSet(
	monitor.NewService,
	webapi.NewWebAPI,
	grpc.NewContainerClient,
	pgrepo.NewContainerRepo,
	webapi.NewGrafanaAPI,
	pgrepo.NewDashboardRepo,
	pgrepo.NewUserRepo,
	webapi.NewPrometheusAPI,
	rabbitmqrepo.NewMonitorMQ,
	
	wire.Bind(new(monitor.MailingWebAPI), new(*webapi.MailingWebAPI)),
	wire.Bind(new(monitor.ContainerServiceClient), new(*grpc.ContainerClient)),
	wire.Bind(new(rest.MonitorService), new(*monitor.Service)),
	wire.Bind(new(monitor.ContainerRepository), new(*pgrepo.ContainerRepository)),
	wire.Bind(new(monitor.GrafanaAPI), new(*webapi.GrafanaAPI)),
	wire.Bind(new(monitor.DashboardRepository), new(*pgrepo.DashboardRepository)),
	wire.Bind(new(monitor.UserRepository), new(*pgrepo.UserRepository)),
	wire.Bind(new(monitor.PrometheusAPI), new(*webapi.PrometheusAPI)),
	wire.Bind(new(monitor.MonitorMQ), new(*rabbitmqrepo.MonitorMQ)),
)

var ProviderSetMonitorGrpcSet wire.ProviderSet = wire.NewSet(
	monitor.NewMonitorServer,
	webapi.NewPrometheusAPI,
	pgrepo.NewContainerRepo,
	rabbitmqrepo.NewMonitorMQ,
	wire.Bind(new(monitor.PrometheusAPI), new(*webapi.PrometheusAPI)),
	wire.Bind(new(monitor.ContainerRepository), new(*pgrepo.ContainerRepository)),
	wire.Bind(new(pb.MonitorServiceServer), new(*monitor.MonitorServerImpl)),
	wire.Bind(new(monitor.MonitorMQ), new(*rabbitmqrepo.MonitorMQ)),
)

func InitMonitorService(rmq *rabbitmq.RabbitMQ, pgRepo *postgres.Postgres, cfg *config.Config,
) *monitor.Service {

	// wire.Build(
	// 	start.InitHTTPandGRPC,
	// )
	wire.Build(
		ProviderSet,
	)
	return nil
}

func InitMonitorGrpcService(rmq *rabbitmq.RabbitMQ, pgRepo *postgres.Postgres, cfg *config.Config) *monitor.MonitorServerImpl {

	wire.Build(
		ProviderSetMonitorGrpcSet,
	)
	return nil
}

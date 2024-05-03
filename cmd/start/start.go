package start

import (
	"dogker/lintang/monitor-service/config"
	"dogker/lintang/monitor-service/internal/grpc"
	"dogker/lintang/monitor-service/internal/repository/pgrepo"
	"dogker/lintang/monitor-service/internal/repository/rabbitmqrepo"
	"dogker/lintang/monitor-service/internal/rest"
	"dogker/lintang/monitor-service/internal/webapi"
	"dogker/lintang/monitor-service/monitor"
	"dogker/lintang/monitor-service/pkg/postgres"
	"dogker/lintang/monitor-service/pkg/rabbitmq"
	"net"

	grpcClient "google.golang.org/grpc"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type InitWireApp struct {
	PG         *postgres.Postgres
	RMQ        *rabbitmq.RabbitMQ
	GRPCServer *grpcClient.Server
}

func InitHTTPandGRPC(cfg *config.Config, handler *gin.Engine) *InitWireApp {
	// Router

	pg := postgres.NewPostgres(cfg)
	containerRepository := pgrepo.NewContainerRepo(pg)
	grf := webapi.NewGrafanaAPI(cfg, cfg.Grafana.GrafanaFileLoc)
	dbRepo := pgrepo.NewDashboardRepo(pg)
	userDb := pgrepo.NewUserRepo(pg)
	prometheusAPI := webapi.NewPrometheusAPI(cfg.Prometheus.URL)

	rmq := rabbitmq.NewRabbitMQ(cfg)
	mtrMq := rabbitmqrepo.NewMonitorMQ(rmq.Channel)

	service := monitor.NewService(containerRepository, grf, dbRepo, userDb, prometheusAPI, mtrMq)
	rest.NewRouter(handler, service)

	address := cfg.GRPC.URLGrpc
	listener, err := net.Listen("tcp", address)
	if err != nil {
		zap.L().Fatal("cannot start server: ", zap.Error(err))
	}

	// GRPC

	monitorServerImpl := monitor.NewMonitorServer(prometheusAPI, containerRepository)
	grpcServerChan := make(chan *grpcClient.Server)

	go func() {
		err := grpc.RunGRPCServer(monitorServerImpl, listener, grpcServerChan)
		if err != nil {
			zap.L().Fatal("cannot start GRPC  Server", zap.Error(err))
		}

	}()

	var grpcServer = <-grpcServerChan

	return &InitWireApp{
		PG:         pg,
		RMQ:        rmq,
		GRPCServer: grpcServer,
	}
}

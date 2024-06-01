package start

import (
	"dogker/lintang/monitor-service/cmd/di"
	"dogker/lintang/monitor-service/config"
	"dogker/lintang/monitor-service/internal/grpc"
	"dogker/lintang/monitor-service/internal/rest"
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
	rmq := rabbitmq.NewRabbitMQ(cfg)
	// containerRepository := pgrepo.NewContainerRepo(pg)
	// grf := webapi.NewGrafanaAPI(cfg)
	// dbRepo := pgrepo.NewDashboardRepo(pg)
	// userDb := pgrepo.NewUserRepo(pg)
	// prometheusAPI := webapi.NewPrometheusAPI(cfg)

	// mtrMq := rabbitmqrepo.NewMonitorMQ(rmq)

	// service := monitor.NewService(containerRepository, grf, dbRepo, userDb, prometheusAPI, mtrMq)

	// cc, err := grpcClient.NewClient(cfg.GRPC.ContainerURL+"?wait=30s", grpcClient.WithTransportCredentials(insecure.NewCredentials()))
	// if err != nil {
	// 	zap.L().Fatal("grpcClient.NewClient", zap.Error(err ))
	// }

	monitorSvc := di.InitMonitorService(rmq, pg, cfg)
	rest.NewRouter(handler, monitorSvc)

	address := cfg.GRPC.URLGrpc
	listener, err := net.Listen("tcp", address)
	if err != nil {
		zap.L().Fatal("cannot start server: ", zap.Error(err))
	}

	// GRPC

	// monitorServerImpl := monitor.NewMonitorServer(prometheusAPI, containerRepository)
	monitorServerImpl := di.InitMonitorGrpcService(rmq, pg, cfg)
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

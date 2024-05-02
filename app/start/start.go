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

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type InitWireApp struct {
	PG  *postgres.Postgres
	RMQ *rabbitmq.RabbitMQ
}

func InitHTTPandGRPC(cfg *config.Config, handler *gin.Engine) *InitWireApp {
	// Router
	pg := postgres.NewPostgres(cfg)
	containerRepository := pgrepo.NewContainerRepo(pg)
	grf := webapi.NewGrafanaAPI(cfg)
	dbRepo := pgrepo.NewDashboardRepo(pg)

	rmq := rabbitmq.NewRabbitMQ(cfg)

	service := monitor.NewService(containerRepository, grf, dbRepo)
	rest.NewRouter(handler, service)

	address := cfg.GRPC.URLGrpc
	listener, err := net.Listen("tcp", address)
	if err != nil {
		zap.L().Fatal("cannot start server: ", zap.Error(err))
	}

	// GRPC
	prometheusAPI := webapi.NewPrometheusAPI(cfg.Prometheus.URL)
	monitorServerImpl := monitor.NewMonitorServer(prometheusAPI, containerRepository)
	err = grpc.RunGRPCServer(monitorServerImpl, listener)
	if err != nil {
		zap.L().Fatal("cannot start GRPC  Server", zap.Error(err))
	}

	// rabbitMQ task
	_, err = rabbitmqrepo.NewMonitor(rmq.Channel)
	if err != nil {
		zap.L().Fatal("cannot start monitor mq: " + err.Error())
	}

	return &InitWireApp{
		PG:  pg,
		RMQ: rmq,
	}
}

package start

import (
	"dogker/lintang/monitor-service/config"
	"dogker/lintang/monitor-service/internal/grpc"
	"dogker/lintang/monitor-service/internal/repository/pgRepo"
	"dogker/lintang/monitor-service/internal/rest"
	"dogker/lintang/monitor-service/internal/webapi"
	"dogker/lintang/monitor-service/monitor"
	"dogker/lintang/monitor-service/pkg/postgres"
	"net"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func InitHTTPandGRPC(cfg *config.Config, handler *gin.Engine) *postgres.Postgres {
	// Router
	pg := postgres.NewPostgres(cfg)
	containerRepository := pgRepo.NewContainerRepo(pg)
	service := monitor.NewService(containerRepository)
	rest.NewRouter(handler, service)

	address := "0.0.0.0:5001"
	listener, err := net.Listen("tcp", address)
	if err != nil {
		zap.L().Fatal("cannot start server: ", zap.Error(err))
	}

	// GRPC
	prometheusAPI := webapi.NewPrometheusAPI(cfg.Prometheus.Url)
	monitorServerImpl := monitor.NewMonitorServer(prometheusAPI, containerRepository)
	err = grpc.RunGRPCServer(monitorServerImpl, listener)
	if err != nil {
		zap.L().Fatal("cannot start GRPC  Server", zap.Error(err))
	}
	return pg
}

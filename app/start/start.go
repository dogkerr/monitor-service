package start

import (
	"dogker/lintang/monitor-service/config"
	"dogker/lintang/monitor-service/internal/grpc"
	"dogker/lintang/monitor-service/internal/repository/pgrepo"
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
	containerRepository := pgrepo.NewContainerRepo(pg)
	grf := webapi.NewGrafanaAPI(cfg)
	dbRepo := pgrepo.NewDashboardRepo(pg)
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
	return pg
}

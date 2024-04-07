package main

import (
	"dogker/lintang/monitor-service/app/di"
	"dogker/lintang/monitor-service/config"
	"dogker/lintang/monitor-service/internal/rest/middleware"
	"dogker/lintang/monitor-service/pkg/httpserver"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

func init() {
	cfg, err := config.NewConfig()

	if err != nil {
		log.Fatalf("Config error: %s", err)
	}
	//
	if err := middleware.InitLogger(cfg.LogConfig); err != nil {
		fmt.Printf("init logger failed, err:%v\n", err)
		return
	}

}

func main() {
	cfg, err := config.NewConfig()

	// init logger

	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// HTTP Server
	handler := gin.New()
	httpServer := httpserver.New(handler, httpserver.Port("5000"))

	// Router
	di.InitRouterApi(cfg, handler)

	address := "0.0.0.0:5001"
	listener, err := net.Listen("tcp", address)
	if err != nil {
		zap.L().Fatal("cannot start server: ", zap.Error(err))
	}

	// GRPC
	err = di.InitGrpcMonitorApi("http://localhost:9090", listener, cfg)
	if err != nil {
		zap.L().Fatal("cannot start GRPC  Server", zap.Error(err))
	}

	
	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		zap.L().Fatal("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		zap.L().Fatal(fmt.Errorf("app - Run - httpServer.Notify: %w", err).Error())
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		zap.L().Fatal(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err).Error())
	}
}

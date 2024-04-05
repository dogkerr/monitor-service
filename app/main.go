package main

import (
	"dogker/lintang/monitor-service/config"
	"dogker/lintang/monitor-service/internal/repository/postgres"
	"dogker/lintang/monitor-service/internal/rest"
	"dogker/lintang/monitor-service/internal/rest/middleware"
	"dogker/lintang/monitor-service/monitor"
	"dogker/lintang/monitor-service/pkg/gorm"
	"dogker/lintang/monitor-service/pkg/httpserver"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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
	//  databae
	gorm, err := gorm.NewGorm(cfg.Postgres.Username, cfg.Postgres.Password)
	if err != nil {
		log.Fatalf("Database Connection error: %s", err)
	}

	// HTTP Server
	handler := gin.New()
	httpServer := httpserver.New(handler, httpserver.Port("3000"))

	// Prepare Repository
	containerRepo := postgres.NewContainerRepo(gorm.Pool)

	// Build service layer
	mSvc := monitor.NewService(containerRepo)

	// router
	handler.Use(middleware.GinLogger(), middleware.GinRecovery(true))
	zap.L().Info("ates zap")
	// Swagger
	swaggerHandler := ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "DISABLE_SWAGGER_HTTP_HANDLER")
	handler.GET("/swagger/*any", swaggerHandler)
	// K8s probe
	handler.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })

	// Prometheus metrics
	handler.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Routers
	h := handler.Group("/api/v1")
	{
		rest.NewMonitorHandler(h, mSvc)
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

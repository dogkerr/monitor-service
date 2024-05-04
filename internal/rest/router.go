package rest

import (
	"dogker/lintang/monitor-service/internal/rest/middleware"
	"net/http"

	_ "dogker/lintang/monitor-service/docs"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

//  NewRouter -.
//  Swagger spec:
//	@title			Dogker Monitor Service
//	@version		1.0
//	@description	Monitor Servicee buat nampilin logs dashboard & conainer metrics dashboard milik user
//	@termsOfService	http://swagger.io/terms/
//	@host			locahost:9191
//	@BasePath		/api/v1/

// @contact.name	Lintang BS
func NewRouter(handler *gin.Engine, mSvc MonitorService) *gin.RouterGroup {
	// router
	handler.Use(middleware.GinLogger(), middleware.GinRecovery(true))

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
		NewMonitorHandler(h, mSvc)
	}
	return h
}

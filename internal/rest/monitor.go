package rest

import (
	"context"
	"dogker/lintang/monitor-service/domain"
	"dogker/lintang/monitor-service/internal/rest/middleware"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// ResponseError represent the response error struct
type ResponseError struct {
	Message string `json:"message"`
}

type MonitorService interface {
	TesDoang(ctx context.Context) (string, error)
	GetAllUserContainerService(ctx context.Context, userID string) (*[]domain.Container, error)
	GetUserMonitorDashboard(ctx context.Context, userID string) (*domain.Dashboard, error)
	GetLogsDashboard(ctx context.Context, userID string) (*domain.Dashboard, error)
	SendAllUsersMetricsToRMQ(ctx context.Context) error
	SendTerminatedInstanceToContainerService(ctx context.Context) error
	AuthorizeGrafanaDashboardAccess(ctx context.Context, ctrID string, userID string) error
}

type MonitorHandler struct {
	service MonitorService
}

func NewMonitorHandler(rg *gin.RouterGroup, svc MonitorService) {
	handler := &MonitorHandler{
		service: svc,
	}
	h := rg.Group("/monitors")
	{
		h.GET("/tes", middleware.AuthMiddleware(), handler.TesDoang)
		h.GET("/services", middleware.AuthMiddleware(), handler.GetAllUserContainerHandler)
		h.GET("/dashboards/monitors", middleware.AuthMiddleware(), handler.GetUserMonitorDashboard)
		h.GET("/dashboards/logs", middleware.AuthMiddleware(), handler.GetUserLogsDashboard)
		h.POST("/cron/usersmetrics", handler.CronAllUsersHandler)
		h.POST("/cron/terminatedAccidentally", handler.ContainerTerminatedAccidentally)
		h.GET("/")
		h.GET("/grafana/auth", middleware.AuthMiddleware(), handler.AuthGrafana)

	}
}

func (m *MonitorHandler) CronAllUsersHandler(c *gin.Context) {
	err := m.service.SendAllUsersMetricsToRMQ(c)
	if err != nil {
		c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, "ok")
}

type ContainerLifecycleRes struct {
	ID          uuid.UUID `json:"id"`
	ContainerID uuid.UUID `json:"containerId"`
	StartTime   time.Time `json:"start_time"`
	StopTime    time.Time `json:"stop_time"`
	CPUCore     float64   `json:"cpu_core"`
	MemCapacity float64   `json:"mem_capacity"`
	Replica     uint64    `json:"replica"`
	Status      string    `json:"status"`
}
type serviceResponse struct {
	ID            uuid.UUID `json:"id"`
	UserID        uuid.UUID `json:"user_id"`
	ImageURL      string    `json:"image_url"`
	Status        string    `json:"status"`
	Name          string    `json:"name"`
	ContainerPort int       `json:"container_port"`
	PublicPort    int       `json:"public_port"`
	CreatedTime   time.Time `json:"created_at"`

	ContainerLifecycles []ContainerLifecycleRes `json:"all_container_lifecycles"`
}

type allUserCtrRes struct {
	Services []serviceResponse `json:"services"`
}

func (m *MonitorHandler) GetAllUserContainerHandler(c *gin.Context) {
	userID := c.Query("userID")
	sv, err := m.service.GetAllUserContainerService(c, userID)
	if err != nil {
		c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
		return
	}
	var res []serviceResponse
	for i := range *sv {
		s := *sv
		var lifes []ContainerLifecycleRes
		for _, l := range s[i].ContainerLifecycles {
			lifes = append(lifes, ContainerLifecycleRes{
				ID:          l.ID,
				ContainerID: l.ContainerID,
				StartTime:   l.StartTime,
				StopTime:    l.StopTime,
				Replica:     l.Replica,
				Status:      l.Status.String(),
			})
		}

		res = append(res, serviceResponse{
			ID:                  s[i].ID,
			UserID:              s[i].UserID,
			ImageURL:            s[i].Image,
			Status:              s[i].Status.String(),
			Name:                s[i].Name,
			ContainerPort:       s[i].ContainerPort,
			PublicPort:          s[i].PublicPort,
			CreatedTime:         s[i].CreatedTime,
			ContainerLifecycles: lifes,
		})
	}
	c.JSON(http.StatusOK, allUserCtrRes{res})
}

//	dashboardRes
//
// @Description	Response saat get metrics dashboard milik user
type dashboardRes struct {
	// data dashboard milik user (isinya uid, owner, type, id)
	Dashboard domain.Dashboard `json:"dashboard"`
	// link dashboard metrics received network per contaainer
	ReceivedNetworkLink string `json:"received_network_link" example:"http://127.0.0.1/d-solo/VWUxxYP3/vwuxxyp3?orgId=1&refresh=5s&from=now-5m&theme=light&to=now&panelId=8"`
	// link dashboard metrics send networks per contaainer
	SendNetworkLink string `json:"send_network_link"`
	// link dashboard cpu usage per contaainer
	CpuUsagePerContainer string `json:"cpu_usage_link"`
	// link memory swap per container
	MemorySwapPerContainer string `json:"memory_swap_per_container_link"`
	// link memory usage per container pake graph
	MemoryUsagePerContainer string `json:"memory_usage_per_container_link"`
	// link memory usage per container gak pake graph cuma angka
	MemoryUsageNotGraph string `json:"memory_usage_not_graph"`
	// link overal cpu usage untuk semua container milik user
	OveralCpuUsage string `json:"overall_cpu_usage"`
	// jumlah container yang dijalankan user di dogker
	TotalContainer string `json:"total_container"`
}

//	GetUserMonitorDashboard godoc
//
// @Summary		Mendapatkan Dashboard Container metrics milik User
// @Description	GetUserMonitorDashboard
// @ID				monitor_dashboard
// @Tags			monitor
// @Accept			json
// @Produce		json
// @Success		200		{object}	dashboardRes	"ok"
// @Failure		500		{object}	ResponseError	"internal server error (bug/error di kode)"
// @Router			/monitors/dashboards/monitors [get]
func (m *MonitorHandler) GetUserMonitorDashboard(c *gin.Context) {
	// userID := c.Query("userID")
	userID, _ := c.Get("userID")
	zap.L().Debug("userID monitor", zap.String("userID", userID.(string)))
	sv, err := m.service.GetUserMonitorDashboard(c, userID.(string))
	if err != nil {
		c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
		return
	}
	dbResult := dashboardRes{
		Dashboard:               *sv,
		ReceivedNetworkLink:     "http://127.0.0.1:3000/d-solo/" + sv.Uid + "/" + strings.ToLower(sv.Uid) + "?orgId=1&refresh=5s&from=now-5m&theme=light&to=now&panelId=8",
		SendNetworkLink:         "http://127.0.0.1:3000/d-solo/" + sv.Uid + "/" + strings.ToLower(sv.Uid) + "?orgId=1&refresh=5s&from=now-5m&theme=light&to=now&panelId=9",
		CpuUsagePerContainer:    "http://127.0.0.1:3000/d-solo/" + sv.Uid + "/" + strings.ToLower(sv.Uid) + "?orgId=1&refresh=5s&from=now-5m&theme=light&to=now&panelId=1",
		MemorySwapPerContainer:  "http://127.0.0.1:3000/d-solo/" + sv.Uid + "/" + strings.ToLower(sv.Uid) + "?orgId=1&refresh=5s&from=now-5m&theme=light&to=now&panelId=34",
		MemoryUsagePerContainer: "http://127.0.0.1:3000/d-solo/" + sv.Uid + "/" + strings.ToLower(sv.Uid) + "?orgId=1&refresh=5s&from=now-5m&theme=light&to=now&panelId=10",
		MemoryUsageNotGraph:     "http://127.0.0.1:3000/d-solo/" + sv.Uid + "/" + strings.ToLower(sv.Uid) + "?orgId=1&refresh=5s&from=now-5m&theme=light&to=now&panelId=37",
		OveralCpuUsage:          "http://127.0.0.1:3000/d-solo/" + sv.Uid + "/" + strings.ToLower(sv.Uid) + "?orgId=1&refresh=5s&from=now-5m&theme=light&to=now&panelId=5",
		TotalContainer:          "http://127.0.0.1:3000/d-solo/" + sv.Uid + "/" + strings.ToLower(sv.Uid) + "?orgId=1&refresh=5s&from=now-5m&theme=light&to=now&panelId=31",
	}
	c.JSON(http.StatusOK, dbResult)
}

//	logsDashboardRes
//
// @Description	Response saat get logs dashboard milik user
type logsDashboardRes struct {
	// data dashboard milik user (isinya uid, owner, type, id)
	Dashboard domain.Dashboard `json:"dashboard"`
	// link dashboard logs yang diembed di frontend
	LogsDashboardLink string `json:"logs_dashboard_link" example:"http://localhost:3000/d/YwXYwNAj/ywxywnaj?orgId=1&var-search_filter=&var-Levels=info&var-container_name=go_container_log2&var-Method=GET&from=1714796971638&to=1714797271638&theme=light"`
}

//	GetUserLogsDashboard godoc
//
// @Summary		Mendapatkan Dashboard Logs containers milik User
// @Description	GetUserLogsDashboard
// @ID				logs_dashboard
// @Tags			monitor
// @Accept			json
// @Produce		json
// @Success		200		{object}	logsDashboardRes	"ok"
// @Failure		500		{object}	ResponseError		"internal server error (bug/error di kode)"
// @Router			/monitors/dashboards/logs [get]
func (m *MonitorHandler) GetUserLogsDashboard(c *gin.Context) {
	// userID := c.Query("userID")
	userID, _ := c.Get("userID")
	zap.L().Debug("userID logs", zap.String("userID", userID.(string)))
	sv, err := m.service.GetLogsDashboard(c, userID.(string))
	if err != nil {
		c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
		return
	}
	dbResult := logsDashboardRes{
		Dashboard:         *sv,
		LogsDashboardLink: "http://localhost:3000/d/" + sv.Uid + "/" + strings.ToLower(sv.Uid) + "?orgId=1&theme=light",
	}
	c.JSON(http.StatusOK, dbResult)
}

func (m *MonitorHandler) TesDoang(c *gin.Context) {
	tesResult, err := m.service.TesDoang(c)
	if err != nil {
		c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, tesResult)
}

type CronTerminatedInstanceResp struct {
	Message string `json:"message"`
}

func (m *MonitorHandler) ContainerTerminatedAccidentally(c *gin.Context) {
	err := m.service.SendTerminatedInstanceToContainerService(c)
	if err != nil {
		c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, CronTerminatedInstanceResp{"ok"})
}

type grafanaAuthReq struct {
	DashboardID string `form:"dashboard_id"`
}

type grafanaAuthRes struct {
	Message string `json:"message"`
}

func (m *MonitorHandler) AuthGrafana(c *gin.Context) {
	userID, _ := c.Get("userID")
	var req grafanaAuthReq
	if err := c.ShouldBindQuery(&req); err != nil {
		// bind query param dashboard_id
		c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
		return
	}
	err := m.service.AuthorizeGrafanaDashboardAccess(c, req.DashboardID, userID.(string))
	if err != nil {
		c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, grafanaAuthRes{Message: "you are authorized"})
}

func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	// logrus.Error(err)
	switch err {
	case domain.ErrInternalServerError:
		return http.StatusInternalServerError
	case domain.ErrNotFound:
		return http.StatusNotFound
	case domain.ErrConflict:
		return http.StatusConflict
	case domain.ErrBadParamInput:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

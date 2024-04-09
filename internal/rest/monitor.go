package rest

import (
	"context"
	"dogker/lintang/monitor-service/domain"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
		h.GET("/tes", handler.TesDoang)
		h.GET("/services", handler.GetAllUserContainerHandler)
		h.GET("/dashboards/monitors", handler.GetUserMonitorDashboard)
		h.GET("/dashboards/logs", handler.GetUserLogsDashboard)
	}
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
			ImageURL:            s[i].ImageURL,
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

type dashboardRes struct {
	Dashboard domain.Dashboard `json:"dashboard"`
}

func (m *MonitorHandler) GetUserMonitorDashboard(c *gin.Context) {
	userID := c.Query("userID")
	sv, err := m.service.GetUserMonitorDashboard(c, userID)
	if err != nil {
		c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	dbResult := dashboardRes{
		*sv,
	}
	c.JSON(http.StatusOK, dbResult)
}

func (m *MonitorHandler) GetUserLogsDashboard(c *gin.Context) {
	userID := c.Query("userID")
	sv, err := m.service.GetLogsDashboard(c, userID)
	if err != nil {
		c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	dbResult := dashboardRes{
		*sv,
	}
	c.JSON(http.StatusOK, dbResult)
}

func (m *MonitorHandler) TesDoang(c *gin.Context) {
	tesResult, err := m.service.TesDoang(c)
	if err != nil {
		c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	c.JSON(http.StatusOK, tesResult)
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

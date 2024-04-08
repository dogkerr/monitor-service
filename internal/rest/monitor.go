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
	GetAllUserContainerService(ctx context.Context, userId string) (*[]domain.Container, error)
}

type MonitorHandler struct {
	Service MonitorService
}

func NewMonitorHandler(rg *gin.RouterGroup, svc MonitorService) {
	handler := &MonitorHandler{
		Service: svc,
	}
	h := rg.Group("/monitors")
	//h.Use(ginkeycloak.Auth(ginkeycloak.AuthCheck(), sbbEndpoint))
	{
		h.GET("/tes", handler.TesDoang)
		h.GET("/services", handler.GetAllUserContainerHandler)
	}
}

type ContainerLifecycleRes struct {
	ID          uuid.UUID `json:"id"`
	ContainerId uuid.UUID `json:"containerId"`
	StartTime   time.Time `json:"start_time"`
	StopTime    time.Time `json:"stop_time"`
	CpuCore     float64   `json:"cpu_core"`
	MemCapacity float64   `json:"mem_capacity"`
	Replica     uint64    `json:"replica"`
	Status      string    `json:"status"`
}
type serviceResponse struct {
	ID            uuid.UUID `json:"id"`
	UserId        uuid.UUID `json:"user_id"`
	ImageUrl      string    `json:"image_url"`
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
	userId := c.Query("userId")
	sv, err := m.Service.GetAllUserContainerService(c, userId)
	if err != nil {
		c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	var res []serviceResponse
	for _, s := range *sv {

		var lifes []ContainerLifecycleRes
		for _, l := range s.ContainerLifecycles {
			lifes = append(lifes, ContainerLifecycleRes{
				ID:          l.ID,
				ContainerId: l.ContainerId,
				StartTime:   l.StartTime,
				StopTime:    l.StopTime,
				Replica:     l.Replica,
				Status:      l.Status.String(),
			})
		}

		res = append(res, serviceResponse{
			ID:                  s.ID,
			UserId:              s.UserId,
			ImageUrl:            s.ImageUrl,
			Status:              s.Status.String(),
			Name:                s.Name,
			ContainerPort:       s.ContainerPort,
			PublicPort:          s.PublicPort,
			CreatedTime:         s.CreatedTime,
			ContainerLifecycles: lifes,
		})

	}
	c.JSON(http.StatusOK, allUserCtrRes{res})
}

func (m *MonitorHandler) TesDoang(c *gin.Context) {
	tesResult, err := m.Service.TesDoang(c)
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

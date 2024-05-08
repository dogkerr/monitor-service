package domain

import (
	"time"

	"github.com/google/uuid"
)

type ContainerStatus int

const (
	RUN ContainerStatus = iota + 1
	STOP
)

func (s ContainerStatus) String() string {
	return [...]string{"RUN", "STOP"}[s-1]
}

// type ContainerStatus string

// const (
// 	RUN  ContainerStatus = "RUN"
// 	STOP ContainerStatus = "STOP"
// )

type Container struct {
	ID                  uuid.UUID            `json:"id"`
	UserID              uuid.UUID            `json:"user_id"`
	Image               string               `json:"image_url"`
	Status              ContainerStatus      `json:"status"`
	Name                string               `json:"name"`
	ContainerPort       int                  `json:"container_port"`
	PublicPort          int                  `json:"public_port"`
	CreatedTime         time.Time            `json:"created_at"`
	TerminatedTime      time.Time            `json:"terminated_time"`
	ContainerLifecycles []ContainerLifecycle `json:"all_container_lifecycles"`
	ServiceID           string               `json:"serviceId"`
}

package domain

import (
	"time"

	"github.com/google/uuid"
)

type Status int

const (
	RUN Status = iota + 1
	STOPPED
)

func (s Status) String() string {
	return [...]string{"RUN", "STOPPED"}[s-1]
}

type Container struct {
	ID                  uuid.UUID            `json:"id"`
	UserID              uuid.UUID            `json:"user_id"`
	ImageURL            string               `json:"image_url"`
	Status              Status               `json:"status"`
	Name                string               `json:"name"`
	ContainerPort       int                  `json:"container_port"`
	PublicPort          int                  `json:"public_port"`
	CreatedTime         time.Time            `json:"created_at"`
	TerminatedTime      time.Time            `json:"terminated_time"`
	ContainerLifecycles []ContainerLifecycle `json:"all_container_lifecycles"`
	ServiceID           string               `json:"serviceId"`
}

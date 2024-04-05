package domain

import (
	"time"
)

type Action int

const (
	START Action = iota + 1
	STOP
)

func (s Status) ActionString() string {
	return [...]string{"START", "STOP"}[s-1]
}

type ContainerAction struct {
	ID          string    `json:"id"`
	ContainerId string    `json:"container_id"`
	Timestamp   time.Time `json:"timestamp"`
	Action      Action    `json:"action"`
}

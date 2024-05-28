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

var GetContainerStatus = map[string]ContainerStatus{ //create a map to link enumeration values with string representation
	"RUN":    RUN,
	"STOP":     STOP,
}


type ContainerLifecycle struct {
	ID          uuid.UUID       `json:"id"`
	ContainerID uuid.UUID       `json:"containerId"`
	StartTime   time.Time       `json:"start_time"`
	StopTime    time.Time       `json:"stop_time"`
	CPUCore     float64         `json:"cpu_core"`
	MemCapacity float64         `json:"mem_capacity"`
	Replica     uint64          `json:"replica"`
	Status      ContainerStatus `json:"status"`
}

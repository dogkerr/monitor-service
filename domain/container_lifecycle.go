package domain

import (
	"time"

	"github.com/google/uuid"
)




type ContainerLifecycle struct {
	ID uuid.UUID `json:"id"`
	ContainerId uuid.UUID `json:"containerId"`
	StartTime time.Time `json:"start_time"`
	StopTime time.Time `json:"stop_time"`
	CpuCore float64 `json:"cpu_core"`
	MemCapacity float64 `json:"mem_capacity"`
	Replica uint64 `json:"replica"`
	Status Status `json:"status"`
}

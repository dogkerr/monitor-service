package domain

import "time"




type ContainerLifecycle struct {
	ID string `json:"id"`
	ContainerId string `json:"containerId"`
	StartTime time.Time `json:"start_time"`
	StopTime time.Time `json:"stop_time"`
	CpuCore float64 `json:"cpu_core"`
	MemCapacity float64 `json:"mem_capacity"`
	Replica uint64 `json:"replica"`
}

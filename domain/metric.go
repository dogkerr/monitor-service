package domain

import "time"

type Metric struct {
	CpuUsage            float32   `json:"cpu_usage"`
	MemoryUsage         float32   `json:"memory_usage"`
	NetworkIngressUsage float32   `json:"network_ingress_usage"`
	NetworkEgressUsage  float32   `json:"network_egress_usage"`
	CreatedTime         time.Time `json:"created_time"`
}

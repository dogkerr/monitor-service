package domain

import "time"


type Prometheus struct {
	CurrentTime time.Time `json:"currentTime"`
	AllCpuUsage float32 `json:"all_cpu_usage"`
	AllMemoryUsage float32 `json:"all_memory_usage"`
	AllNetworkIngressUsage float32 `json:"all_network_ingress_usage"`
	AllNetworkEgressUsage float32 `json:"all_network_egress_usage"`
	FromTime time.Time `json:"fromTime"`
}
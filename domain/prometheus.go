package domain

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type Prometheus struct {
	CurrentTime            *timestamppb.Timestamp `json:"currentTime"`
	AllCPUUsage            float32                `json:"all_cpu_usage"`
	AllMemoryUsage         float32                `json:"all_memory_usage"`
	AllNetworkIngressUsage float32                `json:"all_network_ingress_usage"`
	AllNetworkEgressUsage  float32                `json:"all_network_egress_usage"`
	FromTime               *timestamppb.Timestamp `json:"fromTime"`
}


type PrometheusNotGRPC struct {
	CurrentTime            time.Time `json:"currentTime"`
	AllCPUUsage            float32                `json:"all_cpu_usage"`
	AllMemoryUsage         float32                `json:"all_memory_usage"`
	AllNetworkIngressUsage float32                `json:"all_network_ingress_usage"`
	AllNetworkEgressUsage  float32                `json:"all_network_egress_usage"`
	FromTime               time.Time `json:"fromTime"`
}

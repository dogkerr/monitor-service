package domain





type UserMetricsMessage struct {
	ContainerID string 
	UserID string 
	CpuUsage            float32             
	MemoryUsage         float32               
	NetworkIngressUsage float32                
	NetworkEgressUsage  float32                
}





// type AllUsersMetricsMessage struct {
// 	AllUsersMetrics []UserMetricsMessage `json:"all_users_metrics"`
// }


// type UserMetricsMessage struct {
// 	ContainerID string `json:"containerID"`
// 	UserID string `json:"userID"`
// 	CpuUsage            float32                `json:"cpu_usage"`
// 	MemoryUsage         float32                `json:"memory_usage"`
// 	NetworkIngressUsage float32                `json:"network_ingress_usage"`
// 	NetworkEgressUsage  float32                `json:"network_egress_usage"`
// }



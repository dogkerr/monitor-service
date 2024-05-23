package domain
type CommonLabelsMailing struct {
	Alertname                       string `json:"alertname"`
	ContainerSwarmServiceID         string `json:"container_label_com_docker_swarm_service_id"`
	ContainerDockerSwarmServiceName string `json:"container_label_com_docker_swarm_service_name"`
	ContainerLabelUserID            string `json:"container_label_user_id"`
}
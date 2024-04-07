
### Cara testing
1. bikin client grpc buat monitor service
2. migrate database && insert data dummy ke masing masing table (lihat di migrations/.....sql)

3. jalanin container dg label userId sama dg yg di database
```

 docker service create --name  go_container  --publish 8040:80 --replicas 3 --container-label  user_id="c25ed8a3-cc49-4ba2-9f3d-23c1db70ec24"  generate_user_dashboard_dan_perfomance_testing-go_container_log_user1:latest


  docker service create --name  go_container_user2  --publish 8032:80 --replicas 1 --container-label  user_id="c25ed8a3-cc49-4ba2-9f3d-23c1db70ec24"  generate_user_dashboard_dan_perfomance_testing-go_container_log_user1:latest


```

----
query prome
```
1. buat dapetin metrics per swarm service:
sum(rate(container_cpu_usage_seconds_total{container_label_com_docker_swarm_service_id=~"swarmServiceId"}[1h])) * 100 * 7200 /3600
sum(avg_over_time(container_memory_usage_bytes{container_label_com_docker_swarm_service_id=~"swarmServiceId"}[1h])) * 30*24*3600 / 3600 / (1024^3)
sum(container_network_receive_bytes_total{container_label_com_docker_swarm_service_id=~"swarmServiceId"}) / 1024
sum(container_network_transmit_bytes_total{container_label_com_docker_swarm_service_id=~"swarmServiceId"}) / 1024

ue8bafuvfbtra1yzn3u0kg3gu


```

buat stop service docker swarm (semua replica container ): `docker service scale go_container=0`
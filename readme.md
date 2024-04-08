
### Cara testing
1. ikutin cara nyalain prometheus di readme repo dogker/configs
2. docker compose up [monitor service]
3. migrate database && insert data dummy ke masing masing table (lihat di migrations/.....up.sql)
4.  jalanin docker swarm service
```
 docker service create --name  go_container_4  --publish 8040:80 --replicas 3 --container-label  user_id=<user_id_di_table_container>    generate_user_dashboard_dan_perfomance_testing-go_container_log_user1:latest


  docker service create --name  go_container_user2  --publish 8032:80 --replicas 1 --container-label  user_id=<user_id_di_table_container>    generate_user_dashboard_dan_perfomance_testing-go_container_log_user1:latest

```


4. tambahin serviceID docker swarm ke setiap row di tabel container (harusnya container service yg namabahin)
```
buat dapetin serviceId: `docker service ls`
paling kiri id service ny.
```


6. go run app/main.go
7. jalanin client & kirim request ke
```
http://localhost:5033/api/v1/monitors/metrics?userId=<user_id_di_database>



```


----

### query prome
```
1. buat dapetin metrics per swarm service:
sum(rate(container_cpu_usage_seconds_total{container_label_com_docker_swarm_service_id=~"swarmServiceId"}[1h])) * 100 * 7200 /3600
sum(avg_over_time(container_memory_usage_bytes{container_label_com_docker_swarm_service_id=~"swarmServiceId"}[1h])) * 30*24*3600 / 3600 / (1024^3)
sum(container_network_receive_bytes_total{container_label_com_docker_swarm_service_id=~"swarmServiceId"}) / 1024
sum(container_network_transmit_bytes_total{container_label_com_docker_swarm_service_id=~"swarmServiceId"}) / 1024

ue8bafuvfbtra1yzn3u0kg3gu


```

buat stop service docker swarm (semua replica container ): `docker service scale go_container=0`


### Cara testing
```
1. buat folder pb
2. buat .env isinya PG_URL
3. generate protobuf code `make proto`

```


1. ikutin cara nyalain prometheus di readme repo dogker/configs
2. docker compose up [monitor service]
3. migrate database && insert data dummy ke masing masing table (lihat di migrations/.....up.sql)
```
docker swarm init
docker service ls
insert service_id ke rows tabel container,
```
4.  jalanin docker swarm service
```
 docker service create --name  go_container_4  --publish 8040:80 --replicas 3 --container-label  user_id=<user_id_di_table_container>    generate_user_dashboard_dan_perfomance_testing-go_container_log_user1:latest


  docker service create --name  go_container_user2  --publish 8032:80 --replicas 1 --container-label  user_id=<user_id_di_table_container>    generate_user_dashboard_dan_perfomance_testing-go_container_log_user1:latest



docker service create --name  go_container_log1  --publish 8036:80 --replicas 2 --container-label  user_id=<user_id_di_table_container> --log-driver=loki \
    --log-opt loki-url="http://localhost:3100/loki/api/v1/push" \
    --log-opt loki-retries=5 \
    --log-opt loki-batch-size=400 \
    --log-opt loki-external-labels="job=docker,container_name=go_container_api_user2,userId=<user_id_di_table_container>" generate_user_dashboard_dan_perfomance_testing-go_container_log_user1:latest 


## contoh:
docker service create --name  go_container_log2  --publish 8038:80 --replicas 2 --container-label  user_id=eff92b7f-3f90-405b-9fb8-1ff12eb72431 --log-driver=loki \
    --log-opt loki-url="http://localhost:3100/loki/api/v1/push" \
    --log-opt loki-retries=5 \
    --log-opt loki-batch-size=400 \
    --log-opt loki-external-labels="job=docker,container_name=go_container_api_user2,userId=eff92b7f-3f90-405b-9fb8-1ff12eb72431" generate_user_dashboard_dan_perfomance_testing-go_container_log_user1:latest 




# harus localhost:3100/loki/api/v1/push biar bisa kedetect loki (pake loki:3100 gakbisa)


docker service create --name  go_container_log2  --publish 8037:80 --replicas 1 --container-label  user_id=<user_id_di_table_container>    generate_user_dashboard_dan_perfomance_testing-go_container_log_user1:latest 

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
http://localhost:5033/api/v1/monitors/ctrMetrics?userId=<user_id_di_database>&serviceId=<docker_swarm_serviceId>


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


### Coba generate dashboard

-- get service account token
```
nama service account: adm
id serviceAcc: sa-adm
mending grafananya dipasangin volume biar gaperlu buat serviceAcc setiap setup grafana

harus masukin api key serviceAcc sa-adm ke config.yaml
```


### genereate loki dashboard

```
curl localhost:3000/

```

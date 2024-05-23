docker build .  -f Dockerfile.debug -t lintangbirdas/monitor-service-dev:v1
docker tag lintangbirdas/monitor-service-dev:v1 lintangbirdas/monitor-service-dev:v1
docker push lintangbirdas/monitor-service-dev:v1

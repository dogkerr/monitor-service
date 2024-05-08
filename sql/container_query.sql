-- name: GetAllUserContainers :many
SELECT c.id, c.user_id, c.image, c.status, c.name, c.container_port, c.public_port, c.created_time,c.service_id, c.terminated_time,
			cl.id as lifecycleId, cl.start_time as lifecycleStartTime, cl.stop_time as lifecycleStopTime, 
			cl.replica as lifecycleReplica, cl.status as lifecycleStatus FROM containers c  LEFT JOIN container_lifecycles cl ON cl.container_id=c.id
			WHERE c.user_id=$1;



-- name: GetSpecificContainerMetrics :one
SELECT m.id, m.cpus, m.memory, m.network_ingress, m.network_egress
	FROM container_metrics m 
	WHERE m.container_id=$1;


-- name: GetContainer :many
SELECT c.id, c.user_id, c.image, c.status, c.name, c.container_port, c.public_port,c.created_time,
	c.service_id,c.terminated_time, cl.id as lifeId, cl.start_time as lifecycleStartTime, cl.stop_time as lifecycleStopTime, cl.replica  as lifecycleReplica, cl.status as lifecycleStatus 
	FROM containers c LEFT JOIN container_lifecycles cl ON cl.container_id=c.id
	WHERE c.service_id=$1;


-- name: GetAllUserContainer :many
SELECT c.id, c.user_id, c.image, c.status, c.name, c.container_port, c.public_port, c.created_time,c.service_id, c.terminated_time,
	cl.id as lifecycleId, cl.start_time as lifecycleStartTime, cl.stop_time as lifecycleStopTime, 
	cl.replica as lifecycleReplica, cl.status as lifecycleStatus FROM containers c  LEFT JOIN container_lifecycles cl ON cl.container_id=c.id
	WHERE c.user_id=$1;


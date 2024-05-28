-- name: GetAllUserContainers :many
SELECT c.id, c.user_id, c.image, c.status, c.name, c.container_port, c.public_port, c.created_time,c.service_id, c.terminated_time,
			cl.id as lifecycleId, cl.start_time as lifecycleStartTime, cl.stop_time as lifecycleStopTime, 
			cl.replica as lifecycleReplica, cl.status as lifecycleStatus FROM containers c  LEFT JOIN container_lifecycles cl ON cl.container_id=c.id
			WHERE c.user_id=$1;



-- name: GetSpecificContainerMetrics :many
SELECT m.id, m.cpus, m.memory, m.network_ingress, m.network_egress, m.created_time
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


-- name: GetContainerOwnerByID :one
SELECT d.id, d.owner, d.uid
	FROM dashboards d 
	WHERE d.uid=$1;


-- name: InsertTerminatedContainer :exec 
INSERT INTO processed_terminated_container(
	container_id, down_time
) VALUES (
	$1, $2
);





-- name: GetProcessedContainers :many
SELECT c.container_id, c.down_time
	FROM processed_terminated_container c 
	WHERE c.container_id  = ANY($1::UUID[]);
	
	-- IN ($1::UUID[]); -- ini gakbisa

-- name: GetContainerByServiceIDs :many
SELECT c.id, c.service_id, cl.replica as lifecycleReplica, cl.start_time
	FROM containers c 
	LEFT JOIN container_lifecycles cl ON cl.container_id=c.id
	WHERE c.service_id = ANY($1::varchar[]);


-- name: GetSwarmServiceDetailByServiceIDs :many
SELECT c.id, c.service_id, c.name, c.user_id
	FROM containers c 
	LEFT JOIN container_lifecycles cl ON cl.container_id=c.id
	WHERE c.service_id = ANY($1::varchar[]);



-- name: GetContainerLifecycleByCtrID :many
SELECT cl.start_time, cl.status, cl.id
	FROM container_lifecycles cl
	WHERE cl.container_id = $1;




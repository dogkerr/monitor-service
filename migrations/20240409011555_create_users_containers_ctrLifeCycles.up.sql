-- CREATE EXTENSION IF NOT EXISTS 'uuid-ossp';
CREATE TYPE container_status AS ENUM ('RUN', 'STOPPED');


CREATE TABLE users (
    id  UUID DEFAULT gen_random_uuid() PRIMARY KEY ,
    username VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL, 
    password VARCHAR(255) NOT NULL
);

CREATE TABLE containers (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY ,
    user_id UUID NOT NULL,
    image_url VARCHAR(255) NOT NULL,
    status container_status NOT NULL,
    name VARCHAR(255) NOT NULL,
    container_port int NOT NULL,
    public_port int,
    terminated_time timestamp with time zone ,
    created_time timestamp with time zone  DEFAULT NOW() not null,
    service_id VARCHAR(255)
);




CREATE TABLE container_lifecycles (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY ,
    container_id UUID,
    start_time TIMESTAMP with time zone NOT NULL,
    stop_time TIMESTAMP with time zone,
    status container_status NOT NULL,
    replica INT NOT NULL
);




ALTER TABLE container_lifecycles ADD  CONSTRAINT fk_lifecycles_containers
    FOREIGN KEY (container_id)
    REFERENCES containers (id);

ALTER TABLE containers ADD CONSTRAINT fk_containers_users
    FOREIGN KEY (user_id) 
    REFERENCES users(id);


INSERT INTO users( username, email, password) 
    VALUES ('asda', 'sadasd@gmail.com', 'asdad');

/*
14b16002-c5ec-4906-90c9-8e3749864408
INSERT INTO containers(user_id, image_url, status, name, container_port, public_port) 
    VALUES('<user_id_di_tabel_user>', 'prome', 'RUN', 'prome1', 9090, 9090 );



SELECT * FROM containers;

8bd49dfb-ef83-4315-a90d-e6275f5871b4
INSERT INTO container_lifecycles(container_id, start_time, stop_time, replica, status)
    VALUES('<container_id_di_tabel_container>', NOW(), NOW(), 3, 'RUN');


INSERT INTO container_lifecycles(container_id, start_time, stop_time, replica, status)
    VALUES('<container_id_di_tabel_container>', NOW(), NOW(), 3, 'STOPPED');

INSERT INTO container_lifecycles(container_id, start_time, stop_time, replica, status)
    VALUES('<container_id_di_tabel_container>', NOW(), NOW(), 3, 'RUN');


SELECT c.id, c.user_id, c.image_url, c.status, c.name, c.container_port, c.public_port, c.created_time,
			cl.id as lifecycleId, cl.start_time as lifecycleStartTime, cl.stop_time as lifecycleStopTime, 
			cl.replica as lifecycleReplica, cl.status FROM containers c  LEFT JOIN container_lifecycles cl ON cl.container_id=c.id
			WHERE c.user_id='<user_id_di_tabel_user>';


SELECT * FROM containers c 
 LEFT JOIN container_lifecycles cl ON cl.container_id=c.id
 WHERE c.user_id='<user_id_di_tabel_user>';
*/




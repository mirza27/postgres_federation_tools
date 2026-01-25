# run container
run-base:
	docker compose -f docker-compose.base.yml --env-file .base.env up -d 

run-pub:
	docker compose -f docker-compose.publication.yml --env-file .publication.env up -d

run-northwind:
	docker compose -f docker-compose.northwind.yml --env-file .northwind.env up -d

run-chinook:
	docker compose -f docker-compose.chinook.yml --env-file .chinook.env up -d

# RUN worker
executor:
	go run ./cmd/executor/main.go >  logs/executor.log 2>&1

parser:
	go run ./cmd/parser/main.go >  logs/parser.log 2>&1

checker:
	go run ./cmd/checker/main.go >  logs/checker.log 2>&1

joiner:
	go run ./cmd/joiner/main.go >  logs/joiner.log 2>&1

api:
	go run main.go


# Debezium Config
conn-publication:
	curl -X POST -H "Content-Type: application/json"   -d @internal/debezium/connector-publication.json   http://localhost:8083/connectors
dis-publication:
	curl -X DELETE http://localhost:8083/connectors/publication-connector

#kafka
reset-kafka:
	docker stop debezium
	docker exec -it kafka kafka-consumer-groups --bootstrap-server kafka:9092 --delete --group 1
	docker exec -it kafka kafka-topics --bootstrap-server kafka:9092 --delete --topic '.*'
	docker start debezium

# pivotable 
up-pivot:
	docker exec -i pivot_db psql -U pivot_user -d pivot < ./internal/config/pivot_db.sql

drop-pivot:
	docker exec -i pivot_db  psql -U pivot_user -d pivot < ./migration/pivot/down.sql

empty-pivot:
	docker exec -i pivot_db  psql -U pivot_user -d pivot < ./migration/pivot/empty.sql
	

# PUBLICATION CASE
add-new-publication:
	docker exec -i new_publication psql -U new_user -d new_publication_db < ./migration/publication/new/up.sql
	docker exec -i new_publication psql -U new_user -d new_publication_db < ./migration/publication/new/seed.sql

drop-new-publication:
	docker exec -i new_publication psql -U new_user -d new_publication_db < ./migration/publication/new/down.sql

empty-new-publication:
	docker exec -i new_publication psql -U new_user -d new_publication_db < ./migration/publication/new/empty.sql

add-old-publication:
	docker exec -i old_publication psql -U old_user -d old_publication_db < ./migration/publication/old/up.sql
	docker exec -i old_publication psql -U old_user -d old_publication_db < ./migration/publication/old/seed.sql

drop-old-publication:
	docker exec -i old_publication psql -U old_user -d old_publication_db < ./migration/publication/old/down.sql

empty-old-publication:
	docker exec -i old_publication psql -U old_user -d old_publication_db < ./migration/publication/old/empty.sql

pub-seed:
	docker exec -i old_publication psql -U old_user -d old_publication_db < ./migration/publication/old/empty.sql
	go run ./cmd/seeder/*.go

# CHINOOK CASE
add-new-chinook:
	docker exec -i new_chinook psql -U new_user -d new_chinook_db < ./migration/chinook/new/up.sql

drop-new-chinook:
	docker exec -i new_chinook psql -U new_user -d new_chinook_db < ./migration/chinook/new/down.sql

add-old-chinook:
	docker exec -i old_chinook psql -U old_user -d old_chinook_db < ./migration/chinook/old/up.sql
	docker exec -i old_chinook psql -U old_user -d old_chinook_db < ./migration/chinook/old/seed.sql

drop-old-chinook:
	docker exec -i old_chinook psql -U old_user -d old_chinook_db < ./migration/chinook/old/down.sql


# NORTHWIND CASE
add-new-northwind:
	docker exec -i new_northwind psql -U new_user -d new_northwind_db < ./migration/northwind/new/up.sql

drop-new-northwind:
	docker exec -i new_northwind psql -U new_user -d new_northwind_db < ./migration/northwind/new/down.sql

add-old-northwind:
	docker exec -i old_northwind psql -U old_user -d old_northwind_db < ./migration/northwind/old/up.sql
	docker exec -i old_northwind psql -U old_user -d old_northwind_db < ./migration/northwind/old/seed.sql

drop-old-northwind:
	docker exec -i old_northwind psql -U old_user -d old_northwind_db < ./migration/northwind/old/down.sql


.PHONY: api
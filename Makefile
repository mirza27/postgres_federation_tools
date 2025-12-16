# run container
run-base:
	docker compose -f docker-compose.base.yml --env-file .base.env up -d 

run-pub:
	docker compose -f docker-compose.publication.yml --env-file .publication.env up -d

# RUN worker
executor:
	go run ./cmd/executor/main.go >  logs/executor.log 2>&1

parser:
	go run ./cmd/parser/main.go >  logs/parser.log 2>&1

checker:
	go run ./cmd/checker/main.go >  logs/checker.log 2>&1

joiner:
	go run ./cmd/joiner/main.go >  logs/joiner.log 2>&1


# Debezium Config
conn-publication:
	curl -X POST -H "Content-Type: application/json"   -d @connector-publication.json   http://localhost:8083/connectors
dis-publication:
	curl -X DELETE http://localhost:8083/connectors/publication-connector

conn-ojol:
	curl -X POST -H "Content-Type: application/json"   -d @connector-ojol.json   http://localhost:8083/connectors
dis-ojol:
	curl -X DELETE http://localhost:8083/connectors/ojol-connector

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

add-old-publication:
	docker exec -i old_publication psql -U old_user -d old_publication_db < ./migration/publication/old/up.sql
	docker exec -i old_publication psql -U old_user -d old_publication_db < ./migration/publication/old/seed.sql

drop-old-publication:
	docker exec -i old_publication psql -U old_user -d old_publication_db < ./migration/publication/old/down.sql

# run container
run-base:
	docker compose -f docker-compose.base.yml --env-file .base.env up -d

run-pub:
	docker compose -f docker-compose.publication.yml --env-file .publication.env up -d

# RUN worker
executor:
	go run ./cmd/executor/main.go

parser:
	go run ./cmd/parser/main.go

checker:
	go run ./cmd/checker/main.go

joiner:
	go run ./cmd/joiner/main.go


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

# CRM ECOM CASE
db1:
	docker exec -i fdw_test_ecom psql -U ecom_user -d ecom_db < ./migration/ecom_crm/ecom/up.sql
	docker exec -i fdw_test_ecom psql -U ecom_user -d ecom_db < ./migration/ecom_crm/ecom/seed.sql

drop1:
	docker exec -i fdw_test_ecom psql -U ecom_user -d ecom_db < ./migration/ecom_crm/ecom/down.sql
	
db2:
	docker exec -i fdw_test_crm psql -U crm_user -d crm_db < ./migration/ecom_crm/crm/up.sql
	docker exec -i fdw_test_crm psql -U crm_user -d crm_db < ./migration/ecom_crm/crm/seed.sql

drop2:
	docker exec -i fdw_test_crm psql -U crm_user -d crm_db < ./migration/ecom_crm/crm/down.sql

# EV CASE
add-new-ev:
	docker exec -i new_ev psql -U new_user -d new_ev_db < ./migration/ev/new/up.sql
	docker exec -i new_ev psql -U new_user -d new_ev_db < ./migration/ev/new/seed.sql

drop-new-ev:
	docker exec -i new_ev psql -U new_user -d new_ev_db < ./migration/ev/new/down.sql

add-old-ev:
	docker exec -i old_ev psql -U old_user -d old_ev_db < ./migration/ev/old/up.sql
	docker exec -i old_ev psql -U old_user -d old_ev_db < ./migration/ev/old/seed.sql

drop-old-ev:
	docker exec -i old_ev psql -U old_user -d old_ev_db < ./migration/ev/old/down.sql


# OJOL CASE
add-new-ojol:
	docker exec -i new_ojol psql -U new_user -d new_ojol_db < ./migration/ojol/new/up.sql
	docker exec -i new_ojol psql -U new_user -d new_ojol_db < ./migration/ojol/new/seed.sql

drop-new-ojol:
	docker exec -i new_ojol psql -U new_user -d new_ojol_db < ./migration/ojol/new/down.sql

add-old-ojol:
	docker exec -i old_ojol psql -U old_user -d old_ojol_db < ./migration/ojol/old/up.sql
	docker exec -i old_ojol psql -U old_user -d old_ojol_db < ./migration/ojol/old/seed.sql

drop-old-ojol:
	docker exec -i old_ojol psql -U old_user -d old_ojol_db < ./migration/ojol/old/down.sql
	

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

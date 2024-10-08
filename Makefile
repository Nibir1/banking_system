postgres:
	docker run --name postgres17.0 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:17.0-alpine3.20

createdb: 
	docker exec -it postgres17.0 createdb --username=root --owner=root banking_system

dropdb: 
	docker exec -it postgres17.0 dropdb banking_systems

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/banking_system?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/banking_system?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

.PHONY: postgres createdb dropdb migrateUp migratedown sqlc test
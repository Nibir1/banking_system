postgres:
	docker run --name postgres17.0 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:17.0-alpine3.20

createdb: 
	docker exec -it postgres17.0 createdb --username=root --owner=root banking_system

dropdb: 
	docker exec -it postgres17.0 dropdb banking_systems

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/banking_system?sslmode=disable" -verbose up

migrateup1:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/banking_system?sslmode=disable" -verbose up 1

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/banking_system?sslmode=disable" -verbose down

migratedown1:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/banking_system?sslmode=disable" -verbose down 1

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/nibir1/banking_system/db/sqlc Store

.PHONY: postgres createdb dropdb migrateUp migratedown migrateUp1 migratedown1 sqlc test server mock
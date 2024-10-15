postgres:
	docker run --name postgres17.0 --network bank-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:17.0-alpine3.20

createdb: 
	docker exec -it postgres17.0 createdb --username=root --owner=root banking_system

dropdb: 
	docker exec -it postgres17.0 dropdb banking_systems

# migrate create -ext sql -dir db/migration -seq add_users 

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

dockerImageBuild:
	docker build -t banking_system:latest .

# we need to use the IPAddress of postgres17 after root:secret@"here". So we first create a network where we put the db so that we can use its IP
# docker network create bank-network 
# docker network connect bank-network postgres17.0
# docker network inspect bank-network 
dockerImageRun:
	docker run --name banking_system --network bank-network -p 8080:8080 -e GIN_MODE=release -e DB_SOURCE="postgresql://root:secret@postgres17.0:5432/banking_system?sslmode=disable" banking_system:latest

.PHONY: postgres createdb dropdb migrateUp migratedown migrateUp1 migratedown1 sqlc test server mock dockerImageBuild dockerImageRun
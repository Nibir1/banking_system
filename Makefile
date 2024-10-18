postgres:
	docker run --name postgres17.0 --network bank-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:17.0-alpine3.20

createdb: 
	docker exec -it postgres17.0 createdb --username=root --owner=root banking_system

dropdb: 
	docker exec -it postgres17.0 dropdb banking_systems

# migrate create -ext sql -dir db/migration -seq add_users == For creating migration files

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/banking_system?sslmode=disable" -verbose up

# migrate database to aws postgres - banking_system
# After creating database on aws RDS we need to migrate our database to it
migrateupAWS:
	migrate -path db/migration -database "postgresql://root:OmHiiaKco2zomi5F2FWK@banking-system-id.c92a8wwkqopk.ap-south-1.rds.amazonaws.com:5432/banking_system" -verbose up

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


dockerComposeUp:
	docker compose up

# To remove all existing containers and networks
dockerComposeDown:
	docker compose down


# aws configure = to configure aws cli to access our aws services
# This command retrieves the secrets value from aws secret management and stores into app.env
awsSecretsToappenv:
	aws secretsmanager get-secret-value --secret-id banking_system --query SecretString --output text | jq -r 'to_entries|map("\(.key)=\(.value)")|.[]' > app.env
# ------------------------
awsECRlogin:
	aws ecr get-login-password | docker login --username AWS --password-stdin 339712865282.dkr.ecr.ap-south-1.amazonaws.com

dockerPullImageFromAwsECR:
	docker pull 339712865282.dkr.ecr.ap-south-1.amazonaws.com/banking_system:b478141455a02e65e6cd375f7a41954792ce87eb

dockerRunPulledImageFromAwsECR:
	docker pull 339712865282.dkr.ecr.ap-south-1.amazonaws.com/banking_system:b478141455a02e65e6cd375f7a41954792ce87eb
# ------------------------

.PHONY: postgres createdb dropdb migrateUp migratedown migrateUp1 migratedown1 sqlc test server mock dockerImageBuild dockerImageRun dockerComposeUp dockerComposeDown migrateupAWS awsSecretsToappenv awsECRlogin dockerPullImageFromAwsECR dockerRunPulledImageFromAwsECR

# openssl rand -hex 64 | head -c 32 == To generate random 32 TOKEN_SYMMETRIC_KEY
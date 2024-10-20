postgres:
	docker run --name postgres17.0 --network bank-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:17.0-alpine3.20

createdb: 
	docker exec -it postgres17.0 createdb --username=root --owner=root banking_system

dropdb: 
	docker exec -it postgres17.0 dropdb banking_systems

# migrate create -ext sql -dir db/migration -seq add_users == For creating migration files for db migration

# ------------------------- Start
# For Local Testing
migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/banking_system?sslmode=disable" -verbose up
migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/banking_system?sslmode=disable" -verbose down

# Remove "source /app/app.env" from start.sh when running locally
dockerImageBuild:
	docker build -t banking_system:latest .
# we need to use the IPAddress of postgres17 after root:secret@"here". So we first create a network where we put the db so that we can use its IP
# docker network create bank-network 
# docker network connect bank-network postgres17.0
# docker network inspect bank-network 
dockerImageRun:
	docker run --name banking_system --network bank-network -p 8080:8080 -e GIN_MODE=release -e DB_SOURCE="postgresql://root:secret@postgres17.0:5432/banking_system?sslmode=disable" banking_system:latest
# ------------------------- End

# migrate database to aws postgres - banking_system
# After creating database on aws RDS we need to migrate our database to it
migrateupAWS:
	migrate -path db/migration -database "postgresql://root:pjzo1nxUNsCYum4qyxUR@banking-system-id.c92a8wwkqopk.ap-south-1.rds.amazonaws.com/banking_system" -verbose up

migrateup1:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/banking_system?sslmode=disable" -verbose up 1

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

# ------------------------ Start
# Builds the api and database into one container locally - need to run the postgres database as well to make the db connection successful
dockerComposeUp:
	docker compose up

# To remove all existing containers and networks
dockerComposeDown:
	docker compose down
# ------------------------ End

# aws configure = to configure aws cli to access our aws services
# This command retrieves the secrets value from aws secret management and stores into app.env
awsSecretsToappenv:
	aws secretsmanager get-secret-value --secret-id banking_system --query SecretString --output text | jq -r 'to_entries|map("\(.key)=\(.value)")|.[]' > app.env

# ------------------------ Start
# Need to update the URI each time there is a new push to Main - since there is a new image on aws ecr - also in eks-> deployment.yaml
awsECRlogin:
	aws ecr get-login-password | docker login --username AWS --password-stdin 339712865282.dkr.ecr.ap-south-1.amazonaws.com

dockerPullImageFromAwsECR:
	docker pull 339712865282.dkr.ecr.ap-south-1.amazonaws.com/banking_system:58bd2384fd28086eed63a794c965692b4b6b4057

dockerRunPulledImageFromAwsECR:
	docker run 339712865282.dkr.ecr.ap-south-1.amazonaws.com/banking_system:58bd2384fd28086eed63a794c965692b4b6b4057

dockerRunPulledImageFromAwsECRlocally:
	docker run -p 8080:8080 339712865282.dkr.ecr.ap-south-1.amazonaws.com/banking_system:58bd2384fd28086eed63a794c965692b4b6b4057

# ------------------------ End

# ------------------------ Start
# kubectl version --client
# ls -l ~/.aws
# cat ~/.kube/config 
# cat ~/.aws/credentials
# aws sts get-caller-identity
# vi ~/.aws/credentials = To edit the credentials
# kubectl cluster-info
# kubectl apply -f eks/aws-auth.yaml --> make sure your in the default profile when executing this
# kubectl get service
# kubectl get pods
# use k9s for better kubernetes cluster usage
# kubectl apply -f eks/deployment.yaml
# kubectl apply -f eks/service.yaml
configawsEKS:
	aws eks update-kubeconfig --name banking_system --region ap-south-1

# To connect kubectl to aws eks cluster
connectawsEKSCluster:
	kubectl config use-context arn:aws:eks:ap-south-1:339712865282:cluster/banking_system
# ------------------------ End

.PHONY: postgres createdb dropdb migrateUp migratedown migrateUp1 migratedown1 sqlc test server mock dockerImageBuild dockerImageRun dockerComposeUp dockerComposeDown migrateupAWS awsSecretsToappenv awsECRlogin dockerPullImageFromAwsECR dockerRunPulledImageFromAwsECR

# openssl rand -hex 64 | head -c 32 == To generate random 32 TOKEN_SYMMETRIC_KEY
# A Go-Gin Based Banking System

## Technologies & Terminologies I Learned While Doing This Project:

- Designing the database schema and generating SQL code using dbdiagram.io
- Setting up Docker, PostgreSQL, and TablePlus for database management
- Writing and running database migrations in Golang
- Generating CRUD operations in Golang using various libraries (db/sql, gorm, sqlx, and sqlc)
- Writing unit tests for database interactions with random data generation
- Implementing database transactions and handling potential deadlocks
- Understanding transaction isolation levels and read phenomena
- Setting up automated testing with Github Actions for Golang and PostgreSQL
- Building a RESTful API using Gin framework
- Configuring the application with Viper library
- Mocking the database for API testing and achieving high test coverage
- Implementing money transfer functionality with custom parameter validation
- Adding user authentication with secure password hashing (Bcrypt)
- Writing robust unit tests using gomock custommatchers
- Examining the benefits of PASETO over JWT for token-based authentication
- Creating and verifying JWT and PASETO tokens in Golang
- Building a login API returning PASETO or JWT access tokens
- Implementing authentication middleware and authorization rules with Gin
- Constructing minimal Docker images with multi-stage Dockerfiles
- Connecting Docker containers using Docker networks
- Utilizing docker-compose to manage service dependencies and startup order
- Creating a free AWS account
- Using IAM for granting access to aws resources
- Automating Docker image building and pushing to AWS ECR using Github Actions
- Setting up a production database on AWS RDS
- Securing sensitive data with AWS Secrets Manager
- Makefile usage for console commands
- K9s for kubectl commands
- Github CI/CD for testing, building & deploying to AWS

## Services of Banking System

This project focuses on building a basic bank backend that offers the following functionalities:

- Account Management: Create and manage bank accounts, storing owner name, balance, and currency.
- Transaction Tracking: Record all balance changes (deposits/withdrawals) as detailed account entries.
- Secure Money Transfer: Perform money transfers between accounts using transactions to ensure data consistency.

## Technologies & Softwares Used for Developing Banking System

- Golang, Go-Gin
- Docker, Postgres [ Docker Image ], TablePlus (DB GUI), sqlc, GoMock(Mockgen)
- JWT & PASETO
- testify/require
- Viper
- Makefile
- Github CI/CD -> AWS ( IAM, ECR, EKS, RDS, Secrets Manager)

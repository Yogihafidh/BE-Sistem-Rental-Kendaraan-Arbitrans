.PHONY: help up down migrate-up migrate-down migrate-version sqlc-compile sqlc-generate swag test fmt vet build dev

DATABASE_URL ?= postgres://postgres:postgres@localhost:5432/rental_db?sslmode=disable

help:
	@echo "Available commands:"
	@echo "  make up               - Start Docker infrastructure (Postgres, Redis, MinIO, CloudBeaver)"
	@echo "  make down             - Stop Docker infrastructure"
	@echo "  make migrate-up       - Apply all database migrations"
	@echo "  make migrate-down     - Roll back last database migration"
	@echo "  make migrate-version  - Show current database migration version"
	@echo "  make sqlc-compile     - Verify SQL queries with sqlc"
	@echo "  make sqlc-generate    - Generate Go code with sqlc"
	@echo "  make swag             - Regenerate Swagger documentation"
	@echo "  make test             - Run all unit tests"
	@echo "  make fmt              - Format Go code"
	@echo "  make vet              - Run go vet"
	@echo "  make dev              - Run backend with live reload (Air)"
	@echo "  make build            - Build executable binary"

up:
	docker compose up -d

down:
	docker compose down

migrate-up:
	migrate -path database/migrations -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path database/migrations -database "$(DATABASE_URL)" down 1

migrate-version:
	migrate -path database/migrations -database "$(DATABASE_URL)" version

sqlc-compile:
	sqlc compile

sqlc-generate:
	sqlc generate

swag:
	swag init -g cmd/api/main.go -o docs

test:
	go test -v ./...

fmt:
	go fmt ./...

vet:
	go vet ./...

dev:
	air

build:
	go build -o bin/api ./cmd/api

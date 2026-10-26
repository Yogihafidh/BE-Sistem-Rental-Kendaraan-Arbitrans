# ARBITRANS — Vehicle Rental System Backend

[![Go Version](https://img.shields.io/badge/Go-1.25.0-blue.svg)](https://golang.org)
[![Framework](https://img.shields.io/badge/Framework-Gin-00ADD8.svg)](https://gin-gonic.com/)
[![Database](https://img.shields.io/badge/Database-PostgreSQL_16-336791.svg)](https://www.postgresql.org/)
[![Architecture](https://img.shields.io/badge/Architecture-Modular_Monolith-orange.svg)](#architecture)
[![License](https://img.shields.io/badge/License-Proprietary-red.svg)]()

Backend service for **ARBITRANS**, an enterprise-grade vehicle rental management system designed to coordinate fleet operations, customer identity verification, reservation scheduling, pricing calculations, status lifecycles, and security auditing.

Built in Go with high-performance concurrency, strict compile-time SQL verification, and a domain-driven Modular Monolith structure.

---

## Documentation Index

Canonical project documentation is maintained in [`docs/`](docs/) and [`AGENTS.md`](AGENTS.md):

| Document | Description |
| :--- | :--- |
| **[AI Developer Instructions](AGENTS.md)** | **Primary entry point for AI**: Required reading, workflow, and verified commands. |
| **[Architecture Blueprint](docs/ARCHITECTURE.md)** | System design, directory structure, Modular Monolith boundaries, request lifecycles. |
| **[Product Requirements (PRD)](docs/PRD.md)** | Functional requirements, actor roles, business rules, and acceptance criteria. |
| **[Engineering & Coding Rules](docs/RULES.md)** | Coding conventions, layer flow, error handling, security, and verification checklist. |


---

## Technology Stack

- **Language:** Go `1.25.0`
- **HTTP Routing & Engine:** [Gin](https://github.com/gin-gonic/gin)
- **Database & Pooling:** PostgreSQL 16 via [pgxpool](https://github.com/jackc/pgx)
- **Type-Safe SQL Generation:** [sqlc](https://sqlc.dev/)
- **Database Migrations:** [golang-migrate](https://github.com/golang-migrate/migrate)
- **Authentication:** Stateless HMAC-SHA256 JWT
- **Password Hashing:** `bcrypt`
- **Caching:** Redis 7 via [go-redis](https://github.com/redis/go-redis)
- **Object Storage:** MinIO (S3-compatible) via [minio-go](https://github.com/minio/minio-go)
- **Structured Logging:** Go Standard Library `log/slog`
- **Validation:** [go-playground/validator](https://github.com/go-playground/validator)
- **Live Reload:** [Air](https://github.com/air-verse/air)
- **API Documentation:** OpenAPI 2.0 / [Swagger UI](https://github.com/swaggo/gin-swagger)
- **Local Infrastructure:** Docker Compose

---

## Quickstart

### 1. Prerequisites
Ensure you have installed:
- Go `>= 1.25.0`
- Docker Desktop & Docker Compose
- `air`, `sqlc`, and `migrate` CLI tools (see [AGENTS.md](AGENTS.md#5-verified-repository-commands))

### 2. Configure Environment
```bash
cp .env.example .env
```

### 3. Start Infrastructure
```bash
# Windows PowerShell
.\scripts\dev.ps1 up

# Linux / macOS
make up

# Direct
docker compose up -d
```
Starts PostgreSQL (`5432`), Redis (`6379`), MinIO (`9000`/`9001`), and CloudBeaver (`8978`).

### 4. Run Migrations
```bash
# Windows PowerShell
.\scripts\dev.ps1 migrate-up

# Linux / macOS
make migrate-up

# Direct CLI:
# PowerShell: $env:DATABASE_URL="postgres://postgres:postgres@localhost:5432/rental_db?sslmode=disable"; migrate -path database/migrations -database $env:DATABASE_URL up
# Bash: export DATABASE_URL="postgres://postgres:postgres@localhost:5432/rental_db?sslmode=disable"; migrate -path database/migrations -database "$DATABASE_URL" up
```

### 5. Start Application with Live Reload
```bash
# Windows PowerShell
.\scripts\dev.ps1 dev

# Linux / macOS
make dev

# Or standard execution
go run ./cmd/api
```

### 6. Verify System Health & Swagger
- **Health Check**: `http://localhost:8080/health`
- **Swagger UI**: `http://localhost:8080/swagger/index.html`
- **CloudBeaver Web DB**: `http://localhost:8978`
- **MinIO Console**: `http://localhost:9001` (User: `admin` / Password: `password123`)

---

## Developer Commands & Task Runners

To simplify daily engineering workflows across different operating systems without remembering complex CLI parameters, unified task runners are provided:

| Task | Windows PowerShell | Linux / macOS (`make`) | Direct Command Equivalent |
| :--- | :--- | :--- | :--- |
| **Start Infrastructure** | `.\scripts\dev.ps1 up` | `make up` | `docker compose up -d` |
| **Stop Infrastructure** | `.\scripts\dev.ps1 down` | `make down` | `docker compose down` |
| **Apply Migrations (Up)** | `.\scripts\dev.ps1 migrate-up` | `make migrate-up` | `migrate -path database/migrations -database "$DATABASE_URL" up` |
| **Roll Back Migration** | `.\scripts\dev.ps1 migrate-down` | `make migrate-down` | `migrate -path database/migrations -database "$DATABASE_URL" down 1` |
| **Migration Version** | `.\scripts\dev.ps1 migrate-version` | `make migrate-version` | `migrate -path database/migrations -database "$DATABASE_URL" version` |
| **Compile & Verify SQL** | `.\scripts\dev.ps1 sqlc-compile` | `make sqlc-compile` | `sqlc compile` |
| **Generate sqlc Code** | `.\scripts\dev.ps1 sqlc-generate` | `make sqlc-generate` | `sqlc generate` |
| **Regenerate Swagger** | `.\scripts\dev.ps1 swag` | `make swag` | `swag init -g cmd/api/main.go -o docs` |
| **Run Unit Tests** | `.\scripts\dev.ps1 test` | `make test` | `go test -v ./...` |
| **Format Code** | `.\scripts\dev.ps1 fmt` | `make fmt` | `go fmt ./...` |
| **Static Analysis** | `.\scripts\dev.ps1 vet` | `make vet` | `go vet ./...` |
| **Live Reload (Air)** | `.\scripts\dev.ps1 dev` | `make dev` | `air` |
| **Build Executable** | `.\scripts\dev.ps1 build` | `make build` | `go build -o bin/api ./cmd/api` |

---

## Testing & Quality Assurance

```bash
# Run all unit tests
go test -v ./...

# Verify Go formatting
go fmt ./...

# Run static analysis
go vet ./...

# Verify sqlc queries
sqlc compile
sqlc diff
```

---

## Architecture Overview

ARBITRANS follows a **Modular Monolith** architecture:

```
cmd/api/main.go
  │
  ├──► configs/             (Typed configuration, environment validation)
  ├──► internal/shared/     (Infrastructure: DB pool, Redis, MinIO, JWT, logger, middleware)
  ├──► database/sqlc/       (Compile-time verified type-safe SQL queries)
  └──► internal/modules/    (Domain-driven business modules)
         ├── auth/          (Authentication, login, token issue)
         ├── user/          (User accounts, role management)
         ├── customer/      (Customers, KYC document verification)
         ├── vehicle/       (Vehicles, categories, media uploads)
         ├── booking/       (Reservations, schedule locking, status lifecycle)
         └── audit/         (System mutation audit trail)
```

For full architecture details, transaction boundaries, and concurrency rules, see [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md).

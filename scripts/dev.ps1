<#
.SYNOPSIS
    ARBITRANS Developer Task Runner for Windows PowerShell.
.DESCRIPTION
    Conveniently manage Docker infrastructure, migrations, sqlc, Swagger, and testing.
.EXAMPLE
    .\scripts\dev.ps1 up
    .\scripts\dev.ps1 migrate-up
    .\scripts\dev.ps1 dev
#>

param (
    [Parameter(Position=0, Mandatory=$true)]
    [ValidateSet("help", "up", "down", "migrate-up", "migrate-down", "migrate-version", "sqlc-compile", "sqlc-generate", "swag", "test", "test-race", "lint", "vulncheck", "fmt", "vet", "build", "dev")]
    [string]$Command
)

$DatabaseURL = if ($env:DATABASE_URL) { $env:DATABASE_URL } else { "postgres://postgres:postgres@localhost:5432/rental_db?sslmode=disable" }

switch ($Command) {
    "help" {
        Write-Host "ARBITRANS Development Commands:" -ForegroundColor Cyan
        Write-Host "  .\scripts\dev.ps1 up              - Start Docker infrastructure"
        Write-Host "  .\scripts\dev.ps1 down            - Stop Docker infrastructure"
        Write-Host "  .\scripts\dev.ps1 migrate-up      - Apply all migrations"
        Write-Host "  .\scripts\dev.ps1 migrate-down    - Roll back 1 migration"
        Write-Host "  .\scripts\dev.ps1 migrate-version - Check current schema version"
        Write-Host "  .\scripts\dev.ps1 sqlc-compile    - Statically check queries"
        Write-Host "  .\scripts\dev.ps1 sqlc-generate   - Regenerate Go code with sqlc"
        Write-Host "  .\scripts\dev.ps1 swag            - Regenerate Swagger documentation"
        Write-Host "  .\scripts\dev.ps1 test            - Run all unit tests"
        Write-Host "  .\scripts\dev.ps1 test-race       - Run unit tests with race detector"
        Write-Host "  .\scripts\dev.ps1 lint            - Run golangci-lint (includes gosec)"
        Write-Host "  .\scripts\dev.ps1 vulncheck       - Check dependencies with govulncheck"
        Write-Host "  .\scripts\dev.ps1 fmt             - Format Go source code"
        Write-Host "  .\scripts\dev.ps1 vet             - Run static analysis"
        Write-Host "  .\scripts\dev.ps1 dev             - Run backend with live reload (Air)"
        Write-Host "  .\scripts\dev.ps1 build           - Build executable binary"
    }
    "up" {
        docker compose up -d
    }
    "down" {
        docker compose down
    }
    "migrate-up" {
        migrate -path database/migrations -database $DatabaseURL up
    }
    "migrate-down" {
        migrate -path database/migrations -database $DatabaseURL down 1
    }
    "migrate-version" {
        migrate -path database/migrations -database $DatabaseURL version
    }
    "sqlc-compile" {
        sqlc compile
    }
    "sqlc-generate" {
        sqlc generate
    }
    "swag" {
        swag init -g cmd/api/main.go -o docs
    }
    "test" {
        go test -v ./...
    }
    "test-race" {
        go test -v -race ./...
    }
    "lint" {
        golangci-lint run
    }
    "vulncheck" {
        govulncheck ./...
    }
    "fmt" {
        go fmt ./...
    }
    "vet" {
        go vet ./...
    }
    "dev" {
        air
    }
    "build" {
        go build -o bin/api.exe ./cmd/api
    }
}

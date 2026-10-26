# ARBITRANS — AI Developer & Agent Instructions

This document is the **primary entry point** for AI assistants (Gemini, Antigravity, etc.) working on the **ARBITRANS** backend.

---

## 1. Project Context

**ARBITRANS** is a commercial vehicle rental management backend written in **Go 1.25** using a **Modular Monolith** architecture.
- **HTTP Engine:** Gin (`github.com/gin-gonic/gin`)
- **Database & Pooling:** PostgreSQL 16 via `pgxpool` (`jackc/pgx/v5`)
- **Type-Safe SQL:** `sqlc` (`database/queries/` $\rightarrow$ `database/sqlc/`)
- **Database Migrations:** `golang-migrate` (`database/migrations/`)
- **Authentication & Security:** Stateless JWT (HMAC-SHA256) & bcrypt
- **Cache & Storage:** Redis 7 & MinIO (S3-compatible)
- **Local Infrastructure:** Docker Compose (`docker-compose.yml`)

---

## 2. Required Reading Before Implementation

Before proposing or making changes, consult the relevant canonical document:
1. **Engineering & Coding Rules**: Read [`docs/RULES.md`](docs/RULES.md) for Go conventions, error handling, layer flow, and test criteria.
2. **System Architecture**: Read [`docs/ARCHITECTURE.md`](docs/ARCHITECTURE.md) for module boundaries, startup flow, and layer responsibilities.
3. **Business Requirements**: Read [`docs/PRD.md`](docs/PRD.md) when working on business features, statuses, or pricing calculations.
4. **Active Code Inspection**: Inspect the actual Go source code, migrations, and queries affected by the task.

> *Rule: Consult only the relevant sections needed for the task rather than rereading every document for minor edits.*

---

## 3. Standard Development Workflow

1. **Understand**: Clarify requirements and verify acceptance criteria.
2. **Inspect**: Read existing code, migrations, queries, and tests before modifying.
3. **Plan**: Outline changes and identify module boundaries before editing.
4. **Implement**: Make small, focused changes following the layer flow:
   $$\text{Route} \longrightarrow \text{Handler} \longrightarrow \text{Service} \longrightarrow \text{Repository} \longrightarrow \text{Database}$$
5. **Test**: Add or update unit tests covering success, validation, and error paths.
6. **Verify**: Execute formatting, static analysis, unit tests, and sqlc checks.
7. **Document**: Update `docs/ARCHITECTURE.md`, `docs/PRD.md`, or `docs/RULES.md` if behavior changes.
8. **Report**: Summarize changes, executed checks, observed outputs, and assumptions.

---

## 4. Non-Negotiable Rules

- **Respect Module Boundaries**: Modules live in `internal/modules/<domain>/`. Never execute direct cross-module SQL mutations.
- **Never Invent Requirements**: If business rules or pricing formulas are unclear, mark them as **Needs clarification** and ask.
- **Never Expose Secrets**: Never commit `.env` or hardcode credentials, tokens, or passwords.
- **No Destructive Database Commands**: Never run `migrate drop`, database reset, or volume deletion without explicit approval.
- **Never Claim Unverified Results**: Only report test passes or command successes that were actually executed.
- **Actual vs. Planned**: Never treat database tables or sqlc code as proof that business features are complete. Verify active handlers in `internal/modules/`.
- **Keep Documentation Minimal**: Never create extra documentation files without explicit user approval.

---

## 5. Verified Repository Commands

Execute commands from the repository root:

| Task | Command |
| :--- | :--- |
| **Start Infrastructure** | `docker compose up -d` |
| **Stop Infrastructure** | `docker compose down` |
| **Run Migrations (Up)** | `migrate -path database/migrations -database "$DATABASE_URL" up` |
| **Roll Back Migration** | `migrate -path database/migrations -database "$DATABASE_URL" down 1` |
| **Verify sqlc Queries** | `sqlc compile` and `sqlc diff` |
| **Regenerate sqlc Code**| `sqlc generate` |
| **Regenerate Swagger** | `swag init -g cmd/api/main.go -o docs` |
| **Run Unit Tests** | `go test -v ./...` |
| **Run Linter** | `golangci-lint run` |
| **Audit Vulnerabilities** | `govulncheck ./...` |
| **Format Code** | `go fmt ./...` |
| **Static Analysis** | `go vet ./...` |
| **Compile Project** | `go build ./...` |
| **Live Reload (Air)** | `air` |

*PowerShell equivalent: Use `.\scripts\dev.ps1 <command>` (e.g., `.\scripts\dev.ps1 test`).*

---

## 6. Documentation Maintenance

When changes affect architecture, business rules, or engineering standards, update the corresponding canonical document:
- **Structure / Components / Layers** $\rightarrow$ [`docs/ARCHITECTURE.md`](docs/ARCHITECTURE.md)
- **Business Features / Rules / Statuses** $\rightarrow$ [`docs/PRD.md`](docs/PRD.md)
- **Coding Conventions / Verification** $\rightarrow$ [`docs/RULES.md`](docs/RULES.md)

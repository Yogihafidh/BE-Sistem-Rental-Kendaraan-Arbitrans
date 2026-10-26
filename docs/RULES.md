# ARBITRANS — Engineering & AI Coding Rules

This document establishes the binding engineering standards, coding conventions, architectural rules, and verification criteria for developing the **ARBITRANS** backend. All human developers and AI assistants must follow these rules.

---

## 1. General Development Rules

1. **Inspect Before Modifying**: Always read and understand the existing code, migrations, and queries before proposing or making changes. Never assume an implementation exists or works without verifying it.
2. **Implementation Status vs. Planned Architecture**:
   - Never assume a feature, business rule, or endpoint is implemented merely because a database table, SQL query, DTO struct, or folder exists.
   - The presence of `database/queries/*.sql` and generated `database/sqlc/*.go` proves only that SQL query definitions exist. It does **not** mean the corresponding business module, domain validation, service logic, or HTTP endpoint is implemented.
   - Always verify active handlers and services in `internal/modules/` and route registrations in `cmd/api/main.go`.
3. **Make Small, Focused Changes**: Implement only what is necessary to fulfill the immediate task. Never introduce broad, unsolicited refactoring or stylistic reformatting.
4. **Avoid Unnecessary Abstractions**: Do not create generic interfaces, reflection-based utilities, event buses, or design patterns unless they solve an active requirement or simplify unit testing.
5. **Never Invent Requirements**: If business rules, pricing logic, role permissions, or customer flows are ambiguous, do not guess or silently decide them. Mark them as **Needs clarification** and consult `docs/PRD.md`.
6. **Never Claim Unverified Results**: Do not report that a test passed, a command succeeded, or a service is running unless you actually executed the command and observed the result.
7. **Identify Assumptions**: State all assumptions explicitly in the final response.

---

## 2. Go Coding Conventions

### 2.1 Package Organization and Naming
- Package names must be lowercase, single-word, and reflect their purpose (e.g., `vehicle`, `booking`, `database`).
- File names must use snake_case (e.g., `vehicle_status.go`, `password_test.go`).
- Exported functions, structs, and types must use CamelCase with an explanatory Go doc comment:
  ```go
  // Service provides operations for managing the vehicle fleet.
  type Service struct { ... }
  ```

### 2.2 Error Handling and Propagation
- Always handle errors explicitly. Never use `_` to discard errors returned by standard library, database, or external service calls.
- Wrap errors with context using `%w` to preserve the error chain:
  ```go
  if err != nil {
      return fmt.Errorf("failed to retrieve vehicle by id %s: %w", id, err)
  }
  ```
- Use `errors.Is()` or `errors.As()` for error inspection; never compare error strings directly with `err.Error() == "..."`.
- Define sentinel errors in package scope when callers need to branch on specific failure cases:
  ```go
  var ErrVehicleNotFound = errors.New("vehicle not found")
  ```

### 2.3 Context Propagation
- `ctx context.Context` must always be the **first parameter** of any function performing I/O, database queries, Redis calls, or external service requests:
  ```go
  func (s *Service) GetVehicle(ctx context.Context, id string) (*VehicleDTO, error)
  ```
- Never store a `context.Context` inside a struct.
- Respect cancellation: pass `ctx` directly to `sqlc.Queries`, `pgxpool.Pool`, and HTTP clients.

### 2.4 Dependency Injection & Idiomatic Go Contracts
- Follow the Go principle: **"Accept interfaces, return structs"**.
- Use explicit constructor injection with a `New...` function returning concrete structs:
  ```go
  func NewService(repo *Repository, minio *sharedMinIO.Client, logger *slog.Logger) *Service
  ```
- Do not create 1:1 interfaces for every struct upfront. Let the *consumer* package define interfaces where needed for test mocking or multiple implementations.
- Avoid global mutable variables, except for stateless package-level sentinel errors or constants.

### 2.5 Structured Logging (`log/slog`)
- Never use `fmt.Println`, `fmt.Printf`, or standard `log.Print` in business logic or handlers.
- Use structured logging via `log/slog` with strongly-typed attributes:
  ```go
  logger.Info("Vehicle status updated",
      slog.String("vehicle_id", vehicleID),
      slog.String("new_status", status),
      slog.String("changed_by", userID),
  )
  ```
- Use `Error` level for server failures (HTTP 5xx), `Warn` level for degraded states or client rejection (HTTP 4xx), and `Info`/`Debug` for lifecycle events.

### 2.6 Validation and DTOs
- Validate request payloads at the handler boundary using `internal/shared/validator`:
  ```go
  var req CreateVehicleRequest
  if err := c.ShouldBindJSON(&req); err != nil {
      fieldErrors := validator.FormatValidationErrors(err)
      response.BadRequest(c, "Validation failed", fieldErrors)
      return
  }
  ```
- DTO structs must include `validate` tags (e.g., `validate:"required,min=3"`).

### 2.7 API Response Consistency
- All HTTP endpoints must return the standardized response envelope from `internal/shared/response`:
  - `response.Success(c, statusCode, message, data)`
  - `response.SuccessWithMeta(c, statusCode, message, data, meta)`
  - `response.BadRequest(c, message, errors)`
  - `response.Unauthorized(c, message)`
  - `response.Forbidden(c, message)`
  - `response.NotFound(c, message)`
  - `response.InternalServerError(c, message)`

---

## 3. Modular Monolith Rules

### 3.1 Strict Layered Request Flow
All business features must adhere to the unidirectional flow:

$$\text{Route} \longrightarrow \text{Handler} \longrightarrow \text{Service} \longrightarrow \text{Repository} \longrightarrow \text{Database}$$

- **Handlers**: Only parse input, validate DTOs, extract auth claims, call Service, and write response. Handlers must never contain business rules or direct SQL queries.
- **Services**: Enforce business rules, check permissions, calculate prices, coordinate transactions, and call repositories. Services must remain independent of Gin (`*gin.Context`).
- **Repositories**: Wrap `*sqlc.Queries` or database transactions. Never contain business logic.

### 3.2 Module Isolation & Communication Rules
1. **No Direct Table Modification**: A module must never execute SQL queries that update or delete data in tables owned by another module.
2. **Explicit Contracts Without Forced Interfaces**:
   - Inter-module communication must use explicit public contracts (exported concrete methods, structs, and DTOs).
   - If Module A needs data or actions from Module B (e.g., `booking` checking `vehicle` availability), Module A calls Module B's exported Service methods directly or defines a local consumer interface if mocking is needed for unit testing.
3. **No Circular Package Imports**: If Module A imports Module B, Module B must never import Module A. If two modules require shared types, extract them into `internal/shared/` or common DTO packages.
4. **Shared Infrastructure vs. Module Code**:
   - `internal/shared/`: Generic infrastructure (database connection, Redis, MinIO, JWT, middlewares). Contains zero business logic.
   - `internal/modules/<domain>/`: Domain-specific business logic.

---

## 4. Database and sqlc Rules

### 4.1 Migration Rules & Limitations
1. **Never Edit Already-Applied Migrations**: Once a migration file has been applied or committed, never modify its contents.
2. **Create New Migrations for Schema Changes**:
   ```bash
   migrate create -ext sql -dir database/migrations -seq <migration_name>
   ```
3. **Reversible Migrations with Documented Limitations**:
   - Migrations should be reversible whenever practical (`.up.sql` and corresponding `.down.sql`).
   - Acknowledge realistic PostgreSQL limitations:
     - Enum additions (`ALTER TYPE ... ADD VALUE`) cannot be undone via standard SQL commands in PostgreSQL.
     - Destructive operations (dropping columns or tables) cause permanent data loss if rolled back.
     - One-way data transformations cannot be automatically restored.
   - When a migration cannot be safely or cleanly reversed, document the rollback limitation, risks, and manual recovery procedure as comments in the `.down.sql` file.
4. **Never Run Destructive Database Commands**: Never execute `migrate drop`, database reset, or volume deletion without explicit user approval.

### 4.2 Query and sqlc Rules
1. **No Raw SQL in Application Code**: All SQL must be authored in `database/queries/*.sql`.
2. **Never Edit Generated Code**: Files in `database/sqlc/*.go` are overwritten on generation. Do not edit them manually.
3. **Always Verify and Regenerate**:
   ```bash
   sqlc compile    # Verify SQL syntax against migrations
   sqlc diff       # Verify generated code matches SQL files
   sqlc generate   # Regenerate Go code after changing queries
   ```
4. **Transactions**: Use `sharedDatabase.WithTransaction(ctx, pool, fn)` whenever two or more database operations must succeed or fail together.

---

## 5. Security Rules

1. **Never Commit Secrets**: Never commit `.env` files, private keys, API credentials, or database passwords to version control.
2. **Enforce Authentication**: Protected endpoints must use `middleware.AuthRequired()`.
3. **Enforce Role-Based Access Control**: Sensitive actions must check roles via `middleware.RequireRoles("admin", "staff")`.
4. **Password Security**: Always hash passwords using `utils.HashPassword` (bcrypt standard cost). Never store plaintext passwords.
5. **Token Security**: JWT tokens must use HMAC-SHA256 with the verified TokenManager in `internal/shared/jwt`. `alg: none` and algorithm switching attacks are strictly blocked.
6. **File Upload Sanitization**: Always sanitize multipart upload filenames using `filepath.Base` before creating MinIO object keys to prevent directory traversal attacks.
7. **Never Log Sensitive Data**: Do not log plaintext passwords, JWT tokens, customer NIKs, or identity documents in structured logs.

---

## 6. Testing and Completion Criteria

A task is considered **complete** only when all of the following criteria are verified:

1. **Unit Tests**:
   - New or modified logic has corresponding unit tests (`*_test.go`).
   - All tests pass:
     ```bash
     go test -v ./...
     ```
2. **Static Analysis & Formatting**:
   - Code is formatted cleanly using the standard Go formatting command:
     ```bash
     go fmt ./...
     ```
   - Static analysis passes with zero warnings:
     ```bash
     go vet ./...
     ```
3. **Compilation**:
   - Entire workspace compiles cleanly:
     ```bash
     go build ./...
     ```
4. **Database & SQL Consistency**:
   - If SQL queries or migrations were added or modified:
     ```bash
     sqlc compile && sqlc diff
     ```
5. **Documentation Maintenance**:
   - If architecture, business requirements, or workflows changed, update the relevant canonical file (`ARCHITECTURE.md`, `PRD.md`, or `RULES.md`).
6. **Clean Diff**:
   - Review `git diff` to ensure no stray debug statements, unused imports, or accidental changes were made.

---

## 7. Required AI Final Response Format

When completing an implementation task, the AI must conclude with a structured summary containing:

1. **Summary of Changes**: High-level explanation of what was implemented.
2. **Files Created or Modified**: Explicit list of file paths with 1-sentence explanations.
3. **Important Implementation Decisions**: Key trade-offs, architecture choices, or security safeguards made.
4. **Executed Verification Checks & Results**: Exact commands run (`go test`, `go vet`, `sqlc compile`) and their observed outputs.
5. **Remaining Issues, Assumptions, and Risks**: Any items that require clarification or follow-up.

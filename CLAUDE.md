# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

---

## Behavioral Guidelines

**Tradeoff:** These guidelines bias toward caution over speed. For trivial tasks, use judgment.

### Think Before Coding

**Don't assume. Don't hide confusion. Surface tradeoffs.**

Before implementing:
- State your assumptions explicitly. If uncertain, ask.
- If multiple interpretations exist, present them — don't pick silently.
- If a simpler approach exists, say so. Push back when warranted.
- If something is unclear, stop. Name what's confusing. Ask.

### Simplicity First

**Minimum code that solves the problem. Nothing speculative.**

- No features beyond what was asked.
- No abstractions for single-use code.
- No "flexibility" or "configurability" that wasn't requested.
- No error handling for impossible scenarios.
- If you write 200 lines and it could be 50, rewrite it.

Ask yourself: "Would a senior engineer say this is overcomplicated?" If yes, simplify.

### Surgical Changes

**Touch only what you must. Clean up only your own mess.**

When editing existing code:
- Don't "improve" adjacent code, comments, or formatting.
- Don't refactor things that aren't broken.
- Match existing style, even if you'd do it differently.
- If you notice unrelated dead code, mention it — don't delete it.

When your changes create orphans:
- Remove imports/variables/functions that YOUR changes made unused.
- Don't remove pre-existing dead code unless asked.

The test: Every changed line should trace directly to the user's request.

### Goal-Driven Execution

**Define success criteria. Loop until verified.**

Transform tasks into verifiable goals:
- "Add validation" → "Write tests for invalid inputs, then make them pass"
- "Fix the bug" → "Write a test that reproduces it, then make it pass"
- "Refactor X" → "Ensure tests pass before and after"

For multi-step tasks, state a brief plan:
```
1. [Step] → verify: [check]
2. [Step] → verify: [check]
3. [Step] → verify: [check]
```

Strong success criteria let you loop independently. Weak criteria ("make it work") require constant clarification.

---

## Commands

```bash
make serve              # go run . serve — start server without rebuilding binary
make build              # compile to bin/digimap-backend with version ldflags
make test               # go test ./... -count=1
make test-coverage      # run tests + print total coverage percentage
make test-integration   # go test ./tests/integration/... -tags integration
make migrate-up         # apply all pending migrations (calls go run . migrate up)
make migrate-down       # revert last migration
make migrate-create NAME=create_foo_table  # scaffold next numbered .up/.down pair
make docker-up          # start postgres + redis, then run migrate container
```

**Run a single test:**
```bash
go test ./internal/service/... -run TestAuthService_Login -v
```

**Local setup:** copy `config/env/local.env.example` to `.env` in the repo root before running.

---

## Architecture

The codebase is a clean-layered Go API (Gin + pgx/v5 + PostgreSQL):

```
Handler → Service → Repository → Domain
               ↓
          Platform adapters
```

- **`internal/domain/`** — pure Go structs and sentinel errors. No DB or JSON tags. All errors are wrapped with `domain.AppError` (`NewNotFound`, `NewConflict`, etc.); handlers map these to HTTP codes via `respondError()`.
- **`internal/repository/interfaces.go`** — all repository interfaces in one file. Services depend only on these interfaces.
- **`internal/repository/postgres/`** — pgx implementations. Each file owns one resource. `db.go` provides shared helpers: `newID()` (UUIDv7), `softDelete()`, `IsUniqueViolation()`, and the `scanner` interface used by all `scanX()` helpers.
- **`internal/service/`** — business logic. Each file declares an interface + unexported struct. Constructors return the interface.
- **`internal/handler/`** — thin HTTP layer. `router.go` wires all routes and DI. `helpers.go` has `parseVenueID`, `paginationFromQuery`, `respondError`, `bindingErrors`. `errors.go` maps `domain.AppError` to HTTP codes.
- **`internal/dto/`** — JSON request/response types with `binding:` tags. `common.go` defines the response envelope: `OK(data)`, `Paginated(items, total, page, pageSize)`, `Fail(code, msg)`.
- **`internal/platform/`** — adapters: `database/`, `cache/`, `storage/`, `email/`, `firebase/`. Each exposes an interface + real implementation + `Log*` dev stub that prints instead of calling the real service.

### Adding a new resource

Follow this pattern for every new domain object:

1. **Migration** — `make migrate-create NAME=create_foo_table`, write SQL in `migrations/NNNNNN_create_foo_table.{up,down}.sql`. Every table needs `id UUID`, `created_at`, `updated_at` (auto-trigger), `deleted_at` (soft-delete).
2. **Domain** — struct in `internal/domain/foo.go`, no tags.
3. **Repository** — add interface methods to `internal/repository/interfaces.go`, implement in `internal/repository/postgres/foo_repo.go`. Use `softDelete()` helper, `newID()` for PK, the `scanner` interface for scan helpers.
4. **Service** — `internal/service/foo_service.go`: interface + struct + constructor returning the interface.
5. **DTO** — `internal/dto/foo_dto.go`: `FooResponse`, `FooRequest`, `FooToResponse()` converter.
6. **Handler** — `internal/handler/foo_handler.go`: unexported struct + `newFooHandler(svc)` constructor.
7. **Wire** — add to `Dependencies` struct in `router.go`, register routes, init repo + service in `cmd/serve.go`.

### RBAC

Every venue-scoped route goes through `middleware.VenueAccess(userRepo, minRole)`. Roles are numeric: `RoleViewer=1 < RoleEditor=2 < RoleOwner=3 < RoleSystemAdmin=4`. System admins bypass the check. The `protected` group (JWT only, no venue check) and `adminOnly` group (system admin only) are pre-built in `router.go`.

### Router groups

`router.go` defines five route groups:
- `auth` — unauthenticated, auth-rate-limited (`authRL`)
- `protected` — JWT only, no venue check
- `adminOnly` — JWT + system admin required, API-rate-limited (`apiRL`)
- `venues` — JWT + per-route `VenueAccess(minRole)`, API-rate-limited
- `public` — no auth, public-rate-limited (`publicRL`) — used for analytics tracking

### Response envelope

All responses use `dto.Response{Code, Data, Messages}`. HTTP `2xx` always carries `Code: 0`. Errors carry a non-zero `AppCode` (1000–1007 range). Do not use raw `gin.H{}` for responses — always use `dto.OK`, `dto.Fail`, or `dto.Paginated`.

### Migrations

SQL files are embedded via `//go:embed migrations/*.sql` in `cmd/migrate.go` and applied with `golang-migrate`. The `create_updated_at_trigger(table)` helper function (defined in migration `000001`) must be called for every new table.

### Testing

Service tests use hand-written mocks in `internal/repository/mocks/mocks.go` (not generated — extend manually when adding new repository interface methods). S3 tests use `StorerMock` defined in `internal/platform/storage/s3.go`.

### Platform stubs

In local/dev mode all external services use `Log*` stubs (`LogSender` for email, `LogStorer` for S3, `LogPusher` for Firebase). The `cmd/serve.go` DI wiring is where real implementations get swapped in when credentials are present.

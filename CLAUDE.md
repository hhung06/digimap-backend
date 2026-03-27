# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

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

### Response envelope

All responses use `dto.Response{Code, Data, Messages}`. HTTP `2xx` always carries `Code: 0`. Errors carry a non-zero `AppCode` (1000–1007 range). Do not use raw `gin.H{}` for responses — always use `dto.OK`, `dto.Fail`, or `dto.Paginated`.

### Migrations

SQL files are embedded via `//go:embed migrations/*.sql` in `cmd/migrate.go` and applied with `golang-migrate`. The `create_updated_at_trigger(table)` helper function (defined in migration `000001`) must be called for every new table.

### Platform stubs

In local/dev mode all external services use `Log*` stubs (`LogSender` for email, `LogStorer` for S3, `LogPusher` for Firebase). The `cmd/serve.go` DI wiring is where real implementations get swapped in when credentials are present.

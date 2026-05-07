# digimap-backend

Indoor mapping REST API built with Go (Gin + pgx + PostgreSQL).

## Requirements

- Go 1.22+
- Docker + Docker Compose
- [`swag`](https://github.com/swaggo/swag) CLI for regenerating Swagger docs

## Quick Start

```bash
# Copy environment config
cp config/env/local.env.example .env

# Start PostgreSQL, Redis, and run migrations
make docker-up

# Start the API server (hot-reload via Air)
make docker-up-app
```

The API is available at `http://localhost:8080`.
Swagger UI: `http://localhost:8080/swagger/index.html`

## Development (local Go, Docker infra)

```bash
# Start only postgres + redis
make docker-up

# Run the server locally with hot-reload
make serve
```

## Commands

| Command | Description |
|---------|-------------|
| `make serve` | Run server with `go run` (no rebuild) |
| `make build` | Compile binary to `bin/` |
| `make test` | Run unit tests |
| `make test-coverage` | Run tests and print total coverage |
| `make test-integration` | Run integration tests (requires Docker) |
| `make load-test-local` | Run local k6 public load test (requires Docker) |
| `make load-test-staging` | Run staging k6 load test (requires Docker) |
| `make load-test-scenario ENV=staging SCENARIO=app PROFILE=peak` | Run one load scenario/profile |
| `make migrate-up` | Apply all pending migrations |
| `make migrate-down` | Revert last migration |
| `make migrate-create NAME=xxx` | Scaffold next numbered migration pair |
| `make docker-up` | Start postgres + redis + run migrations |
| `make docker-up-app` | Start full stack (postgres + redis + app) detached |
| `make docker-down` | Stop all containers |

## Run a single test

```bash
go test ./internal/service/... -run TestAuthService_Login -v
```

## Swagger

After modifying handler annotations, regenerate the docs:

```bash
swag init
```

Then rebuild the app container:

```bash
docker-compose up -d --build app
```

## Create a system admin

```bash
go run . create-admin
```

## Architecture

```
Handler → Service → Repository → Domain
               ↓
          Platform adapters (DB, Redis, S3, Email, Firebase)
```

- **`internal/domain/`** — pure Go structs and sentinel errors
- **`internal/repository/postgres/`** — pgx implementations
- **`internal/service/`** — business logic
- **`internal/handler/`** — Gin HTTP layer, routing, middleware
- **`internal/dto/`** — JSON request/response types
- **`internal/platform/`** — adapters for external services
- **`migrations/`** — SQL migration files (golang-migrate)

## Environment Variables

See `config/env/local.env.example` for all available configuration options.

## EC2 Monitoring (Prometheus + Grafana)

Infrastructure monitoring is defined in `deploy/stack-infra.yml` and includes:
- `node-exporter` (host metrics per EC2 node)
- `prometheus` (metrics scrape + storage)
- `grafana` (dashboards + alerting UI)

Setup:

```bash
cp deploy/.env.deploy.example deploy/.env.deploy
# edit deploy/.env.deploy:
# - PROMETHEUS_HOST
# - GRAFANA_HOST
# - GRAFANA_ADMIN_PASSWORD
```

Deploy/update infra stack:

```bash
docker network create --driver overlay --attachable traefik-public
docker stack deploy -c deploy/stack-infra.yml digimap-infra
```

After deploy:
- Prometheus: `https://$PROMETHEUS_HOST`
- Grafana: `https://$GRAFANA_HOST` (user: `admin`, password: `GRAFANA_ADMIN_PASSWORD`)

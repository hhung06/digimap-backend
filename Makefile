.PHONY: build build-alpine run serve clean test test-integration \
        migrate-up migrate-down migrate-version migrate-create \
        docker-build docker-up docker-down \
        swarm-init swarm-deploy-infra swarm-deploy swarm-rollback swarm-status \
        help

BIN_NAME   = digimap-backend
VERSION   := $(shell grep "const Version " version/version.go | sed -E 's/.*"(.+)"$$/\1/')
GIT_COMMIT = $(shell git rev-parse HEAD)
GIT_DIRTY  = $(shell test -n "`git status --porcelain`" && echo "+CHANGES" || true)
BUILD_DATE = $(shell date '+%Y-%m-%d-%H:%M:%S')
IMAGE_NAME = hungnh14/digimap-backend

LDFLAGS = -X github.com/hhung06/digimap-backend/version.GitCommit=$(GIT_COMMIT)$(GIT_DIRTY) \
          -X github.com/hhung06/digimap-backend/version.BuildDate=$(BUILD_DATE)

default: build

help:
	@echo 'Usage:'
	@echo '  make build            Compile binary to ./bin/'
	@echo '  make build-alpine     Compile static binary for Alpine'
	@echo '  make run              Run the server locally (requires .env)'
	@echo '  make test             Run unit tests'
	@echo '  make test-integration Run integration tests (requires Docker)'
	@echo '  make migrate-up       Apply all pending migrations'
	@echo '  make migrate-down     Revert last migration'
	@echo '  make migrate-version  Show current migration version'
	@echo '  make migrate-create NAME=xxx  Create new migration pair'
	@echo '  make docker-up        Start local docker-compose stack'
	@echo '  make docker-down      Stop local docker-compose stack'
	@echo '  make docker-build     Build production Docker image'
	@echo '  make swarm-init       One-time Swarm + network setup'
	@echo '  make swarm-deploy-infra  Deploy Traefik + Postgres + Redis stack'
	@echo '  make swarm-deploy TAG=x  Blue-green deploy with image tag'
	@echo '  make swarm-rollback   Rollback to previous slot instantly'
	@echo '  make swarm-status     Show current blue-green state'
	@echo '  make clean            Remove build artifacts'

build:
	@echo "building $(BIN_NAME) $(VERSION)"
	go build -ldflags "$(LDFLAGS)" -o bin/$(BIN_NAME) .

build-alpine:
	@echo "building $(BIN_NAME) $(VERSION) for alpine"
	CGO_ENABLED=0 GOOS=linux go build \
		-ldflags '-w -s $(LDFLAGS)' \
		-o bin/$(BIN_NAME) .

run: build
	./bin/$(BIN_NAME) serve

serve:
	go run . serve

test:
	go test ./... -count=1

test-coverage:
	go test ./... -count=1 -coverprofile=coverage.out
	go tool cover -func=coverage.out | grep total

test-integration:
	go test ./tests/integration/... -count=1 -tags integration -v

# ── Migrations ────────────────────────────────────────────────────────────────
migrate-up:
	go run . migrate up

migrate-down:
	go run . migrate down 1

migrate-version:
	go run . migrate version

migrate-create:
	@[ "$(NAME)" ] || (echo "usage: make migrate-create NAME=create_foo_table" && exit 1)
	@NEXT=$$(ls migrations/*.up.sql 2>/dev/null | wc -l | tr -d ' '); \
	NEXT=$$((NEXT + 1)); \
	SEQ=$$(printf "%06d" $$NEXT); \
	touch migrations/$${SEQ}_$(NAME).up.sql migrations/$${SEQ}_$(NAME).down.sql; \
	echo "created migrations/$${SEQ}_$(NAME).{up,down}.sql"

# ── Docker ────────────────────────────────────────────────────────────────────
docker-up:
	docker-compose up -d postgres redis
	docker-compose --profile migrate run --rm migrate

docker-down:
	docker-compose down

docker-up-app:
	docker-compose up

docker-build:
	docker build \
		--build-arg GIT_COMMIT=$(GIT_COMMIT) \
		--build-arg BUILD_DATE=$(BUILD_DATE) \
		-t $(IMAGE_NAME):$(GIT_COMMIT) \
		-t $(IMAGE_NAME):$(VERSION) \
		-t $(IMAGE_NAME):latest \
		.

clean:
	@test ! -e bin/$(BIN_NAME) || rm bin/$(BIN_NAME)
	@rm -f coverage.out
	@rm -rf tmp/

# ── Docker Swarm ───────────────────────────────────────────────────────────────
# One-time Swarm setup: creates the overlay network Traefik uses
swarm-init:
	docker swarm init || true
	docker network create --driver overlay --attachable traefik-public || true
	docker node update --label-add postgres=true $$(docker node ls -q | head -1)
	@echo "Swarm initialised. Copy deploy/.env.deploy.example to deploy/.env.deploy and fill values."

# Deploy the shared infrastructure stack (Traefik, Postgres, Redis)
swarm-deploy-infra:
	docker stack deploy \
		--with-registry-auth \
		--compose-file deploy/stack-infra.yml \
		digimap-infra

# Blue-green deploy: build + push + deploy to idle slot + cutover
# Usage: make swarm-deploy TAG=abc1234
swarm-deploy:
	@[ "$(TAG)" ] || (echo "usage: make swarm-deploy TAG=<image-tag>" && exit 1)
	$(MAKE) docker-build GIT_COMMIT=$(TAG)
	docker push $(IMAGE_NAME):$(TAG)
	chmod +x deploy/deploy.sh
	deploy/deploy.sh $(IMAGE_NAME):$(TAG)

# Instant rollback to previous slot
swarm-rollback:
	chmod +x deploy/rollback.sh
	deploy/rollback.sh

# Show current deployment state
swarm-status:
	@echo "=== Active slot ==="
	@cat /var/lib/digimap/active-slot 2>/dev/null || echo "(state file not found)"
	@echo ""
	@echo "=== Services ==="
	@docker service ls --filter name=digimap
	@echo ""
	@echo "=== Blue tasks ==="
	@docker service ps digimap-blue_app 2>/dev/null || echo "(not deployed)"
	@echo ""
	@echo "=== Green tasks ==="
	@docker service ps digimap-green_app 2>/dev/null || echo "(not deployed)"

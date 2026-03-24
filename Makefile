.PHONY: build build-alpine run serve clean test test-integration \
        migrate-up migrate-down migrate-version migrate-create \
        docker-build docker-up docker-down help

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

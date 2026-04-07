# ── Build stage ──────────────────────────────────────────────────────────────
FROM golang:1.25-alpine AS builder

ARG GIT_COMMIT=unknown
ARG BUILD_DATE=unknown

RUN apk --no-cache add ca-certificates tzdata git

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s \
      -X github.com/hhung06/digimap-backend/version.GitCommit=${GIT_COMMIT} \
      -X github.com/hhung06/digimap-backend/version.BuildDate=${BUILD_DATE}" \
    -o bin/digimap-backend .

# ── Runtime stage ─────────────────────────────────────────────────────────────
FROM alpine:3.21

RUN apk --no-cache add ca-certificates tzdata dumb-init

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /opt/app
COPY --from=builder /build/bin/digimap-backend .
COPY --from=builder /build/migrations ./migrations

USER appuser

EXPOSE 8080

ENTRYPOINT ["/usr/bin/dumb-init", "--"]
CMD ["/opt/app/digimap-backend", "serve"]

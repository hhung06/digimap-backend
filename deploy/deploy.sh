#!/usr/bin/env bash
# deploy.sh — Zero-downtime blue-green deployment for Docker Swarm + Traefik
#
# Usage:
#   ./deploy.sh <image-tag>           # e.g. ./deploy.sh hungnh14/digimap-backend:abc1234
#   ./deploy.sh --force-slot blue     # force deploy to a specific slot
#   ./deploy.sh --skip-health-check   # skip health verification (dangerous)
#
# Environment (must be set or sourced from .env.deploy):
#   API_DOMAIN, DB_NAME, DB_USER, DB_PASSWORD, JWT_SECRET,
#   S3_BUCKET, S3_REGION, AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY
#
# State file:
#   /var/lib/digimap/active-slot  — contains "blue" or "green"; persists across deploys

set -euo pipefail

# ── Constants ──────────────────────────────────────────────────────────────────
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
STATE_FILE="${STATE_FILE:-/var/lib/digimap/active-slot}"
HEALTH_URL_TEMPLATE="http://localhost:8080/health"
HEALTH_TIMEOUT="${HEALTH_TIMEOUT:-120}"   # seconds to wait for new slot to become healthy
HEALTH_INTERVAL=5                         # seconds between health check polls
DRAIN_DELAY="${DRAIN_DELAY:-30}"          # seconds to wait after cutover before draining old slot

# ── Colors for output ──────────────────────────────────────────────────────────
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

log()   { echo -e "${GREEN}[$(date '+%H:%M:%S')] $*${NC}"; }
warn()  { echo -e "${YELLOW}[$(date '+%H:%M:%S')] WARN: $*${NC}"; }
error() { echo -e "${RED}[$(date '+%H:%M:%S')] ERROR: $*${NC}" >&2; }
step()  { echo -e "${BLUE}[$(date '+%H:%M:%S')] ── $* ──${NC}"; }

# ── Argument parsing ───────────────────────────────────────────────────────────
APP_IMAGE=""
FORCE_SLOT=""
SKIP_HEALTH=false

while [[ $# -gt 0 ]]; do
  case $1 in
    --force-slot)   FORCE_SLOT="$2"; shift 2 ;;
    --skip-health-check) SKIP_HEALTH=true; shift ;;
    -*)             error "Unknown flag: $1"; exit 1 ;;
    *)              APP_IMAGE="$1"; shift ;;
  esac
done

if [[ -z "$APP_IMAGE" ]]; then
  error "Usage: $0 <image-tag> [--force-slot blue|green] [--skip-health-check]"
  exit 1
fi

# ── Load environment ───────────────────────────────────────────────────────────
if [[ -f "${SCRIPT_DIR}/.env.deploy" ]]; then
  set -a
  # shellcheck disable=SC1091
  source "${SCRIPT_DIR}/.env.deploy"
  set +a
fi

required_vars=(API_DOMAIN DB_NAME DB_USER DB_PASSWORD JWT_SECRET)
for var in "${required_vars[@]}"; do
  if [[ -z "${!var:-}" ]]; then
    error "Required environment variable '$var' is not set."
    exit 1
  fi
done

# ── Determine active/target slot ───────────────────────────────────────────────
determine_active_slot() {
  # 1. Check state file (fastest path)
  if [[ -f "$STATE_FILE" ]]; then
    cat "$STATE_FILE"
    return
  fi
  # 2. Query Swarm services (handles the case where state file was lost)
  if docker service inspect digimap-blue_app &>/dev/null; then
    active=$(docker service inspect digimap-blue_app \
      --format '{{index .Spec.Labels "com.digimap.active"}}' 2>/dev/null || echo "false")
    if [[ "$active" == "true" ]]; then
      echo "blue"
      return
    fi
  fi
  if docker service inspect digimap-green_app &>/dev/null; then
    active=$(docker service inspect digimap-green_app \
      --format '{{index .Spec.Labels "com.digimap.active"}}' 2>/dev/null || echo "false")
    if [[ "$active" == "true" ]]; then
      echo "green"
      return
    fi
  fi
  # 3. First deployment — no active slot yet
  echo "none"
}

ACTIVE_SLOT=$(determine_active_slot)

if [[ -n "$FORCE_SLOT" ]]; then
  TARGET_SLOT="$FORCE_SLOT"
else
  case "$ACTIVE_SLOT" in
    blue)  TARGET_SLOT="green" ;;
    green) TARGET_SLOT="blue" ;;
    none)  TARGET_SLOT="blue" ;;  # First deploy always goes to blue
    *)     error "Unknown active slot: $ACTIVE_SLOT"; exit 1 ;;
  esac
fi

log "Deploying $APP_IMAGE"
log "Active slot: ${ACTIVE_SLOT} → Target slot: ${TARGET_SLOT}"

# ── Helper: wait for all replicas in a service to be running ──────────────────
wait_for_service_running() {
  local service="$1"
  local timeout="${2:-$HEALTH_TIMEOUT}"
  local elapsed=0

  step "Waiting for ${service} replicas to be running"

  while true; do
    if [[ $elapsed -ge $timeout ]]; then
      error "Timed out waiting for ${service} to become healthy after ${timeout}s"
      return 1
    fi

    # Get replica counts: "running/desired"
    running=$(docker service ls --filter "name=${service}" \
      --format "{{.Replicas}}" 2>/dev/null | head -1 || echo "0/0")
    desired=$(echo "$running" | cut -d/ -f2)
    actual=$(echo "$running" | cut -d/ -f1)

    if [[ "$desired" -gt 0 && "$actual" -eq "$desired" ]]; then
      log "Service ${service} healthy: ${actual}/${desired} replicas running"
      return 0
    fi

    warn "Waiting... ${actual}/${desired} replicas running (${elapsed}s elapsed)"
    sleep $HEALTH_INTERVAL
    elapsed=$((elapsed + HEALTH_INTERVAL))
  done
}

# ── Helper: verify app-level health endpoint ──────────────────────────────────
verify_app_health() {
  local service="$1"
  local timeout="${2:-60}"
  local elapsed=0

  step "Verifying /health endpoint on ${service}"

  # Get one task ID from the service and exec into it
  while true; do
    if [[ $elapsed -ge $timeout ]]; then
      error "App health check timed out after ${timeout}s"
      return 1
    fi

    task_id=$(docker service ps "$service" \
      --filter "desired-state=running" \
      --format "{{.ID}}" 2>/dev/null | head -1)

    if [[ -n "$task_id" ]]; then
      # Find the container ID for this task
      container_id=$(docker inspect --format "{{.Status.ContainerStatus.ContainerID}}" \
        "$task_id" 2>/dev/null | head -c 12)

      if [[ -n "$container_id" ]]; then
        if docker exec "$container_id" wget -qO- http://localhost:8080/health &>/dev/null; then
          log "App health check passed on ${service}"
          return 0
        fi
      fi
    fi

    warn "Health check pending... (${elapsed}s elapsed)"
    sleep $HEALTH_INTERVAL
    elapsed=$((elapsed + HEALTH_INTERVAL))
  done
}

# ── Helper: activate a slot via Traefik label swap ────────────────────────────
activate_slot() {
  local slot="$1"       # blue or green
  local deactivate="$2" # green or blue (the one being turned off)

  step "Cutting over traffic: ${deactivate} → ${slot}"

  # Step 1: Enable target slot at HIGH priority (Traefik immediately begins routing)
  # Both slots are briefly active — Traefik will prefer the higher priority one
  docker service update \
    --label-add "traefik.enable=true" \
    --label-add "traefik.http.routers.api-${slot}.priority=10" \
    --label-add "com.digimap.active=true" \
    "digimap-${slot}_app"

  # Brief wait for Traefik to pick up the new service (Swarm poll interval ~1s)
  sleep 3

  # Step 2: Disable old slot (Traefik stops routing to it immediately)
  docker service update \
    --label-add "traefik.enable=false" \
    --label-add "traefik.http.routers.api-${deactivate}.priority=1" \
    --label-add "com.digimap.active=false" \
    "digimap-${deactivate}_app"

  log "Traffic cutover complete. ${slot} is now ACTIVE."
}

# ── Step 1: Run database migrations BEFORE deploying new app ─────────────────
step "Running database migrations"
docker run --rm \
  --network digimap-infra_internal \
  -e DB_HOST=postgres \
  -e DB_PORT=5432 \
  -e DB_NAME="${DB_NAME}" \
  -e DB_USER="${DB_USER}" \
  -e DB_PASSWORD="${DB_PASSWORD}" \
  -e DB_SSLMODE="${DB_SSLMODE:-disable}" \
  "${APP_IMAGE}" migrate up
log "Migrations complete"

# ── Step 2: Deploy new image to target slot ────────────────────────────────────
step "Deploying ${APP_IMAGE} to slot: ${TARGET_SLOT}"

# Build env for docker stack deploy substitution
export APP_IMAGE
export BLUE_ACTIVE="false"
export GREEN_ACTIVE="false"
export BLUE_PRIORITY="1"
export GREEN_PRIORITY="1"

# If this is not the first deploy, keep the current active slot active during deployment
if [[ "$ACTIVE_SLOT" == "blue" ]]; then
  export BLUE_ACTIVE="true"
  export BLUE_PRIORITY="10"
elif [[ "$ACTIVE_SLOT" == "green" ]]; then
  export GREEN_ACTIVE="true"
  export GREEN_PRIORITY="10"
fi

docker stack deploy \
  --with-registry-auth \
  --compose-file "${SCRIPT_DIR}/stack-${TARGET_SLOT}.yml" \
  "digimap-${TARGET_SLOT}"

log "Stack digimap-${TARGET_SLOT} deployed"

# ── Step 3: Wait for target slot to become healthy ────────────────────────────
if [[ "$SKIP_HEALTH" == "true" ]]; then
  warn "--skip-health-check set: skipping health verification (dangerous!)"
else
  wait_for_service_running "digimap-${TARGET_SLOT}_app" "$HEALTH_TIMEOUT"
  verify_app_health "digimap-${TARGET_SLOT}_app" 60
fi

# ── Step 4: Atomic traffic cutover ────────────────────────────────────────────
if [[ "$ACTIVE_SLOT" != "none" ]]; then
  activate_slot "$TARGET_SLOT" "$ACTIVE_SLOT"
else
  # First deployment — just enable the target slot
  step "First deployment: enabling ${TARGET_SLOT} slot"
  docker service update \
    --label-add "traefik.enable=true" \
    --label-add "traefik.http.routers.api-${TARGET_SLOT}.priority=10" \
    --label-add "com.digimap.active=true" \
    "digimap-${TARGET_SLOT}_app"
  log "${TARGET_SLOT} slot is now ACTIVE (first deployment)"
fi

# ── Step 5: Drain the old slot (keep containers for rollback) ────────────────
if [[ "$ACTIVE_SLOT" != "none" ]]; then
  step "Draining old slot: ${ACTIVE_SLOT} (waiting ${DRAIN_DELAY}s for in-flight requests)"
  sleep "$DRAIN_DELAY"
  docker service scale "digimap-${ACTIVE_SLOT}_app=0"
  log "Slot ${ACTIVE_SLOT} scaled to 0 (kept for rollback; use rollback.sh to reactivate)"
fi

# ── Step 6: Persist active slot state ─────────────────────────────────────────
mkdir -p "$(dirname "$STATE_FILE")"
echo "$TARGET_SLOT" > "$STATE_FILE"
log "State saved: active slot = ${TARGET_SLOT}"

# ── Done ──────────────────────────────────────────────────────────────────────
echo ""
log "✓ Deployment complete"
log "  Image:       ${APP_IMAGE}"
log "  Active slot: ${TARGET_SLOT}"
log "  Rollback:    ./deploy/rollback.sh"

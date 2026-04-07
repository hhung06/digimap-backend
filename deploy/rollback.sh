#!/usr/bin/env bash
# rollback.sh — Instant rollback to the previous deployment slot
#
# Rollback is instant because the old slot's containers are still running
# (just scaled to 0 / traffic disabled). We simply reverse the cutover.
#
# Usage:
#   ./deploy/rollback.sh              # rollback to previous slot
#   ./deploy/rollback.sh --dry-run    # show what would happen without doing it

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
STATE_FILE="${STATE_FILE:-/var/lib/digimap/active-slot}"
APP_REPLICAS="${APP_REPLICAS:-2}"
DRAIN_DELAY="${DRAIN_DELAY:-15}"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

log()   { echo -e "${GREEN}[$(date '+%H:%M:%S')] $*${NC}"; }
warn()  { echo -e "${YELLOW}[$(date '+%H:%M:%S')] WARN: $*${NC}"; }
error() { echo -e "${RED}[$(date '+%H:%M:%S')] ERROR: $*${NC}" >&2; }
step()  { echo -e "${BLUE}[$(date '+%H:%M:%S')] ── $* ──${NC}"; }

DRY_RUN=false
[[ "${1:-}" == "--dry-run" ]] && DRY_RUN=true

# ── Determine current and previous slots ──────────────────────────────────────
if [[ ! -f "$STATE_FILE" ]]; then
  error "State file not found: $STATE_FILE"
  error "Cannot determine active slot. Check Swarm service labels manually."
  exit 1
fi

ACTIVE_SLOT=$(cat "$STATE_FILE")
case "$ACTIVE_SLOT" in
  blue)  ROLLBACK_SLOT="green" ;;
  green) ROLLBACK_SLOT="blue" ;;
  *)     error "Invalid state in $STATE_FILE: $ACTIVE_SLOT"; exit 1 ;;
esac

log "Current active slot:   ${ACTIVE_SLOT}"
log "Rollback to slot:      ${ROLLBACK_SLOT}"

# Verify rollback slot has containers (were not fully removed)
rollback_service="digimap-${ROLLBACK_SLOT}_app"
if ! docker service inspect "$rollback_service" &>/dev/null; then
  error "Rollback service '${rollback_service}' does not exist in Swarm."
  error "The previous slot may have been removed. Check 'docker service ls'."
  exit 1
fi

if [[ "$DRY_RUN" == "true" ]]; then
  warn "DRY RUN — no changes will be made"
  log "Would scale up:  ${rollback_service} to ${APP_REPLICAS} replicas"
  log "Would cut over:  traffic → ${ROLLBACK_SLOT}"
  log "Would scale down: digimap-${ACTIVE_SLOT}_app to 0"
  exit 0
fi

# ── Step 1: Scale up rollback slot ────────────────────────────────────────────
step "Scaling up rollback slot: ${ROLLBACK_SLOT}"
docker service scale "${rollback_service}=${APP_REPLICAS}"

# Wait for rollback slot to be ready
elapsed=0
while true; do
  running=$(docker service ls --filter "name=${rollback_service}" \
    --format "{{.Replicas}}" | head -1)
  desired=$(echo "$running" | cut -d/ -f2)
  actual=$(echo "$running" | cut -d/ -f1)
  if [[ "$desired" -gt 0 && "$actual" -eq "$desired" ]]; then
    log "Rollback slot ready: ${actual}/${desired} replicas"
    break
  fi
  if [[ $elapsed -ge 120 ]]; then
    error "Timed out waiting for rollback slot to become ready"
    exit 1
  fi
  warn "Waiting for replicas... ${actual}/${desired} (${elapsed}s)"
  sleep 5
  elapsed=$((elapsed + 5))
done

# ── Step 2: Atomic traffic cutover back to rollback slot ─────────────────────
step "Cutting over traffic: ${ACTIVE_SLOT} → ${ROLLBACK_SLOT}"

docker service update \
  --label-add "traefik.enable=true" \
  --label-add "traefik.http.routers.api-${ROLLBACK_SLOT}.priority=10" \
  --label-add "com.digimap.active=true" \
  "${rollback_service}"

sleep 3

docker service update \
  --label-add "traefik.enable=false" \
  --label-add "traefik.http.routers.api-${ACTIVE_SLOT}.priority=1" \
  --label-add "com.digimap.active=false" \
  "digimap-${ACTIVE_SLOT}_app"

log "Traffic cutover complete — ${ROLLBACK_SLOT} is now ACTIVE"

# ── Step 3: Drain failed deployment slot ─────────────────────────────────────
step "Draining failed slot: ${ACTIVE_SLOT} (waiting ${DRAIN_DELAY}s for in-flight requests)"
sleep "$DRAIN_DELAY"
docker service scale "digimap-${ACTIVE_SLOT}_app=0"
log "Slot ${ACTIVE_SLOT} drained"

# ── Step 4: Update state ──────────────────────────────────────────────────────
echo "$ROLLBACK_SLOT" > "$STATE_FILE"

echo ""
log "✓ Rollback complete"
log "  Active slot: ${ROLLBACK_SLOT}"
log "  Failed slot: ${ACTIVE_SLOT} (containers idle, image preserved)"

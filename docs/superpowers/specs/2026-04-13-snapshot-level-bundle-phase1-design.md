# Snapshot & LevelBundle — Phase 1 Design

**Date:** 2026-04-13
**Scope:** CRUD layer only — no publish pipeline (Phase 2)
**Status:** Approved

---

## Background

The original `indoormap-backend` (Django/Python) has a snapshot and sync-data system for managing versioned venue map data. `digimap-backend` (Go) is missing this entirely. This document covers Phase 1: the domain, storage, repository, service, and HTTP layers needed to create and manage snapshots and level bundles.

The publish pipeline (assembling per-language bundles, encrypting, CDN invalidation) is Phase 2.

---

## Decisions

| Decision | Choice | Rationale |
|---|---|---|
| Upload flow | Pass-through (backend receives bundle, uploads to S3) | Low-frequency admin operation; bundles are a few MB; single API call is simpler for the client |
| State representation | Integer constants | Consistent with all other domain objects in this codebase |
| SyncData rename | `LevelBundle` | Describes what it is: the bundled data artifact for a specific map level |
| LevelBundle versioning | Cascade from Snapshot pruning | Adding `SnapshotID` FK + `ON DELETE CASCADE` makes bundle cleanup automatic — no separate version limit needed |
| S3 ordering | S3-first, DB-second | Orphaned S3 objects are harmless (deterministic keys, overwritten on retry); avoids DB records with no S3 backing |
| `pending` state (0) | Reserved, never written in Phase 1 | Pass-through flow goes straight to `draft`; reserved for future async upload path and cleanup job detection |
| RBAC | SystemAdmin only | All snapshot/bundle routes restricted to system admins; goes under existing `adminOnly` group in `router.go` |
| FK design | DB-level FK constraints | `snapshots.venue_id → venues(id)`, `level_bundles.snapshot_id → snapshots(id) ON DELETE CASCADE`, `level_bundles.level_id → levels(id)` |

---

## Domain

### `internal/domain/snapshot.go`

```go
const (
    SnapshotStatePending = 0  // reserved — never written in Phase 1
    SnapshotStateDraft   = 1
    SnapshotStatePublic  = 2

    SnapshotMethodManual = 1
    SnapshotMethodAuto   = 2

    MaxSnapshotVersions = 5
)

type Snapshot struct {
    ID        uuid.UUID
    VenueID   uuid.UUID   // FK → venues(id)
    State     int         // 0=pending 1=draft 2=public
    Method    int         // 1=manual 2=auto
    CreatedBy *uuid.UUID  // FK → users(id)
    PublishAt *time.Time
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt *time.Time
}
```

### `internal/domain/level_bundle.go`

```go
const (
    LevelBundleStatePending = 0  // reserved
    LevelBundleStateDraft   = 1
    LevelBundleStatePublic  = 2
)

type LevelBundle struct {
    ID         uuid.UUID
    SnapshotID uuid.UUID  // FK → snapshots(id) ON DELETE CASCADE
    VenueID    uuid.UUID  // FK → venues(id)
    LevelID    uuid.UUID  // FK → levels(id)
    State      int        // 0=pending 1=draft 2=public
    CreatedAt  time.Time
    UpdatedAt  time.Time
    DeletedAt  *time.Time
}
```

No `Version` field on `LevelBundle` — version control is handled by pruning the parent `Snapshot`.

---

## Migration

**File:** `migrations/000024_create_snapshots.up.sql`

```sql
CREATE TABLE snapshots (
    id          UUID        PRIMARY KEY DEFAULT uuidv7(),
    venue_id    UUID        NOT NULL REFERENCES venues(id),
    state       SMALLINT    NOT NULL DEFAULT 0,
    method      SMALLINT    NOT NULL DEFAULT 1,
    created_by  UUID        REFERENCES users(id),
    publish_at  TIMESTAMPTZ,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);
SELECT create_updated_at_trigger('snapshots');
CREATE INDEX idx_snapshots_venue_state ON snapshots (venue_id, state) WHERE deleted_at IS NULL;

CREATE TABLE level_bundles (
    id          UUID        PRIMARY KEY DEFAULT uuidv7(),
    snapshot_id UUID        NOT NULL REFERENCES snapshots(id) ON DELETE CASCADE,
    venue_id    UUID        NOT NULL REFERENCES venues(id),
    level_id    UUID        NOT NULL REFERENCES levels(id),
    state       SMALLINT    NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);
SELECT create_updated_at_trigger('level_bundles');
CREATE INDEX idx_level_bundles_snapshot ON level_bundles (snapshot_id) WHERE deleted_at IS NULL;
```

**Down migration** drops both tables in reverse order (level_bundles first due to FK).

---

## S3 Storage

### Interface extension (`internal/platform/storage/s3.go`)

```go
type Storer interface {
    PresignUpload(ctx, key, contentType string, ttl time.Duration) (string, error)
    PresignDownload(ctx, key string, ttl time.Duration) (string, error)
    PutObject(ctx context.Context, key string, data []byte) error   // new
    GetObject(ctx context.Context, key string) ([]byte, error)       // new
}
```

`LogStorer` gets no-op implementations that log the key and return nil/empty.

### S3 Key patterns

| Object | Key |
|---|---|
| Snapshot bundle | `{env}/{venueID}/snapshots/draft/{snapshotID}.json` |
| LevelBundle | `{env}/{venueID}/level-bundles/{snapshotID}/{levelID}.json` |

Keys are deterministic — derived from IDs, never stored in DB. On retry, the upload overwrites the previous attempt.

---

## Repository

### `internal/repository/interfaces.go` additions

```go
type SnapshotRepository interface {
    FindByID(ctx context.Context, id uuid.UUID) (*domain.Snapshot, error)
    List(ctx context.Context, venueID uuid.UUID, p domain.Pagination) ([]*domain.Snapshot, int64, error)
    LatestPublished(ctx context.Context, venueID uuid.UUID) (*domain.Snapshot, error)
    Create(ctx context.Context, s *domain.Snapshot) error
    UpdateState(ctx context.Context, id uuid.UUID, state int, publishAt *time.Time) error
    Delete(ctx context.Context, id uuid.UUID) error
    CountDraftsByVenue(ctx context.Context, venueID uuid.UUID) (int64, error)
    DeleteOldestDraft(ctx context.Context, venueID uuid.UUID) error
}

type LevelBundleRepository interface {
    FindByID(ctx context.Context, id uuid.UUID) (*domain.LevelBundle, error)
    ListBySnapshot(ctx context.Context, snapshotID uuid.UUID) ([]*domain.LevelBundle, error)
    Create(ctx context.Context, b *domain.LevelBundle) error
    Delete(ctx context.Context, id uuid.UUID) error
}
```

Implementations in:
- `internal/repository/postgres/snapshot_repo.go`
- `internal/repository/postgres/level_bundle_repo.go`

`DeleteOldestDraft` hard-deletes (not soft-delete) the draft snapshot with the oldest `created_at` for the given venue, so the `ON DELETE CASCADE` fires and removes its associated `LevelBundle` rows.

---

## Service

### `internal/service/snapshot_service.go`

```go
type SnapshotService interface {
    List(ctx context.Context, venueID uuid.UUID, p domain.Pagination) ([]*domain.Snapshot, int64, error)
    Get(ctx context.Context, id uuid.UUID) (*domain.Snapshot, error)
    LatestPublished(ctx context.Context, venueID uuid.UUID) (*domain.Snapshot, error)
    CreateDraft(ctx context.Context, venueID, createdBy uuid.UUID, bundle []byte) (*domain.Snapshot, error)
    Delete(ctx context.Context, id uuid.UUID) error
}
```

**`CreateDraft` execution order:**
1. Generate `snapshotID = newID()`
2. Upload `bundle` → S3 key `{env}/{venueID}/snapshots/draft/{snapshotID}.json`
   - Failure → return error; DB untouched
3. `INSERT` snapshot with `state = SnapshotStateDraft`
   - Failure → S3 object orphaned but harmless (overwritten on retry with same key)
4. `CountDraftsByVenue` → if `count >= MaxSnapshotVersions (5)`: `DeleteOldestDraft` (hard-delete, cascades to LevelBundles)
5. Return snapshot

### `internal/service/level_bundle_service.go`

```go
type LevelBundleService interface {
    ListBySnapshot(ctx context.Context, snapshotID uuid.UUID) ([]*domain.LevelBundle, error)
    Create(ctx context.Context, snapshotID, venueID, levelID uuid.UUID, bundle []byte) (*domain.LevelBundle, error)
    Delete(ctx context.Context, id uuid.UUID) error
}
```

**`Create` execution order:** same S3-first pattern as `CreateDraft`. The new `LevelBundle` is written with `State` matching its parent `Snapshot`'s state at creation time (always `draft` in Phase 1).

---

## HTTP Layer

### DTOs (`internal/dto/`)

- `snapshot_dto.go` — `SnapshotResponse`, `CreateSnapshotRequest`, `SnapshotToResponse()`
- `level_bundle_dto.go` — `LevelBundleResponse`, `CreateLevelBundleRequest`, `LevelBundleToResponse()`

`CreateSnapshotRequest` carries `bundle json.RawMessage` (the full map bundle payload).
`CreateLevelBundleRequest` carries `level_id uuid.UUID` and `bundle json.RawMessage`.

### Handlers

- `internal/handler/snapshot_handler.go`
- `internal/handler/level_bundle_handler.go`

### Routes — all under `adminOnly` group in `router.go`

```
GET    /venues/{venueID}/snapshots                                list
POST   /venues/{venueID}/snapshots                                createDraft
GET    /venues/{venueID}/snapshots/recent                         latestPublished
GET    /venues/{venueID}/snapshots/{snapshotID}                   get
DELETE /venues/{venueID}/snapshots/{snapshotID}                   delete

POST   /venues/{venueID}/snapshots/{snapshotID}/bundles           create
GET    /venues/{venueID}/snapshots/{snapshotID}/bundles           list
DELETE /venues/{venueID}/snapshots/{snapshotID}/bundles/{id}      delete
```

All routes use the existing `adminOnly` group — no new middleware needed.

---

## Files to Create / Modify

| Action | File |
|---|---|
| Create | `migrations/000024_create_snapshots.up.sql` |
| Create | `migrations/000024_create_snapshots.down.sql` |
| Create | `internal/domain/snapshot.go` |
| Create | `internal/domain/level_bundle.go` |
| Modify | `internal/platform/storage/s3.go` — add `PutObject`, `GetObject` |
| Modify | `internal/repository/interfaces.go` — add two interfaces |
| Create | `internal/repository/postgres/snapshot_repo.go` |
| Create | `internal/repository/postgres/level_bundle_repo.go` |
| Create | `internal/service/snapshot_service.go` |
| Create | `internal/service/level_bundle_service.go` |
| Create | `internal/dto/snapshot_dto.go` |
| Create | `internal/dto/level_bundle_dto.go` |
| Create | `internal/handler/snapshot_handler.go` |
| Create | `internal/handler/level_bundle_handler.go` |
| Modify | `internal/handler/router.go` — register routes + DI |
| Modify | `cmd/serve.go` — init repos, services, wire storer |

---

## Out of Scope (Phase 2)

- Publish (`draft → public`) and revert logic
- Async publish pipeline (per-language bundle assembly, encryption, CDN invalidation)
- `SnapshotStatePublic` and `UpdateState` are defined but not used until Phase 2
- `SnapshotMethodAuto` and `VenueSyncFlag` (webhook dirty flag)

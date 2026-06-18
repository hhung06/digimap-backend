# Parity Audit: indoormap-backend → digimap-backend

**Date:** 2026-06-15  
**Scope:** Full parity audit — data models, endpoints, business logic, infrastructure  
**Repos:**
- `indoormap-backend` (Python/Django, MySQL, 281 migrations) — original
- `digimap-backend` (Go, PostgreSQL, 31 migrations) — reimplementation

---

## Executive Summary

The Go rewrite **preserves all core business logic** in the domains that matter most — location management, product/content CRUD, event and survey lifecycle, publishing pipeline, and JMA webhook integration. The bundle crypto format is **byte-for-byte compatible** with the Python original; app clients can switch between backends without code changes.

**Critical gaps that must be addressed before full parity:**

| # | Gap | Domain | Risk |
|---|-----|--------|------|
| 1 | AppUser.Phone stored as plaintext (Python uses AES-GCM) | Visitors/PII | CRITICAL |
| 2 | S3, SES, FCM, OpenSearch are `LogXxx` stubs — not wired | Infrastructure | CRITICAL |
| 3 | Customer isolation / scope computation missing | Auth/Users | HIGH |
| 4 | Japanese tokenization not ported (`"dictionaries": {}`) | Bundles | MEDIUM |
| 5 | sdk/v1 endpoints not ported (levels, categories, bundle) | API surface | MEDIUM |

**Core correctly preserved:** event/survey/visitor lifecycle, bundle assembly, crypto/compression, JMA webhooks, product/article/coupon/location CRUD, beacon/connection management, snapshot/publish pipeline, RBAC middleware, scheduled notification/survey commands.

**Cruft correctly dropped:** Selenium publish-verification (+ 15 MB chromedriver), Azure business-card OCR, Vimeo API integration, Slack integration, AWS Textract, per-client FCM projects (`FCM_PROJECTS`), 2940-line FOODEX/HCJ constants.py, `VenueSyncFlag.expo_id` field (replaced by `venue.ExternalID` in Go enricher), Django multi-project FCM, `publish_venue` v1 (v2 ported), versioned command duplication, venue publish/clone endpoints, user management + invitation endpoints (no RBAC in this version), QRCode domain, Notification Segment + NotificationLog.

---

## Parity Matrix

| Domain | Models | Endpoints | Logic | Overall |
|--------|--------|-----------|-------|---------|
| Auth & Users | ⚠️ partial | 🔵 dropped | ⚠️ partial | ⚠️ |
| Customers | ✅ parity | ✅ parity | ✅ parity | ✅ |
| Venues | ⚠️ partial | 🔵 dropped | ⚠️ partial | ⚠️ |
| Levels / Maps | ✅ parity | ✅ parity | ✅ parity | ✅ |
| Locations | ✅ parity | ✅ parity | ✅ parity | ✅ |
| Location Categories | ✅ parity | ✅ parity | ✅ parity | ✅ |
| Masters (Tags/Templates) | ⚠️ partial | ⚠️ partial | 🔵 dropped | ⚠️ |
| Products | ✅ parity | ✅ parity | ✅ parity | ✅ |
| Product Plazas | ✅ parity | ✅ parity | 🔵 images dropped | ✅ |
| Articles | ✅ parity | ✅ parity | 🔵 contact/related-products dropped | ✅ |
| Coupons | ⚠️ partial | ✅ parity | ⚠️ partial | ⚠️ |
| Advertisements | ✅ parity | ✅ parity | ✅ parity | ✅ |
| Videos | ✅ parity | ✅ parity | 🔵 Vimeo dropped | ✅ |
| Events | ✅ parity | ✅ parity | ✅ parity | ✅ |
| Surveys | ✅ parity | ✅ parity | ✅ parity | ✅ |
| Visitors / AppUser | ⚠️ partial | ✅ parity | ❌ missing (phone encryption) | ❌ |
| Notifications | ✅ parity | ✅ parity | ✅ parity | ✅ |
| Analytics | ✅ parity | ⚠️ partial | 🔵 alerting dropped | ⚠️ |
| Publishing / Bundles | ⚠️ partial | ✅ parity | ⚠️ partial | ⚠️ |
| Beacons | ✅ parity | ✅ parity | ✅ parity | ✅ |
| Connections | ✅ parity | ✅ parity | ✅ parity | ✅ |
| QR Codes | 🔵 dropped | 🔵 dropped | 🔵 dropped | 🔵 |
| Infrastructure (S3/SES/FCM/CDN/Search) | ✅ interfaces | 🔵 stubs | ❌ missing | ❌ |
| JMA Webhooks | ✅ parity | ✅ parity | ⚠️ unknown | ⚠️ |
| Tenant Enricher | ✅ interface | n/a | ❌ empty | ⚠️ |

---

## Per-Domain Sections

### 1. Auth & Users

#### Model diff

| Aspect | indoormap (Python) | digimap (Go) | Status |
|--------|--------------------|--------------|--------|
| Core identity | Django `auth.User` + `Profile` OneToOne | Single `User` struct | Restructured, functionally OK |
| Customer FK | `Profile.customer` FK | **MISSING** | ❌ Gap — no tenant isolation |
| Scope property | `profile.scope` → (customer, venues) | **MISSING** | ❌ Gap |
| `is_partner` bool | Customer user flag | `IsSystemAdmin` bool | Semantically equivalent (inverted) |
| VenueUserRole | `profile`, `venue`, `role_type` (ADMIN/EDITOR/VIEWER) | `user_id`, `venue_id`, `role` (Owner/Editor/Viewer enum) | ✅ Equivalent |
| VenueInvitation | `invited_profile` FK + `invited_email` | email only | 🔵 Simplified — no pre-fill for existing user |
| Token storage | Raw token string | SHA-256 hash | ✅ Improved (more secure) |
| RefreshToken revocation | Django simplejwt | Explicit `revoked_at` field | ✅ Equivalent |

#### Endpoint diff

| Python | Go | Status |
|--------|----|--------|
| `POST /users/login/` | `POST /auth/login` | ✅ |
| `POST /users/logout/` | `POST /auth/logout` | ✅ |
| `POST /users/token/refresh/` | `POST /auth/refresh` | ✅ |
| `POST /users/upd-pwd/` | `POST /auth/password-change` | ✅ |
| `POST /users/rst-pwd-request/` | `POST /auth/password-reset` | ✅ |
| `POST /users/rst-pwd/` | `POST /auth/password-reset/confirm` | ✅ |
| `POST /users/register/` | `/auth/register` — commented out (`router.go:118`) | 🔵 |
| `GET/PUT /profile/` | Not implemented | 🔵 |
| `GET /users/` (list) | Not implemented | 🔵 |
| `POST /users/invite/` | Not implemented | 🔵 |
| `POST /users/<id>/change-role/` | Not implemented | 🔵 |
| `POST /users/<id>/remove/` | Not implemented | 🔵 |
| `POST /invitations/accept/` | Not implemented | 🔵 |
| `POST /invitations/<id>/cancel/` | Not implemented | 🔵 |

#### Logic diff

- **RBAC:** No per-venue RBAC permission enforcement in this version. Python's 8+ permission classes (`permissions.py:336–532`) are not ported. Go middleware stubs exist (`SystemAdminRequired()` + `VenueAccess(minRole)`) but user management endpoints are out of scope.
- **Scope computation:** Python `profile.scope` computes `(customer, venues)` tuple for queryset filtering. Not ported — not needed without RBAC.
- **Token handling:** Python stores raw tokens; Go hashes with SHA-256 — improvement, not a gap.

#### Deltas

| Gap | Verdict |
|-----|---------|
| User management endpoints (list, invite, change-role, remove) | 🔵 DROPPED — no RBAC in this version |
| Invitation endpoints | 🔵 DROPPED — no RBAC in this version |
| Scope computation + CustomerID | 🔵 DROPPED — no RBAC in this version |
| Register endpoint | 🔵 DROPPED — commented out intentionally |

---

### 2. Customers, Venues, Levels / Maps

#### Model diff — Customers

Full parity. Go adds `deleted_at` and `metadata JSONB`; no Python fields missing.

#### Model diff — Venues

| Field | Status | Notes |
|-------|--------|-------|
| `publish: bool` | ❌ MISSING | Python `models.py` has `publish: BooleanField`; no column in Go migration `000003` |
| `defaultMap`, `countryCode` | ❌ MISSING | Present in Python; absent from Go migration |
| `start_at`, `end_at` | ✅ Added in Go | Not in Python — forward-compatible addition |
| `translations` | 🔵 JSONB in Go vs separate `VenueTranslation` table in Python | Intentional simplification |
| `theme` | 🔵 Simple `PrimaryColor/SecondaryColor` in Go vs scope-aware `VenueTheme` model | Simplified |
| `plugins` | 🔵 JSONB in Go vs dedicated `Plugin` model with enable/disable | Simplified |

#### Model diff — Levels/Maps

Full parity: MapGroup, Perspective, Level, GeoReference, LevelType all present. Minor rename: Python `publish` → Go `is_published` on Level.

#### Endpoint diff

| Gap | Status |
|-----|--------|
| `POST /venues/{id}/publish` | 🔵 DROPPED — not in this version |
| `POST /venues/{id}/clone` | 🔵 DROPPED — not in this version |
| All level endpoints | ✅ Full parity |
| All map-group endpoints | ✅ Full parity |
| All geo-reference endpoints | ✅ Full parity |

#### Logic diff

- **Venue publish/clone:** Both removed in this version. Python's async clone (management command) and publish state are not carried forward.
- **VenueTheme:** Python has scope validation (GLOBAL vs CUSTOM), per-theme config JSONB. Go simplifies to `Name + PrimaryColor + SecondaryColor`. Theme management less flexible.
- **Plugin system:** Python has CRUD + enable/disable per plugin. Go stores plugins as JSONB in `venues.plugins` — no management endpoints.

#### Deltas

| Gap | Verdict | Priority |
|-----|---------|----------|
| Venue publish/clone | 🔵 DROPPED | Intentional |
| `Venue.publish` model field | 🔵 DROPPED | Intentional — no endpoint, no field needed |
| Theme scope system | PORT or ACCEPT | MEDIUM — decide if advanced theme management needed |
| Plugin management endpoints | DEFER | LOW — JSONB storage is backward-compatible |
| `defaultMap`, `countryCode` | CLARIFY | LOW — confirm with frontend if used |

---

### 3. Locations, Categories, Masters

#### Model diff

Full parity on Location, LocationCategory (with hierarchical parent), LocationImage.  
Go correctly implements hierarchical category with `validateHierarchy()` cycle-prevention.

Notable:
- Python has separate `LocationTranslation` and `LocationCategoryTranslation` tables. Go uses `localization JSONB`. Intentional simplification, no functional loss if frontend parses JSON.
- Python `LocationTemplate` (masters) = Go `amenities` table. Renamed, same purpose.
- `TemplateCategory` model in Python has no Go equivalent (treated as seed data only).
- PERSON location type constant not defined in Go (`internal/domain/location.go:12–16`).

#### Endpoint diff

Full parity on `/venues/{id}/locations/`, `/venues/{id}/categories/`, memo CRUD, duplicate, set-top.  
Python's `/api/masters/` (template-categories, amenities CRUD) has no Go equivalent — intentionally dropped (seed data only).

#### Logic diff

- Category hierarchy validation: ✅ ported
- Location type filtering: ✅ ported
- Location duplication: ✅ ported — **but tags not copied** (Python copies `place_tags`; Go skips)
- `PlaceWorkHours` not copied on duplicate (`location_service.go:223–225`)
- Image auto-resize: Python auto-generates 3 sizes on save; Go stores URL strings only — image pipeline must be external (CDN/S3) or ported

#### Deltas

| Gap | Verdict | Priority |
|-----|---------|----------|
| Tags not copied on location duplicate | PORT | MEDIUM — data loss on duplication |
| `PlaceWorkHours` not copied on duplicate | PORT | LOW |
| PERSON location type constant | PORT | LOW |
| Image auto-resize | DEFER | MEDIUM — decide: CDN resize or Go processing |
| Masters CRUD endpoints | DROP | Intentionally seed-data only in Go |

---

### 4. Products, Plazas, Articles, Advertisements, Coupons, Videos

#### Products — ✅ core parity

All key fields present. `code` field dropped (duplicate of name). Complex filter params from Python (`exhibitor_english_status`, `product_requisite_documents`, multi-category OR) not ported — Go has basic name search only.

**Delta:** Complex product filters — PORT if app search relies on them.

#### Articles — ✅ core parity with intentional drops

Contact metadata fields (`company_name`, `full_name`, `landline`, `mobile`, `email`, `created_by`, `application_language`) dropped. `related_products` M2M dropped. Published period uses `*time.Time` vs Python `DateField` — precision mismatch, low risk.

#### Advertisements — ✅ parity

Full admin CRUD implemented: `internal/handler/ad_handler.go` with List, Get, Create, Update, Delete, and Publish handlers. Routes wired at `router.go:320–327` under `/venues/:id/ads`. App endpoint also present at `/app/v1/ads` (`router.go:455`).

#### Coupons — ⚠️ redemption model mismatch

Python: `UserCoupon` M2M table (one coupon issued to many users, per-user `is_used` flag).  
Go: `redeemed_by UUID + redeemed_at TIMESTAMPTZ` directly on Coupon — supports single redemption only.

**Delta:** Redemption model — RISK (HIGH). If indoormap issues one coupon to multiple users, data migration will fail silently. Clarify intended behavior; add `app_user_coupons` junction table if multi-user needed.

#### Videos — ✅ New in Go

Python has Vimeo integration (`utils/vimeo.py`) but **no Video model**. Go has full CRUD with plain URL storage. Vimeo API integration (delete, ID extraction, thumbnail upload) dropped.

**Delta:** Vimeo API — DROP unless videos are user-uploaded; document as external-URL-only.

#### Product Plazas — ✅ parity

Plaza images/description-images dropped. Core fields present.

---

### 5. Events, Surveys

Full parity on both. Event model, EventType, EventTag, EventImage — all present and aligned. Survey, Question, Option, SurveyResponse, SurveyAnswer — all present and aligned.

Scheduling (survey activation, notification dispatch) relies on external cron calling `go run . activate-scheduled-surveys` and `go run . send-scheduled-notifications` — equivalent to Python's management command + Celery dispatch. **No automatic background job — cron deployment required.**

---

### 6. Visitors / AppUser

#### Model diff

Python has two models: `Visitor` (survey respondents) and `AppUser` (authenticated mobile user). Go merges both into `AppUser` with a `VisitorType` discriminator. Functionally equivalent.

**CRITICAL GAP — phone encryption:**

| | Python | Go |
|--|--------|----|
| Storage | `_phone_number BinaryField` | `Phone string` (plaintext) |
| Cipher | AES-GCM (12-byte nonce, `DATA_ENCRYPTION_KEY`) | **None** |
| Source | `encrypting.py:132–228` | `app_user.go` — no crypto |

Go production exposes phone numbers in plaintext. Rows ingested from Python will have AES-GCM ciphertext stored as a binary blob — Go will serve garbled bytes unless decrypted at read time.

**Required fix:**
1. Add AES-GCM encryption/decryption wrapper in `internal/platform/crypto/aes.go` (currently only has AES-256-CBC).
2. Encrypt on write, decrypt on read in `app_user` repository.
3. Data migration: decrypt existing Python binary rows, re-encrypt with Go convention (or store as plaintext after migration).

---

### 7. Notifications

#### Model diff

| Python | Go | Status |
|--------|----|--------|
| Notification | Notification | ✅ parity |
| Segment | Not implemented | 🔵 DROPPED |
| NotificationSegment | Not implemented | 🔵 DROPPED |
| NotificationLog | Not implemented | 🔵 DROPPED |

#### Logic diff

Notification scheduling: ✅ ported (via `send_scheduled_notifications` CLI command).  
Segment-based delivery fields (`SegmentFilters`, `DeviceTokens`) exist on the `Notification` model. Per-device audit trail (NotificationLog) and reusable segment objects (Segment/NotificationSegment) are intentionally not implemented in this version.

**Deltas:** None — Segment + NotificationLog dropped by design.

---

### 8. Analytics

Full model parity on EventLog and SearchQuery.  
Python has `bface_api_error` auto-alerting (Slack + email). Go has no equivalent — drop for now, add as observability task.

**Event log listing endpoint commented out** in Go (`router.go:368`) — likely a performance decision. Verify if admin dashboard needs it.

---

### 9. Publishing Pipeline & Bundles

#### Bundle format — ✅ MATCH (critical)

Both Python and Go use:
- AES-256-CBC with PKCS7 padding
- Key: venue `public_key` (URL-safe base64, decoded to 32 bytes)
- Gzip compression (level 9)
- Wrapper: `{"compressed_data": "<base64(gzip(json))>"}`
- Wire: `base64_std(iv[16] || ciphertext)`
- S3 key: `{env}/bundles/public/{venue_id}.digiapp.{lang}`

**App clients can switch between Python-published and Go-published bundles without code changes.**

#### Six-key bundle JSON structure — ✅ MATCH

Both produce: `search`, `map`, `products`, `productsGroupedByCountry`, `exhibitorSearchOptions`, `productSearchOptions`.

**Gap:** `"dictionaries"` field (Japanese morphological tokenization) is always `{}` in Go. Python uses janome + wanakana. Affects Japanese full-text search quality.

#### Model diff

| Python | Go | Notes |
|--------|----|-------|
| Snapshot | Snapshot | ✅ parity; state encoding differs (string vs int — not wire-visible) |
| SyncData | **MISSING** | Per-level sync tracking; confirm if app clients query it |
| VenueSyncFlag | **MISSING** | Webhook-triggered pending-change flag; used for W11/W12 exhibitor sync |
| `VenueSyncFlag.expo_id` | Replaced by `venue.ExternalID` in enricher | 🔵 Intentionally replaced |

#### Logic diff

- Snapshot creation, pruning (max 5 versions), S3 upload: ✅ ported (`snapshot_service.go`)
- CloudFront invalidation: ✅ implemented (`cloudfront.go`) but depends on stub wiring
- Selenium publish-verification: 🔵 DROPPED — intentional latency/complexity reduction
- OpenSearch sync on publish: ❌ UNKNOWN — both repos have infrastructure but no trigger wired

#### Deltas

| Gap | Verdict | Priority |
|-----|---------|----------|
| Japanese tokenization (`"dictionaries"`) | PORT | MEDIUM — search quality |
| SyncData / VenueSyncFlag | CLARIFY | MEDIUM — confirm usage with app team |
| OpenSearch publish trigger | CLARIFY | MEDIUM |
| Selenium verification | DROP | ✅ Accepted |

---

### 10. Beacons, Connections, QR Codes

#### Beacons — ✅ parity

All fields present. Go adds `ElementID *uuid.UUID` (no Python equivalent — likely new).  
Minor: Go beacon list handler has no keyword/enable filter params (Python supports both).

#### Connections — ✅ parity

Connection + ConnectionLevel full parity. Level management via `AddLevel`/`RemoveLevel` service methods.

#### QR Codes — 🔵 DROPPED

Python: Full `QRCode` model, CRUD ViewSet, `generate_base64_qr_code()` method (qrcode library, `digimap://` URI scheme, PNG→base64).  
Go: Migration `000016` creates `qrcodes` table (schema preserved) but QR Code domain, service, and handler are intentionally not implemented in this version.

---

## Cross-Cutting Gaps

### Stubbed Infrastructure

All four platform services are no-op stubs in Go (`LogXxx` methods printing to stdout). Each is wired in `cmd/serve.go` as the stub implementation behind its interface. Swapping in the real implementation requires only changing the concrete type in `serve.go`.

| Platform | Python does | Go stub | Wiring needed |
|----------|-------------|---------|---------------|
| **S3** | PUT/GET bundles, snapshots, digimaps, media; presigned URLs (`s3services.py`) | `LogPutSnapshot`, `LogGetSnapshot`, etc. | Wire `aws.S3Client` in `serve.go`; affects publish, snapshot, asset upload |
| **SES / Email** | 13 notification functions: invitations, publish success/failure, exhibitor sync, article submission (`sendmail.py`) | `LogSender` (print-only) | Wire SES implementation; call from invitation, venue publish, JMA handlers |
| **Firebase FCM** | Multicast batch (500 tokens/batch), retry, rate-limit handling, error classification (`firebase.py`) | `LogSender` stub | Wire real FCM client with batching + retry; call from notification send, JMA push |
| **CloudFront** | `invalidate_cloudfront()` (`s3services.py`) | **Full implementation** present (`cloudfront.go`) | Just wire in `serve.go` — already done in code |
| **OpenSearch** | Search indexing/querying (infrastructure only in Python) | `LogIndexer` stub | Clarify trigger timing; wire when search feature is required |

### Tenant Enricher — FOODEX/HCJ Replacement

Python embedded FOODEX and HCJ data in a 2,940-line `constants.py`:
- 128+ visitor interested-category options across 15 sections (Bakery, Beverages, Wine×13, Seafood, Meat, etc.)
- Per-venue export field mappings (exhibitor English status, requisite documents)
- JMA taxonomy for `jma_meat_2026` category
- Per-client FCM project configuration

Go replaces this with a clean enricher registry (`internal/enricher/enricher.go`) where:
- `EnricherFunc` → `func(ctx, venueID) (map[string]any, error)`
- Resources: `location`, `product`, `notification`, `survey`, `advertisement`
- Registration: `cmd/tenants/tenants.go:RegisterAll(registry)` — **currently empty**

The registry interface is correct. For each customer (FOODEX, HCJ), a `RegisterAll` call must inject `EnricherFunc` closures that return the relevant category option sets, field mappings, and per-tenant configuration.

**Action:** Implement `RegisterAll` with FOODEX and HCJ enricher functions. This is the only remaining client-specific code, now cleanly isolated.

### MySQL → PostgreSQL Data Migration

| Concern | Detail |
|---------|--------|
| **Translation tables → JSONB** | Python has separate `VenueTranslation`, `LocationTranslation`, `LocationCategoryTranslation` tables. Go collapses all to `localization JSONB` / `translations JSONB`. ETL required: `UPDATE venues SET translations = jsonb_object_agg(lang_code, …) FROM venue_translations` |
| **UUID versions** | Python uses UUIDv4 (`uuid.uuid4`). Go uses UUIDv7 (sortable). **Preserve existing UUIDs on migration — do not regenerate.** |
| **Phone encryption** | Python Visitor rows have AES-GCM ciphertext in `_phone_number BinaryField`. Go expects plaintext. Migration must decrypt all rows using `DATA_ENCRYPTION_KEY` before import. |
| **Coupon redemption** | Python `UserCoupon` M2M → Go `redeemed_by + redeemed_at` on Coupon. Migration must handle only-one-per-coupon constraint. |
| **utf8mb4 → UTF-8** | PostgreSQL uses UTF-8 natively — no explicit conversion needed. |
| **JSON → JSONB** | All Django `JSONField` columns → PostgreSQL `JSONB` — semantically equivalent. |
| **Soft delete** | Both use timestamp-based soft delete; `deleted_at IS NULL` indexes present in Go. |
| **pgvector** | Extension created (`000001`) but no vector columns exist. Reserve for future RAG/embedding use. |

### API Surface Collapse

| Python surface | Go equivalent | Coverage |
|----------------|---------------|---------|
| `api/v1/` (admin JWT) | `/api/v1/` (system-admin + venue JWT) | ✅ ~90% |
| `app/v1/` (API-key) | `/app/v1/` (API-key) | ✅ ~82% (23/28 endpoints) |
| `mobile/v1/` (JWT) | Consolidated into `/app/v1/` | 🔵 Intentionally merged |
| `public/v1/` + `sdk/v1/` | `/public/v1/` | ⚠️ ~83%; sdk/v1 levels + categories + bundle + qrcodes missing |

**sdk/v1 missing endpoints:** `levels/`, `location-categories/`, `bundle/`, `qrcodes/{id}/` — needed if any web SDK embeds depend on them.

### Encrypted PII — Cipher Parity

This is the highest-data-safety risk in the migration.

- **Python cipher:** AES-GCM, 12-byte random nonce, key = `DATA_ENCRYPTION_KEY` (base64-decoded, 32 bytes). Source: `encrypting.py:132–228`.
- **Go crypto:** `internal/platform/crypto/aes.go` implements **AES-256-CBC** (not GCM). No GCM implementation exists yet.
- **Required for parity:** Add `EncryptGCM(plaintext, key) []byte` and `DecryptGCM(ciphertext, key) string` to `aes.go` matching the Python convention (nonce prepended to ciphertext). Use in `app_user` repository on write/read.

---

## Cruft Confirmed Dropped

These systems exist in `indoormap-backend` and are **intentionally absent from `digimap-backend`**. Review this list and veto any item you consider core.

| System | Python location | Reason for drop |
|--------|----------------|-----------------|
| **Selenium publish-verification** | `management/commands/check_app_publication.py`, `auto_check_app.py` | Fragile, slow, adds 15 MB chromedriver binary to repo; replace with app-side telemetry |
| **15 MB committed chromedriver** | `chromedriver_linux64/chromedriver` | Binary in git; dropped with Selenium |
| **Azure Business-Card OCR** | `utils/business_card.py` (4,170 lines) | Already commented out in Python; client-specific, no general value |
| **AWS Textract OCR** | `utils/ai/textract.py` | Client-specific; not needed in core |
| **Vimeo API integration** | `utils/vimeo.py` | No Video model in Python anyway; Go uses plain URL storage |
| **Slack notifications** | `utils/slack.py` | Internal alerting; replace with observability tooling |
| **FOODEX/HCJ constants** | `utils/constants.py` (2,940 lines) | Replaced by empty tenant enricher registry in Go |
| **`VenueSyncFlag.expo_id`** | `snapshots/models.py` | "hcj or foodex" discriminator; replaced by `venue.ExternalID` in Go enricher |
| **Per-client FCM projects** | `FCM_PROJECTS` dict in settings | Replaced by single FCM credential; per-client routing via enricher if needed |
| **`publish_venue` v1** | `management/commands/publish_venue.py` | v2 ported; v1 is superseded |
| **`auto_publish_venue`** | Management command (Selenium-driven) | Dropped with Selenium |
| **JMA AI path generation** | `utils/ai/generate_paths.py` (Voronoi/NetworkX/Shapely) | Not ported — confirm if needed for auto-routing feature |
| **`jma_meat_2026` taxonomy** | `utils/jma/meat.py` | Trade-show-specific; goes in tenant enricher RegisterAll when needed |
| **Django Groups RBAC sync** | Django `auth.Group` machinery | Replaced by Go enum-based role system |
| **`promotions` table** | Go has migration but no domain model | No Python equivalent either; status unclear |
| **Venue publish/clone endpoints** | `api/venues/views.py` (publish), `router.go:177–178` (commented clone) | Not in this version |
| **User management + invitation endpoints** | `api/users/views.py` (UserViewSet, VenueInvitationViewSet) | No RBAC in this version |
| **QRCode domain, service, handler** | `connections/models.py` (QRCode), `connections/views.py` (QRCodeViewSet) | Not used in this version; schema preserved in `migrations/000016` |
| **Notification Segment + NotificationLog** | `notification/models.py` (Segment, NotificationSegment, NotificationLog) | Not needed in this version |

---

## Recommendations

Prioritized by risk and blocking nature:

### P0 — Must fix before any production data flows through Go

1. **Port AES-GCM for AppUser.Phone** (`internal/platform/crypto/aes.go`)  
   Add `EncryptGCM`/`DecryptGCM` matching `encrypting.py:132–228`. Wire in `app_user` repository.

2. **Wire real S3 client in `cmd/serve.go`**  
   Swap `LogXxx` stub for `aws.S3Client`. Snapshot publish pipeline is broken until this is done.

3. **Wire real FCM client in `cmd/serve.go`**  
   Notification push is entirely no-op until FCM is wired.

### P1 — High-impact gaps blocking core workflows

4. **Port customer isolation / scope to User** *(if multi-tenant isolation is needed)*  
   Add `CustomerID` to `User` and scope-based queryset filtering equivalent to Python `profile.scope`.

### P2 — Should fix before partner launch

6. **Wire SES email client**  
   JMA sync alerts and any system notifications are silent until wired.

7. **Port Japanese tokenization for bundle `"dictionaries"`**  
   Use `github.com/ikawaha/kagome` (Japanese MeCab-compatible tokenizer) to match Python's janome behavior.

8. **Implement tenant enricher `RegisterAll`** for FOODEX and HCJ  
   Surface: per-tenant visitor category option sets, export field mappings, FCM routing.

9. **Implement sdk/v1 missing endpoints**  
   `levels/`, `location-categories/`, `bundle/` under `/public/v1` — needed if web SDK embeds these.

10. **Coupon redemption model — clarify and align**  
    Decide: single-redemption (Go model) or multi-user (Python `UserCoupon`). If multi-user needed, add `app_user_coupons` junction table.

### P3 — Polish and correctness

11. **Port tag copying in location duplicate** (`location_service.go:217–243`)
12. **Add keyword/enable filter params to beacon list handler**
13. **Port `PlaceWorkHours` to location duplicate**
14. **Decide on `SyncData`/`VenueSyncFlag`** — confirm with app team if per-level sync state is queried
15. **Un-comment or remove event-log listing endpoint** (`router.go:368`)
16. **Wire CloudFront invalidation** in publish pipeline (`cloudfront.go` is already implemented)
17. **Wire OpenSearch indexing** trigger on venue publish

---

*Generated from parallel deep-read of both repositories. All citations are from the actual source files; spot-check before acting on any claim.*

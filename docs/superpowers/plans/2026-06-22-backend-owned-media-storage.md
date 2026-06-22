# Backend-Owned Media Storage Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Make the Go backend exclusively upload, replace, delete, and presign upload-owned media while PostgreSQL stores collision-safe S3 object keys and article periods remain date-only.

**Architecture:** Add a shared `MediaService` over the S3 adapter, then migrate entity endpoints to multipart requests containing a JSON `data` part and named files. Implement article as the first complete backend/frontend slice, apply the contract to remaining upload-owned fields, and finish with an idempotent legacy-media migration.

**Tech Stack:** Go, Gin, pgx/v5, AWS SDK v2 S3, PostgreSQL, React 18, TypeScript, TanStack Query, Jest.

---

### Task 1: Add Collision-Safe Storage Primitives

**Files:**
- Modify: `internal/platform/storage/keys.go`
- Create: `internal/platform/storage/keys_test.go`
- Modify: `internal/platform/storage/s3.go`
- Create: `internal/platform/storage/s3_test.go`

- [ ] **Step 1: Write the failing key tests**

```go
func TestMediaKeyIsUniqueAndRecordScoped(t *testing.T) {
	recordID := uuid.New()
	first := MediaKey("develop", "articles", recordID, "images", uuid.New(), ".png")
	second := MediaKey("develop", "articles", recordID, "images", uuid.New(), ".png")
	require.NotEqual(t, first, second)
	require.True(t, OwnsMediaKey("develop", "articles", recordID, "images", first))
	require.False(t, OwnsMediaKey("develop", "articles", uuid.New(), "images", first))
}
```

- [ ] **Step 2: Verify RED**

Run: `env GOCACHE=/private/tmp/go-build-cache go test ./internal/platform/storage -run 'TestMediaKey' -v`

Expected: compilation fails because `MediaKey` and `OwnsMediaKey` do not exist.

- [ ] **Step 3: Implement immutable keys**

```go
func MediaKey(env, entity string, recordID uuid.UUID, field string, uploadID uuid.UUID, ext string) string {
	return path.Join(env, "media", entity, recordID.String(), field, uploadID.String()+strings.ToLower(ext))
}

func OwnsMediaKey(env, entity string, recordID uuid.UUID, field, key string) bool {
	prefix := path.Join(env, "media", entity, recordID.String(), field) + "/"
	return strings.HasPrefix(key, prefix)
}
```

- [ ] **Step 4: Write failing S3 tests**

Test that `PutMedia` forwards body, key, content type, and length. Test construction without static keys and verify the default AWS credential chain remains active.

- [ ] **Step 5: Implement S3 media upload and task-role credentials**

```go
type Storer interface {
	PutMedia(context.Context, string, string, int64, io.Reader) error
	PresignDownload(context.Context, string, time.Duration) (string, error)
	DeleteObject(context.Context, string) error
	// Retain existing bundle methods.
}
```

Implement `PutMedia` with `s3.PutObject`. Remove `credentials.NewStaticCredentialsProvider`; load region plus the AWS default credential chain. Extend `LogStorer` consistently.

- [ ] **Step 6: Verify GREEN and commit**

Run: `env GOCACHE=/private/tmp/go-build-cache go test ./internal/platform/storage -v`

```bash
git add internal/platform/storage
git commit -m "fix: add collision-safe backend media storage"
```

### Task 2: Build and Wire the Shared Media Service

**Files:**
- Create: `internal/service/media_service.go`
- Create: `internal/service/media_service_test.go`
- Modify: `cmd/serve.go`

- [ ] **Step 1: Write failing service tests**

```go
func TestMediaServiceUploadUsesUniqueKeysForSameFilename(t *testing.T) {
	storer := newRecordingStorer()
	svc := NewMediaService(storer, "develop")
	target := MediaTarget{Entity: "articles", RecordID: uuid.New(), Field: "images"}
	upload := imageUpload("image01.png")
	first, err := svc.Upload(context.Background(), target, upload)
	require.NoError(t, err)
	second, err := svc.Upload(context.Background(), target, imageUpload("image01.png"))
	require.NoError(t, err)
	require.NotEqual(t, first, second)
}
```

Also test different records, MIME mismatch, oversized input, undecodable images, presign failure returning nil URL, and deletion rejecting another record's key.

- [ ] **Step 2: Verify RED**

Run: `env GOCACHE=/private/tmp/go-build-cache go test ./internal/service -run TestMediaService -v`

- [ ] **Step 3: Implement the service contract**

```go
type MediaTarget struct { Entity string; RecordID uuid.UUID; Field string }
type MediaUpload struct { Filename string; ContentType string; Size int64; Reader io.Reader }
type MediaService interface {
	Upload(context.Context, MediaTarget, MediaUpload) (string, error)
	URL(context.Context, string) *string
	DeleteOwned(context.Context, MediaTarget, string) error
}
```

Bound reads, detect MIME from content, decode image configuration for image policies, derive extensions from validated content, generate a new UUID per upload, and presign reads for 15 minutes.

- [ ] **Step 4: Verify GREEN, wire once in `cmd/serve.go`, and commit**

Run: `env GOCACHE=/private/tmp/go-build-cache go test ./internal/service -run TestMediaService -v`

```bash
git add internal/service/media_service.go internal/service/media_service_test.go cmd/serve.go
git commit -m "feat: add backend media service"
```

### Task 3: Add Date-Only and Multipart Contracts

**Files:**
- Create: `internal/dto/date.go`
- Create: `internal/dto/date_test.go`
- Create: `internal/handler/multipart.go`
- Create: `internal/handler/multipart_test.go`

- [ ] **Step 1: Write failing date tests**

```go
func TestDateJSONIsCalendarOnly(t *testing.T) {
	var d Date
	require.NoError(t, json.Unmarshal([]byte(`"2026-06-18"`), &d))
	require.Equal(t, `"2026-06-18"`, string(mustJSON(d)))
	require.Error(t, json.Unmarshal([]byte(`"2026-06-18T00:00:00Z"`), &d))
}
```

- [ ] **Step 2: Verify RED, then implement `dto.Date`**

Run: `env GOCACHE=/private/tmp/go-build-cache go test ./internal/dto -run TestDate -v`

Implement strict `2006-01-02` JSON parsing/serialization and UTC-midnight conversion only at the pgx boundary.

- [ ] **Step 3: Write failing multipart tests**

Cover valid `data` JSON, one file, repeated files, malformed JSON, bounded size, and conflicting replacement/remove instructions.

- [ ] **Step 4: Implement shared helpers**

```go
func bindMultipartData(c *gin.Context, dst any, maxMemory int64) error
func multipartFile(c *gin.Context, field string) (*multipart.FileHeader, bool, error)
func multipartFiles(c *gin.Context, field string) ([]*multipart.FileHeader, bool, error)
func mediaUpload(*multipart.FileHeader) (service.MediaUpload, io.Closer, error)
```

- [ ] **Step 5: Verify and commit**

Run: `env GOCACHE=/private/tmp/go-build-cache go test ./internal/dto ./internal/handler -run 'TestDate|TestMultipart' -v`

```bash
git add internal/dto/date* internal/handler/multipart*
git commit -m "feat: add date-only and multipart contracts"
```

### Task 4: Implement Atomic Article Media Replacement

**Files:**
- Modify: `internal/domain/article.go`
- Modify/Create tests: `internal/dto/article_dto.go`, `internal/dto/article_dto_test.go`
- Modify: `internal/repository/interfaces.go`, `internal/repository/postgres/article_repo.go`, `internal/repository/mocks/mocks.go`
- Create: `internal/repository/postgres/article_repo_integration_test.go`
- Modify/Create tests: `internal/service/article_service.go`, `internal/service/article_service_test.go`
- Modify/Create tests: `internal/handler/article_handler.go`, `internal/handler/article_handler_test.go`

- [ ] **Step 1: Write failing DTO tests**

Assert article periods accept/emit `YYYY-MM-DD`; image responses contain `image` and nullable `image_url`; omitted images differ from replace/clear.

- [ ] **Step 2: Write failing repository tests**

Insert an article with `old-key`, call `UpdateWithImages` with `new-key`, and assert one active new row. Force an insert failure and assert scalar fields and images both roll back.

- [ ] **Step 3: Implement atomic persistence**

```go
type ArticleMediaChange struct { Replace bool; Keys []string }
UpdateWithImages(context.Context, *domain.Article, ArticleMediaChange) error
```

Use one pgx transaction. Preserve rows if `Replace=false`; otherwise soft-delete active rows and insert ordered keys.

- [ ] **Step 4: Write failing service tests**

Cover S3 failure before DB mutation, DB failure compensation, successful post-commit old-key deletion, same-name uniqueness, and ownership isolation.

- [ ] **Step 5: Implement article orchestration**

Upload all replacement files first, persist article and keys atomically, delete new keys on rollback, then delete verified old keys after commit.

- [ ] **Step 6: Write failing handler tests and implement multipart handling**

Send `data` with date-only values and repeated `images` files. Add JSON-only preserve and `remove_images=true` tests. Enrich reads with presigned URLs without failing the record when one presign fails.

- [ ] **Step 7: Verify and commit**

Run: `env GOCACHE=/private/tmp/go-build-cache go test ./internal/dto ./internal/repository/postgres ./internal/service ./internal/handler -run 'TestArticle|TestDate' -v`

```bash
git add internal/domain/article.go internal/dto/article_dto* internal/repository internal/service/article_service* internal/handler/article_handler*
git commit -m "fix: store and replace article media in s3"
```

### Task 5: Convert the CMS Article Form

**Files (frontend repository):**
- Create: `src/services/mediaFormData.ts`
- Create: `src/services/__tests__/mediaFormData.test.ts`
- Modify: `src/services/index.ts`
- Modify: `src/repositories/promotion.repo.ts`
- Modify: `src/hooks/queries/usePromotionQuery.ts`
- Modify: `src/pages/promotion/article/components/ArticleFormTypes.ts`
- Modify: `src/pages/promotion/article/hooks/useArticleForm.ts`
- Create: `src/pages/promotion/article/hooks/__tests__/useArticleForm.test.tsx`
- Delete: `src/services/S3Service.ts`

- [ ] **Step 1: Write failing FormData tests**

```ts
it('sends date-only data and the original File', () => {
  const file = new File(['png'], 'image01.png', { type: 'image/png' });
  const body = buildMediaFormData(
    { title: 'Updated', published_period_start: '2026-06-18' },
    { images: [file] }
  );
  expect(JSON.parse(String(body.get('data'))).published_period_start).toBe('2026-06-18');
  expect(body.getAll('images')).toEqual([file]);
});
```

Also test no-file JSON fallback, explicit clear, and that existing keys/URLs are not appended as files.

- [ ] **Step 2: Verify RED and implement FormData support**

Run: `yarn test --watchAll=false --runTestsByPath src/services/__tests__/mediaFormData.test.ts`

Create `buildMediaFormData`; let the browser set multipart boundaries; allow article repository/mutations to accept `FormData`.

- [ ] **Step 3: Write failing form-hook tests**

Assert `{image,image_url}` normalizes to `{key,data_url,existing:true}`, replacement files remain `File`, dates stay date-only, omission preserves, and clear sends `remove_images`.

- [ ] **Step 4: Remove base64/direct AWS behavior**

Delete `getBase64` use and `S3Service.ts`. Submit only actual replacement files, use `image_url` for previews, and retain `image` as identity.

- [ ] **Step 5: Verify and commit**

Run: `yarn test --watchAll=false --runTestsByPath src/services/__tests__/mediaFormData.test.ts src/pages/promotion/article/hooks/__tests__/useArticleForm.test.tsx`

Run: `yarn lint`

```bash
git add src/services src/repositories/promotion.repo.ts src/hooks/queries/usePromotionQuery.ts src/pages/promotion/article
git commit -m "fix: upload article media through backend"
```

### Task 6: Apply Singular Media Semantics Repository-Wide

**Files (backend):**
- Modify: `internal/domain/venue.go`, `user.go`, `location.go`, `product.go`, `event.go`, `advertisement.go`, `product_plaza.go`, `video.go`, `asset.go`, `level_type.go`
- Modify: `internal/dto/venue_dto.go`, `user_dto.go`, `location_dto.go`, `product_dto.go`, `event_dto.go`, `ad_dto.go`, `product_plaza_dto.go`, `video_dto.go`, `asset_dto.go`, `level_type_dto.go`
- Modify: `internal/handler/customer_handler.go`, `user_handler.go`, `venue_handler.go`, `location_handler.go`, `product_handler.go`, `event_handler.go`, `ad_handler.go`, `product_plaza_handler.go`, `video_handler.go`, `asset_handler.go`, `level_type_handler.go`
- Modify: `internal/service/customer_service.go`, `user_service.go`, `venue_service.go`, `location_service.go`, `product_service.go`, `event_service.go`, `ad_service.go`, `product_plaza_service.go`, `video_service.go`, `asset_service.go`, `level_type_service.go`
- Modify: `internal/repository/interfaces.go`, `internal/repository/mocks/mocks.go`, and PostgreSQL repositories `customer_repo.go`, `user_repo.go`, `venue_repo.go`, `location_category_repo.go`, `location_repo.go`, `product_repo.go`, `event_repo.go`, `ad_repo.go`, `product_plaza_repo.go`, `video_repo.go`, `asset_repo.go`, `level_type_repo.go`
- Create focused `*_media_test.go` files beside each changed handler/service/repository package; modify existing entity tests when that package already owns equivalent CRUD test setup.
- Modify: `cmd/serve.go`

**Files (frontend):**
- Modify: `src/services/index.ts`; repositories `customer.repo.ts`, `user.repo.ts`, `venue.repo.ts`, `location.repo.ts`, `product.repo.ts`, `event.repo.ts`, and `promotion.repo.ts`; query hooks `useCustomerQuery.ts`, `useUserQuery.ts`, `useVenueQuery.ts`, `useLocationQuery.ts`, `useProductQuery.ts`, `useEventQuery.ts`, `usePremiumQuery.ts`, and `usePromotionQuery.ts`.
- Modify: `src/pages/settings/index.tsx`, `src/components/layouts/header/settings.tsx`, `src/components/location/components/LocationFormPopup.tsx`, `src/components/premium/MarkerFormPopup.tsx`, `src/components/premium/AdFormMediaSection.tsx`, `src/components/event/components/popup.tsx`, `src/components/event/components/event.tsx`, `src/pages/promotion/top-logo/top-logo-detail.tsx`, `src/pages/promotion/top-logo/components/TopLogoImage.tsx`, `src/pages/promotion/product-plaza/components/MediaUploadSection.tsx`, `src/pages/promotion/product-plaza/useEditProductPlazaForm.ts`, and `src/pages/promotion/product-plaza/useSubmitProductForm.ts`.
- Create focused tests adjacent to each changed form/hook and extend `src/services/__tests__/mediaFormData.test.ts` with each named file field.

- [ ] **Step 1: Add failing table-driven tests per entity**

Each entity test must prove omitted preserves, supplied replaces, explicit remove clears, same filenames produce different keys, rollback compensates, and another record's key is never deleted.

- [ ] **Step 2: Implement backend batches**

Batch A: customer, user, venue, level type. Batch B: location category and location. Batch C: product, event, advertisement, product plaza, video thumbnail, and asset.

Store only keys and expose a sibling URL property using the persisted field name plus `_url`. Do not convert CTA/navigation links, Vimeo/video URLs, attachment `source_url`, generated QR/base64, bundles, or theme storage paths.

- [ ] **Step 3: Add frontend FormData tests and convert forms**

For every form, assert files are appended, existing URLs are preview-only, remove flags are explicit, and no-file JSON updates preserve values. Reuse `buildMediaFormData`.

- [ ] **Step 4: Verify and commit each batch**

Backend: `env GOCACHE=/private/tmp/go-build-cache go test ./internal/service ./internal/handler ./internal/repository/postgres -run 'Test(Customer|User|Venue|LevelType|LocationCategory|Location|Product|Event|Advertisement|ProductPlaza|Video|Asset).*Media' -v`

Frontend: `yarn test --watchAll=false` then `yarn lint`.

```bash
git commit -m "feat: store account and venue media in s3"
git commit -m "feat: store location media in s3"
git commit -m "feat: store promotion media in s3"
```

### Task 7: Apply Collection and Attachment Semantics

**Files:**
- Modify backend domain files: `internal/domain/article.go`, `location.go`, `event.go`, `product.go`.
- Modify backend DTO/handler/service files: `article_dto.go`, `location_dto.go`, `event_dto.go`, `product_dto.go`; `article_handler.go`, `location_handler.go`, `event_handler.go`, `product_handler.go`; `article_service.go`, `location_service.go`, `event_service.go`, `product_service.go`.
- Modify persistence: `internal/repository/interfaces.go`, `internal/repository/mocks/mocks.go`, `internal/repository/postgres/article_repo.go`, `location_repo.go`, `event_repo.go`, `product_repo.go`.
- Modify focused tests: existing `article_service_test.go`, `location_service_test.go`, `event_service_test.go`, and `crud_services_test.go`; create handler and PostgreSQL integration tests for each collection transaction.
- Modify frontend: `src/pages/promotion/article/hooks/useArticleForm.ts`, `src/pages/promotion/article/submit-article.tsx`, `src/pages/promotion/article/components/ArticleImages.tsx`, `src/components/location/components/LocationFormPopup.tsx`, `src/components/event/components/popup.tsx`, `src/components/product/detail.tsx`, and `src/services/index.ts`.
- Extend: `src/services/__tests__/mediaFormData.test.ts` and create focused collection tests beside each changed form hook.

- [ ] **Step 1: Write failing collection tests**

Prove omitted collection preservation, ordered full replacement, explicit clear, transaction rollback, compensation of all new uploads, and post-commit deletion of all replaced owned keys.

- [ ] **Step 2: Implement transactional collection methods**

Use repository methods such as `UpdateWithImages`/`UpdateWithAttachments`; handlers must not loop over independent persistence calls. Preserve attachment `source_url`; expose uploaded `file` key plus `file_url`.

- [ ] **Step 3: Convert frontend repeated multipart parts and verify**

Preserve server media unless changed or cleared. Never convert returned URLs to base64.

```bash
git commit -m "feat: replace media collections atomically"
```

### Task 8: Migrate Existing Media Values

**Files:**
- Create: `internal/service/media_migration_service.go`
- Create: `internal/service/media_migration_service_test.go`
- Create: `cmd/migrate_media.go`
- Create: `cmd/migrate_media_test.go`
- Modify: `cmd/root.go`

- [ ] **Step 1: Write failing migration tests**

Use an HTTP test server, fake storer, and test rows. Prove reachable URLs upload then store keys, configured S3 URLs convert without copying, `example.com`/unreachable URLs become null, external fields are skipped, valid keys are skipped, reruns make no changes, and failures report record IDs without reverting successes.

- [ ] **Step 2: Implement an explicit field registry**

Descriptors contain table, primary key, media column, entity, field, nullable representation, and parent-record column. Do not classify columns dynamically by name.

- [ ] **Step 3: Implement resumable CLI behavior**

Provide `go run . migrate-media --dry-run`, `--apply`, and `--batch-size`. Update a row only after its destination object exists. Report migrated, converted, cleared, skipped, and failed totals without logging credentials or presigned queries.

- [ ] **Step 4: Verify and commit**

Run: `env GOCACHE=/private/tmp/go-build-cache go test ./internal/service ./cmd -run 'TestMediaMigration|TestMigrateMedia' -v`

```bash
git add internal/service/media_migration* cmd/migrate_media*
git commit -m "feat: migrate legacy media to s3 keys"
```

### Task 9: Remove Obsolete Upload Presigning and Refresh Contracts

**Files:**
- Modify: `internal/service/product_service.go`, `internal/handler/product_handler.go`, `internal/handler/router.go`, `internal/dto/product_dto.go`
- Regenerate: `docs/swagger.json`, `docs/swagger.yaml`, `docs/docs.go`
- Modify applicable backend/frontend `AGENTS.md` files.

- [ ] **Step 1: Add a failing router test**

Assert `POST /storage/presign-upload` is absent while multipart entity endpoints retain current authentication and RBAC.

- [ ] **Step 2: Remove client-controlled presigning**

Delete upload presign service/DTO/route and frontend service. Keep download presigning backend-only through `MediaService`.

- [ ] **Step 3: Regenerate Swagger and perform DOX updates**

Run: `$(go env GOPATH)/bin/swag init --generalInfo cmd/serve.go --output docs`

Update nearest owning backend docs for storage, services, handlers, DTOs, repositories, and migration command; update frontend services/pages guidance. Refresh Child DOX indexes only where structure changes.

- [ ] **Step 4: Commit**

```bash
git commit -m "refactor: centralize media storage in backend"
```

### Task 10: Full Verification and Dry Run

**Files:**
- Verify all changes from Tasks 1-9.

- [ ] **Step 1: Run backend verification**

Run: `env GOCACHE=/private/tmp/go-build-cache go test ./...`

Run: `git diff --check HEAD~9`

Expected: all tests pass and no whitespace errors.

- [ ] **Step 2: Run frontend verification**

Run: `yarn test --watchAll=false`

Run: `yarn lint`

Run: `yarn build`

Run: `git diff --check HEAD~4`

Expected: tests/build pass, lint has no new errors, and no whitespace errors.

- [ ] **Step 3: Run migration dry-run**

Run: `env GOCACHE=/private/tmp/go-build-cache go run . migrate-media --dry-run --batch-size 100`

Expected: no DB/S3 writes; classified fields are reported; placeholders count as cleared; external/generated fields are absent.

- [ ] **Step 4: Exercise the original regression**

Create two records with files named `image01.png`; verify different keys. Replace one with another `image01.png`; verify only that record changes, the other still renders, the replaced owned object is removed, dates remain `YYYY-MM-DD`, and retrieval returns temporary URLs.

- [ ] **Step 5: Final DOX and working-tree review**

Walk every changed path through its root/child `AGENTS.md`, remove stale direct-S3/base64 guidance, and report migration/config prerequisites.

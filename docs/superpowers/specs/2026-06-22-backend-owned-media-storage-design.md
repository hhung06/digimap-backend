# Backend-Owned Media Storage Design

## Goal

Make the Go backend the sole owner of uploaded-media persistence across Digimap. Clients send files to entity create/update endpoints; the backend stores unique objects in S3, stores only object keys in PostgreSQL, and returns temporary presigned URLs on reads.

Article publication start and end remain date-only values represented as `YYYY-MM-DD`.

## Scope

Apply the media contract to upload-owned fields across the backend and CMS frontend:

- customer and user images
- venue logos
- location and location-category icons, logos, and images
- product images and uploaded attachment files
- event banner, icon, and child images
- article images
- advertisement and promotion images
- uploaded thumbnails and level-type icons

Do not treat external resources as uploaded media. CTA and navigation links, Vimeo/video URLs, attachment `source_url`, generated QR base64 data, generated bundles, and theme storage paths retain their current purpose.

Before implementation, audit every media-like field and classify it as upload-owned, external, or generated. This classification is the source of truth for endpoint and migration coverage.

## Request Contract

Entity create and update endpoints that accept files use `multipart/form-data`:

- `data` contains the entity's normal JSON payload.
- Named file parts contain singular media such as `image` or `banner_image`.
- Repeated file parts contain collections such as `images`.

Update semantics are explicit:

- An omitted file part preserves the existing object key.
- A supplied singular file replaces the existing object.
- `remove_<field>: true` removes a singular object.
- An omitted repeated file part preserves the existing collection.
- Supplied repeated files replace the complete collection.
- `remove_<field>: true` clears a collection.

Base64 file payloads are not part of the normal contract because they add bandwidth and memory overhead. JSON-only requests remain valid for endpoints and updates that contain no media changes.

## Storage Contract

The backend generates every object key. Original filenames are retained only as optional metadata and never determine object identity.

Keys follow this shape:

```text
media/{entity-type}/{record-id}/{field}/{upload-uuid}.{extension}
```

Every upload, including a replacement with the same filename, receives a new UUID. Existing S3 objects are never overwritten. This isolates files across entity types, records, fields, and repeated replacements.

The shared media service owns:

- file validation
- object-key generation
- S3 upload and deletion
- presigned download URL generation
- replacement cleanup and compensating cleanup

S3 uses the AWS default credential chain so ECS task-role credentials work without static keys.

## Write Flow

For create operations, the backend allocates the record ID before generating object keys. For updates, it loads the existing entity and media ownership before accepting replacements.

The write sequence is:

1. Parse entity data and media parts.
2. Validate file size, declared MIME type, extension, and decodability where applicable.
3. Generate unique object keys and upload new objects to S3.
4. In one database transaction, persist entity changes and replace relevant object keys or child media rows.
5. If the transaction fails, delete newly uploaded objects as compensation.
6. After commit, delete replaced objects whose ownership was verified against the entity and field prefix.

Old-object cleanup failures do not revert committed business data. They are logged with enough context for bounded retry. The implementation must not delete an object referenced by another record; unique per-upload keys make sharing invalid by construction.

## Read Contract

PostgreSQL stores only object keys. Read responses expose both the durable key and an expiring URL. For child media objects:

```json
{
  "image": "media/articles/<article-id>/images/<upload-id>.png",
  "image_url": "https://presigned-s3-url"
}
```

Singular media fields follow the same `<field>` and `<field>_url` convention. The frontend uses the URL for preview/display and never treats it as the persistent identifier.

If presigning one object fails, the API returns its key with a null URL and records the failure. A single presign failure does not make the entire entity response unavailable.

## Date-Only Contract

Article `published_period_start` and `published_period_end` accept and return `YYYY-MM-DD`. They use a dedicated date-only DTO type rather than `time.Time`, avoiding implicit timezone conversion. Domain and persistence mapping must preserve calendar dates exactly.

## Legacy Migration

Provide a resumable, idempotent migration command that processes every upload-owned field and child media row:

- Convert identifiable URLs for objects already in the configured S3 bucket to object keys without copying.
- Download other reachable legacy URLs, upload each file under a new unique key, and replace the stored value.
- Set unreachable values, including placeholder `example.com` media, to `NULL` or the field's established empty representation.
- Never modify external or generated URL fields.
- Record migrated, converted, cleared, skipped, and failed identifiers and counts.

The migration updates a database value only after its destination object exists. Reruns skip already-valid object keys and safely resume incomplete work.

## Frontend Changes

The CMS removes direct browser AWS SDK usage and base64 conversion for persisted media. Forms send files and entity data to backend multipart endpoints. Existing media state uses returned object keys for identity and presigned URLs for previews.

Frontend request adapters must preserve JSON-only behavior when no file change is requested and must distinguish preserve, replace, and remove operations.

## Error Handling

- Invalid media returns app code `1000` with field-specific messages.
- S3 upload failure leaves database state unchanged.
- Database failure triggers cleanup of newly uploaded objects.
- Post-commit old-object cleanup failure is observable and retryable.
- Unsupported request combinations, such as supplying a file and its remove flag together, are rejected.
- Media errors must not be reported as successful entity updates.

## Verification

Automated coverage must include:

- identical filenames uploaded to different entities and records produce different keys
- repeated same-name replacement on one record creates a new key without affecting other records
- omitted, replaced, and explicitly removed singular media
- omitted, replaced, and cleared media collections
- S3 failure before database mutation
- database rollback and new-object compensation
- old-object deletion only after commit
- presigned read URLs with durable keys retained
- null URL behavior when presigning fails
- date-only parse, serialization, persistence, and timezone independence
- migration of reachable URLs, conversion of existing S3 URLs, clearing unreachable placeholders, and safe reruns
- frontend multipart request construction and preview normalization

Run focused tests per changed package and frontend module, then the complete Go test suite, frontend test suite, lint, build, and `git diff --check` in both repositories.

## DOX Impact

Implementation must update the nearest owning `AGENTS.md` contracts for handlers, DTOs, services, repositories, storage adapters, migrations/commands, and frontend services/pages where the durable media workflow changes. This design document alone does not change those runtime contracts.

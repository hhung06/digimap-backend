# Purpose

Adapters for all external services. Each sub-package wraps a third-party dependency behind an interface so business logic never imports SDK types directly.

# Ownership

Owned by engineers integrating or maintaining external service dependencies.

# Local Contracts

- Each adapter exposes an interface; the implementation is unexported.
- When credentials are absent locally, adapters fall back to `Log*` stubs that print instead of calling the real service.
- No business logic here — adapters translate between service contracts and external APIs only.
- Backend-owned media uses immutable `{env}/media/{entity}/{record_id}/{field}/{upload_id}.{ext}` keys; verify the full record-field prefix before deletion.
- S3 storage loads the AWS default credential chain with the configured region so task-role credentials remain available.
- Sub-packages: `cache` (Redis), `cdn` (CDN URL signing), `crypto` (hashing/encryption), `database` (pgx pool), `email` (SMTP/SES), `firebase` (FCM push), `search` (Elasticsearch/Typesense), `storage` (S3/compatible).

# Work Guidance

When adding a new adapter:
1. Define the interface in the sub-package.
2. Provide a real implementation and a `Log*` stub for local dev.
3. Wire into `Dependencies` in `internal/handler/router.go`.

# Verification

```bash
go build ./internal/platform/...
```

# Child DOX Index

All sub-packages (`cache`, `cdn`, `crypto`, `database`, `email`, `firebase`, `search`, `storage`) are governed by this doc. No deeper AGENTS.md files needed at current complexity.

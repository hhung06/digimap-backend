# Load Testing

This directory contains k6 load tests for non-production Digimap API environments. The setup is read-heavy by default and requires explicit opt-in for mutating requests.

## Structure

- `scenarios/main.js` - k6 scenario runner.
- `data/test-data.json` - reusable IDs and request payloads.
- `env/*.env.example` - environment templates for local and staging runs.
- `reports/` - generated `summary.json` and `summary.html` output.

## Authentication

The API uses three access modes:

- Public endpoints: no auth, for `/health`, `/version`, `/public/v1/*`, and analytics tracking under `/api/v1/public/*`.
- Admin API: `Authorization: Bearer <jwt>` for `/api/v1/*`.
- App API: `Authorization: ApiKey <public_key>:<private_key>` for `/app/v1/*`.

Use staging/local credentials only. The runner refuses production-looking URLs unless `ALLOW_PRODUCTION_LOAD_TEST=true` is set.

## Load Profiles

- `normal`: steady baseline traffic for public, app, and admin read flows.
- `peak`: ramping virtual users for sustained high traffic.
- `burst`: arrival-rate spikes for short traffic surges.
- `error_retry`: invalid auth, invalid payloads, not-found/validation checks, and retry wrapper behavior.

## Endpoint Coverage

| Flow | Method | Path | Auth | Payload | Success | Failure |
|---|---:|---|---|---|---|---|
| Health | GET | `/health` | none | none | `200`, health body | `5xx` |
| Version | GET | `/version` | none | none | `200`, `code:0` | `5xx` |
| Public venue | GET | `/public/v1/venue-info?public_key=...` | none | none | `200`, `code:0` | `400`, `404` |
| Visitor surveys | GET | `/public/v1/visitor-surveys?public_key=...` | optional `Token` | none | `200`, `code:0` | `400`, `401`, `404` |
| Track event | POST | `/api/v1/public/venues/:id/events` | none | `track_event_valid` | `200`, `code:0` | `400`, `404` |
| Track search | POST | `/api/v1/public/venues/:id/searches` | none | `track_search_valid` | `200`, `code:0` | `400`, `404` |
| App locations | GET | `/app/v1/locations` | API key | none | `200`, paginated `code:0` | `401` |
| App products | GET | `/app/v1/products` | API key | none | `200`, paginated `code:0` | `401` |
| App events | GET | `/app/v1/events` | API key | none | `200`, paginated `code:0` | `401` |
| App articles | GET | `/app/v1/articles` | API key | none | `200`, paginated `code:0` | `401` |
| App search | GET | `/app/v1/search-options?q=coffee` | API key | none | `200`, `code:0` | `401` |
| Promotions | GET | `/app/v1/promotions` | API key | none | `200`, `code:0` | `401` |
| Profile | GET | `/api/v1/profile` | Bearer JWT | none | `200`, `code:0` | `401` |
| Venues list | GET | `/api/v1/venues` | Bearer JWT | none | `200`, paginated `code:0` | `401` |
| Venue detail | GET | `/api/v1/venues/:id` | Bearer JWT | none | `200`, `code:0` | `401`, `403`, `404` |
| Venue locations | GET | `/api/v1/venues/:id/locations` | Bearer JWT + venue role | none | `200`, paginated `code:0` | `401`, `403`, `404` |
| Venue products | GET | `/api/v1/venues/:id/products` | Bearer JWT + venue role | none | `200`, paginated `code:0` | `401`, `403`, `404` |
| Analytics searches | GET | `/api/v1/venues/:id/analytics/searches` | Bearer JWT + venue role | none | `200`, paginated `code:0` | `401`, `403`, `404` |
| Invalid login | POST | `/api/v1/auth/login` | none | `login_invalid` | n/a | `400`, `401`, non-zero `code` |

Optional write smoke coverage is disabled by default. Set `ALLOW_WRITES=true` to include `POST /api/v1/venues/:id/categories` and `POST /api/v1/venues/:id/locations` with payloads from `test-data.json`.

## Running Tests

Copy an env template and fill non-production values:

```bash
cp load-tests/env/local.env.example load-tests/env/local.env
cp load-tests/env/staging.env.example load-tests/env/staging.env
```

Run local baseline:

```bash
make load-test-local
```

Run staging baseline:

```bash
make load-test-staging
```

Run one scenario:

```bash
make load-test-scenario ENV=staging SCENARIO=app PROFILE=peak
```

Run the full suite:

```bash
make load-test-full ENV=staging
```

Without Make, run k6 directly:

```bash
k6 run --env TARGET_BASE_URL=http://localhost:8080 --env PROFILE=normal --env SCENARIO=full load-tests/scenarios/main.js
```

## Reports

Each run writes:

- `load-tests/reports/<profile>-<scenario>-summary.json`
- `load-tests/reports/<profile>-<scenario>-summary.html`

The summary includes request count, requests per second, average latency, P95/P99 latency, error rate, failed requests, and k6 tag data for endpoint-level breakdown. Use `summary.json` for CI parsing and `summary.html` for review.

## Thresholds

Default thresholds:

- Overall failed request rate: `< 5%`
- Average latency: `< 500ms`
- P95 latency: `< 1000ms`
- P99 latency: `< 2000ms`
- Public P95: `< 800ms`
- App P95: `< 1000ms`
- Admin P95: `< 1200ms`

Adjust thresholds in `scenarios/main.js` only after confirming the target environment capacity.

## Adding Scenarios

1. Add reusable payloads or IDs to `data/test-data.json`.
2. Add a flow function in `scenarios/main.js`.
3. Tag each request with `flow` and `endpoint`.
4. Add the flow to one or more profiles in `options`.
5. Document the endpoint, auth, payload, and expected responses in this file.

CREATE TABLE app_versions (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    venue_id   UUID NOT NULL UNIQUE REFERENCES venues(id) ON DELETE CASCADE,
    version    UUID NOT NULL DEFAULT gen_random_uuid(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

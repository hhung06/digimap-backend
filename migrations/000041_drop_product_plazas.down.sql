CREATE TABLE IF NOT EXISTS product_plazas (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    venue_id    UUID        NOT NULL REFERENCES venues(id) ON DELETE CASCADE,
    name        TEXT        NOT NULL,
    description TEXT        NOT NULL DEFAULT '',
    location_id UUID        REFERENCES locations(id) ON DELETE SET NULL,
    localization JSONB,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_product_plazas_venue ON product_plazas (venue_id) WHERE deleted_at IS NULL;

SELECT create_updated_at_trigger('product_plazas');

-- Product categories (venue-scoped)
CREATE TABLE product_categories (
    id           UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    venue_id     UUID        NOT NULL REFERENCES venues (id) ON DELETE CASCADE,
    external_id  TEXT,
    name         TEXT        NOT NULL,
    source       TEXT        NOT NULL DEFAULT 'internal',
    localization JSONB,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at   TIMESTAMPTZ
);

CREATE INDEX idx_product_categories_venue_id ON product_categories (venue_id) WHERE deleted_at IS NULL;

CREATE TRIGGER set_product_categories_updated_at
    BEFORE UPDATE ON product_categories
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Products
CREATE TABLE products (
    id               UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    venue_id         UUID        REFERENCES venues (id) ON DELETE CASCADE,
    location_id      UUID        REFERENCES locations (id) ON DELETE SET NULL,
    main_category_id UUID        REFERENCES product_categories (id) ON DELETE SET NULL,
    image            TEXT,
    name             TEXT,
    code             TEXT,
    size             TEXT,
    price            TEXT,
    origin_country   TEXT,
    expiration       TEXT,
    description      TEXT,
    custom           JSONB,
    localization     JSONB,
    source           TEXT        NOT NULL DEFAULT 'internal',
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at       TIMESTAMPTZ
);

CREATE INDEX idx_products_venue_id    ON products (venue_id)    WHERE deleted_at IS NULL;
CREATE INDEX idx_products_location_id ON products (location_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_products_name_trgm   ON products USING GIN (name gin_trgm_ops) WHERE deleted_at IS NULL;

CREATE TRIGGER set_products_updated_at
    BEFORE UPDATE ON products
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Product ↔ category M2M
CREATE TABLE product_category_links (
    product_id  UUID NOT NULL REFERENCES products (id) ON DELETE CASCADE,
    category_id UUID NOT NULL REFERENCES product_categories (id) ON DELETE CASCADE,
    PRIMARY KEY (product_id, category_id)
);

CREATE INDEX idx_product_category_links_category_id ON product_category_links (category_id);

-- Product attachments
CREATE TABLE product_attachments (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id  UUID        NOT NULL REFERENCES products (id) ON DELETE CASCADE,
    title       TEXT,
    file_type   TEXT        NOT NULL DEFAULT 'document',
    file        TEXT,
    source_url  TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);

CREATE INDEX idx_product_attachments_product_id ON product_attachments (product_id) WHERE deleted_at IS NULL;

CREATE TRIGGER set_product_attachments_updated_at
    BEFORE UPDATE ON product_attachments
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

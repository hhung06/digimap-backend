CREATE TABLE customers (
    id          UUID        PRIMARY KEY DEFAULT uuidv7(),
    name        TEXT        NOT NULL,
    email       TEXT        UNIQUE NOT NULL,
    phone       TEXT,
    address     TEXT,
    image       TEXT,
    is_active   BOOLEAN     NOT NULL DEFAULT TRUE,
    metadata    JSONB,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ,
    url         TEXT,
    description TEXT
);

CREATE INDEX idx_customers_email     ON customers (email)      WHERE deleted_at IS NULL;
CREATE INDEX idx_customers_is_active ON customers (is_active)  WHERE deleted_at IS NULL;

CREATE TRIGGER set_customers_updated_at
    BEFORE UPDATE ON customers
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

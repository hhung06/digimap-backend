CREATE TABLE users (
    id              UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    email           TEXT        UNIQUE NOT NULL,
    password_hash   TEXT        NOT NULL,
    first_name      TEXT        NOT NULL DEFAULT '',
    last_name       TEXT        NOT NULL DEFAULT '',
    phone           TEXT,
    avatar_url      TEXT,
    is_active       BOOLEAN     NOT NULL DEFAULT TRUE,
    is_system_admin BOOLEAN     NOT NULL DEFAULT FALSE,
    last_login_at   TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

CREATE INDEX idx_users_email          ON users (email)          WHERE deleted_at IS NULL;
CREATE INDEX idx_users_is_active      ON users (is_active)      WHERE deleted_at IS NULL;
CREATE INDEX idx_users_is_system_admin ON users (is_system_admin) WHERE deleted_at IS NULL AND is_system_admin = TRUE;

CREATE TRIGGER set_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

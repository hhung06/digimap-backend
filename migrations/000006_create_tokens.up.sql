-- ── Refresh tokens ────────────────────────────────────────────────────────────
-- The raw token is issued to the client once; only its SHA-256 hash is stored.

CREATE TABLE refresh_tokens (
    id          UUID        PRIMARY KEY DEFAULT uuidv7(),
    user_id     UUID        NOT NULL REFERENCES users (id),
    token_hash  TEXT        UNIQUE NOT NULL,
    expires_at  TIMESTAMPTZ NOT NULL,
    revoked_at  TIMESTAMPTZ,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_refresh_tokens_user_id    ON refresh_tokens (user_id)    WHERE revoked_at IS NULL;
CREATE INDEX idx_refresh_tokens_token_hash ON refresh_tokens (token_hash) WHERE revoked_at IS NULL;

-- ── Password reset tokens ─────────────────────────────────────────────────────

CREATE TABLE reset_password_tokens (
    id          UUID        PRIMARY KEY DEFAULT uuidv7(),
    user_id     UUID        NOT NULL REFERENCES users (id),
    token_hash  TEXT        UNIQUE NOT NULL,
    expires_at  TIMESTAMPTZ NOT NULL,
    used_at     TIMESTAMPTZ,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_reset_tokens_token_hash ON reset_password_tokens (token_hash) WHERE used_at IS NULL;

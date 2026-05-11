CREATE TYPE venue_role AS ENUM ('owner', 'editor', 'viewer');

CREATE TABLE venue_user_roles (
    id          UUID        PRIMARY KEY DEFAULT uuidv7(),
    venue_id    UUID        NOT NULL REFERENCES venues (id),
    user_id     UUID        NOT NULL REFERENCES users (id),
    role        venue_role  NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ,
    UNIQUE (venue_id, user_id)
);

CREATE INDEX idx_venue_user_roles_venue_id ON venue_user_roles (venue_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_venue_user_roles_user_id  ON venue_user_roles (user_id)  WHERE deleted_at IS NULL;

CREATE TRIGGER set_venue_user_roles_updated_at
    BEFORE UPDATE ON venue_user_roles
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ── Invitations ───────────────────────────────────────────────────────────────

CREATE TYPE invitation_status AS ENUM ('pending', 'accepted', 'cancelled');

CREATE TABLE venue_invitations (
    id           UUID              PRIMARY KEY DEFAULT uuidv7(),
    venue_id     UUID              NOT NULL REFERENCES venues (id),
    email        TEXT              NOT NULL,
    role         venue_role        NOT NULL,
    token        TEXT              UNIQUE NOT NULL,
    invited_by   UUID              NOT NULL REFERENCES users (id),
    status       invitation_status NOT NULL DEFAULT 'pending',
    accepted_at  TIMESTAMPTZ,
    cancelled_at TIMESTAMPTZ,
    expires_at   TIMESTAMPTZ       NOT NULL,
    created_at   TIMESTAMPTZ       NOT NULL DEFAULT NOW(),
    deleted_at   TIMESTAMPTZ
);

CREATE INDEX idx_venue_invitations_venue_id ON venue_invitations (venue_id)   WHERE deleted_at IS NULL;
CREATE INDEX idx_venue_invitations_email    ON venue_invitations (email)       WHERE deleted_at IS NULL;
CREATE INDEX idx_venue_invitations_token    ON venue_invitations (token)       WHERE deleted_at IS NULL;

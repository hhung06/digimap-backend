-- Script 04: Transform customers, users, venue_user_roles, venue_invitations

-- ── customers ────────────────────────────────────────────────────────────────
ALTER TABLE customers RENAME COLUMN image TO logo_url;
ALTER TABLE customers DROP COLUMN IF EXISTS url;
ALTER TABLE customers DROP COLUMN IF EXISTS description;
ALTER TABLE customers DROP COLUMN IF EXISTS restored_at;
ALTER TABLE customers DROP COLUMN IF EXISTS transaction_id;
ALTER TABLE customers DROP COLUMN IF EXISTS address;
ALTER TABLE customers ADD COLUMN IF NOT EXISTS is_active BOOLEAN NOT NULL DEFAULT TRUE;
ALTER TABLE customers ADD COLUMN IF NOT EXISTS metadata JSONB;

-- ── users (was auth_user) ────────────────────────────────────────────────────
-- Build UUID mapping: old integer id → new UUID
CREATE TABLE IF NOT EXISTS _user_id_map (
    old_id BIGINT PRIMARY KEY,
    new_id UUID NOT NULL DEFAULT gen_random_uuid()
);
INSERT INTO _user_id_map (old_id)
SELECT id FROM users
ON CONFLICT DO NOTHING;

-- Add new UUID column alongside old integer id
ALTER TABLE users ADD COLUMN IF NOT EXISTS new_id UUID;
UPDATE users u SET new_id = m.new_id FROM _user_id_map m WHERE m.old_id = u.id;

-- Rename / add / drop columns
ALTER TABLE users RENAME COLUMN password    TO password_hash;
ALTER TABLE users RENAME COLUMN last_login  TO last_login_at;
ALTER TABLE users RENAME COLUMN is_superuser TO is_system_admin;
ALTER TABLE users ADD COLUMN IF NOT EXISTS avatar_url  TEXT;
ALTER TABLE users ADD COLUMN IF NOT EXISTS phone       TEXT;
ALTER TABLE users ADD COLUMN IF NOT EXISTS deleted_at  TIMESTAMPTZ;
ALTER TABLE users ADD COLUMN IF NOT EXISTS created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW();
ALTER TABLE users ADD COLUMN IF NOT EXISTS updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW();
ALTER TABLE users DROP COLUMN IF EXISTS is_staff;
ALTER TABLE users DROP COLUMN IF EXISTS username;
ALTER TABLE users DROP COLUMN IF EXISTS date_joined;

-- Swap integer PK for UUID
ALTER TABLE users DROP CONSTRAINT IF EXISTS idx_21110_primary;
ALTER TABLE users DROP COLUMN IF EXISTS id;
ALTER TABLE users RENAME COLUMN new_id TO id;
ALTER TABLE users ADD PRIMARY KEY (id);
ALTER TABLE users ALTER COLUMN email TYPE TEXT;
ALTER TABLE users ADD CONSTRAINT users_email_unique UNIQUE (email);

-- ── venue_user_roles ─────────────────────────────────────────────────────────
-- NOTE: role column is named role_type (varchar(20)) in the source DB.
-- Values are already strings ('owner', 'editor', 'viewer'), so we just rename the column.
ALTER TABLE venue_user_roles RENAME COLUMN role_type TO role;
ALTER TABLE venue_user_roles DROP COLUMN IF EXISTS restored_at;
ALTER TABLE venue_user_roles DROP COLUMN IF EXISTS transaction_id;
-- NOTE: profile_id is an integer FK to indoormap_api_profile (not a UUID).
-- UUID remapping for profile → user will be handled after appuser/profile mapping is established.

-- ── venue_invitations ────────────────────────────────────────────────────────
-- NOTE: role and status columns are already varchar(20) strings in the source DB.
-- No type conversion needed; just clean up housekeeping columns.
ALTER TABLE venue_invitations DROP COLUMN IF EXISTS restored_at;
ALTER TABLE venue_invitations DROP COLUMN IF EXISTS transaction_id;
-- NOTE: invited_by_id and invited_profile_id are integer FKs to indoormap_api_profile.
-- UUID remapping will be handled after appuser/profile mapping is established.

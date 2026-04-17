-- Script 04: Transform customers, users, venue_user_roles, venue_invitations

-- ── customers ────────────────────────────────────────────────────────────────
DO $$ BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns
             WHERE table_name='customers' AND column_name='image') THEN
    ALTER TABLE customers RENAME COLUMN image TO logo_url;
  END IF;
END $$;
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
DO $$ BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns
             WHERE table_name='users' AND column_name='password') THEN
    ALTER TABLE users RENAME COLUMN password TO password_hash;
  END IF;
END $$;
DO $$ BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns
             WHERE table_name='users' AND column_name='last_login') THEN
    ALTER TABLE users RENAME COLUMN last_login TO last_login_at;
  END IF;
END $$;
DO $$ BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns
             WHERE table_name='users' AND column_name='is_superuser') THEN
    ALTER TABLE users RENAME COLUMN is_superuser TO is_system_admin;
  END IF;
END $$;
ALTER TABLE users ADD COLUMN IF NOT EXISTS avatar_url  TEXT;
ALTER TABLE users ADD COLUMN IF NOT EXISTS phone       TEXT;
ALTER TABLE users ADD COLUMN IF NOT EXISTS deleted_at  TIMESTAMPTZ;
ALTER TABLE users ADD COLUMN IF NOT EXISTS created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW();
ALTER TABLE users ADD COLUMN IF NOT EXISTS updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW();
ALTER TABLE users DROP COLUMN IF EXISTS is_staff;
ALTER TABLE users DROP COLUMN IF EXISTS username;
ALTER TABLE users DROP COLUMN IF EXISTS date_joined;

-- Swap integer PK for UUID
DO $$
DECLARE cname TEXT;
BEGIN
  SELECT constraint_name INTO cname
  FROM information_schema.table_constraints
  WHERE table_name = 'users' AND constraint_type = 'PRIMARY KEY';
  IF cname IS NOT NULL THEN
    EXECUTE format('ALTER TABLE users DROP CONSTRAINT %I', cname);
  END IF;
END$$;
ALTER TABLE users DROP COLUMN IF EXISTS id;
DO $$ BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns
             WHERE table_name='users' AND column_name='new_id') THEN
    ALTER TABLE users RENAME COLUMN new_id TO id;
  END IF;
END $$;
ALTER TABLE users ADD PRIMARY KEY (id);
ALTER TABLE users ALTER COLUMN email TYPE TEXT;
ALTER TABLE users ADD CONSTRAINT users_email_unique UNIQUE (email);

-- ── venue_user_roles ─────────────────────────────────────────────────────────
-- NOTE: role column is named role_type (varchar(20)) in the source DB.
-- Values are already strings ('owner', 'editor', 'viewer'), so we just rename the column.
DO $$ BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns
             WHERE table_name='venue_user_roles' AND column_name='role_type') THEN
    ALTER TABLE venue_user_roles RENAME COLUMN role_type TO role;
  END IF;
END $$;
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

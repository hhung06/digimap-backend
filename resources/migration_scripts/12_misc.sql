-- Script 12: Transform tags, search_queries, snapshots
SET search_path TO digimap_db, public;

-- ── tags ─────────────────────────────────────────────────────────────────────
-- Actual columns: name, localization, restored_at, transaction_id
ALTER TABLE tags DROP COLUMN IF EXISTS restored_at;
ALTER TABLE tags DROP COLUMN IF EXISTS transaction_id;

-- ── search_queries ───────────────────────────────────────────────────────────
-- Actual columns: search_term, search_count (bigint), last_searched, origin, app_id,
--                 restored_at, transaction_id, venue_id, is_promoted, reference, status
ALTER TABLE search_queries DROP COLUMN IF EXISTS restored_at;
ALTER TABLE search_queries DROP COLUMN IF EXISTS transaction_id;
ALTER TABLE search_queries ALTER COLUMN search_count TYPE INTEGER USING search_count::INTEGER;

-- ── snapshots ────────────────────────────────────────────────────────────────
-- Actual columns: state (varchar(6)), venue_id, publish_at, publish_by_id (integer FK to auth_user),
--                 method (integer), restored_at, transaction_id
-- state values in source: 'draft', 'publish', 'published'
-- Map to SMALLINT: draft→0, publish/published→1
ALTER TABLE snapshots ADD COLUMN IF NOT EXISTS state_new SMALLINT;
UPDATE snapshots SET state_new = CASE state
    WHEN 'draft'     THEN 0
    WHEN 'publish'   THEN 1
    WHEN 'published' THEN 1
    ELSE 0
END;
ALTER TABLE snapshots DROP COLUMN state;
DO $$ BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns
             WHERE table_name='snapshots' AND column_name='state_new') THEN
    ALTER TABLE snapshots RENAME COLUMN state_new TO state;
  END IF;
END $$;
ALTER TABLE snapshots ALTER COLUMN state SET NOT NULL;

-- Remap publish_by_id (old integer) → created_by (UUID) using _user_id_map
ALTER TABLE snapshots ADD COLUMN IF NOT EXISTS created_by UUID;
UPDATE snapshots s
SET created_by = m.new_id
FROM _user_id_map m
WHERE s.publish_by_id = m.old_id;
ALTER TABLE snapshots DROP COLUMN IF EXISTS publish_by_id;

ALTER TABLE snapshots DROP COLUMN IF EXISTS restored_at;
ALTER TABLE snapshots DROP COLUMN IF EXISTS transaction_id;

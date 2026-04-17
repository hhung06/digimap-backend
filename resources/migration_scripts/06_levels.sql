-- Script 06: Transform map_groups, perspectives, levels, geo_references

-- ── map_groups ───────────────────────────────────────────────────────────────
-- Actual columns: type, name, shortname, sortindex (bigint), venue_id, restored_at, transaction_id
ALTER TABLE map_groups RENAME COLUMN sortindex  TO sort_index;
ALTER TABLE map_groups RENAME COLUMN shortname  TO short_name;
ALTER TABLE map_groups DROP COLUMN IF EXISTS restored_at;
ALTER TABLE map_groups DROP COLUMN IF EXISTS transaction_id;
-- Cast sort_index from bigint to INTEGER
ALTER TABLE map_groups ALTER COLUMN sort_index TYPE INTEGER USING sort_index::INTEGER;

-- ── perspectives ─────────────────────────────────────────────────────────────
-- Actual columns (all lowercased by pgloader):
--   camerazoom, cameratype, cameramaxzoom, cameraminzoom,
--   cameratargetbearing, cameratargetcenterlat, cameratargetcenterlng,
--   cameratargetpitch, cameratargetzoom
ALTER TABLE perspectives RENAME COLUMN camerazoom            TO camera_zoom;
ALTER TABLE perspectives RENAME COLUMN cameratype            TO camera_type;
ALTER TABLE perspectives RENAME COLUMN cameramaxzoom         TO camera_max_zoom;
ALTER TABLE perspectives RENAME COLUMN cameraminzoom         TO camera_min_zoom;
ALTER TABLE perspectives RENAME COLUMN cameratargetbearing   TO camera_target_bearing;
ALTER TABLE perspectives RENAME COLUMN cameratargetcenterlat TO camera_target_center_lat;
ALTER TABLE perspectives RENAME COLUMN cameratargetcenterlng TO camera_target_center_lng;
ALTER TABLE perspectives RENAME COLUMN cameratargetpitch     TO camera_target_pitch;
ALTER TABLE perspectives RENAME COLUMN cameratargetzoom      TO camera_target_zoom;
ALTER TABLE perspectives DROP COLUMN IF EXISTS restored_at;
ALTER TABLE perspectives DROP COLUMN IF EXISTS transaction_id;

-- ── levels ───────────────────────────────────────────────────────────────────
-- Actual columns: externalid, shortname, publish (boolean), type_id, restored_at, transaction_id
ALTER TABLE levels RENAME COLUMN externalid TO external_id;
ALTER TABLE levels RENAME COLUMN shortname  TO short_name;
ALTER TABLE levels RENAME COLUMN publish    TO is_published;
ALTER TABLE levels DROP COLUMN IF EXISTS type_id;
ALTER TABLE levels DROP COLUMN IF EXISTS restored_at;
ALTER TABLE levels DROP COLUMN IF EXISTS transaction_id;

-- ── geo_references ───────────────────────────────────────────────────────────
-- Actual columns: controlx (integer), controly (integer), targetx (double precision), targety (double precision)
ALTER TABLE geo_references RENAME COLUMN controlx TO control_x;
ALTER TABLE geo_references RENAME COLUMN controly TO control_y;
ALTER TABLE geo_references RENAME COLUMN targetx  TO target_x;
ALTER TABLE geo_references RENAME COLUMN targety  TO target_y;
ALTER TABLE geo_references DROP COLUMN IF EXISTS restored_at;
ALTER TABLE geo_references DROP COLUMN IF EXISTS transaction_id;

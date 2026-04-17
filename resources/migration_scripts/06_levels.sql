-- Script 06: Transform map_groups, perspectives, levels, geo_references

-- ── map_groups ───────────────────────────────────────────────────────────────
-- Actual columns: type, name, shortname, sortindex (bigint), venue_id, restored_at, transaction_id
DO $$ BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns
             WHERE table_name='map_groups' AND column_name='sortindex') THEN
    ALTER TABLE map_groups RENAME COLUMN sortindex TO sort_index;
  END IF;
END $$;
DO $$ BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns
             WHERE table_name='map_groups' AND column_name='shortname') THEN
    ALTER TABLE map_groups RENAME COLUMN shortname TO short_name;
  END IF;
END $$;
ALTER TABLE map_groups DROP COLUMN IF EXISTS restored_at;
ALTER TABLE map_groups DROP COLUMN IF EXISTS transaction_id;
-- Cast sort_index from bigint to INTEGER
ALTER TABLE map_groups ALTER COLUMN sort_index TYPE INTEGER USING sort_index::INTEGER;

-- ── perspectives ─────────────────────────────────────────────────────────────
-- Actual columns (all lowercased by pgloader):
--   camerazoom, cameratype, cameramaxzoom, cameraminzoom,
--   cameratargetbearing, cameratargetcenterlat, cameratargetcenterlng,
--   cameratargetpitch, cameratargetzoom
DO $$ BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns
             WHERE table_name='perspectives' AND column_name='camerazoom') THEN
    ALTER TABLE perspectives RENAME COLUMN camerazoom TO camera_zoom;
  END IF;
END $$;
DO $$ BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns
             WHERE table_name='perspectives' AND column_name='cameratype') THEN
    ALTER TABLE perspectives RENAME COLUMN cameratype TO camera_type;
  END IF;
END $$;
DO $$ BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns
             WHERE table_name='perspectives' AND column_name='cameramaxzoom') THEN
    ALTER TABLE perspectives RENAME COLUMN cameramaxzoom TO camera_max_zoom;
  END IF;
END $$;
DO $$ BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns
             WHERE table_name='perspectives' AND column_name='cameraminzoom') THEN
    ALTER TABLE perspectives RENAME COLUMN cameraminzoom TO camera_min_zoom;
  END IF;
END $$;
DO $$ BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns
             WHERE table_name='perspectives' AND column_name='cameratargetbearing') THEN
    ALTER TABLE perspectives RENAME COLUMN cameratargetbearing TO camera_target_bearing;
  END IF;
END $$;
DO $$ BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns
             WHERE table_name='perspectives' AND column_name='cameratargetcenterlat') THEN
    ALTER TABLE perspectives RENAME COLUMN cameratargetcenterlat TO camera_target_center_lat;
  END IF;
END $$;
DO $$ BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns
             WHERE table_name='perspectives' AND column_name='cameratargetcenterlng') THEN
    ALTER TABLE perspectives RENAME COLUMN cameratargetcenterlng TO camera_target_center_lng;
  END IF;
END $$;
DO $$ BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns
             WHERE table_name='perspectives' AND column_name='cameratargetpitch') THEN
    ALTER TABLE perspectives RENAME COLUMN cameratargetpitch TO camera_target_pitch;
  END IF;
END $$;
DO $$ BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns
             WHERE table_name='perspectives' AND column_name='cameratargetzoom') THEN
    ALTER TABLE perspectives RENAME COLUMN cameratargetzoom TO camera_target_zoom;
  END IF;
END $$;
ALTER TABLE perspectives DROP COLUMN IF EXISTS restored_at;
ALTER TABLE perspectives DROP COLUMN IF EXISTS transaction_id;

-- ── levels ───────────────────────────────────────────────────────────────────
-- Actual columns: externalid, shortname, publish (boolean), type_id, restored_at, transaction_id
DO $$ BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns
             WHERE table_name='levels' AND column_name='externalid') THEN
    ALTER TABLE levels RENAME COLUMN externalid TO external_id;
  END IF;
END $$;
DO $$ BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns
             WHERE table_name='levels' AND column_name='shortname') THEN
    ALTER TABLE levels RENAME COLUMN shortname TO short_name;
  END IF;
END $$;
DO $$ BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns
             WHERE table_name='levels' AND column_name='publish') THEN
    ALTER TABLE levels RENAME COLUMN publish TO is_published;
  END IF;
END $$;
ALTER TABLE levels DROP COLUMN IF EXISTS type_id;
ALTER TABLE levels DROP COLUMN IF EXISTS restored_at;
ALTER TABLE levels DROP COLUMN IF EXISTS transaction_id;

-- ── geo_references ───────────────────────────────────────────────────────────
-- Actual columns: controlx (integer), controly (integer), targetx (double precision), targety (double precision)
DO $$ BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns
             WHERE table_name='geo_references' AND column_name='controlx') THEN
    ALTER TABLE geo_references RENAME COLUMN controlx TO control_x;
  END IF;
END $$;
DO $$ BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns
             WHERE table_name='geo_references' AND column_name='controly') THEN
    ALTER TABLE geo_references RENAME COLUMN controly TO control_y;
  END IF;
END $$;
DO $$ BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns
             WHERE table_name='geo_references' AND column_name='targetx') THEN
    ALTER TABLE geo_references RENAME COLUMN targetx TO target_x;
  END IF;
END $$;
DO $$ BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns
             WHERE table_name='geo_references' AND column_name='targety') THEN
    ALTER TABLE geo_references RENAME COLUMN targety TO target_y;
  END IF;
END $$;
ALTER TABLE geo_references DROP COLUMN IF EXISTS restored_at;
ALTER TABLE geo_references DROP COLUMN IF EXISTS transaction_id;

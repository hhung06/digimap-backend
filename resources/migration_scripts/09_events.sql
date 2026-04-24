-- Script 09: Transform event_types, event_tags, events, event_tag_links, event_location_links, event_images
SET search_path TO digimap_db, public;

-- ── event_types ──────────────────────────────────────────────────────────────
-- Actual columns: name, venue_id, localization, restored_at, transaction_id
ALTER TABLE event_types DROP COLUMN IF EXISTS restored_at;
ALTER TABLE event_types DROP COLUMN IF EXISTS transaction_id;

-- ── event_tags ───────────────────────────────────────────────────────────────
-- Actual columns: name, localization, restored_at, transaction_id
ALTER TABLE event_tags DROP COLUMN IF EXISTS restored_at;
ALTER TABLE event_tags DROP COLUMN IF EXISTS transaction_id;

-- ── events ───────────────────────────────────────────────────────────────────
-- Actual columns: bannerimage, title, starttime, endtime, type_id, contentdetail,
--                 contenturl, description, iconimage, showendtime, showstarttime,
--                 venue_id, localization, restored_at, transaction_id
DO $$ BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns
             WHERE table_name='events' AND column_name='bannerimage') THEN
    ALTER TABLE events RENAME COLUMN bannerimage TO banner_image;
  END IF;
END $$;
DO $$ BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns
             WHERE table_name='events' AND column_name='starttime') THEN
    ALTER TABLE events RENAME COLUMN starttime TO start_time;
  END IF;
END $$;
DO $$ BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns
             WHERE table_name='events' AND column_name='endtime') THEN
    ALTER TABLE events RENAME COLUMN endtime TO end_time;
  END IF;
END $$;
DO $$ BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns
             WHERE table_name='events' AND column_name='contentdetail') THEN
    ALTER TABLE events RENAME COLUMN contentdetail TO content_detail;
  END IF;
END $$;
DO $$ BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns
             WHERE table_name='events' AND column_name='contenturl') THEN
    ALTER TABLE events RENAME COLUMN contenturl TO content_url;
  END IF;
END $$;
DO $$ BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns
             WHERE table_name='events' AND column_name='iconimage') THEN
    ALTER TABLE events RENAME COLUMN iconimage TO icon_image;
  END IF;
END $$;
DO $$ BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns
             WHERE table_name='events' AND column_name='showendtime') THEN
    ALTER TABLE events RENAME COLUMN showendtime TO show_end_time;
  END IF;
END $$;
DO $$ BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns
             WHERE table_name='events' AND column_name='showstarttime') THEN
    ALTER TABLE events RENAME COLUMN showstarttime TO show_start_time;
  END IF;
END $$;
ALTER TABLE events DROP COLUMN IF EXISTS restored_at;
ALTER TABLE events DROP COLUMN IF EXISTS transaction_id;

-- ── event_tag_links ──────────────────────────────────────────────────────────
-- Actual columns: id (bigserial PK), event_id, eventtag_id
DO $$ BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns
             WHERE table_name='event_tag_links' AND column_name='eventtag_id') THEN
    ALTER TABLE event_tag_links RENAME COLUMN eventtag_id TO tag_id;
  END IF;
END $$;
ALTER TABLE event_tag_links DROP COLUMN IF EXISTS id;
ALTER TABLE event_tag_links ADD PRIMARY KEY (event_id, tag_id);

-- ── event_location_links ─────────────────────────────────────────────────────
-- Actual columns: id (bigserial PK), event_id, location_id
ALTER TABLE event_location_links DROP COLUMN IF EXISTS id;
ALTER TABLE event_location_links ADD PRIMARY KEY (event_id, location_id);

-- ── event_images ─────────────────────────────────────────────────────────────
-- Actual columns: image, event_id, restored_at, transaction_id
ALTER TABLE event_images DROP COLUMN IF EXISTS restored_at;
ALTER TABLE event_images DROP COLUMN IF EXISTS transaction_id;

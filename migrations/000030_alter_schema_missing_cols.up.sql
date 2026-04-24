-- locations: add columns missing from the base DB schema
ALTER TABLE locations
    ADD COLUMN IF NOT EXISTS top_logo      TEXT,
    ADD COLUMN IF NOT EXISTS top_logo_type TEXT,
    ADD COLUMN IF NOT EXISTS start_time    TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS end_time      TIMESTAMPTZ;

-- levels: add type column missing from the base DB schema
ALTER TABLE levels
    ADD COLUMN IF NOT EXISTS type SMALLINT NOT NULL DEFAULT 0;

-- location_categories: rename shortname → short_name to match the codebase
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_schema = 'digimap_db'
          AND table_name   = 'location_categories'
          AND column_name  = 'shortname'
    ) THEN
        ALTER TABLE location_categories RENAME COLUMN shortname TO short_name;
    END IF;
END $$;

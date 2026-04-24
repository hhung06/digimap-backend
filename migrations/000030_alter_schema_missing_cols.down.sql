ALTER TABLE locations
    DROP COLUMN IF EXISTS top_logo,
    DROP COLUMN IF EXISTS top_logo_type,
    DROP COLUMN IF EXISTS start_time,
    DROP COLUMN IF EXISTS end_time;

ALTER TABLE levels
    DROP COLUMN IF EXISTS type;

DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_schema = 'digimap_db'
          AND table_name   = 'location_categories'
          AND column_name  = 'short_name'
    ) THEN
        ALTER TABLE location_categories RENAME COLUMN short_name TO shortname;
    END IF;
END $$;

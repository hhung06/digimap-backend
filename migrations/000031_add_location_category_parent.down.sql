DROP INDEX IF EXISTS idx_location_categories_parent_id;

ALTER TABLE location_categories
    DROP CONSTRAINT IF EXISTS chk_location_categories_not_self_parent;

ALTER TABLE location_categories
    DROP COLUMN IF EXISTS parent_id;

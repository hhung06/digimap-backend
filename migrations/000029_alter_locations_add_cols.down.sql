ALTER TABLE locations
    DROP COLUMN IF EXISTS common_location_sub_type,
    DROP COLUMN IF EXISTS common_show_short_name;

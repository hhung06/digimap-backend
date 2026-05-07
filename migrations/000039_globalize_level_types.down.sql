UPDATE level_types
SET venue_id = (
    SELECT id
    FROM venues
    WHERE deleted_at IS NULL
    ORDER BY created_at
    LIMIT 1
)
WHERE venue_id IS NULL;

ALTER TABLE level_types
    ALTER COLUMN venue_id SET NOT NULL;

-- Restore NOT NULL only if all rows have a venue_id (may fail if global assets exist).
ALTER TABLE assets ALTER COLUMN venue_id SET NOT NULL;

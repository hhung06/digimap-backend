-- Languages previously had no way to be disabled independently of soft-delete
-- or is_default. Add a dedicated enabled flag so a venue can toggle a
-- language off without losing its row.
ALTER TABLE languages ADD COLUMN enabled BOOLEAN NOT NULL DEFAULT true;

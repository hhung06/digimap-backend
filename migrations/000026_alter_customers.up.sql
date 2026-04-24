ALTER TABLE customers
    RENAME COLUMN logo_url TO image;

ALTER TABLE customers
    ADD COLUMN IF NOT EXISTS address     TEXT,
    ADD COLUMN IF NOT EXISTS url         TEXT,
    ADD COLUMN IF NOT EXISTS description TEXT;

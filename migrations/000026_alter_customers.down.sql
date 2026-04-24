ALTER TABLE customers
    DROP COLUMN IF EXISTS description,
    DROP COLUMN IF EXISTS url,
    DROP COLUMN IF EXISTS address;

ALTER TABLE customers
    RENAME COLUMN image TO logo_url;

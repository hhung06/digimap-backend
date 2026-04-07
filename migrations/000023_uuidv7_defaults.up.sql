-- Migration: replace gen_random_uuid() (UUIDv4) defaults with uuidv7()
-- Rationale: the Go layer already generates UUIDv7 via uuid.NewV7(); this keeps
-- the database default consistent so raw SQL / direct inserts also produce v7.
-- UUIDv7 is time-ordered (48-bit ms timestamp prefix), which reduces B-tree
-- index fragmentation compared to random UUIDv4.
--
-- Requires: PostgreSQL 18+ (native uuidv7() built-in, no extension needed).

-- ─── Swap DEFAULT on every table ───────────────────────────────────────────
ALTER TABLE customers            ALTER COLUMN id SET DEFAULT uuidv7();
ALTER TABLE venues               ALTER COLUMN id SET DEFAULT uuidv7();
ALTER TABLE users                ALTER COLUMN id SET DEFAULT uuidv7();
ALTER TABLE venue_user_roles     ALTER COLUMN id SET DEFAULT uuidv7();
ALTER TABLE venue_invitations    ALTER COLUMN id SET DEFAULT uuidv7();
ALTER TABLE refresh_tokens       ALTER COLUMN id SET DEFAULT uuidv7();
ALTER TABLE reset_password_tokens ALTER COLUMN id SET DEFAULT uuidv7();
ALTER TABLE map_groups           ALTER COLUMN id SET DEFAULT uuidv7();
ALTER TABLE perspectives         ALTER COLUMN id SET DEFAULT uuidv7();
ALTER TABLE levels               ALTER COLUMN id SET DEFAULT uuidv7();
ALTER TABLE geo_references       ALTER COLUMN id SET DEFAULT uuidv7();
ALTER TABLE location_categories  ALTER COLUMN id SET DEFAULT uuidv7();
ALTER TABLE amenities            ALTER COLUMN id SET DEFAULT uuidv7();
ALTER TABLE venue_amenities      ALTER COLUMN id SET DEFAULT uuidv7();
ALTER TABLE locations            ALTER COLUMN id SET DEFAULT uuidv7();
ALTER TABLE location_images      ALTER COLUMN id SET DEFAULT uuidv7();
ALTER TABLE promotions           ALTER COLUMN id SET DEFAULT uuidv7();
ALTER TABLE product_categories   ALTER COLUMN id SET DEFAULT uuidv7();
ALTER TABLE products             ALTER COLUMN id SET DEFAULT uuidv7();
ALTER TABLE product_attachments  ALTER COLUMN id SET DEFAULT uuidv7();
ALTER TABLE event_tags           ALTER COLUMN id SET DEFAULT uuidv7();
ALTER TABLE event_types          ALTER COLUMN id SET DEFAULT uuidv7();
ALTER TABLE events               ALTER COLUMN id SET DEFAULT uuidv7();
ALTER TABLE event_images         ALTER COLUMN id SET DEFAULT uuidv7();
ALTER TABLE notifications        ALTER COLUMN id SET DEFAULT uuidv7();
ALTER TABLE surveys              ALTER COLUMN id SET DEFAULT uuidv7();
ALTER TABLE questions            ALTER COLUMN id SET DEFAULT uuidv7();
ALTER TABLE options              ALTER COLUMN id SET DEFAULT uuidv7();
ALTER TABLE survey_responses     ALTER COLUMN id SET DEFAULT uuidv7();
ALTER TABLE survey_answers       ALTER COLUMN id SET DEFAULT uuidv7();
ALTER TABLE beacons              ALTER COLUMN id SET DEFAULT uuidv7();
ALTER TABLE connections          ALTER COLUMN id SET DEFAULT uuidv7();
ALTER TABLE connection_levels    ALTER COLUMN id SET DEFAULT uuidv7();
ALTER TABLE qrcodes              ALTER COLUMN id SET DEFAULT uuidv7();
ALTER TABLE advertisements       ALTER COLUMN id SET DEFAULT uuidv7();
ALTER TABLE articles             ALTER COLUMN id SET DEFAULT uuidv7();
ALTER TABLE article_images       ALTER COLUMN id SET DEFAULT uuidv7();
ALTER TABLE coupons              ALTER COLUMN id SET DEFAULT uuidv7();
ALTER TABLE videos               ALTER COLUMN id SET DEFAULT uuidv7();
ALTER TABLE tags                 ALTER COLUMN id SET DEFAULT uuidv7();
ALTER TABLE entity_tags          ALTER COLUMN id SET DEFAULT uuidv7();
ALTER TABLE event_logs           ALTER COLUMN id SET DEFAULT uuidv7();
ALTER TABLE search_queries       ALTER COLUMN id SET DEFAULT uuidv7();

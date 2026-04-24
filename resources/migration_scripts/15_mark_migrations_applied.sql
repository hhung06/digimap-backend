SET search_path TO digimap_db, public;

CREATE TABLE IF NOT EXISTS schema_migrations (
    version  BIGINT  NOT NULL PRIMARY KEY,
    dirty    BOOLEAN NOT NULL
);

-- golang-migrate expects exactly ONE row: the current (highest) applied version
TRUNCATE schema_migrations;
INSERT INTO schema_migrations (version, dirty) VALUES (25, false);

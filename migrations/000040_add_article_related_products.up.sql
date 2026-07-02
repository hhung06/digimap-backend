ALTER TABLE articles
    ADD COLUMN related_products JSONB NOT NULL DEFAULT '[]';

ALTER TABLE location_categories
    ADD COLUMN parent_id UUID REFERENCES location_categories (id) ON DELETE RESTRICT;

ALTER TABLE location_categories
    ADD CONSTRAINT chk_location_categories_not_self_parent
    CHECK (parent_id IS NULL OR parent_id <> id);

CREATE INDEX idx_location_categories_parent_id
    ON location_categories (parent_id)
    WHERE deleted_at IS NULL;

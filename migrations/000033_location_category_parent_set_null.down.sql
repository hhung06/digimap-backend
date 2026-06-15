ALTER TABLE location_categories
    DROP CONSTRAINT location_categories_parent_id_fkey;

ALTER TABLE location_categories
    ADD CONSTRAINT location_categories_parent_id_fkey
    FOREIGN KEY (parent_id) REFERENCES location_categories (id) ON DELETE RESTRICT;

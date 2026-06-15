-- Change parent_id FK from RESTRICT to SET NULL so that deleting a parent
-- category orphans its subcategories (sets parent_id = NULL) rather than
-- blocking the delete. Matches Django's ForeignKey(SET_NULL) behaviour.
ALTER TABLE location_categories
    DROP CONSTRAINT location_categories_parent_id_fkey;

ALTER TABLE location_categories
    ADD CONSTRAINT location_categories_parent_id_fkey
    FOREIGN KEY (parent_id) REFERENCES location_categories (id) ON DELETE SET NULL;

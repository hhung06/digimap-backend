-- =============================================================================
-- Script 13: Restore FK constraints, updated_at triggers, uuidv7() defaults
-- =============================================================================

-- =============================================================================
-- SECTION 1: Re-add FK constraints
-- =============================================================================

DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints
                 WHERE constraint_name = 'venues_customer_id_fk' AND table_schema = current_schema()) THEN
    ALTER TABLE venues ADD CONSTRAINT venues_customer_id_fk
      FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE SET NULL;
  END IF;
END $$;

-- venue_user_roles: venue_id only (user_id FK skipped — not yet remapped to UUID)
DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints
                 WHERE constraint_name = 'venue_user_roles_venue_id_fk' AND table_schema = current_schema()) THEN
    ALTER TABLE venue_user_roles ADD CONSTRAINT venue_user_roles_venue_id_fk
      FOREIGN KEY (venue_id) REFERENCES venues(id) ON DELETE CASCADE;
  END IF;
END $$;

-- venue_invitations: venue_id only (invited_by FK skipped — not yet remapped to UUID)
DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints
                 WHERE constraint_name = 'venue_invitations_venue_id_fk' AND table_schema = current_schema()) THEN
    ALTER TABLE venue_invitations ADD CONSTRAINT venue_invitations_venue_id_fk
      FOREIGN KEY (venue_id) REFERENCES venues(id) ON DELETE CASCADE;
  END IF;
END $$;

DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints
                 WHERE constraint_name = 'map_groups_venue_id_fk' AND table_schema = current_schema()) THEN
    ALTER TABLE map_groups ADD CONSTRAINT map_groups_venue_id_fk
      FOREIGN KEY (venue_id) REFERENCES venues(id) ON DELETE CASCADE;
  END IF;
END $$;

DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints
                 WHERE constraint_name = 'levels_venue_id_fk' AND table_schema = current_schema()) THEN
    ALTER TABLE levels ADD CONSTRAINT levels_venue_id_fk
      FOREIGN KEY (venue_id) REFERENCES venues(id) ON DELETE CASCADE;
  END IF;
END $$;

DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints
                 WHERE constraint_name = 'levels_map_group_id_fk' AND table_schema = current_schema()) THEN
    ALTER TABLE levels ADD CONSTRAINT levels_map_group_id_fk
      FOREIGN KEY (map_group_id) REFERENCES map_groups(id) ON DELETE SET NULL;
  END IF;
END $$;

DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints
                 WHERE constraint_name = 'levels_perspective_id_fk' AND table_schema = current_schema()) THEN
    ALTER TABLE levels ADD CONSTRAINT levels_perspective_id_fk
      FOREIGN KEY (perspective_id) REFERENCES perspectives(id) ON DELETE SET NULL;
  END IF;
END $$;

DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints
                 WHERE constraint_name = 'geo_references_level_id_fk' AND table_schema = current_schema()) THEN
    ALTER TABLE geo_references ADD CONSTRAINT geo_references_level_id_fk
      FOREIGN KEY (level_id) REFERENCES levels(id) ON DELETE CASCADE;
  END IF;
END $$;

DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints
                 WHERE constraint_name = 'location_categories_venue_id_fk' AND table_schema = current_schema()) THEN
    ALTER TABLE location_categories ADD CONSTRAINT location_categories_venue_id_fk
      FOREIGN KEY (venue_id) REFERENCES venues(id) ON DELETE CASCADE;
  END IF;
END $$;

DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints
                 WHERE constraint_name = 'locations_venue_id_fk' AND table_schema = current_schema()) THEN
    ALTER TABLE locations ADD CONSTRAINT locations_venue_id_fk
      FOREIGN KEY (venue_id) REFERENCES venues(id) ON DELETE CASCADE;
  END IF;
END $$;

DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints
                 WHERE constraint_name = 'locations_level_id_fk' AND table_schema = current_schema()) THEN
    ALTER TABLE locations ADD CONSTRAINT locations_level_id_fk
      FOREIGN KEY (level_id) REFERENCES levels(id) ON DELETE SET NULL;
  END IF;
END $$;

DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints
                 WHERE constraint_name = 'locations_main_category_id_fk' AND table_schema = current_schema()) THEN
    ALTER TABLE locations ADD CONSTRAINT locations_main_category_id_fk
      FOREIGN KEY (main_category_id) REFERENCES location_categories(id) ON DELETE SET NULL;
  END IF;
END $$;

DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints
                 WHERE constraint_name = 'location_category_links_location_id_fk' AND table_schema = current_schema()) THEN
    ALTER TABLE location_category_links ADD CONSTRAINT location_category_links_location_id_fk
      FOREIGN KEY (location_id) REFERENCES locations(id) ON DELETE CASCADE;
  END IF;
END $$;

DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints
                 WHERE constraint_name = 'location_category_links_category_id_fk' AND table_schema = current_schema()) THEN
    ALTER TABLE location_category_links ADD CONSTRAINT location_category_links_category_id_fk
      FOREIGN KEY (category_id) REFERENCES location_categories(id) ON DELETE CASCADE;
  END IF;
END $$;

DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints
                 WHERE constraint_name = 'location_images_location_id_fk' AND table_schema = current_schema()) THEN
    ALTER TABLE location_images ADD CONSTRAINT location_images_location_id_fk
      FOREIGN KEY (location_id) REFERENCES locations(id) ON DELETE CASCADE;
  END IF;
END $$;

DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints
                 WHERE constraint_name = 'product_categories_venue_id_fk' AND table_schema = current_schema()) THEN
    ALTER TABLE product_categories ADD CONSTRAINT product_categories_venue_id_fk
      FOREIGN KEY (venue_id) REFERENCES venues(id) ON DELETE CASCADE;
  END IF;
END $$;

DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints
                 WHERE constraint_name = 'products_venue_id_fk' AND table_schema = current_schema()) THEN
    ALTER TABLE products ADD CONSTRAINT products_venue_id_fk
      FOREIGN KEY (venue_id) REFERENCES venues(id) ON DELETE CASCADE;
  END IF;
END $$;

DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints
                 WHERE constraint_name = 'products_location_id_fk' AND table_schema = current_schema()) THEN
    ALTER TABLE products ADD CONSTRAINT products_location_id_fk
      FOREIGN KEY (location_id) REFERENCES locations(id) ON DELETE SET NULL;
  END IF;
END $$;

DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints
                 WHERE constraint_name = 'products_main_category_id_fk' AND table_schema = current_schema()) THEN
    ALTER TABLE products ADD CONSTRAINT products_main_category_id_fk
      FOREIGN KEY (main_category_id) REFERENCES product_categories(id) ON DELETE SET NULL;
  END IF;
END $$;

DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints
                 WHERE constraint_name = 'product_category_links_product_id_fk' AND table_schema = current_schema()) THEN
    ALTER TABLE product_category_links ADD CONSTRAINT product_category_links_product_id_fk
      FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE;
  END IF;
END $$;

DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints
                 WHERE constraint_name = 'product_category_links_category_id_fk' AND table_schema = current_schema()) THEN
    ALTER TABLE product_category_links ADD CONSTRAINT product_category_links_category_id_fk
      FOREIGN KEY (category_id) REFERENCES product_categories(id) ON DELETE CASCADE;
  END IF;
END $$;

DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints
                 WHERE constraint_name = 'product_attachments_product_id_fk' AND table_schema = current_schema()) THEN
    ALTER TABLE product_attachments ADD CONSTRAINT product_attachments_product_id_fk
      FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE;
  END IF;
END $$;

DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints
                 WHERE constraint_name = 'event_types_venue_id_fk' AND table_schema = current_schema()) THEN
    ALTER TABLE event_types ADD CONSTRAINT event_types_venue_id_fk
      FOREIGN KEY (venue_id) REFERENCES venues(id) ON DELETE CASCADE;
  END IF;
END $$;

DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints
                 WHERE constraint_name = 'events_venue_id_fk' AND table_schema = current_schema()) THEN
    ALTER TABLE events ADD CONSTRAINT events_venue_id_fk
      FOREIGN KEY (venue_id) REFERENCES venues(id) ON DELETE CASCADE;
  END IF;
END $$;

DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints
                 WHERE constraint_name = 'events_type_id_fk' AND table_schema = current_schema()) THEN
    ALTER TABLE events ADD CONSTRAINT events_type_id_fk
      FOREIGN KEY (type_id) REFERENCES event_types(id) ON DELETE SET NULL;
  END IF;
END $$;

DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints
                 WHERE constraint_name = 'event_tag_links_event_id_fk' AND table_schema = current_schema()) THEN
    ALTER TABLE event_tag_links ADD CONSTRAINT event_tag_links_event_id_fk
      FOREIGN KEY (event_id) REFERENCES events(id) ON DELETE CASCADE;
  END IF;
END $$;

DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints
                 WHERE constraint_name = 'event_tag_links_tag_id_fk' AND table_schema = current_schema()) THEN
    ALTER TABLE event_tag_links ADD CONSTRAINT event_tag_links_tag_id_fk
      FOREIGN KEY (tag_id) REFERENCES event_tags(id) ON DELETE CASCADE;
  END IF;
END $$;

DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints
                 WHERE constraint_name = 'event_location_links_event_id_fk' AND table_schema = current_schema()) THEN
    ALTER TABLE event_location_links ADD CONSTRAINT event_location_links_event_id_fk
      FOREIGN KEY (event_id) REFERENCES events(id) ON DELETE CASCADE;
  END IF;
END $$;

DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints
                 WHERE constraint_name = 'event_location_links_location_id_fk' AND table_schema = current_schema()) THEN
    ALTER TABLE event_location_links ADD CONSTRAINT event_location_links_location_id_fk
      FOREIGN KEY (location_id) REFERENCES locations(id) ON DELETE CASCADE;
  END IF;
END $$;

DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints
                 WHERE constraint_name = 'event_images_event_id_fk' AND table_schema = current_schema()) THEN
    ALTER TABLE event_images ADD CONSTRAINT event_images_event_id_fk
      FOREIGN KEY (event_id) REFERENCES events(id) ON DELETE CASCADE;
  END IF;
END $$;

DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints
                 WHERE constraint_name = 'surveys_venue_id_fk' AND table_schema = current_schema()) THEN
    ALTER TABLE surveys ADD CONSTRAINT surveys_venue_id_fk
      FOREIGN KEY (venue_id) REFERENCES venues(id) ON DELETE CASCADE;
  END IF;
END $$;

DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints
                 WHERE constraint_name = 'surveys_created_by_fk' AND table_schema = current_schema()) THEN
    ALTER TABLE surveys ADD CONSTRAINT surveys_created_by_fk
      FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL;
  END IF;
END $$;

DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints
                 WHERE constraint_name = 'questions_survey_id_fk' AND table_schema = current_schema()) THEN
    ALTER TABLE questions ADD CONSTRAINT questions_survey_id_fk
      FOREIGN KEY (survey_id) REFERENCES surveys(id) ON DELETE CASCADE;
  END IF;
END $$;

DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints
                 WHERE constraint_name = 'options_question_id_fk' AND table_schema = current_schema()) THEN
    ALTER TABLE options ADD CONSTRAINT options_question_id_fk
      FOREIGN KEY (question_id) REFERENCES questions(id) ON DELETE CASCADE;
  END IF;
END $$;

DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints
                 WHERE constraint_name = 'survey_responses_survey_id_fk' AND table_schema = current_schema()) THEN
    ALTER TABLE survey_responses ADD CONSTRAINT survey_responses_survey_id_fk
      FOREIGN KEY (survey_id) REFERENCES surveys(id) ON DELETE CASCADE;
  END IF;
END $$;

DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints
                 WHERE constraint_name = 'survey_answers_response_id_fk' AND table_schema = current_schema()) THEN
    ALTER TABLE survey_answers ADD CONSTRAINT survey_answers_response_id_fk
      FOREIGN KEY (response_id) REFERENCES survey_responses(id) ON DELETE CASCADE;
  END IF;
END $$;

DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints
                 WHERE constraint_name = 'survey_answers_question_id_fk' AND table_schema = current_schema()) THEN
    ALTER TABLE survey_answers ADD CONSTRAINT survey_answers_question_id_fk
      FOREIGN KEY (question_id) REFERENCES questions(id) ON DELETE CASCADE;
  END IF;
END $$;

DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints
                 WHERE constraint_name = 'survey_answers_option_id_fk' AND table_schema = current_schema()) THEN
    ALTER TABLE survey_answers ADD CONSTRAINT survey_answers_option_id_fk
      FOREIGN KEY (option_id) REFERENCES options(id) ON DELETE SET NULL;
  END IF;
END $$;

DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints
                 WHERE constraint_name = 'advertisements_venue_id_fk' AND table_schema = current_schema()) THEN
    ALTER TABLE advertisements ADD CONSTRAINT advertisements_venue_id_fk
      FOREIGN KEY (venue_id) REFERENCES venues(id) ON DELETE CASCADE;
  END IF;
END $$;

DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints
                 WHERE constraint_name = 'advertisements_location_id_fk' AND table_schema = current_schema()) THEN
    ALTER TABLE advertisements ADD CONSTRAINT advertisements_location_id_fk
      FOREIGN KEY (location_id) REFERENCES locations(id) ON DELETE SET NULL;
  END IF;
END $$;

DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints
                 WHERE constraint_name = 'articles_venue_id_fk' AND table_schema = current_schema()) THEN
    ALTER TABLE articles ADD CONSTRAINT articles_venue_id_fk
      FOREIGN KEY (venue_id) REFERENCES venues(id) ON DELETE CASCADE;
  END IF;
END $$;

DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints
                 WHERE constraint_name = 'articles_location_id_fk' AND table_schema = current_schema()) THEN
    ALTER TABLE articles ADD CONSTRAINT articles_location_id_fk
      FOREIGN KEY (location_id) REFERENCES locations(id) ON DELETE SET NULL;
  END IF;
END $$;

DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints
                 WHERE constraint_name = 'article_images_article_id_fk' AND table_schema = current_schema()) THEN
    ALTER TABLE article_images ADD CONSTRAINT article_images_article_id_fk
      FOREIGN KEY (article_id) REFERENCES articles(id) ON DELETE CASCADE;
  END IF;
END $$;

DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints
                 WHERE constraint_name = 'coupons_venue_id_fk' AND table_schema = current_schema()) THEN
    ALTER TABLE coupons ADD CONSTRAINT coupons_venue_id_fk
      FOREIGN KEY (venue_id) REFERENCES venues(id) ON DELETE CASCADE;
  END IF;
END $$;

DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints
                 WHERE constraint_name = 'beacons_venue_id_fk' AND table_schema = current_schema()) THEN
    ALTER TABLE beacons ADD CONSTRAINT beacons_venue_id_fk
      FOREIGN KEY (venue_id) REFERENCES venues(id) ON DELETE CASCADE;
  END IF;
END $$;

DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints
                 WHERE constraint_name = 'beacons_level_id_fk' AND table_schema = current_schema()) THEN
    ALTER TABLE beacons ADD CONSTRAINT beacons_level_id_fk
      FOREIGN KEY (level_id) REFERENCES levels(id) ON DELETE SET NULL;
  END IF;
END $$;

DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints
                 WHERE constraint_name = 'connections_venue_id_fk' AND table_schema = current_schema()) THEN
    ALTER TABLE connections ADD CONSTRAINT connections_venue_id_fk
      FOREIGN KEY (venue_id) REFERENCES venues(id) ON DELETE CASCADE;
  END IF;
END $$;

DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints
                 WHERE constraint_name = 'connection_levels_connection_id_fk' AND table_schema = current_schema()) THEN
    ALTER TABLE connection_levels ADD CONSTRAINT connection_levels_connection_id_fk
      FOREIGN KEY (connection_id) REFERENCES connections(id) ON DELETE CASCADE;
  END IF;
END $$;

DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints
                 WHERE constraint_name = 'connection_levels_level_id_fk' AND table_schema = current_schema()) THEN
    ALTER TABLE connection_levels ADD CONSTRAINT connection_levels_level_id_fk
      FOREIGN KEY (level_id) REFERENCES levels(id) ON DELETE CASCADE;
  END IF;
END $$;

DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints
                 WHERE constraint_name = 'snapshots_venue_id_fk' AND table_schema = current_schema()) THEN
    ALTER TABLE snapshots ADD CONSTRAINT snapshots_venue_id_fk
      FOREIGN KEY (venue_id) REFERENCES venues(id) ON DELETE CASCADE;
  END IF;
END $$;

DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints
                 WHERE constraint_name = 'snapshots_created_by_fk' AND table_schema = current_schema()) THEN
    ALTER TABLE snapshots ADD CONSTRAINT snapshots_created_by_fk
      FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL;
  END IF;
END $$;

-- =============================================================================
-- SECTION 2: Install updated_at triggers
-- =============================================================================

DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.triggers
                 WHERE trigger_name = 'set_updated_at' AND event_object_table = 'customers' AND trigger_schema = current_schema()) THEN
    PERFORM create_updated_at_trigger('customers');
  END IF;
END $$;

DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.triggers
                 WHERE trigger_name = 'set_updated_at' AND event_object_table = 'users' AND trigger_schema = current_schema()) THEN
    PERFORM create_updated_at_trigger('users');
  END IF;
END $$;

DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.triggers
                 WHERE trigger_name = 'set_updated_at' AND event_object_table = 'venues' AND trigger_schema = current_schema()) THEN
    PERFORM create_updated_at_trigger('venues');
  END IF;
END $$;

DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.triggers
                 WHERE trigger_name = 'set_updated_at' AND event_object_table = 'venue_user_roles' AND trigger_schema = current_schema()) THEN
    PERFORM create_updated_at_trigger('venue_user_roles');
  END IF;
END $$;

DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.triggers
                 WHERE trigger_name = 'set_updated_at' AND event_object_table = 'venue_invitations' AND trigger_schema = current_schema()) THEN
    PERFORM create_updated_at_trigger('venue_invitations');
  END IF;
END $$;

DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.triggers
                 WHERE trigger_name = 'set_updated_at' AND event_object_table = 'map_groups' AND trigger_schema = current_schema()) THEN
    PERFORM create_updated_at_trigger('map_groups');
  END IF;
END $$;

DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.triggers
                 WHERE trigger_name = 'set_updated_at' AND event_object_table = 'perspectives' AND trigger_schema = current_schema()) THEN
    PERFORM create_updated_at_trigger('perspectives');
  END IF;
END $$;

DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.triggers
                 WHERE trigger_name = 'set_updated_at' AND event_object_table = 'levels' AND trigger_schema = current_schema()) THEN
    PERFORM create_updated_at_trigger('levels');
  END IF;
END $$;

DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.triggers
                 WHERE trigger_name = 'set_updated_at' AND event_object_table = 'geo_references' AND trigger_schema = current_schema()) THEN
    PERFORM create_updated_at_trigger('geo_references');
  END IF;
END $$;

DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.triggers
                 WHERE trigger_name = 'set_updated_at' AND event_object_table = 'location_categories' AND trigger_schema = current_schema()) THEN
    PERFORM create_updated_at_trigger('location_categories');
  END IF;
END $$;

DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.triggers
                 WHERE trigger_name = 'set_updated_at' AND event_object_table = 'locations' AND trigger_schema = current_schema()) THEN
    PERFORM create_updated_at_trigger('locations');
  END IF;
END $$;

DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.triggers
                 WHERE trigger_name = 'set_updated_at' AND event_object_table = 'location_images' AND trigger_schema = current_schema()) THEN
    PERFORM create_updated_at_trigger('location_images');
  END IF;
END $$;

DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.triggers
                 WHERE trigger_name = 'set_updated_at' AND event_object_table = 'product_categories' AND trigger_schema = current_schema()) THEN
    PERFORM create_updated_at_trigger('product_categories');
  END IF;
END $$;

DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.triggers
                 WHERE trigger_name = 'set_updated_at' AND event_object_table = 'products' AND trigger_schema = current_schema()) THEN
    PERFORM create_updated_at_trigger('products');
  END IF;
END $$;

DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.triggers
                 WHERE trigger_name = 'set_updated_at' AND event_object_table = 'product_attachments' AND trigger_schema = current_schema()) THEN
    PERFORM create_updated_at_trigger('product_attachments');
  END IF;
END $$;

DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.triggers
                 WHERE trigger_name = 'set_updated_at' AND event_object_table = 'event_types' AND trigger_schema = current_schema()) THEN
    PERFORM create_updated_at_trigger('event_types');
  END IF;
END $$;

DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.triggers
                 WHERE trigger_name = 'set_updated_at' AND event_object_table = 'event_tags' AND trigger_schema = current_schema()) THEN
    PERFORM create_updated_at_trigger('event_tags');
  END IF;
END $$;

DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.triggers
                 WHERE trigger_name = 'set_updated_at' AND event_object_table = 'events' AND trigger_schema = current_schema()) THEN
    PERFORM create_updated_at_trigger('events');
  END IF;
END $$;

DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.triggers
                 WHERE trigger_name = 'set_updated_at' AND event_object_table = 'event_images' AND trigger_schema = current_schema()) THEN
    PERFORM create_updated_at_trigger('event_images');
  END IF;
END $$;

DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.triggers
                 WHERE trigger_name = 'set_updated_at' AND event_object_table = 'surveys' AND trigger_schema = current_schema()) THEN
    PERFORM create_updated_at_trigger('surveys');
  END IF;
END $$;

DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.triggers
                 WHERE trigger_name = 'set_updated_at' AND event_object_table = 'questions' AND trigger_schema = current_schema()) THEN
    PERFORM create_updated_at_trigger('questions');
  END IF;
END $$;

DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.triggers
                 WHERE trigger_name = 'set_updated_at' AND event_object_table = 'options' AND trigger_schema = current_schema()) THEN
    PERFORM create_updated_at_trigger('options');
  END IF;
END $$;

DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.triggers
                 WHERE trigger_name = 'set_updated_at' AND event_object_table = 'survey_responses' AND trigger_schema = current_schema()) THEN
    PERFORM create_updated_at_trigger('survey_responses');
  END IF;
END $$;

DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.triggers
                 WHERE trigger_name = 'set_updated_at' AND event_object_table = 'survey_answers' AND trigger_schema = current_schema()) THEN
    PERFORM create_updated_at_trigger('survey_answers');
  END IF;
END $$;

DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.triggers
                 WHERE trigger_name = 'set_updated_at' AND event_object_table = 'beacons' AND trigger_schema = current_schema()) THEN
    PERFORM create_updated_at_trigger('beacons');
  END IF;
END $$;

DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.triggers
                 WHERE trigger_name = 'set_updated_at' AND event_object_table = 'connections' AND trigger_schema = current_schema()) THEN
    PERFORM create_updated_at_trigger('connections');
  END IF;
END $$;

DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.triggers
                 WHERE trigger_name = 'set_updated_at' AND event_object_table = 'connection_levels' AND trigger_schema = current_schema()) THEN
    PERFORM create_updated_at_trigger('connection_levels');
  END IF;
END $$;

DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.triggers
                 WHERE trigger_name = 'set_updated_at' AND event_object_table = 'advertisements' AND trigger_schema = current_schema()) THEN
    PERFORM create_updated_at_trigger('advertisements');
  END IF;
END $$;

DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.triggers
                 WHERE trigger_name = 'set_updated_at' AND event_object_table = 'articles' AND trigger_schema = current_schema()) THEN
    PERFORM create_updated_at_trigger('articles');
  END IF;
END $$;

DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.triggers
                 WHERE trigger_name = 'set_updated_at' AND event_object_table = 'article_images' AND trigger_schema = current_schema()) THEN
    PERFORM create_updated_at_trigger('article_images');
  END IF;
END $$;

DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.triggers
                 WHERE trigger_name = 'set_updated_at' AND event_object_table = 'coupons' AND trigger_schema = current_schema()) THEN
    PERFORM create_updated_at_trigger('coupons');
  END IF;
END $$;

DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.triggers
                 WHERE trigger_name = 'set_updated_at' AND event_object_table = 'tags' AND trigger_schema = current_schema()) THEN
    PERFORM create_updated_at_trigger('tags');
  END IF;
END $$;

DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.triggers
                 WHERE trigger_name = 'set_updated_at' AND event_object_table = 'search_queries' AND trigger_schema = current_schema()) THEN
    PERFORM create_updated_at_trigger('search_queries');
  END IF;
END $$;

DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.triggers
                 WHERE trigger_name = 'set_updated_at' AND event_object_table = 'snapshots' AND trigger_schema = current_schema()) THEN
    PERFORM create_updated_at_trigger('snapshots');
  END IF;
END $$;

-- =============================================================================
-- SECTION 3: Set uuidv7() defaults on id columns
-- =============================================================================

DO $$ BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns
             WHERE table_name = 'customers' AND column_name = 'id' AND data_type = 'uuid' AND table_schema = current_schema()) THEN
    ALTER TABLE customers ALTER COLUMN id SET DEFAULT uuidv7();
  END IF;
END $$;

DO $$ BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns
             WHERE table_name = 'users' AND column_name = 'id' AND data_type = 'uuid' AND table_schema = current_schema()) THEN
    ALTER TABLE users ALTER COLUMN id SET DEFAULT uuidv7();
  END IF;
END $$;

DO $$ BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns
             WHERE table_name = 'venues' AND column_name = 'id' AND data_type = 'uuid' AND table_schema = current_schema()) THEN
    ALTER TABLE venues ALTER COLUMN id SET DEFAULT uuidv7();
  END IF;
END $$;

DO $$ BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns
             WHERE table_name = 'venue_user_roles' AND column_name = 'id' AND data_type = 'uuid' AND table_schema = current_schema()) THEN
    ALTER TABLE venue_user_roles ALTER COLUMN id SET DEFAULT uuidv7();
  END IF;
END $$;

DO $$ BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns
             WHERE table_name = 'venue_invitations' AND column_name = 'id' AND data_type = 'uuid' AND table_schema = current_schema()) THEN
    ALTER TABLE venue_invitations ALTER COLUMN id SET DEFAULT uuidv7();
  END IF;
END $$;

DO $$ BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns
             WHERE table_name = 'map_groups' AND column_name = 'id' AND data_type = 'uuid' AND table_schema = current_schema()) THEN
    ALTER TABLE map_groups ALTER COLUMN id SET DEFAULT uuidv7();
  END IF;
END $$;

DO $$ BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns
             WHERE table_name = 'perspectives' AND column_name = 'id' AND data_type = 'uuid' AND table_schema = current_schema()) THEN
    ALTER TABLE perspectives ALTER COLUMN id SET DEFAULT uuidv7();
  END IF;
END $$;

DO $$ BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns
             WHERE table_name = 'levels' AND column_name = 'id' AND data_type = 'uuid' AND table_schema = current_schema()) THEN
    ALTER TABLE levels ALTER COLUMN id SET DEFAULT uuidv7();
  END IF;
END $$;

DO $$ BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns
             WHERE table_name = 'geo_references' AND column_name = 'id' AND data_type = 'uuid' AND table_schema = current_schema()) THEN
    ALTER TABLE geo_references ALTER COLUMN id SET DEFAULT uuidv7();
  END IF;
END $$;

DO $$ BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns
             WHERE table_name = 'location_categories' AND column_name = 'id' AND data_type = 'uuid' AND table_schema = current_schema()) THEN
    ALTER TABLE location_categories ALTER COLUMN id SET DEFAULT uuidv7();
  END IF;
END $$;

DO $$ BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns
             WHERE table_name = 'locations' AND column_name = 'id' AND data_type = 'uuid' AND table_schema = current_schema()) THEN
    ALTER TABLE locations ALTER COLUMN id SET DEFAULT uuidv7();
  END IF;
END $$;

DO $$ BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns
             WHERE table_name = 'location_images' AND column_name = 'id' AND data_type = 'uuid' AND table_schema = current_schema()) THEN
    ALTER TABLE location_images ALTER COLUMN id SET DEFAULT uuidv7();
  END IF;
END $$;

DO $$ BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns
             WHERE table_name = 'product_categories' AND column_name = 'id' AND data_type = 'uuid' AND table_schema = current_schema()) THEN
    ALTER TABLE product_categories ALTER COLUMN id SET DEFAULT uuidv7();
  END IF;
END $$;

DO $$ BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns
             WHERE table_name = 'products' AND column_name = 'id' AND data_type = 'uuid' AND table_schema = current_schema()) THEN
    ALTER TABLE products ALTER COLUMN id SET DEFAULT uuidv7();
  END IF;
END $$;

DO $$ BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns
             WHERE table_name = 'product_attachments' AND column_name = 'id' AND data_type = 'uuid' AND table_schema = current_schema()) THEN
    ALTER TABLE product_attachments ALTER COLUMN id SET DEFAULT uuidv7();
  END IF;
END $$;

DO $$ BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns
             WHERE table_name = 'event_types' AND column_name = 'id' AND data_type = 'uuid' AND table_schema = current_schema()) THEN
    ALTER TABLE event_types ALTER COLUMN id SET DEFAULT uuidv7();
  END IF;
END $$;

DO $$ BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns
             WHERE table_name = 'event_tags' AND column_name = 'id' AND data_type = 'uuid' AND table_schema = current_schema()) THEN
    ALTER TABLE event_tags ALTER COLUMN id SET DEFAULT uuidv7();
  END IF;
END $$;

DO $$ BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns
             WHERE table_name = 'events' AND column_name = 'id' AND data_type = 'uuid' AND table_schema = current_schema()) THEN
    ALTER TABLE events ALTER COLUMN id SET DEFAULT uuidv7();
  END IF;
END $$;

DO $$ BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns
             WHERE table_name = 'event_images' AND column_name = 'id' AND data_type = 'uuid' AND table_schema = current_schema()) THEN
    ALTER TABLE event_images ALTER COLUMN id SET DEFAULT uuidv7();
  END IF;
END $$;

DO $$ BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns
             WHERE table_name = 'surveys' AND column_name = 'id' AND data_type = 'uuid' AND table_schema = current_schema()) THEN
    ALTER TABLE surveys ALTER COLUMN id SET DEFAULT uuidv7();
  END IF;
END $$;

DO $$ BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns
             WHERE table_name = 'questions' AND column_name = 'id' AND data_type = 'uuid' AND table_schema = current_schema()) THEN
    ALTER TABLE questions ALTER COLUMN id SET DEFAULT uuidv7();
  END IF;
END $$;

DO $$ BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns
             WHERE table_name = 'options' AND column_name = 'id' AND data_type = 'uuid' AND table_schema = current_schema()) THEN
    ALTER TABLE options ALTER COLUMN id SET DEFAULT uuidv7();
  END IF;
END $$;

DO $$ BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns
             WHERE table_name = 'survey_responses' AND column_name = 'id' AND data_type = 'uuid' AND table_schema = current_schema()) THEN
    ALTER TABLE survey_responses ALTER COLUMN id SET DEFAULT uuidv7();
  END IF;
END $$;

DO $$ BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns
             WHERE table_name = 'survey_answers' AND column_name = 'id' AND data_type = 'uuid' AND table_schema = current_schema()) THEN
    ALTER TABLE survey_answers ALTER COLUMN id SET DEFAULT uuidv7();
  END IF;
END $$;

DO $$ BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns
             WHERE table_name = 'beacons' AND column_name = 'id' AND data_type = 'uuid' AND table_schema = current_schema()) THEN
    ALTER TABLE beacons ALTER COLUMN id SET DEFAULT uuidv7();
  END IF;
END $$;

DO $$ BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns
             WHERE table_name = 'connections' AND column_name = 'id' AND data_type = 'uuid' AND table_schema = current_schema()) THEN
    ALTER TABLE connections ALTER COLUMN id SET DEFAULT uuidv7();
  END IF;
END $$;

DO $$ BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns
             WHERE table_name = 'connection_levels' AND column_name = 'id' AND data_type = 'uuid' AND table_schema = current_schema()) THEN
    ALTER TABLE connection_levels ALTER COLUMN id SET DEFAULT uuidv7();
  END IF;
END $$;

DO $$ BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns
             WHERE table_name = 'advertisements' AND column_name = 'id' AND data_type = 'uuid' AND table_schema = current_schema()) THEN
    ALTER TABLE advertisements ALTER COLUMN id SET DEFAULT uuidv7();
  END IF;
END $$;

DO $$ BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns
             WHERE table_name = 'articles' AND column_name = 'id' AND data_type = 'uuid' AND table_schema = current_schema()) THEN
    ALTER TABLE articles ALTER COLUMN id SET DEFAULT uuidv7();
  END IF;
END $$;

DO $$ BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns
             WHERE table_name = 'article_images' AND column_name = 'id' AND data_type = 'uuid' AND table_schema = current_schema()) THEN
    ALTER TABLE article_images ALTER COLUMN id SET DEFAULT uuidv7();
  END IF;
END $$;

DO $$ BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns
             WHERE table_name = 'coupons' AND column_name = 'id' AND data_type = 'uuid' AND table_schema = current_schema()) THEN
    ALTER TABLE coupons ALTER COLUMN id SET DEFAULT uuidv7();
  END IF;
END $$;

DO $$ BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns
             WHERE table_name = 'tags' AND column_name = 'id' AND data_type = 'uuid' AND table_schema = current_schema()) THEN
    ALTER TABLE tags ALTER COLUMN id SET DEFAULT uuidv7();
  END IF;
END $$;

DO $$ BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns
             WHERE table_name = 'search_queries' AND column_name = 'id' AND data_type = 'uuid' AND table_schema = current_schema()) THEN
    ALTER TABLE search_queries ALTER COLUMN id SET DEFAULT uuidv7();
  END IF;
END $$;

DO $$ BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.columns
             WHERE table_name = 'snapshots' AND column_name = 'id' AND data_type = 'uuid' AND table_schema = current_schema()) THEN
    ALTER TABLE snapshots ALTER COLUMN id SET DEFAULT uuidv7();
  END IF;
END $$;

-- =============================================================================
-- SECTION 4: Drop temporary migration helper table
-- =============================================================================

DROP TABLE IF EXISTS _user_id_map;

-- Script 03: Rename all indoormap_api_* tables to new schema names
ALTER TABLE indoormap_api_customer               RENAME TO customers;
ALTER TABLE auth_user                            RENAME TO users;
ALTER TABLE indoormap_api_venue                  RENAME TO venues;
ALTER TABLE indoormap_api_venueuserrole          RENAME TO venue_user_roles;
ALTER TABLE indoormap_api_venueinvitation        RENAME TO venue_invitations;
ALTER TABLE indoormap_api_mapgroup               RENAME TO map_groups;
ALTER TABLE indoormap_api_perspective            RENAME TO perspectives;
ALTER TABLE indoormap_api_level                  RENAME TO levels;
ALTER TABLE indoormap_api_georeference           RENAME TO geo_references;
ALTER TABLE indoormap_api_locationcategory       RENAME TO location_categories;
ALTER TABLE indoormap_api_location               RENAME TO locations;
ALTER TABLE indoormap_api_location_common_categories RENAME TO location_category_links;
ALTER TABLE indoormap_api_locationimage          RENAME TO location_images;
-- NOTE: indoormap_api_locationpromotion does not exist in the source DB; skipping promotions rename
ALTER TABLE indoormap_api_productcategory        RENAME TO product_categories;
ALTER TABLE indoormap_api_product                RENAME TO products;
ALTER TABLE indoormap_api_product_categories     RENAME TO product_category_links;
ALTER TABLE indoormap_api_productattachment      RENAME TO product_attachments;
ALTER TABLE indoormap_api_eventtype              RENAME TO event_types;
ALTER TABLE indoormap_api_eventtag               RENAME TO event_tags;
ALTER TABLE indoormap_api_event                  RENAME TO events;
ALTER TABLE indoormap_api_event_tags             RENAME TO event_tag_links;
ALTER TABLE indoormap_api_event_locations        RENAME TO event_location_links;
ALTER TABLE indoormap_api_eventimage             RENAME TO event_images;
ALTER TABLE indoormap_api_survey                 RENAME TO surveys;
ALTER TABLE indoormap_api_question               RENAME TO questions;
ALTER TABLE indoormap_api_option                 RENAME TO options;
ALTER TABLE indoormap_api_surveyresponse         RENAME TO survey_responses;
ALTER TABLE indoormap_api_surveyanswer           RENAME TO survey_answers;
ALTER TABLE indoormap_api_beacon                 RENAME TO beacons;
ALTER TABLE indoormap_api_connection             RENAME TO connections;
ALTER TABLE indoormap_api_connectionlevel        RENAME TO connection_levels;
ALTER TABLE indoormap_api_advertisement          RENAME TO advertisements;
ALTER TABLE indoormap_api_article                RENAME TO articles;
ALTER TABLE indoormap_api_articleimage           RENAME TO article_images;
ALTER TABLE indoormap_api_coupon                 RENAME TO coupons;
ALTER TABLE indoormap_api_tag                    RENAME TO tags;
ALTER TABLE indoormap_api_searchquery            RENAME TO search_queries;
ALTER TABLE indoormap_api_snapshot               RENAME TO snapshots;
-- NOTE: indoormap_api_qrcode does not exist in the source DB; skipping qrcodes rename
-- NOTE: indoormap_api_resetpasswordtoken does not exist in the source DB; skipping reset_password_tokens rename

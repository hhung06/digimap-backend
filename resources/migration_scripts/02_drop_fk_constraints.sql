-- Script 02: Drop all FK constraints before renaming/transforming tables
SET search_path TO digimap_db, public;

-- Helper: drop a FK constraint, silently ignoring missing tables/constraints
CREATE OR REPLACE FUNCTION _drop_fk(tbl TEXT, con TEXT) RETURNS VOID LANGUAGE plpgsql AS $$
BEGIN
  EXECUTE format('ALTER TABLE %I DROP CONSTRAINT IF EXISTS %I', tbl, con);
EXCEPTION WHEN undefined_table THEN NULL;
END $$;

-- indoormap_api_profile (may not exist)
SELECT _drop_fk('indoormap_api_profile', 'indoormap_api_profile_user_id_5a841cea_fk_auth_user_id');
SELECT _drop_fk('indoormap_api_profile', 'indoormap_api_profil_customer_id_f1024b17_fk_indoormap');

-- indoormap_api_appuserdevice (may not exist)
SELECT _drop_fk('indoormap_api_appuserdevice', 'indoormap_api_appuse_user_id_99c19ad8_fk_indoormap');

-- indoormap_api_venue
SELECT _drop_fk('indoormap_api_venue', 'indoormap_api_venue_customer_id_0de13a77_fk_indoormap');

-- indoormap_api_venueuserrole
SELECT _drop_fk('indoormap_api_venueuserrole', 'indoormap_api_venueu_profile_id_37cc1c83_fk_indoormap');
SELECT _drop_fk('indoormap_api_venueuserrole', 'indoormap_api_venueu_venue_id_3ab36cdc_fk_indoormap');

-- indoormap_api_venueinvitation
SELECT _drop_fk('indoormap_api_venueinvitation', 'indoormap_api_venuei_invited_by_id_2915f3a7_fk_indoormap');
SELECT _drop_fk('indoormap_api_venueinvitation', 'indoormap_api_venuei_invited_profile_id_203f3f41_fk_indoormap');
SELECT _drop_fk('indoormap_api_venueinvitation', 'indoormap_api_venuei_venue_id_4fb38cff_fk_indoormap');

-- indoormap_api_mapgroup
SELECT _drop_fk('indoormap_api_mapgroup', 'indoormap_api_mapgro_venue_id_ae910a39_fk_indoormap');

-- indoormap_api_level
SELECT _drop_fk('indoormap_api_level', 'indoormap_api_level_venue_id_2413c37c_fk');
SELECT _drop_fk('indoormap_api_level', 'indoormap_api_level_perspective_id_2e9d9229_fk_indoormap');
SELECT _drop_fk('indoormap_api_level', 'indoormap_api_level_type_id_a18364cf_fk_indoormap');
SELECT _drop_fk('indoormap_api_level', 'indoormap_api_level_mapgroup_id_9d651189_fk_indoormap');

-- indoormap_api_georeference
SELECT _drop_fk('indoormap_api_georeference', 'indoormap_api_georef_level_id_d043778c_fk_indoormap');

-- indoormap_api_locationcategory
SELECT _drop_fk('indoormap_api_locationcategory', 'indoormap_api_locati_venue_id_6a5feeb8_fk_indoormap');

-- indoormap_api_location
SELECT _drop_fk('indoormap_api_location', 'indoormap_api_locati_level_id_efc07011_fk_indoormap');
SELECT _drop_fk('indoormap_api_location', 'indoormap_api_locati_venue_id_ea11932a_fk_indoormap');
SELECT _drop_fk('indoormap_api_location', 'indoormap_api_locati_main_category_id_62f50b77_fk_indoormap');

-- indoormap_api_location_common_categories
SELECT _drop_fk('indoormap_api_location_common_categories', 'indoormap_api_locati_location_id_5431b8ac_fk_indoormap');
SELECT _drop_fk('indoormap_api_location_common_categories', 'indoormap_api_locati_locationcategory_id_b3cb6515_fk_indoormap');

-- indoormap_api_locationimage
SELECT _drop_fk('indoormap_api_locationimage', 'indoormap_api_locati_location_id_74cf1f87_fk_indoormap');

-- indoormap_api_productcategory
SELECT _drop_fk('indoormap_api_productcategory', 'indoormap_api_produc_venue_id_348a595a_fk_indoormap');

-- indoormap_api_product
SELECT _drop_fk('indoormap_api_product', 'indoormap_api_produc_location_id_edf95cb2_fk_indoormap');
SELECT _drop_fk('indoormap_api_product', 'indoormap_api_produc_main_category_id_7221dfe3_fk_indoormap');
SELECT _drop_fk('indoormap_api_product', 'indoormap_api_produc_venue_id_927295fe_fk_indoormap');

-- indoormap_api_product_categories
SELECT _drop_fk('indoormap_api_product_categories', 'indoormap_api_produc_product_id_5c419598_fk_indoormap');
SELECT _drop_fk('indoormap_api_product_categories', 'indoormap_api_produc_productcategory_id_678f4361_fk_indoormap');

-- indoormap_api_productattachment
SELECT _drop_fk('indoormap_api_productattachment', 'indoormap_api_produc_product_id_8a524e4e_fk_indoormap');

-- indoormap_api_eventtype
SELECT _drop_fk('indoormap_api_eventtype', 'indoormap_api_eventt_venue_id_aa1577e0_fk_indoormap');

-- indoormap_api_event
SELECT _drop_fk('indoormap_api_event', 'indoormap_api_event_venue_id_08a7d602_fk_indoormap_api_venue_id');

-- indoormap_api_event_tags
SELECT _drop_fk('indoormap_api_event_tags', 'indoormap_api_event__event_id_4582f9f1_fk_indoormap');
SELECT _drop_fk('indoormap_api_event_tags', 'indoormap_api_event__eventtag_id_60216ac7_fk_indoormap');

-- indoormap_api_event_locations
SELECT _drop_fk('indoormap_api_event_locations', 'indoormap_api_event__event_id_35d8fa44_fk_indoormap');
SELECT _drop_fk('indoormap_api_event_locations', 'indoormap_api_event__location_id_9e57c39f_fk_indoormap');

-- indoormap_api_eventimage
SELECT _drop_fk('indoormap_api_eventimage', 'indoormap_api_eventi_event_id_7c1f1f54_fk_indoormap');

-- indoormap_api_survey
SELECT _drop_fk('indoormap_api_survey', 'indoormap_api_survey_created_by_id_e3b2ab61_fk_indoormap');
SELECT _drop_fk('indoormap_api_survey', 'indoormap_api_survey_venue_id_ceb6ac5c_fk_indoormap_api_venue_i');

-- indoormap_api_question
SELECT _drop_fk('indoormap_api_question', 'indoormap_api_questi_survey_id_d0aa8913_fk_indoormap');

-- indoormap_api_option
SELECT _drop_fk('indoormap_api_option', 'indoormap_api_option_question_id_32136110_fk_indoormap');

-- indoormap_api_surveyresponse
SELECT _drop_fk('indoormap_api_surveyresponse', 'indoormap_api_survey_survey_id_cf938a9c_fk_indoormap');
SELECT _drop_fk('indoormap_api_surveyresponse', 'indoormap_api_survey_user_id_c6490791_fk_indoormap');

-- indoormap_api_surveyanswer
SELECT _drop_fk('indoormap_api_surveyanswer', 'indoormap_api_survey_option_id_1f032bc1_fk_indoormap');
SELECT _drop_fk('indoormap_api_surveyanswer', 'indoormap_api_survey_question_id_344d0ff0_fk_indoormap');
SELECT _drop_fk('indoormap_api_surveyanswer', 'indoormap_api_survey_response_id_c45e3d07_fk_indoormap');

-- indoormap_api_beacon
SELECT _drop_fk('indoormap_api_beacon', 'indoormap_api_beacon_venue_id_08aabc2f_fk_indoormap_api_venue_i');
SELECT _drop_fk('indoormap_api_beacon', 'indoormap_api_beacon_level_id_029d21be_fk_indoormap_api_level_i');

-- indoormap_api_connection
SELECT _drop_fk('indoormap_api_connection', 'indoormap_api_connec_venue_id_3002b720_fk_indoormap');

-- indoormap_api_connectionlevel
SELECT _drop_fk('indoormap_api_connectionlevel', 'indoormap_api_connec_connection_id_4f11bb58_fk_indoormap');
SELECT _drop_fk('indoormap_api_connectionlevel', 'indoormap_api_connec_level_id_91d5e466_fk_indoormap');

-- indoormap_api_advertisement
SELECT _drop_fk('indoormap_api_advertisement', 'indoormap_api_advert_venue_id_20557289_fk_indoormap');
SELECT _drop_fk('indoormap_api_advertisement', 'indoormap_api_advert_location_id_c8af4309_fk_indoormap');
SELECT _drop_fk('indoormap_api_advertisement', 'indoormap_api_advert_article_id_5edbc569_fk_indoormap');

-- indoormap_api_article
SELECT _drop_fk('indoormap_api_article', 'indoormap_api_articl_product_id_e8f3d65f_fk_indoormap');
SELECT _drop_fk('indoormap_api_article', 'indoormap_api_articl_product_category_id_395d7522_fk_indoormap');
SELECT _drop_fk('indoormap_api_article', 'indoormap_api_articl_venue_id_337016f5_fk_indoormap');
SELECT _drop_fk('indoormap_api_article', 'indoormap_api_articl_location_id_aa68488f_fk_indoormap');

-- indoormap_api_articleimage
SELECT _drop_fk('indoormap_api_articleimage', 'indoormap_api_articl_article_id_2c093e47_fk_indoormap');

-- indoormap_api_coupon
SELECT _drop_fk('indoormap_api_coupon', 'indoormap_api_coupon_survey_id_5c6684cd_fk_indoormap');
SELECT _drop_fk('indoormap_api_coupon', 'indoormap_api_coupon_venue_id_4e51a9b3_fk_indoormap_api_venue_i');

-- indoormap_api_snapshot
SELECT _drop_fk('indoormap_api_snapshot', 'indoormap_api_snapshot_publish_by_id_4e82db52_fk_auth_user_id');

-- Cleanup helper
DROP FUNCTION IF EXISTS _drop_fk(TEXT, TEXT);

-- Script 01: Prerequisites — extensions, functions, framework cleanup
SET search_path TO digimap_db, public;

-- =============================================================================
-- Extensions & functions
-- =============================================================================

CREATE EXTENSION IF NOT EXISTS "pgcrypto";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";
CREATE EXTENSION IF NOT EXISTS "btree_gist";

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION create_updated_at_trigger(tbl TEXT)
RETURNS VOID LANGUAGE plpgsql AS $$
BEGIN
  EXECUTE format(
    'CREATE TRIGGER set_updated_at BEFORE UPDATE ON %I
     FOR EACH ROW EXECUTE FUNCTION update_updated_at_column()', tbl);
END;
$$;

-- =============================================================================
-- Drop Django framework tables
-- =============================================================================

DROP TABLE IF EXISTS auth_group_permissions CASCADE;
DROP TABLE IF EXISTS auth_user_groups CASCADE;
DROP TABLE IF EXISTS auth_user_user_permissions CASCADE;
DROP TABLE IF EXISTS authtoken_token CASCADE;
DROP TABLE IF EXISTS auth_group CASCADE;
DROP TABLE IF EXISTS auth_permission CASCADE;
DROP TABLE IF EXISTS debug_toolbar_historyentry CASCADE;
DROP TABLE IF EXISTS django_admin_log CASCADE;
DROP TABLE IF EXISTS django_content_type CASCADE;
DROP TABLE IF EXISTS django_migrations CASCADE;
DROP TABLE IF EXISTS django_session CASCADE;
DROP TABLE IF EXISTS rest_framework_api_key_apikey CASCADE;
DROP TABLE IF EXISTS token_blacklist_blacklistedtoken CASCADE;
DROP TABLE IF EXISTS token_blacklist_outstandingtoken CASCADE;

-- =============================================================================
-- Drop no-equivalent indoormap_api_* tables
-- =============================================================================

DROP TABLE IF EXISTS indoormap_api_accesshistory CASCADE;
DROP TABLE IF EXISTS indoormap_api_adtrack CASCADE;
DROP TABLE IF EXISTS indoormap_api_articlerelatedproducts CASCADE;
DROP TABLE IF EXISTS indoormap_api_libraryasset CASCADE;
DROP TABLE IF EXISTS indoormap_api_categoryad CASCADE;
DROP TABLE IF EXISTS indoormap_api_coupon_gifts CASCADE;
DROP TABLE IF EXISTS indoormap_api_promogift CASCADE;
DROP TABLE IF EXISTS indoormap_api_forcesyncversion CASCADE;
DROP TABLE IF EXISTS indoormap_api_forceversion CASCADE;
DROP TABLE IF EXISTS indoormap_api_location_place_tags CASCADE;
drop table if exists indoormap_api_locationamenity cascade;
DROP TABLE IF EXISTS indoormap_api_locationcategorytranslation CASCADE;
DROP TABLE IF EXISTS indoormap_api_locationtranslation CASCADE;
DROP TABLE IF EXISTS indoormap_api_locationtemplate CASCADE;
DROP TABLE IF EXISTS indoormap_api_locationtemplate_common_categories CASCADE;
DROP TABLE IF EXISTS indoormap_api_locationtemplate_place_tags CASCADE;
DROP TABLE IF EXISTS indoormap_api_locationtemplateimage CASCADE;
DROP TABLE IF EXISTS indoormap_api_locationtemplatetranslation CASCADE;
DROP TABLE IF EXISTS indoormap_api_locationpromotion CASCADE;
DROP TABLE IF EXISTS indoormap_api_new CASCADE;
DROP TABLE IF EXISTS indoormap_api_plugin CASCADE;
DROP TABLE IF EXISTS indoormap_api_profile CASCADE;
DROP TABLE IF EXISTS indoormap_api_segment CASCADE;
DROP TABLE IF EXISTS indoormap_api_syncdata CASCADE;
DROP TABLE IF EXISTS indoormap_api_venue_sync_flag CASCADE;
DROP TABLE IF EXISTS indoormap_api_venuetheme CASCADE;
DROP TABLE IF EXISTS indoormap_api_venuetranslation CASCADE;
DROP TABLE IF EXISTS indoormap_api_tagtranslation CASCADE;
DROP TABLE IF EXISTS indoormap_api_usercoupon CASCADE;
DROP TABLE IF EXISTS indoormap_api_templatecategory CASCADE;
DROP TABLE IF EXISTS indoormap_api_templatecategorytranslation CASCADE;
-- NOTE: indoormap_api_appuser and indoormap_api_visitor are kept here;
-- their data is migrated in script 14 before they are dropped.

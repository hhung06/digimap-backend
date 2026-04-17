-- Script 01: Prerequisites — extensions, functions, framework cleanup

-- =============================================================================
-- Extensions & functions
-- =============================================================================

-- Required PostgreSQL extensions
CREATE EXTENSION IF NOT EXISTS "pgcrypto";   -- gen_random_bytes() and other crypto helpers (uuidv7 is native in PG18+)
CREATE EXTENSION IF NOT EXISTS "pg_trgm";    -- trigram similarity for fuzzy text search
CREATE EXTENSION IF NOT EXISTS "btree_gist"; -- exclusion constraints with ranges
CREATE EXTENSION IF NOT EXISTS "vector";     -- pgvector: dense vector similarity search (RAG / embeddings)

-- Reusable trigger: auto-update updated_at on every UPDATE
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

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
DROP TABLE IF EXISTS indoormap_api_asset CASCADE;
DROP TABLE IF EXISTS indoormap_api_libraryasset CASCADE;
DROP TABLE IF EXISTS indoormap_api_categoryad CASCADE;
DROP TABLE IF EXISTS indoormap_api_coupon_gifts CASCADE;
DROP TABLE IF EXISTS indoormap_api_promogift CASCADE;
DROP TABLE IF EXISTS indoormap_api_forcesyncversion CASCADE;
DROP TABLE IF EXISTS indoormap_api_forceversion CASCADE;
DROP TABLE IF EXISTS indoormap_api_language CASCADE;
DROP TABLE IF EXISTS indoormap_api_leveltype CASCADE;
DROP TABLE IF EXISTS indoormap_api_location_place_tags CASCADE;
DROP TABLE IF EXISTS indoormap_api_locationcategorytranslation CASCADE;
DROP TABLE IF EXISTS indoormap_api_locationtranslation CASCADE;
DROP TABLE IF EXISTS indoormap_api_locationtemplate CASCADE;
DROP TABLE IF EXISTS indoormap_api_locationtemplate_common_categories CASCADE;
DROP TABLE IF EXISTS indoormap_api_locationtemplate_place_tags CASCADE;
DROP TABLE IF EXISTS indoormap_api_locationtemplateimage CASCADE;
DROP TABLE IF EXISTS indoormap_api_locationtemplatetranslation CASCADE;
DROP TABLE IF EXISTS indoormap_api_new CASCADE;
DROP TABLE IF EXISTS indoormap_api_plugin CASCADE;
DROP TABLE IF EXISTS indoormap_api_productplaza CASCADE;
DROP TABLE IF EXISTS indoormap_api_productplazadescriptionimage CASCADE;
DROP TABLE IF EXISTS indoormap_api_productplazaimage CASCADE;
DROP TABLE IF EXISTS indoormap_api_productplazapendingvideo CASCADE;
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

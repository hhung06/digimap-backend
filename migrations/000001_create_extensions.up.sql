-- Required PostgreSQL extensions
CREATE EXTENSION IF NOT EXISTS "pgcrypto";   -- gen_random_uuid()
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

CREATE TABLE articles (
    id                     UUID PRIMARY KEY DEFAULT uuidv7(),
    venue_id               UUID REFERENCES venues(id) ON DELETE CASCADE,
    external_id            VARCHAR(255),
    location_id            UUID REFERENCES locations(id) ON DELETE SET NULL,
    placement              VARCHAR(50) NOT NULL DEFAULT 'article',
    navigate               VARCHAR(50),
    title                  VARCHAR(255) NOT NULL DEFAULT '',
    label                  VARCHAR(255),
    content                TEXT,
    status                 VARCHAR(50) NOT NULL DEFAULT 'draft',
    published_at           TIMESTAMPTZ,
    published_period_start DATE,
    published_period_end   DATE,
    localization           JSONB,
    created_at             TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at             TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at             TIMESTAMPTZ
);

CREATE TABLE article_images (
    id         UUID PRIMARY KEY DEFAULT uuidv7(),
    article_id UUID NOT NULL REFERENCES articles(id) ON DELETE CASCADE,
    image      VARCHAR(1000) NOT NULL DEFAULT '',
    sort_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_articles_venue_id ON articles(venue_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_article_images_article_id ON article_images(article_id) WHERE deleted_at IS NULL;

SELECT create_updated_at_trigger('articles');
SELECT create_updated_at_trigger('article_images');

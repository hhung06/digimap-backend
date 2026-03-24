CREATE TABLE IF NOT EXISTS surveys (
    id               UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    venue_id         UUID        REFERENCES venues(id) ON DELETE CASCADE,
    external_id      TEXT,
    title            TEXT,
    content          TEXT,
    start_date       TIMESTAMPTZ,
    end_date         TIMESTAMPTZ,
    status           SMALLINT    NOT NULL DEFAULT 1,  -- 1=draft 2=active 3=inactive 4=closed
    publish_type     SMALLINT    NOT NULL DEFAULT 1,  -- 1=push 2=in_app 3=both 4=promo_banner
    is_forced        BOOLEAN     NOT NULL DEFAULT FALSE,
    source           SMALLINT    NOT NULL DEFAULT 1,  -- 1=cms 2=app 3=import
    app              TEXT        NOT NULL DEFAULT 'all',
    segment_filters  JSONB,
    created_by       UUID        REFERENCES users(id) ON DELETE SET NULL,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at       TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_surveys_venue  ON surveys (venue_id)         WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_surveys_status ON surveys (status)           WHERE deleted_at IS NULL;

CREATE TRIGGER set_surveys_updated_at
    BEFORE UPDATE ON surveys
    FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();

-- Add FK from notifications to surveys now that surveys table exists
ALTER TABLE notifications ADD CONSTRAINT fk_notifications_survey
    FOREIGN KEY (survey_id) REFERENCES surveys(id) ON DELETE SET NULL;

CREATE TABLE IF NOT EXISTS questions (
    id              UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    survey_id       UUID        NOT NULL REFERENCES surveys(id) ON DELETE CASCADE,
    question_number INT         NOT NULL DEFAULT 1,
    question_type   TEXT        NOT NULL DEFAULT 'paragraph',
    question_text   TEXT        NOT NULL,
    is_required     BOOLEAN     NOT NULL DEFAULT FALSE,
    is_other        BOOLEAN     NOT NULL DEFAULT FALSE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_questions_survey ON questions (survey_id) WHERE deleted_at IS NULL;

CREATE TRIGGER set_questions_updated_at
    BEFORE UPDATE ON questions
    FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();

CREATE TABLE IF NOT EXISTS options (
    id            UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    question_id   UUID        NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
    option_number INT         NOT NULL DEFAULT 1,
    option_text   TEXT        NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at    TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_options_question ON options (question_id) WHERE deleted_at IS NULL;

CREATE TRIGGER set_options_updated_at
    BEFORE UPDATE ON options
    FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();

CREATE TABLE IF NOT EXISTS survey_responses (
    id           UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    survey_id    UUID        NOT NULL REFERENCES surveys(id) ON DELETE CASCADE,
    submitted_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at   TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_survey_responses_survey ON survey_responses (survey_id) WHERE deleted_at IS NULL;

CREATE TRIGGER set_survey_responses_updated_at
    BEFORE UPDATE ON survey_responses
    FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();

CREATE TABLE IF NOT EXISTS survey_answers (
    id          UUID  PRIMARY KEY DEFAULT gen_random_uuid(),
    response_id UUID  NOT NULL REFERENCES survey_responses(id) ON DELETE CASCADE,
    question_id UUID  NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
    option_id   UUID  REFERENCES options(id) ON DELETE SET NULL,
    answer_text TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_survey_answers_response  ON survey_answers (response_id)  WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_survey_answers_question  ON survey_answers (question_id)  WHERE deleted_at IS NULL;

CREATE TRIGGER set_survey_answers_updated_at
    BEFORE UPDATE ON survey_answers
    FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();

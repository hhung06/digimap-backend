ALTER TABLE survey_responses
    ADD COLUMN IF NOT EXISTS external_id TEXT;

CREATE UNIQUE INDEX IF NOT EXISTS uq_survey_responses_survey_external_id
    ON survey_responses (survey_id, external_id)
    WHERE external_id IS NOT NULL AND deleted_at IS NULL;

DROP INDEX IF EXISTS uq_survey_responses_survey_external_id;

ALTER TABLE survey_responses
    DROP COLUMN IF EXISTS external_id;

-- Script 10: Transform surveys, questions, options, survey_responses, survey_answers

-- ── surveys ──────────────────────────────────────────────────────────────────
-- Actual columns: created_by_id (bigint FK to indoormap_api_appuser), end_date, is_forced,
--                 publish_type, start_date, status, title, app, restored_at, transaction_id,
--                 content, source, venue_id, segment_filters, external_id
ALTER TABLE surveys DROP COLUMN IF EXISTS created_by_id;
ALTER TABLE surveys DROP COLUMN IF EXISTS restored_at;
ALTER TABLE surveys DROP COLUMN IF EXISTS transaction_id;
ALTER TABLE surveys ADD COLUMN IF NOT EXISTS created_by UUID;

-- ── questions ────────────────────────────────────────────────────────────────
-- Actual columns: question_number (bigint), question_type, question_text, survey_id,
--                 is_required, is_other, restored_at, transaction_id
ALTER TABLE questions DROP COLUMN IF EXISTS restored_at;
ALTER TABLE questions DROP COLUMN IF EXISTS transaction_id;
ALTER TABLE questions ALTER COLUMN question_number TYPE INTEGER USING question_number::INTEGER;

-- ── options ──────────────────────────────────────────────────────────────────
-- Actual columns: option_number (bigint), option_text, question_id, restored_at, transaction_id
ALTER TABLE options DROP COLUMN IF EXISTS restored_at;
ALTER TABLE options DROP COLUMN IF EXISTS transaction_id;
ALTER TABLE options ALTER COLUMN option_number TYPE INTEGER USING option_number::INTEGER;

-- ── survey_responses ─────────────────────────────────────────────────────────
-- Actual columns: submitted_at, survey_id, user_id (bigint FK to indoormap_api_appuser),
--                 created_at, deleted_at, restored_at, transaction_id, updated_at
ALTER TABLE survey_responses DROP COLUMN IF EXISTS user_id;
ALTER TABLE survey_responses DROP COLUMN IF EXISTS restored_at;
ALTER TABLE survey_responses DROP COLUMN IF EXISTS transaction_id;

-- ── survey_answers ───────────────────────────────────────────────────────────
-- Actual columns: answer_text, option_id, question_id, response_id,
--                 created_at, deleted_at, restored_at, transaction_id, updated_at
ALTER TABLE survey_answers DROP COLUMN IF EXISTS restored_at;
ALTER TABLE survey_answers DROP COLUMN IF EXISTS transaction_id;

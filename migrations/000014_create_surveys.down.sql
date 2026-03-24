ALTER TABLE notifications DROP CONSTRAINT IF EXISTS fk_notifications_survey;
DROP TABLE IF EXISTS survey_answers;
DROP TABLE IF EXISTS survey_responses;
DROP TABLE IF EXISTS options;
DROP TABLE IF EXISTS questions;
DROP TABLE IF EXISTS surveys;

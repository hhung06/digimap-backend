package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/repository"
)

type surveyRepo struct {
	pool *pgxpool.Pool
}

// NewSurveyRepository returns a SurveyRepository backed by PostgreSQL.
func NewSurveyRepository(pool *pgxpool.Pool) repository.SurveyRepository {
	return &surveyRepo{pool: pool}
}

const surveySelectCols = `id, venue_id, external_id, title, content, start_date, end_date,
	status, publish_type, is_forced, source, app, segment_filters, created_by, created_at, updated_at`

func (r *surveyRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.Survey, error) {
	q := `SELECT ` + surveySelectCols + ` FROM surveys WHERE id = $1 AND deleted_at IS NULL`
	s, err := scanSurvey(r.pool.QueryRow(ctx, q, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.NewNotFound("survey not found")
	}
	if err != nil {
		return nil, err
	}
	// Eagerly load questions + options
	qs, err := r.ListQuestions(ctx, s.ID)
	if err != nil {
		return nil, err
	}
	s.Questions = qs
	return s, nil
}

func (r *surveyRepo) List(ctx context.Context, venueID uuid.UUID, p domain.Pagination) ([]*domain.Survey, int64, error) {
	var total int64
	if err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM surveys WHERE venue_id = $1 AND deleted_at IS NULL`, venueID,
	).Scan(&total); err != nil {
		return nil, 0, err
	}

	q := `SELECT ` + surveySelectCols + ` FROM surveys WHERE venue_id = $1 AND deleted_at IS NULL ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	rows, err := r.pool.Query(ctx, q, venueID, p.PageSize, p.Offset())
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var surveys []*domain.Survey
	for rows.Next() {
		s, err := scanSurvey(rows)
		if err != nil {
			return nil, 0, err
		}
		surveys = append(surveys, s)
	}
	return surveys, total, rows.Err()
}

func (r *surveyRepo) ListDueActivation(ctx context.Context, now time.Time) ([]*domain.Survey, error) {
	q := `SELECT ` + surveySelectCols + ` FROM surveys
		WHERE status = $1 AND start_date IS NOT NULL AND start_date <= $2 AND deleted_at IS NULL
		ORDER BY start_date ASC`
	rows, err := r.pool.Query(ctx, q, domain.SurveyStatusInactive, now)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var surveys []*domain.Survey
	for rows.Next() {
		s, err := scanSurvey(rows)
		if err != nil {
			return nil, err
		}
		surveys = append(surveys, s)
	}
	return surveys, rows.Err()
}

func (r *surveyRepo) ListDueClosure(ctx context.Context, now time.Time) ([]*domain.Survey, error) {
	q := `SELECT ` + surveySelectCols + ` FROM surveys
		WHERE status = $1 AND end_date IS NOT NULL AND end_date <= $2 AND deleted_at IS NULL
		ORDER BY end_date ASC`
	rows, err := r.pool.Query(ctx, q, domain.SurveyStatusActive, now)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var surveys []*domain.Survey
	for rows.Next() {
		s, err := scanSurvey(rows)
		if err != nil {
			return nil, err
		}
		surveys = append(surveys, s)
	}
	return surveys, rows.Err()
}

func (r *surveyRepo) Create(ctx context.Context, s *domain.Survey) error {
	if s.ID == uuid.Nil {
		s.ID = newID()
	}
	const q = `
		INSERT INTO surveys (id, venue_id, external_id, title, content, start_date, end_date,
		                     status, publish_type, is_forced, source, app, segment_filters, created_by)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)
		RETURNING created_at, updated_at`

	return r.pool.QueryRow(ctx, q,
		s.ID, s.VenueID, nullStr(s.ExternalID), nullStr(s.Title), nullStr(s.Content),
		s.StartDate, s.EndDate, s.Status, s.PublishType, s.IsForced,
		s.Source, s.App, jsonOrNil(s.SegmentFilters), s.CreatedBy,
	).Scan(&s.CreatedAt, &s.UpdatedAt)
}

func (r *surveyRepo) Update(ctx context.Context, s *domain.Survey) error {
	const q = `
		UPDATE surveys
		SET external_id=$2, title=$3, content=$4, start_date=$5, end_date=$6,
		    status=$7, publish_type=$8, is_forced=$9, source=$10, app=$11, segment_filters=$12
		WHERE id=$1 AND deleted_at IS NULL
		RETURNING updated_at`

	err := r.pool.QueryRow(ctx, q,
		s.ID, nullStr(s.ExternalID), nullStr(s.Title), nullStr(s.Content),
		s.StartDate, s.EndDate, s.Status, s.PublishType, s.IsForced,
		s.Source, s.App, jsonOrNil(s.SegmentFilters),
	).Scan(&s.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.NewNotFound("survey not found")
	}
	return err
}

func (r *surveyRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return softDelete(ctx, r.pool, "surveys", id.String())
}

// ── Questions ─────────────────────────────────────────────────────────────────

func (r *surveyRepo) FindQuestionByID(ctx context.Context, id uuid.UUID) (*domain.Question, error) {
	const q = `SELECT id, survey_id, question_number, question_type, question_text, is_required, is_other, created_at, updated_at
		FROM questions WHERE id = $1 AND deleted_at IS NULL`
	qn, err := scanQuestion(r.pool.QueryRow(ctx, q, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.NewNotFound("question not found")
	}
	if err != nil {
		return nil, err
	}
	opts, err := r.listOptions(ctx, qn.ID)
	if err != nil {
		return nil, err
	}
	qn.Options = opts
	return qn, nil
}

func (r *surveyRepo) ListQuestions(ctx context.Context, surveyID uuid.UUID) ([]*domain.Question, error) {
	const q = `SELECT id, survey_id, question_number, question_type, question_text, is_required, is_other, created_at, updated_at
		FROM questions WHERE survey_id = $1 AND deleted_at IS NULL ORDER BY question_number`
	rows, err := r.pool.Query(ctx, q, surveyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions []*domain.Question
	for rows.Next() {
		qn, err := scanQuestion(rows)
		if err != nil {
			return nil, err
		}
		opts, err := r.listOptions(ctx, qn.ID)
		if err != nil {
			return nil, err
		}
		qn.Options = opts
		questions = append(questions, qn)
	}
	return questions, rows.Err()
}

func (r *surveyRepo) CreateQuestion(ctx context.Context, q *domain.Question) error {
	if q.ID == uuid.Nil {
		q.ID = newID()
	}
	const sql = `INSERT INTO questions (id, survey_id, question_number, question_type, question_text, is_required, is_other)
		VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING created_at, updated_at`
	return r.pool.QueryRow(ctx, sql,
		q.ID, q.SurveyID, q.QuestionNumber, q.QuestionType, q.QuestionText, q.IsRequired, q.IsOther,
	).Scan(&q.CreatedAt, &q.UpdatedAt)
}

func (r *surveyRepo) UpdateQuestion(ctx context.Context, q *domain.Question) error {
	const sql = `UPDATE questions SET question_number=$2, question_type=$3, question_text=$4, is_required=$5, is_other=$6
		WHERE id=$1 AND deleted_at IS NULL RETURNING updated_at`
	err := r.pool.QueryRow(ctx, sql,
		q.ID, q.QuestionNumber, q.QuestionType, q.QuestionText, q.IsRequired, q.IsOther,
	).Scan(&q.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.NewNotFound("question not found")
	}
	return err
}

func (r *surveyRepo) DeleteQuestion(ctx context.Context, id uuid.UUID) error {
	return softDelete(ctx, r.pool, "questions", id.String())
}

// ── Options ───────────────────────────────────────────────────────────────────

func (r *surveyRepo) listOptions(ctx context.Context, questionID uuid.UUID) ([]*domain.Option, error) {
	const q = `SELECT id, question_id, option_number, option_text, created_at, updated_at
		FROM options WHERE question_id = $1 AND deleted_at IS NULL ORDER BY option_number`
	rows, err := r.pool.Query(ctx, q, questionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var opts []*domain.Option
	for rows.Next() {
		o, err := scanOption(rows)
		if err != nil {
			return nil, err
		}
		opts = append(opts, o)
	}
	return opts, rows.Err()
}

func (r *surveyRepo) CreateOption(ctx context.Context, o *domain.Option) error {
	if o.ID == uuid.Nil {
		o.ID = newID()
	}
	const q = `INSERT INTO options (id, question_id, option_number, option_text)
		VALUES ($1,$2,$3,$4) RETURNING created_at, updated_at`
	return r.pool.QueryRow(ctx, q, o.ID, o.QuestionID, o.OptionNumber, o.OptionText).Scan(&o.CreatedAt, &o.UpdatedAt)
}

func (r *surveyRepo) UpdateOption(ctx context.Context, o *domain.Option) error {
	const q = `UPDATE options SET option_number=$2, option_text=$3 WHERE id=$1 AND deleted_at IS NULL RETURNING updated_at`
	err := r.pool.QueryRow(ctx, q, o.ID, o.OptionNumber, o.OptionText).Scan(&o.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.NewNotFound("option not found")
	}
	return err
}

func (r *surveyRepo) DeleteOption(ctx context.Context, id uuid.UUID) error {
	return softDelete(ctx, r.pool, "options", id.String())
}

// ── Responses ─────────────────────────────────────────────────────────────────

func (r *surveyRepo) ListResponses(ctx context.Context, surveyID uuid.UUID, p domain.Pagination) ([]*domain.SurveyResponse, int64, error) {
	var total int64
	if err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM survey_responses WHERE survey_id = $1 AND deleted_at IS NULL`, surveyID,
	).Scan(&total); err != nil {
		return nil, 0, err
	}

	const q = `SELECT id, survey_id, external_id, submitted_at, created_at FROM survey_responses WHERE survey_id = $1 AND deleted_at IS NULL ORDER BY submitted_at DESC LIMIT $2 OFFSET $3`
	rows, err := r.pool.Query(ctx, q, surveyID, p.PageSize, p.Offset())
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var responses []*domain.SurveyResponse
	for rows.Next() {
		var sr domain.SurveyResponse
		var externalID *string
		if err := rows.Scan(&sr.ID, &sr.SurveyID, &externalID, &sr.SubmittedAt, &sr.CreatedAt); err != nil {
			return nil, 0, err
		}
		derefStr(&sr.ExternalID, externalID)
		responses = append(responses, &sr)
	}
	return responses, total, rows.Err()
}

func (r *surveyRepo) CreateResponse(ctx context.Context, resp *domain.SurveyResponse) error {
	if resp.ID == uuid.Nil {
		resp.ID = newID()
	}
	const q = `INSERT INTO survey_responses (id, survey_id, external_id, submitted_at) VALUES ($1,$2,$3,NOW()) RETURNING submitted_at, created_at`
	if err := r.pool.QueryRow(ctx, q, resp.ID, resp.SurveyID, nullStr(resp.ExternalID)).Scan(&resp.SubmittedAt, &resp.CreatedAt); err != nil {
		if IsUniqueViolation(err) && resp.ExternalID != "" {
			return domain.NewConflict("survey response already exists for external_id")
		}
		return err
	}
	for _, ans := range resp.Answers {
		if ans.ID == uuid.Nil {
			ans.ID = newID()
		}
		ans.ResponseID = resp.ID
		const aq = `INSERT INTO survey_answers (id, response_id, question_id, option_id, answer_text) VALUES ($1,$2,$3,$4,$5)`
		if _, err := r.pool.Exec(ctx, aq, ans.ID, ans.ResponseID, ans.QuestionID, ans.OptionID, nullStr(ans.AnswerText)); err != nil {
			return err
		}
	}
	return nil
}

func (r *surveyRepo) ListActive(ctx context.Context, venueID uuid.UUID, publishTypes []int) ([]*domain.Survey, error) {
	q := `SELECT ` + surveySelectCols + `
		FROM surveys
		WHERE venue_id = $1 AND status = $2 AND publish_type = ANY($3) AND deleted_at IS NULL
		ORDER BY created_at DESC`
	rows, err := r.pool.Query(ctx, q, venueID, domain.SurveyStatusActive, publishTypes)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*domain.Survey
	for rows.Next() {
		s, err := scanSurvey(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

// ── scan helpers ──────────────────────────────────────────────────────────────

func scanSurvey(row scanner) (*domain.Survey, error) {
	var s domain.Survey
	var extID, title, content, app *string
	err := row.Scan(
		&s.ID, &s.VenueID, &extID, &title, &content,
		&s.StartDate, &s.EndDate, &s.Status, &s.PublishType, &s.IsForced,
		&s.Source, &app, &s.SegmentFilters, &s.CreatedBy, &s.CreatedAt, &s.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	derefStr(&s.ExternalID, extID)
	derefStr(&s.Title, title)
	derefStr(&s.Content, content)
	if app != nil {
		s.App = *app
	}
	return &s, nil
}

func scanQuestion(row scanner) (*domain.Question, error) {
	var q domain.Question
	err := row.Scan(
		&q.ID, &q.SurveyID, &q.QuestionNumber, &q.QuestionType,
		&q.QuestionText, &q.IsRequired, &q.IsOther, &q.CreatedAt, &q.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &q, nil
}

func scanOption(row scanner) (*domain.Option, error) {
	var o domain.Option
	err := row.Scan(&o.ID, &o.QuestionID, &o.OptionNumber, &o.OptionText, &o.CreatedAt, &o.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &o, nil
}

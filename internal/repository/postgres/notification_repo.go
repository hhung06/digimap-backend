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

type notificationRepo struct {
	pool *pgxpool.Pool
}

// NewNotificationRepository returns a NotificationRepository backed by PostgreSQL.
func NewNotificationRepository(pool *pgxpool.Pool) repository.NotificationRepository {
	return &notificationRepo{pool: pool}
}

const notifSelectCols = `id, venue_id, survey_id, title, content, topic, type, status,
	send_status, send_type, data, link_url, scheduled_at, target_app,
	segment_filters, device_tokens, error_infos, retry_count, retry_at,
	published_at, created_by, created_at, updated_at`

func (r *notificationRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.Notification, error) {
	q := `SELECT ` + notifSelectCols + ` FROM notifications WHERE id = $1 AND deleted_at IS NULL`
	n, err := scanNotification(r.pool.QueryRow(ctx, q, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.NewNotFound("notification not found")
	}
	return n, err
}

func (r *notificationRepo) List(ctx context.Context, venueID uuid.UUID, p domain.Pagination) ([]*domain.Notification, int64, error) {
	var total int64
	if err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM notifications WHERE venue_id = $1 AND deleted_at IS NULL`, venueID,
	).Scan(&total); err != nil {
		return nil, 0, err
	}

	q := `SELECT ` + notifSelectCols + ` FROM notifications WHERE venue_id = $1 AND deleted_at IS NULL ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	rows, err := r.pool.Query(ctx, q, venueID, p.PageSize, p.Offset())
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var ns []*domain.Notification
	for rows.Next() {
		n, err := scanNotification(rows)
		if err != nil {
			return nil, 0, err
		}
		ns = append(ns, n)
	}
	return ns, total, rows.Err()
}

func (r *notificationRepo) Create(ctx context.Context, n *domain.Notification) error {
	if n.ID == uuid.Nil {
		n.ID = newID()
	}
	const q = `
		INSERT INTO notifications (id, venue_id, survey_id, title, content, topic, type, status,
		                           send_status, send_type, data, link_url, scheduled_at, target_app,
		                           segment_filters, device_tokens, error_infos, retry_count, created_by)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19)
		RETURNING created_at, updated_at`

	return r.pool.QueryRow(ctx, q,
		n.ID, n.VenueID, n.SurveyID, nullStr(n.Title), nullStr(n.Content), nullStr(n.Topic),
		n.Kind, n.Status, n.SendStatus, n.SendType,
		jsonOrNil(n.Data), nullStr(n.LinkURL), n.ScheduledAt, n.TargetApp,
		jsonOrNil(n.SegmentFilters), jsonOrNil(n.DeviceTokens), jsonOrNil(n.ErrorInfos),
		n.RetryCount, n.CreatedBy,
	).Scan(&n.CreatedAt, &n.UpdatedAt)
}

func (r *notificationRepo) Update(ctx context.Context, n *domain.Notification) error {
	const q = `
		UPDATE notifications
		SET survey_id=$2, title=$3, content=$4, topic=$5, type=$6, send_type=$7,
		    data=$8, link_url=$9, scheduled_at=$10, target_app=$11, segment_filters=$12, device_tokens=$13
		WHERE id=$1 AND deleted_at IS NULL
		RETURNING updated_at`

	err := r.pool.QueryRow(ctx, q,
		n.ID, n.SurveyID, nullStr(n.Title), nullStr(n.Content), nullStr(n.Topic),
		n.Kind, n.SendType, jsonOrNil(n.Data), nullStr(n.LinkURL), n.ScheduledAt,
		n.TargetApp, jsonOrNil(n.SegmentFilters), jsonOrNil(n.DeviceTokens),
	).Scan(&n.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.NewNotFound("notification not found")
	}
	return err
}

func (r *notificationRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return softDelete(ctx, r.pool, "notifications", id.String())
}

func (r *notificationRepo) MarkSent(ctx context.Context, id uuid.UUID, publishedAt time.Time) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE notifications SET status=$2, send_status=$3, published_at=$4, updated_at=NOW() WHERE id=$1 AND deleted_at IS NULL`,
		id, domain.NotifStatusSent, domain.NotifSendSuccess, publishedAt,
	)
	return err
}

func (r *notificationRepo) MarkFailed(ctx context.Context, id uuid.UUID, errInfos []byte) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE notifications SET send_status=$2, error_infos=$3, retry_count=retry_count+1, updated_at=NOW() WHERE id=$1 AND deleted_at IS NULL`,
		id, domain.NotifSendFailed, errInfos,
	)
	return err
}

func scanNotification(row scanner) (*domain.Notification, error) {
	var n domain.Notification
	var title, content, topic, linkURL, targetApp *string
	err := row.Scan(
		&n.ID, &n.VenueID, &n.SurveyID, &title, &content, &topic,
		&n.Kind, &n.Status, &n.SendStatus, &n.SendType,
		&n.Data, &linkURL, &n.ScheduledAt, &targetApp,
		&n.SegmentFilters, &n.DeviceTokens, &n.ErrorInfos,
		&n.RetryCount, &n.RetryAt, &n.PublishedAt, &n.CreatedBy,
		&n.CreatedAt, &n.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	derefStr(&n.Title, title)
	derefStr(&n.Content, content)
	derefStr(&n.Topic, topic)
	derefStr(&n.LinkURL, linkURL)
	if targetApp != nil {
		n.TargetApp = *targetApp
	}
	return &n, nil
}

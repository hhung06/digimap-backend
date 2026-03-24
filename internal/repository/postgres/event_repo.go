package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/repository"
)

type eventRepo struct {
	pool *pgxpool.Pool
}

// NewEventRepository returns an EventRepository backed by PostgreSQL.
func NewEventRepository(pool *pgxpool.Pool) repository.EventRepository {
	return &eventRepo{pool: pool}
}

// ── Event tags ────────────────────────────────────────────────────────────────

func (r *eventRepo) FindTagByID(ctx context.Context, id uuid.UUID) (*domain.EventTag, error) {
	const q = `SELECT id, name, localization, created_at, updated_at FROM event_tags WHERE id = $1 AND deleted_at IS NULL`
	t, err := scanEventTag(r.pool.QueryRow(ctx, q, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.NewNotFound("event tag not found")
	}
	return t, err
}

func (r *eventRepo) ListTags(ctx context.Context) ([]*domain.EventTag, error) {
	const q = `SELECT id, name, localization, created_at, updated_at FROM event_tags WHERE deleted_at IS NULL ORDER BY name`
	rows, err := r.pool.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tags []*domain.EventTag
	for rows.Next() {
		t, err := scanEventTag(rows)
		if err != nil {
			return nil, err
		}
		tags = append(tags, t)
	}
	return tags, rows.Err()
}

func (r *eventRepo) CreateTag(ctx context.Context, t *domain.EventTag) error {
	if t.ID == uuid.Nil {
		t.ID = newID()
	}
	const q = `INSERT INTO event_tags (id, name, localization) VALUES ($1, $2, $3) RETURNING created_at, updated_at`
	return r.pool.QueryRow(ctx, q, t.ID, nullStr(t.Name), jsonOrNil(t.Localization)).Scan(&t.CreatedAt, &t.UpdatedAt)
}

func (r *eventRepo) UpdateTag(ctx context.Context, t *domain.EventTag) error {
	const q = `UPDATE event_tags SET name = $2, localization = $3 WHERE id = $1 AND deleted_at IS NULL RETURNING updated_at`
	err := r.pool.QueryRow(ctx, q, t.ID, nullStr(t.Name), jsonOrNil(t.Localization)).Scan(&t.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.NewNotFound("event tag not found")
	}
	return err
}

func (r *eventRepo) DeleteTag(ctx context.Context, id uuid.UUID) error {
	return softDelete(ctx, r.pool, "event_tags", id.String())
}

// ── Event types ───────────────────────────────────────────────────────────────

func (r *eventRepo) FindEventTypeByID(ctx context.Context, id uuid.UUID) (*domain.EventType, error) {
	const q = `SELECT id, venue_id, name, localization, created_at, updated_at FROM event_types WHERE id = $1 AND deleted_at IS NULL`
	t, err := scanEventType(r.pool.QueryRow(ctx, q, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.NewNotFound("event type not found")
	}
	return t, err
}

func (r *eventRepo) ListEventTypes(ctx context.Context, venueID uuid.UUID) ([]*domain.EventType, error) {
	const q = `SELECT id, venue_id, name, localization, created_at, updated_at FROM event_types WHERE venue_id = $1 AND deleted_at IS NULL ORDER BY created_at`
	rows, err := r.pool.Query(ctx, q, venueID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var types []*domain.EventType
	for rows.Next() {
		t, err := scanEventType(rows)
		if err != nil {
			return nil, err
		}
		types = append(types, t)
	}
	return types, rows.Err()
}

func (r *eventRepo) CreateEventType(ctx context.Context, t *domain.EventType) error {
	if t.ID == uuid.Nil {
		t.ID = newID()
	}
	const q = `INSERT INTO event_types (id, venue_id, name, localization) VALUES ($1, $2, $3, $4) RETURNING created_at, updated_at`
	return r.pool.QueryRow(ctx, q, t.ID, t.VenueID, nullStr(t.Name), jsonOrNil(t.Localization)).Scan(&t.CreatedAt, &t.UpdatedAt)
}

func (r *eventRepo) UpdateEventType(ctx context.Context, t *domain.EventType) error {
	const q = `UPDATE event_types SET name = $2, localization = $3 WHERE id = $1 AND deleted_at IS NULL RETURNING updated_at`
	err := r.pool.QueryRow(ctx, q, t.ID, nullStr(t.Name), jsonOrNil(t.Localization)).Scan(&t.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.NewNotFound("event type not found")
	}
	return err
}

func (r *eventRepo) DeleteEventType(ctx context.Context, id uuid.UUID) error {
	return softDelete(ctx, r.pool, "event_types", id.String())
}

// ── Events ────────────────────────────────────────────────────────────────────

const eventSelectCols = `e.id, e.venue_id, e.type_id, e.title, e.description,
	e.banner_image, e.icon_image, e.start_time, e.end_time,
	e.show_start_time, e.show_end_time, e.content_detail, e.content_url,
	e.localization, e.created_at, e.updated_at`

func (r *eventRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.Event, error) {
	q := `SELECT ` + eventSelectCols + ` FROM events e WHERE e.id = $1 AND e.deleted_at IS NULL`
	e, err := scanEvent(r.pool.QueryRow(ctx, q, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.NewNotFound("event not found")
	}
	if err != nil {
		return nil, err
	}
	if err := r.loadEventRelations(ctx, e); err != nil {
		return nil, err
	}
	return e, nil
}

func (r *eventRepo) List(ctx context.Context, venueID uuid.UUID, p domain.Pagination) ([]*domain.Event, int64, error) {
	const cq = `SELECT COUNT(*) FROM events WHERE venue_id = $1 AND deleted_at IS NULL`
	var total int64
	if err := r.pool.QueryRow(ctx, cq, venueID).Scan(&total); err != nil {
		return nil, 0, err
	}

	q := `SELECT ` + eventSelectCols + ` FROM events e WHERE e.venue_id = $1 AND e.deleted_at IS NULL ORDER BY e.created_at DESC LIMIT $2 OFFSET $3`
	rows, err := r.pool.Query(ctx, q, venueID, p.PageSize, p.Offset())
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var events []*domain.Event
	for rows.Next() {
		e, err := scanEvent(rows)
		if err != nil {
			return nil, 0, err
		}
		events = append(events, e)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	return events, total, nil
}

func (r *eventRepo) Create(ctx context.Context, e *domain.Event) error {
	if e.ID == uuid.Nil {
		e.ID = newID()
	}
	const q = `
		INSERT INTO events (id, venue_id, type_id, title, description, banner_image, icon_image,
		                    start_time, end_time, show_start_time, show_end_time,
		                    content_detail, content_url, localization)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)
		RETURNING created_at, updated_at`

	return r.pool.QueryRow(ctx, q,
		e.ID, e.VenueID, e.TypeID, nullStr(e.Title), nullStr(e.Description),
		nullStr(e.BannerImage), nullStr(e.IconImage),
		e.StartTime, e.EndTime, e.ShowStartTime, e.ShowEndTime,
		nullStr(e.ContentDetail), nullStr(e.ContentURL), jsonOrNil(e.Localization),
	).Scan(&e.CreatedAt, &e.UpdatedAt)
}

func (r *eventRepo) Update(ctx context.Context, e *domain.Event) error {
	const q = `
		UPDATE events
		SET type_id=$2, title=$3, description=$4, banner_image=$5, icon_image=$6,
		    start_time=$7, end_time=$8, show_start_time=$9, show_end_time=$10,
		    content_detail=$11, content_url=$12, localization=$13
		WHERE id=$1 AND deleted_at IS NULL
		RETURNING updated_at`

	err := r.pool.QueryRow(ctx, q,
		e.ID, e.TypeID, nullStr(e.Title), nullStr(e.Description),
		nullStr(e.BannerImage), nullStr(e.IconImage),
		e.StartTime, e.EndTime, e.ShowStartTime, e.ShowEndTime,
		nullStr(e.ContentDetail), nullStr(e.ContentURL), jsonOrNil(e.Localization),
	).Scan(&e.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.NewNotFound("event not found")
	}
	return err
}

func (r *eventRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return softDelete(ctx, r.pool, "events", id.String())
}

func (r *eventRepo) SetTags(ctx context.Context, eventID uuid.UUID, tagIDs []uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM event_tag_links WHERE event_id = $1`, eventID)
	if err != nil {
		return err
	}
	for _, tid := range tagIDs {
		if _, err := r.pool.Exec(ctx,
			`INSERT INTO event_tag_links (event_id, tag_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`,
			eventID, tid,
		); err != nil {
			return err
		}
	}
	return nil
}

func (r *eventRepo) SetLocations(ctx context.Context, eventID uuid.UUID, locationIDs []uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM event_location_links WHERE event_id = $1`, eventID)
	if err != nil {
		return err
	}
	for _, lid := range locationIDs {
		if _, err := r.pool.Exec(ctx,
			`INSERT INTO event_location_links (event_id, location_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`,
			eventID, lid,
		); err != nil {
			return err
		}
	}
	return nil
}

// ── Images ────────────────────────────────────────────────────────────────────

func (r *eventRepo) CreateImage(ctx context.Context, img *domain.EventImage) error {
	if img.ID == uuid.Nil {
		img.ID = newID()
	}
	const q = `INSERT INTO event_images (id, event_id, image) VALUES ($1, $2, $3) RETURNING created_at, updated_at`
	return r.pool.QueryRow(ctx, q, img.ID, img.EventID, nullStr(img.Image)).Scan(&img.CreatedAt, &img.UpdatedAt)
}

func (r *eventRepo) DeleteImage(ctx context.Context, id uuid.UUID) error {
	return softDelete(ctx, r.pool, "event_images", id.String())
}

// ── relation loaders ──────────────────────────────────────────────────────────

func (r *eventRepo) loadEventRelations(ctx context.Context, e *domain.Event) error {
	// Load tags
	const tq = `
		SELECT t.id, t.name, t.localization, t.created_at, t.updated_at
		FROM event_tags t
		JOIN event_tag_links l ON l.tag_id = t.id
		WHERE l.event_id = $1 AND t.deleted_at IS NULL`
	tagRows, err := r.pool.Query(ctx, tq, e.ID)
	if err != nil {
		return err
	}
	defer tagRows.Close()
	for tagRows.Next() {
		t, err := scanEventTag(tagRows)
		if err != nil {
			return err
		}
		e.Tags = append(e.Tags, t)
	}
	if err := tagRows.Err(); err != nil {
		return err
	}

	// Load location IDs
	const lq = `SELECT location_id FROM event_location_links WHERE event_id = $1`
	locRows, err := r.pool.Query(ctx, lq, e.ID)
	if err != nil {
		return err
	}
	defer locRows.Close()
	for locRows.Next() {
		var lid uuid.UUID
		if err := locRows.Scan(&lid); err != nil {
			return err
		}
		e.Locations = append(e.Locations, lid)
	}
	if err := locRows.Err(); err != nil {
		return err
	}

	// Load images
	const iq = `SELECT id, event_id, image, created_at, updated_at FROM event_images WHERE event_id = $1 AND deleted_at IS NULL ORDER BY created_at`
	imgRows, err := r.pool.Query(ctx, iq, e.ID)
	if err != nil {
		return err
	}
	defer imgRows.Close()
	for imgRows.Next() {
		img, err := scanEventImage(imgRows)
		if err != nil {
			return err
		}
		e.Images = append(e.Images, img)
	}
	return imgRows.Err()
}

// ── scan helpers ──────────────────────────────────────────────────────────────

func scanEventTag(row scanner) (*domain.EventTag, error) {
	var t domain.EventTag
	var name *string
	err := row.Scan(&t.ID, &name, &t.Localization, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return nil, err
	}
	derefStr(&t.Name, name)
	return &t, nil
}

func scanEventType(row scanner) (*domain.EventType, error) {
	var t domain.EventType
	var name *string
	err := row.Scan(&t.ID, &t.VenueID, &name, &t.Localization, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return nil, err
	}
	derefStr(&t.Name, name)
	return &t, nil
}

func scanEvent(row scanner) (*domain.Event, error) {
	var e domain.Event
	var title, description, bannerImage, iconImage, contentDetail, contentURL *string
	err := row.Scan(
		&e.ID, &e.VenueID, &e.TypeID, &title, &description,
		&bannerImage, &iconImage, &e.StartTime, &e.EndTime,
		&e.ShowStartTime, &e.ShowEndTime, &contentDetail, &contentURL,
		&e.Localization, &e.CreatedAt, &e.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	derefStr(&e.Title, title)
	derefStr(&e.Description, description)
	derefStr(&e.BannerImage, bannerImage)
	derefStr(&e.IconImage, iconImage)
	derefStr(&e.ContentDetail, contentDetail)
	derefStr(&e.ContentURL, contentURL)
	return &e, nil
}

func scanEventImage(row scanner) (*domain.EventImage, error) {
	var img domain.EventImage
	var image *string
	err := row.Scan(&img.ID, &img.EventID, &image, &img.CreatedAt, &img.UpdatedAt)
	if err != nil {
		return nil, err
	}
	derefStr(&img.Image, image)
	return &img, nil
}

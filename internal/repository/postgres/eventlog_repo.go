package postgres

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/hhung06/digimap-backend/internal/domain"
)

type eventLogRepository struct{ pool *pgxpool.Pool }

func NewEventLogRepository(pool *pgxpool.Pool) *eventLogRepository {
	return &eventLogRepository{pool: pool}
}

func (r *eventLogRepository) Create(ctx context.Context, e *domain.EventLog) error {
	e.ID = newID()
	params := e.Params
	if params == nil {
		params = json.RawMessage("{}")
	}
	return r.pool.QueryRow(ctx,
		`INSERT INTO event_logs (id,venue_id,name,params,device_id,user_id,user_agent,ip_address)
         VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING created_at`,
		e.ID, e.VenueID, e.Name, params, e.DeviceID, e.UserID, e.UserAgent, e.IPAddress,
	).Scan(&e.CreatedAt)
}

func (r *eventLogRepository) ListByVenue(ctx context.Context, venueID uuid.UUID, p domain.Pagination) ([]*domain.EventLog, int, error) {
	var total int
	if err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM event_logs WHERE venue_id=$1`, venueID,
	).Scan(&total); err != nil {
		return nil, 0, err
	}
	rows, err := r.pool.Query(ctx,
		`SELECT id,venue_id,name,params,device_id,user_id,user_agent,ip_address,created_at
         FROM event_logs WHERE venue_id=$1
         ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
		venueID, p.PageSize, p.Offset())
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var out []*domain.EventLog
	for rows.Next() {
		var e domain.EventLog
		if err := rows.Scan(
			&e.ID, &e.VenueID, &e.Name, &e.Params,
			&e.DeviceID, &e.UserID, &e.UserAgent, &e.IPAddress, &e.CreatedAt,
		); err != nil {
			return nil, 0, err
		}
		out = append(out, &e)
	}
	return out, total, rows.Err()
}

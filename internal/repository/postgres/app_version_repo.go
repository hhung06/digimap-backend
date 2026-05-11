package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/hhung06/digimap-backend/internal/domain"
)

type appVersionRepo struct {
	pool *pgxpool.Pool
}

func NewAppVersionRepository(pool *pgxpool.Pool) *appVersionRepo {
	return &appVersionRepo{pool: pool}
}

func (r *appVersionRepo) Upsert(ctx context.Context, venueID, version uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO app_versions (venue_id, version, updated_at)
		VALUES ($1, $2, NOW())
		ON CONFLICT (venue_id) DO UPDATE SET version = EXCLUDED.version, updated_at = NOW()
	`, venueID, version)
	return err
}

func (r *appVersionRepo) Get(ctx context.Context, venueID uuid.UUID) (*domain.AppVersion, error) {
	av := &domain.AppVersion{}
	err := r.pool.QueryRow(ctx, `
		SELECT id, venue_id, version, updated_at FROM app_versions WHERE venue_id = $1
	`, venueID).Scan(&av.ID, &av.VenueID, &av.Version, &av.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// Auto-create on miss (mirrors Django's update_or_create with defaults).
			av.ID = uuid.New()
			av.VenueID = venueID
			av.Version = uuid.New()
			if insertErr := r.Upsert(ctx, venueID, av.Version); insertErr != nil {
				return nil, insertErr
			}
			return av, nil
		}
		return nil, err
	}
	return av, nil
}

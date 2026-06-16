package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/repository"
)

type locationCategoryRepo struct {
	pool *pgxpool.Pool
}

func NewLocationCategoryRepository(pool *pgxpool.Pool) repository.LocationCategoryRepository {
	return &locationCategoryRepo{pool: pool}
}

func (r *locationCategoryRepo) FindByNameAndVenue(ctx context.Context, venueID uuid.UUID, name, source string) (*domain.LocationCategory, error) {
	const q = `
		SELECT id, venue_id, parent_id, external_id, name, short_name, color, icon, icon_default,
		       sort_index, visible, description, type, image, localization, source,
		       created_at, updated_at, deleted_at
		FROM location_categories
		WHERE venue_id = $1 AND name = $2 AND source = $3 AND deleted_at IS NULL
		LIMIT 1`

	c, err := scanLocationCategory(r.pool.QueryRow(ctx, q, venueID, name, source))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.NewNotFound("location category not found")
	}
	return c, err
}

func (r *locationCategoryRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.LocationCategory, error) {
	const q = `
		SELECT id, venue_id, parent_id, external_id, name, short_name, color, icon, icon_default,
		       sort_index, visible, description, type, image, localization, source,
		       created_at, updated_at, deleted_at
		FROM location_categories WHERE id = $1 AND deleted_at IS NULL`

	c, err := scanLocationCategory(r.pool.QueryRow(ctx, q, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.NewNotFound("location category not found")
	}
	return c, err
}

func (r *locationCategoryRepo) List(ctx context.Context, venueID uuid.UUID) ([]*domain.LocationCategory, error) {
	const q = `
		SELECT id, venue_id, parent_id, external_id, name, short_name, color, icon, icon_default,
		       sort_index, visible, description, type, image, localization, source,
		       created_at, updated_at, deleted_at
		FROM location_categories WHERE venue_id = $1 AND deleted_at IS NULL ORDER BY sort_index`

	rows, err := r.pool.Query(ctx, q, venueID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cats []*domain.LocationCategory
	for rows.Next() {
		c, err := scanLocationCategory(rows)
		if err != nil {
			return nil, err
		}
		cats = append(cats, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	hydrateLocationCategoryRelations(cats)
	return cats, nil
}

func (r *locationCategoryRepo) Create(ctx context.Context, c *domain.LocationCategory) error {
	if c.ID == uuid.Nil {
		c.ID = newID()
	}
	const q = `
		INSERT INTO location_categories (
			id, venue_id, parent_id, external_id, name, short_name, color, icon, icon_default,
			sort_index, visible, description, type, image, localization, source
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16)
		RETURNING created_at, updated_at`

	return r.pool.QueryRow(ctx, q,
		c.ID, c.VenueID, c.ParentID, nullStr(c.ExternalID), nullStr(c.Name), nullStr(c.ShortName),
		nullStr(c.Color), nullStr(c.Icon), nullStr(c.IconDefault),
		c.SortIndex, c.Visible, nullStr(c.Description), nullStr(c.Type),
		nullStr(c.Image), jsonOrNil(c.Localization), c.Source,
	).Scan(&c.CreatedAt, &c.UpdatedAt)
}

func (r *locationCategoryRepo) Update(ctx context.Context, c *domain.LocationCategory) error {
	const q = `
		UPDATE location_categories
		SET parent_id=$2, external_id=$3, name=$4, short_name=$5, color=$6, icon=$7, icon_default=$8,
		    sort_index=$9, visible=$10, description=$11, type=$12, image=$13,
		    localization=$14, source=$15
		WHERE id=$1 AND deleted_at IS NULL
		RETURNING updated_at`

	err := r.pool.QueryRow(ctx, q,
		c.ID, c.ParentID, nullStr(c.ExternalID), nullStr(c.Name), nullStr(c.ShortName),
		nullStr(c.Color), nullStr(c.Icon), nullStr(c.IconDefault),
		c.SortIndex, c.Visible, nullStr(c.Description), nullStr(c.Type),
		nullStr(c.Image), jsonOrNil(c.Localization), c.Source,
	).Scan(&c.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.NewNotFound("location category not found")
	}
	return err
}

func (r *locationCategoryRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return softDelete(ctx, r.pool, "location_categories", id.String())
}

func scanLocationCategory(row pgx.Row) (*domain.LocationCategory, error) {
	var c domain.LocationCategory
	var (
		extID, name, shortName, color, icon, iconDefault *string
		desc, typ, image                                 *string
		localization                                     []byte
		deletedAt                                        *time.Time
	)
	err := row.Scan(
		&c.ID, &c.VenueID, &c.ParentID, &extID, &name, &shortName, &color, &icon, &iconDefault,
		&c.SortIndex, &c.Visible, &desc, &typ, &image, &localization, &c.Source,
		&c.CreatedAt, &c.UpdatedAt, &deletedAt,
	)
	if err != nil {
		return nil, err
	}
	derefStr(&c.ExternalID, extID)
	derefStr(&c.Name, name)
	derefStr(&c.ShortName, shortName)
	derefStr(&c.Color, color)
	derefStr(&c.Icon, icon)
	derefStr(&c.IconDefault, iconDefault)
	derefStr(&c.Description, desc)
	derefStr(&c.Type, typ)
	derefStr(&c.Image, image)
	c.Localization = json.RawMessage(localization)
	c.DeletedAt = deletedAt
	return &c, nil
}

func hydrateLocationCategoryRelations(cats []*domain.LocationCategory) {
	byID := make(map[uuid.UUID]*domain.LocationCategory, len(cats))
	for _, c := range cats {
		c.Parent = nil
		c.Subcategories = nil
		byID[c.ID] = c
	}
	for _, c := range cats {
		if c.ParentID == nil {
			continue
		}
		parent, ok := byID[*c.ParentID]
		if !ok {
			continue
		}
		c.Parent = parent
		parent.Subcategories = append(parent.Subcategories, c)
	}
}

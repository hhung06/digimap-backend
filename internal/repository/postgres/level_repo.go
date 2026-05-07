package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/platform/database"
	"github.com/hhung06/digimap-backend/internal/repository"
)

type levelRepo struct {
	pool *pgxpool.Pool
}

// NewLevelRepository returns a LevelRepository backed by PostgreSQL.
func NewLevelRepository(pool *pgxpool.Pool) repository.LevelRepository {
	return &levelRepo{pool: pool}
}

// ── Map groups ────────────────────────────────────────────────────────────────

func (r *levelRepo) FindMapGroupByID(ctx context.Context, id uuid.UUID) (*domain.MapGroup, error) {
	const q = `
		SELECT id, venue_id, type, name, short_name, sort_index, created_at, updated_at, deleted_at
		FROM map_groups WHERE id = $1 AND deleted_at IS NULL`

	mg, err := scanMapGroup(r.pool.QueryRow(ctx, q, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.NewNotFound("map group not found")
	}
	return mg, err
}

func (r *levelRepo) ListMapGroups(ctx context.Context, venueID uuid.UUID) ([]*domain.MapGroup, error) {
	const q = `
		SELECT id, venue_id, type, name, short_name, sort_index, created_at, updated_at, deleted_at
		FROM map_groups WHERE venue_id = $1 AND deleted_at IS NULL ORDER BY sort_index`

	rows, err := r.pool.Query(ctx, q, venueID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []*domain.MapGroup
	for rows.Next() {
		mg, err := scanMapGroup(rows)
		if err != nil {
			return nil, err
		}
		groups = append(groups, mg)
	}
	return groups, rows.Err()
}

func (r *levelRepo) CreateMapGroup(ctx context.Context, mg *domain.MapGroup) error {
	if mg.ID == uuid.Nil {
		mg.ID = newID()
	}
	const q = `
		INSERT INTO map_groups (id, venue_id, type, name, short_name, sort_index)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING created_at, updated_at`

	return r.pool.QueryRow(ctx, q,
		mg.ID, mg.VenueID, nullStr(mg.Type), nullStr(mg.Name), nullStr(mg.ShortName), mg.SortIndex,
	).Scan(&mg.CreatedAt, &mg.UpdatedAt)
}

func (r *levelRepo) UpdateMapGroup(ctx context.Context, mg *domain.MapGroup) error {
	const q = `
		UPDATE map_groups SET type=$2, name=$3, short_name=$4, sort_index=$5
		WHERE id=$1 AND deleted_at IS NULL
		RETURNING updated_at`

	err := r.pool.QueryRow(ctx, q,
		mg.ID, nullStr(mg.Type), nullStr(mg.Name), nullStr(mg.ShortName), mg.SortIndex,
	).Scan(&mg.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.NewNotFound("map group not found")
	}
	return err
}

func (r *levelRepo) DeleteMapGroup(ctx context.Context, id uuid.UUID) error {
	return database.WithTransaction(ctx, r.pool, func(ctx context.Context, tx pgx.Tx) error {
		if _, err := tx.Exec(ctx,
			`UPDATE levels SET map_group_id = NULL, updated_at = NOW() WHERE map_group_id = $1 AND deleted_at IS NULL`,
			id,
		); err != nil {
			return err
		}
		_, err := tx.Exec(ctx,
			`UPDATE map_groups SET deleted_at = NOW(), updated_at = NOW() WHERE id = $1 AND deleted_at IS NULL`,
			id,
		)
		return err
	})
}

// ── Levels ────────────────────────────────────────────────────────────────────

func (r *levelRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.Level, error) {
	const q = `
		SELECT l.id, l.venue_id, l.map_group_id, l.perspective_id,
		       l.name, l.short_name, l.external_id, l.type,
		       l.latitude, l.longitude, l.bearing,
		       l.width, l.height, l.scale, l.level_width, l.level_height,
		       l.file_ids, l.elevation, l.is_published,
		       l.created_at, l.updated_at, l.deleted_at
		FROM levels l WHERE l.id = $1 AND l.deleted_at IS NULL`

	l, err := scanLevel(r.pool.QueryRow(ctx, q, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.NewNotFound("level not found")
	}
	return l, err
}

func (r *levelRepo) List(ctx context.Context, venueID uuid.UUID) ([]*domain.Level, error) {
	const q = `
		SELECT id, venue_id, map_group_id, perspective_id,
		       name, short_name, external_id, type,
		       latitude, longitude, bearing,
		       width, height, scale, level_width, level_height,
		       file_ids, elevation, is_published,
		       created_at, updated_at, deleted_at
		FROM levels WHERE venue_id = $1 AND deleted_at IS NULL ORDER BY created_at`

	rows, err := r.pool.Query(ctx, q, venueID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var levels []*domain.Level
	for rows.Next() {
		l, err := scanLevel(rows)
		if err != nil {
			return nil, err
		}
		levels = append(levels, l)
	}
	return levels, rows.Err()
}

func (r *levelRepo) Create(ctx context.Context, l *domain.Level) error {
	if l.ID == uuid.Nil {
		l.ID = newID()
	}
	const q = `
		INSERT INTO levels (
			id, venue_id, map_group_id, perspective_id,
			name, short_name, external_id, type,
			latitude, longitude, bearing, width, height, scale,
			level_width, level_height, file_ids, elevation, is_published
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19)
		RETURNING created_at, updated_at`

	return r.pool.QueryRow(ctx, q,
		l.ID, uuidOrNil(l.VenueID), l.MapGroupID, l.PerspectiveID,
		nullStr(l.Name), nullStr(l.ShortName), nullStr(l.ExternalID), l.Type,
		l.Latitude, l.Longitude, l.Bearing, l.Width, l.Height, l.Scale,
		l.LevelWidth, l.LevelHeight, nullStr(l.FileIDs), l.Elevation, l.IsPublished,
	).Scan(&l.CreatedAt, &l.UpdatedAt)
}

func (r *levelRepo) Update(ctx context.Context, l *domain.Level) error {
	const q = `
		UPDATE levels SET
			map_group_id=$2, perspective_id=$3,
			name=$4, short_name=$5, external_id=$6, type=$7,
			latitude=$8, longitude=$9, bearing=$10,
			width=$11, height=$12, scale=$13, level_width=$14, level_height=$15,
			file_ids=$16, elevation=$17, is_published=$18
		WHERE id=$1 AND deleted_at IS NULL
		RETURNING updated_at`

	err := r.pool.QueryRow(ctx, q,
		l.ID, l.MapGroupID, l.PerspectiveID,
		nullStr(l.Name), nullStr(l.ShortName), nullStr(l.ExternalID), l.Type,
		l.Latitude, l.Longitude, l.Bearing, l.Width, l.Height, l.Scale,
		l.LevelWidth, l.LevelHeight, nullStr(l.FileIDs), l.Elevation, l.IsPublished,
	).Scan(&l.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.NewNotFound("level not found")
	}
	return err
}

func (r *levelRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return softDelete(ctx, r.pool, "levels", id.String())
}

// ── Perspectives ──────────────────────────────────────────────────────────────

func (r *levelRepo) UpsertPerspective(ctx context.Context, p *domain.Perspective) error {
	if p.ID == uuid.Nil {
		p.ID = newID()
	}
	const q = `
		INSERT INTO perspectives (
			id, name, camera_zoom, camera_type, camera_max_zoom, camera_min_zoom,
			camera_target_center_lng, camera_target_center_lat, camera_target_zoom,
			camera_target_bearing, camera_target_pitch
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
		ON CONFLICT (id) DO UPDATE SET
			name=EXCLUDED.name, camera_zoom=EXCLUDED.camera_zoom,
			camera_type=EXCLUDED.camera_type, camera_max_zoom=EXCLUDED.camera_max_zoom,
			camera_min_zoom=EXCLUDED.camera_min_zoom,
			camera_target_center_lng=EXCLUDED.camera_target_center_lng,
			camera_target_center_lat=EXCLUDED.camera_target_center_lat,
			camera_target_zoom=EXCLUDED.camera_target_zoom,
			camera_target_bearing=EXCLUDED.camera_target_bearing,
			camera_target_pitch=EXCLUDED.camera_target_pitch,
			updated_at=NOW()
		RETURNING created_at, updated_at`

	return r.pool.QueryRow(ctx, q,
		p.ID, nullStr(p.Name), p.CameraZoom, p.CameraType, p.CameraMaxZoom, p.CameraMinZoom,
		p.CameraTargetCenterLng, p.CameraTargetCenterLat, p.CameraTargetZoom,
		p.CameraTargetBearing, p.CameraTargetPitch,
	).Scan(&p.CreatedAt, &p.UpdatedAt)
}

// ── Geo references ────────────────────────────────────────────────────────────

func (r *levelRepo) ListGeoReferences(ctx context.Context, levelID uuid.UUID) ([]*domain.GeoReference, error) {
	const q = `
		SELECT id, level_id, control_x, control_y, target_x, target_y, created_at, updated_at, deleted_at
		FROM geo_references WHERE level_id = $1 AND deleted_at IS NULL`

	rows, err := r.pool.Query(ctx, q, levelID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var refs []*domain.GeoReference
	for rows.Next() {
		g, err := scanGeoRef(rows)
		if err != nil {
			return nil, err
		}
		refs = append(refs, g)
	}
	return refs, rows.Err()
}

func (r *levelRepo) CreateGeoReference(ctx context.Context, g *domain.GeoReference) error {
	if g.ID == uuid.Nil {
		g.ID = newID()
	}
	const q = `
		INSERT INTO geo_references (id, level_id, control_x, control_y, target_x, target_y)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING created_at, updated_at`

	return r.pool.QueryRow(ctx, q,
		g.ID, g.LevelID, g.ControlX, g.ControlY, g.TargetX, g.TargetY,
	).Scan(&g.CreatedAt, &g.UpdatedAt)
}

func (r *levelRepo) DeleteGeoReference(ctx context.Context, id uuid.UUID) error {
	return softDelete(ctx, r.pool, "geo_references", id.String())
}

// ── scan helpers ──────────────────────────────────────────────────────────────

func scanMapGroup(row pgx.Row) (*domain.MapGroup, error) {
	var mg domain.MapGroup
	var typ, name, shortName *string
	var deletedAt *time.Time
	err := row.Scan(
		&mg.ID, &mg.VenueID, &typ, &name, &shortName, &mg.SortIndex,
		&mg.CreatedAt, &mg.UpdatedAt, &deletedAt,
	)
	if err != nil {
		return nil, err
	}
	derefStr(&mg.Type, typ)
	derefStr(&mg.Name, name)
	derefStr(&mg.ShortName, shortName)
	mg.DeletedAt = deletedAt
	return &mg, nil
}

func scanLevel(row pgx.Row) (*domain.Level, error) {
	var l domain.Level
	var name, shortName, externalID, fileIDs *string
	var deletedAt *time.Time
	err := row.Scan(
		&l.ID, &l.VenueID, &l.MapGroupID, &l.PerspectiveID,
		&name, &shortName, &externalID, &l.Type,
		&l.Latitude, &l.Longitude, &l.Bearing,
		&l.Width, &l.Height, &l.Scale, &l.LevelWidth, &l.LevelHeight,
		&fileIDs, &l.Elevation, &l.IsPublished,
		&l.CreatedAt, &l.UpdatedAt, &deletedAt,
	)
	if err != nil {
		return nil, err
	}
	derefStr(&l.Name, name)
	derefStr(&l.ShortName, shortName)
	derefStr(&l.ExternalID, externalID)
	derefStr(&l.FileIDs, fileIDs)
	l.DeletedAt = deletedAt
	return &l, nil
}

func scanGeoRef(row pgx.Row) (*domain.GeoReference, error) {
	var g domain.GeoReference
	var deletedAt *time.Time
	err := row.Scan(
		&g.ID, &g.LevelID, &g.ControlX, &g.ControlY, &g.TargetX, &g.TargetY,
		&g.CreatedAt, &g.UpdatedAt, &deletedAt,
	)
	if err != nil {
		return nil, err
	}
	g.DeletedAt = deletedAt
	return &g, nil
}

// uuidOrNil returns nil for a zero UUID (avoids inserting uuid.Nil into FK columns).
func uuidOrNil(id uuid.UUID) interface{} {
	if id == uuid.Nil {
		return nil
	}
	return id
}

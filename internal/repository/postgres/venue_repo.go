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

type venueRepo struct {
	pool *pgxpool.Pool
}

// NewVenueRepository returns a VenueRepository backed by PostgreSQL.
func NewVenueRepository(pool *pgxpool.Pool) repository.VenueRepository {
	return &venueRepo{pool: pool}
}

const venueSelectCols = `
	id, customer_id, name, external_id, type, public_key, private_key,
	address, city, state, country, postal, lat, lng, timezone, telephone, work_hours,
	description,
	theme, plugins, translations, localization, custom_data, app_configs, app_domains,
	sub_domains, seo_title, seo_description, seo_keywords, head_tag, body_tag,
	original_logo, small_logo, medium_logo, large_logo,
	start_at, end_at, created_at, updated_at, deleted_at`

func (r *venueRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.Venue, error) {
	q := `SELECT ` + venueSelectCols + ` FROM venues WHERE id = $1 AND deleted_at IS NULL`
	v, err := scanVenue(r.pool.QueryRow(ctx, q, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.NewNotFound("venue not found")
	}
	return v, err
}

func (r *venueRepo) FindByPublicKey(ctx context.Context, publicKey string) (*domain.Venue, error) {
	q := `SELECT ` + venueSelectCols + ` FROM venues WHERE public_key = $1 AND deleted_at IS NULL`
	v, err := scanVenue(r.pool.QueryRow(ctx, q, publicKey))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.NewNotFound("venue not found")
	}
	return v, err
}

func (r *venueRepo) FindByPrivateKey(ctx context.Context, privateKey string) (*domain.Venue, error) {
	q := `SELECT ` + venueSelectCols + ` FROM venues WHERE private_key = $1 AND deleted_at IS NULL`
	v, err := scanVenue(r.pool.QueryRow(ctx, q, privateKey))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.NewNotFound("venue not found")
	}
	return v, err
}

func (r *venueRepo) List(ctx context.Context, customerID uuid.UUID, p domain.Pagination) ([]*domain.Venue, int64, error) {
	const countQ = `SELECT COUNT(*) FROM venues WHERE customer_id = $1 AND deleted_at IS NULL`
	q := `SELECT ` + venueSelectCols + ` FROM venues WHERE customer_id = $1 AND deleted_at IS NULL ORDER BY created_at DESC LIMIT $2 OFFSET $3`

	var total int64
	if err := r.pool.QueryRow(ctx, countQ, customerID).Scan(&total); err != nil {
		return nil, 0, err
	}
	return r.queryVenues(ctx, q, customerID, p.PageSize, p.Offset(), total)
}

func (r *venueRepo) ListAll(ctx context.Context, p domain.Pagination) ([]*domain.Venue, int64, error) {
	const countQ = `SELECT COUNT(*) FROM venues WHERE deleted_at IS NULL`
	q := `SELECT ` + venueSelectCols + ` FROM venues WHERE deleted_at IS NULL ORDER BY created_at DESC LIMIT $1 OFFSET $2`

	var total int64
	if err := r.pool.QueryRow(ctx, countQ).Scan(&total); err != nil {
		return nil, 0, err
	}
	return r.queryVenues(ctx, q, p.PageSize, p.Offset(), total)
}

func (r *venueRepo) queryVenues(ctx context.Context, q string, args ...any) ([]*domain.Venue, int64, error) {
	// last two args are always the result of count + page params; separate them
	total := args[len(args)-1].(int64)
	queryArgs := args[:len(args)-1]

	rows, err := r.pool.Query(ctx, q, queryArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var venues []*domain.Venue
	for rows.Next() {
		v, err := scanVenue(rows)
		if err != nil {
			return nil, 0, err
		}
		venues = append(venues, v)
	}
	return venues, total, rows.Err()
}

func (r *venueRepo) Create(ctx context.Context, v *domain.Venue) error {
	if v.ID == uuid.Nil {
		v.ID = newID()
	}
	const q = `
		INSERT INTO venues (
			id, customer_id, name, external_id, type, public_key, private_key,
			address, city, state, country, postal, lat, lng, timezone, telephone,
			description,
			theme, plugins, localization, app_configs, app_domains,
			sub_domains, seo_title, seo_description, seo_keywords, head_tag, body_tag,
			original_logo, small_logo, medium_logo, large_logo, start_at, end_at
		) VALUES (
			$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,
			$18,$19,$20,$21,$22,$23,$24,$25,$26,$27,$28,$29,$30,$31,$32,$33,$34
		) RETURNING created_at, updated_at`

	return r.pool.QueryRow(ctx, q,
		v.ID, v.CustomerID, v.Name, nullStr(v.ExternalID), v.Type,
		v.PublicKey, v.PrivateKey,
		nullStr(v.Address), nullStr(v.City), nullStr(v.State), nullStr(v.Country),
		nullStr(v.Postal), v.Lat, v.Lng, v.Timezone,
		nullStr(v.Telephone), nullStr(v.Description),
		jsonOrNil(v.Theme), jsonOrNil(v.Plugins),
		jsonOrNil(v.Localization), jsonOrNil(v.AppConfigs),
		jsonOrNil(v.AppDomains), nullStr(v.SubDomains),
		nullStr(v.SEOTitle), nullStr(v.SEODescription), nullStr(v.SEOKeywords),
		nullStr(v.HeadTag), nullStr(v.BodyTag),
		nullStr(v.OriginalLogo), nullStr(v.SmallLogo), nullStr(v.MediumLogo), nullStr(v.LargeLogo),
		v.StartAt, v.EndAt,
	).Scan(&v.CreatedAt, &v.UpdatedAt)
}

func (r *venueRepo) Update(ctx context.Context, v *domain.Venue) error {
	const q = `
		UPDATE venues SET
			name=$2, external_id=$3, type=$4,
			address=$5, city=$6, state=$7, country=$8, postal=$9, lat=$10, lng=$11,
			timezone=$12, telephone=$13, description=$14,
			theme=$15, plugins=$16, localization=$17,
			app_configs=$18, app_domains=$19, sub_domains=$20,
			seo_title=$21, seo_description=$22, seo_keywords=$23, head_tag=$24, body_tag=$25,
			original_logo=$26, small_logo=$27, medium_logo=$28, large_logo=$29,
			start_at=$30, end_at=$31
		WHERE id=$1 AND deleted_at IS NULL
		RETURNING updated_at`

	err := r.pool.QueryRow(ctx, q,
		v.ID, v.Name, nullStr(v.ExternalID), v.Type,
		nullStr(v.Address), nullStr(v.City), nullStr(v.State), nullStr(v.Country),
		nullStr(v.Postal), v.Lat, v.Lng, v.Timezone,
		nullStr(v.Telephone), nullStr(v.Description),
		jsonOrNil(v.Theme), jsonOrNil(v.Plugins),
		jsonOrNil(v.Localization), jsonOrNil(v.AppConfigs),
		jsonOrNil(v.AppDomains), nullStr(v.SubDomains),
		nullStr(v.SEOTitle), nullStr(v.SEODescription), nullStr(v.SEOKeywords),
		nullStr(v.HeadTag), nullStr(v.BodyTag),
		nullStr(v.OriginalLogo), nullStr(v.SmallLogo), nullStr(v.MediumLogo), nullStr(v.LargeLogo),
		v.StartAt, v.EndAt,
	).Scan(&v.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.NewNotFound("venue not found")
	}
	return err
}

func (r *venueRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return softDelete(ctx, r.pool, "venues", id.String())
}

func (r *venueRepo) UpdateKeys(ctx context.Context, id uuid.UUID, publicKey, privateKey string) error {
	tag, err := r.pool.Exec(ctx,
		`UPDATE venues SET public_key = $2, private_key = $3 WHERE id = $1 AND deleted_at IS NULL`,
		id, publicKey, privateKey,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.NewNotFound("venue not found")
	}
	return nil
}


func (r *venueRepo) GetCustomerID(ctx context.Context, venueID uuid.UUID) (uuid.UUID, error) {
	var customerID uuid.UUID
	err := r.pool.QueryRow(ctx,
		`SELECT customer_id FROM venues WHERE id = $1 AND deleted_at IS NULL`,
		venueID,
	).Scan(&customerID)
	if errors.Is(err, pgx.ErrNoRows) {
		return uuid.Nil, domain.NewNotFound("venue not found")
	}
	return customerID, err
}

// ── helpers ───────────────────────────────────────────────────────────────────

func scanVenue(row pgx.Row) (*domain.Venue, error) {
	var v domain.Venue
	var (
		externalID                                                   *string
		address, city, state, country, postal, timezone, telephone   *string
		workHours, description                                       *string
		theme, plugins, translations, localization, customData       []byte
		appConfigs, appDomains                                       []byte
		subDomains, seoTitle, seoDesc, seoKeywords, headTag, bodyTag *string
		origLogo, smallLogo, mediumLogo, largeLogo                   *string
		startAt, endAt, deletedAt                                    *time.Time
	)

	err := row.Scan(
		&v.ID, &v.CustomerID, &v.Name, &externalID, &v.Type,
		&v.PublicKey, &v.PrivateKey,
		&address, &city, &state, &country, &postal, &v.Lat, &v.Lng,
		&timezone, &telephone, &workHours, &description,
		&theme, &plugins, &translations, &localization, &customData,
		&appConfigs, &appDomains,
		&subDomains, &seoTitle, &seoDesc, &seoKeywords, &headTag, &bodyTag,
		&origLogo, &smallLogo, &mediumLogo, &largeLogo,
		&startAt, &endAt, &v.CreatedAt, &v.UpdatedAt, &deletedAt,
	)
	if err != nil {
		return nil, err
	}

	derefStr(&v.ExternalID, externalID)
	derefStr(&v.Address, address)
	derefStr(&v.City, city)
	derefStr(&v.State, state)
	derefStr(&v.Country, country)
	derefStr(&v.Postal, postal)
	derefStr(&v.Timezone, timezone)
	derefStr(&v.Telephone, telephone)
	derefStr(&v.Description, description)
	derefStr(&v.SubDomains, subDomains)
	derefStr(&v.SEOTitle, seoTitle)
	derefStr(&v.SEODescription, seoDesc)
	derefStr(&v.SEOKeywords, seoKeywords)
	derefStr(&v.HeadTag, headTag)
	derefStr(&v.BodyTag, bodyTag)
	derefStr(&v.OriginalLogo, origLogo)
	derefStr(&v.SmallLogo, smallLogo)
	derefStr(&v.MediumLogo, mediumLogo)
	derefStr(&v.LargeLogo, largeLogo)

	v.Theme = json.RawMessage(theme)
	v.Plugins = json.RawMessage(plugins)
	v.Localization = json.RawMessage(localization)
	v.AppConfigs = json.RawMessage(appConfigs)
	v.AppDomains = json.RawMessage(appDomains)
	v.StartAt = startAt
	v.EndAt = endAt
	v.DeletedAt = deletedAt
	return &v, nil
}

// jsonOrNil returns nil for an empty/null RawMessage so PostgreSQL stores NULL.
func jsonOrNil(b json.RawMessage) interface{} {
	if len(b) == 0 || string(b) == "null" {
		return nil
	}
	return []byte(b)
}

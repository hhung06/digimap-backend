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
	id, customer_id, name, slug, external_id, type, public_key, private_key,
	address, city, state, country, postal, lat, lng, timezone, telephone, work_hours,
	description, is_published,
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
	return r.queryVenues(ctx, q, nil, p.PageSize, p.Offset(), total)
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
			id, customer_id, name, slug, external_id, type, public_key, private_key,
			address, city, state, country, postal, lat, lng, timezone, telephone, work_hours,
			description, is_published,
			theme, plugins, translations, localization, custom_data, app_configs, app_domains,
			sub_domains, seo_title, seo_description, seo_keywords, head_tag, body_tag,
			original_logo, small_logo, medium_logo, large_logo, start_at, end_at
		) VALUES (
			$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,
			$21,$22,$23,$24,$25,$26,$27,$28,$29,$30,$31,$32,$33,$34,$35,$36,$37,$38,$39
		) RETURNING created_at, updated_at`

	return r.pool.QueryRow(ctx, q,
		v.ID, v.CustomerID, v.Name, nullStr(v.Slug), nullStr(v.ExternalID), v.Type,
		v.PublicKey, v.PrivateKey,
		nullStr(v.Address), nullStr(v.City), nullStr(v.State), nullStr(v.Country),
		nullStr(v.Postal), v.Lat, v.Lng, v.Timezone,
		nullStr(v.Telephone), nullStr(v.WorkHours), nullStr(v.Description), v.IsPublished,
		jsonOrNil(v.Theme), jsonOrNil(v.Plugins), jsonOrNil(v.Translations),
		jsonOrNil(v.Localization), jsonOrNil(v.CustomData), jsonOrNil(v.AppConfigs),
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
			name=$2, slug=$3, external_id=$4, type=$5,
			address=$6, city=$7, state=$8, country=$9, postal=$10, lat=$11, lng=$12,
			timezone=$13, telephone=$14, work_hours=$15, description=$16, is_published=$17,
			theme=$18, plugins=$19, translations=$20, localization=$21, custom_data=$22,
			app_configs=$23, app_domains=$24, sub_domains=$25,
			seo_title=$26, seo_description=$27, seo_keywords=$28, head_tag=$29, body_tag=$30,
			original_logo=$31, small_logo=$32, medium_logo=$33, large_logo=$34,
			start_at=$35, end_at=$36
		WHERE id=$1 AND deleted_at IS NULL
		RETURNING updated_at`

	err := r.pool.QueryRow(ctx, q,
		v.ID, v.Name, nullStr(v.Slug), nullStr(v.ExternalID), v.Type,
		nullStr(v.Address), nullStr(v.City), nullStr(v.State), nullStr(v.Country),
		nullStr(v.Postal), v.Lat, v.Lng, v.Timezone,
		nullStr(v.Telephone), nullStr(v.WorkHours), nullStr(v.Description), v.IsPublished,
		jsonOrNil(v.Theme), jsonOrNil(v.Plugins), jsonOrNil(v.Translations),
		jsonOrNil(v.Localization), jsonOrNil(v.CustomData), jsonOrNil(v.AppConfigs),
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

func (r *venueRepo) UpdatePublished(ctx context.Context, id uuid.UUID, published bool) error {
	tag, err := r.pool.Exec(ctx,
		`UPDATE venues SET is_published = $2 WHERE id = $1 AND deleted_at IS NULL`,
		id, published,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.NewNotFound("venue not found")
	}
	return nil
}

// ── helpers ───────────────────────────────────────────────────────────────────

func scanVenue(row pgx.Row) (*domain.Venue, error) {
	var v domain.Venue
	var (
		slug, externalID                                                    *string
		address, city, state, country, postal, timezone, telephone         *string
		workHours, description                                              *string
		theme, plugins, translations, localization, customData             []byte
		appConfigs, appDomains                                              []byte
		subDomains, seoTitle, seoDesc, seoKeywords, headTag, bodyTag       *string
		origLogo, smallLogo, mediumLogo, largeLogo                         *string
		startAt, endAt, deletedAt                                          *time.Time
	)

	err := row.Scan(
		&v.ID, &v.CustomerID, &v.Name, &slug, &externalID, &v.Type,
		&v.PublicKey, &v.PrivateKey,
		&address, &city, &state, &country, &postal, &v.Lat, &v.Lng,
		&timezone, &telephone, &workHours, &description, &v.IsPublished,
		&theme, &plugins, &translations, &localization, &customData,
		&appConfigs, &appDomains,
		&subDomains, &seoTitle, &seoDesc, &seoKeywords, &headTag, &bodyTag,
		&origLogo, &smallLogo, &mediumLogo, &largeLogo,
		&startAt, &endAt, &v.CreatedAt, &v.UpdatedAt, &deletedAt,
	)
	if err != nil {
		return nil, err
	}

	derefStr(&v.Slug, slug)
	derefStr(&v.ExternalID, externalID)
	derefStr(&v.Address, address)
	derefStr(&v.City, city)
	derefStr(&v.State, state)
	derefStr(&v.Country, country)
	derefStr(&v.Postal, postal)
	derefStr(&v.Timezone, timezone)
	derefStr(&v.Telephone, telephone)
	derefStr(&v.WorkHours, workHours)
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
	v.Translations = json.RawMessage(translations)
	v.Localization = json.RawMessage(localization)
	v.CustomData = json.RawMessage(customData)
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

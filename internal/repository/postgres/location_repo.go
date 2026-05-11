package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/hhung06/digimap-backend/internal/domain"
	"github.com/hhung06/digimap-backend/internal/repository"
)

type locationRepo struct {
	pool *pgxpool.Pool
}

func NewLocationRepository(pool *pgxpool.Pool) repository.LocationRepository {
	return &locationRepo{pool: pool}
}

const locationSelectCols = `
	id, venue_id, level_id, main_category_id, external_id,
	common_hidden, common_name, common_short_name, common_description, common_color,
	common_location_type, common_location_sub_type, common_latitude, common_longitude, common_address,
	common_location_state, common_location_state_start_date, common_location_state_end_date,
	common_logo, common_large_logo, common_medium_logo, common_small_logo,
	common_social_website, common_social_twitter, common_social_tiktok,
	common_social_facebook, common_social_instagram,
	common_contact_email, common_contact_phone, common_show_short_name,
	top_logo, top_logo_type, place_work_hours,
	booth_number, booth_event_date, booth_size, booth_services_offered, booth_products_showcased,
	person_full_name, person_job_title,
	room_number, room_department, room_bed_count, room_equipment_details,
	is_top_location, top_location_sort_index, icon_default,
	custom, localization, source, start_time, end_time, is_searchable,
	created_at, updated_at, deleted_at`

func (r *locationRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.Location, error) {
	q := `SELECT ` + locationSelectCols + ` FROM locations WHERE id = $1 AND deleted_at IS NULL`
	l, err := scanLocation(r.pool.QueryRow(ctx, q, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.NewNotFound("location not found")
	}
	if err != nil {
		return nil, err
	}
	// Load categories and images
	cats, err := r.loadCategories(ctx, id)
	if err != nil {
		return nil, err
	}
	l.Categories = cats
	imgs, err := r.ListImages(ctx, id)
	if err != nil {
		return nil, err
	}
	l.Images = imgs
	return l, nil
}

func (r *locationRepo) List(ctx context.Context, venueID uuid.UUID, typeFilter *int, p domain.Pagination) ([]*domain.Location, int64, error) {
	var countQ, q string
	var args []interface{}

	// Base condition
	baseWhere := `venue_id = $1 AND deleted_at IS NULL`
	args = append(args, venueID)

	// Add type filter if provided
	if typeFilter != nil {
		baseWhere += ` AND common_location_type = $2`
		args = append(args, *typeFilter)
	}

	countQ = `SELECT COUNT(*) FROM locations WHERE ` + baseWhere

	q = `SELECT ` + locationSelectCols + ` FROM locations WHERE ` + baseWhere + `
		ORDER BY created_at DESC LIMIT $` + strconv.Itoa(len(args)+1) + ` OFFSET $` + strconv.Itoa(len(args)+2)

	var total int64
	if err := r.pool.QueryRow(ctx, countQ, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	queryArgs := append(args, p.PageSize, p.Offset())
	rows, err := r.pool.Query(ctx, q, queryArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var locations []*domain.Location
	for rows.Next() {
		l, err := scanLocation(rows)
		if err != nil {
			return nil, 0, err
		}
		locations = append(locations, l)
	}
	if rows.Err() != nil {
		return nil, 0, rows.Err()
	}
	return locations, total, nil
}

func (r *locationRepo) Create(ctx context.Context, l *domain.Location) error {
	if l.ID == uuid.Nil {
		l.ID = newID()
	}
	const q = `
		INSERT INTO locations (
			id, venue_id, level_id, main_category_id, external_id,
			common_hidden, common_name, common_short_name, common_description, common_color,
			common_location_type, common_location_sub_type, common_latitude, common_longitude, common_address,
			common_location_state, common_location_state_start_date, common_location_state_end_date,
			common_logo, common_social_website, common_social_twitter, common_social_tiktok,
			common_social_facebook, common_social_instagram,
			common_contact_email, common_contact_phone, common_show_short_name,
			top_logo, top_logo_type, place_work_hours,
			booth_number, booth_event_date, booth_size, booth_services_offered, booth_products_showcased,
			person_full_name, person_job_title,
			room_number, room_department, room_bed_count, room_equipment_details,
			is_top_location, top_location_sort_index, icon_default,
			custom, localization, source, start_time, end_time, is_searchable
		) VALUES (
			$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,
			$20,$21,$22,$23,$24,$25,$26,$27,$28,$29,$30,
			$31,$32,$33,$34,$35,$36,$37,$38,$39,$40,$41,$42,$43,$44,$45,$46,$47,$48,$49,$50
		) RETURNING created_at, updated_at`

	return r.pool.QueryRow(ctx, q,
		l.ID, uuidOrNil(l.VenueID), l.LevelID, l.MainCategoryID, nullStr(l.ExternalID),
		l.CommonHidden, l.CommonName, nullStr(l.CommonShortName), nullStr(l.CommonDescription),
		nullStr(l.CommonColor), l.CommonLocationType, l.CommonLocationSubType,
		l.CommonLatitude, l.CommonLongitude, nullStr(l.CommonAddress),
		l.CommonLocationState, l.CommonLocationStateStartDate, l.CommonLocationStateEndDate,
		nullStr(l.CommonLogo),
		nullStr(l.CommonSocialWebsite), nullStr(l.CommonSocialTwitter),
		nullStr(l.CommonSocialTiktok), nullStr(l.CommonSocialFacebook),
		nullStr(l.CommonSocialInstagram),
		nullStr(l.CommonContactEmail), nullStr(l.CommonContactPhone), l.CommonShowShortName,
		nullStr(l.TopLogo), nullStr(l.TopLogoType), jsonOrNil(l.PlaceWorkHours),
		nullStr(l.BoothNumber), l.BoothEventDate, nullStr(l.BoothSize),
		nullStr(l.BoothServicesOffered), nullStr(l.BoothProductsShowcased),
		nullStr(l.PersonFullName), nullStr(l.PersonJobTitle),
		nullStr(l.RoomNumber), nullStr(l.RoomDepartment),
		l.RoomBedCount, nullStr(l.RoomEquipmentDetails),
		l.IsTopLocation, l.TopLocationSortIndex, nullStr(l.IconDefault),
		jsonOrNil(l.Custom), jsonOrNil(l.Localization), l.Source,
		l.StartTime, l.EndTime, l.IsSearchable,
	).Scan(&l.CreatedAt, &l.UpdatedAt)
}

func (r *locationRepo) Update(ctx context.Context, l *domain.Location) error {
	const q = `
		UPDATE locations SET
			level_id=$2, main_category_id=$3, external_id=$4,
			common_hidden=$5, common_name=$6, common_short_name=$7, common_description=$8,
			common_color=$9, common_location_type=$10, common_location_sub_type=$11,
			common_latitude=$12, common_longitude=$13, common_address=$14,
			common_location_state=$15, common_location_state_start_date=$16,
			common_location_state_end_date=$17,
			common_logo=$18, common_social_website=$19, common_social_twitter=$20,
			common_social_tiktok=$21, common_social_facebook=$22, common_social_instagram=$23,
			common_contact_email=$24, common_contact_phone=$25, common_show_short_name=$26,
			top_logo=$27, top_logo_type=$28, place_work_hours=$29,
			booth_number=$30, booth_event_date=$31, booth_size=$32,
			booth_services_offered=$33, booth_products_showcased=$34,
			person_full_name=$35, person_job_title=$36,
			room_number=$37, room_department=$38, room_bed_count=$39, room_equipment_details=$40,
			is_top_location=$41, top_location_sort_index=$42, icon_default=$43,
			custom=$44, localization=$45, source=$46, start_time=$47, end_time=$48, is_searchable=$49
		WHERE id=$1 AND deleted_at IS NULL
		RETURNING updated_at`

	err := r.pool.QueryRow(ctx, q,
		l.ID, l.LevelID, l.MainCategoryID, nullStr(l.ExternalID),
		l.CommonHidden, l.CommonName, nullStr(l.CommonShortName), nullStr(l.CommonDescription),
		nullStr(l.CommonColor), l.CommonLocationType, l.CommonLocationSubType,
		l.CommonLatitude, l.CommonLongitude, nullStr(l.CommonAddress),
		l.CommonLocationState, l.CommonLocationStateStartDate, l.CommonLocationStateEndDate,
		nullStr(l.CommonLogo),
		nullStr(l.CommonSocialWebsite), nullStr(l.CommonSocialTwitter),
		nullStr(l.CommonSocialTiktok), nullStr(l.CommonSocialFacebook),
		nullStr(l.CommonSocialInstagram),
		nullStr(l.CommonContactEmail), nullStr(l.CommonContactPhone), l.CommonShowShortName,
		nullStr(l.TopLogo), nullStr(l.TopLogoType), jsonOrNil(l.PlaceWorkHours),
		nullStr(l.BoothNumber), l.BoothEventDate, nullStr(l.BoothSize),
		nullStr(l.BoothServicesOffered), nullStr(l.BoothProductsShowcased),
		nullStr(l.PersonFullName), nullStr(l.PersonJobTitle),
		nullStr(l.RoomNumber), nullStr(l.RoomDepartment),
		l.RoomBedCount, nullStr(l.RoomEquipmentDetails),
		l.IsTopLocation, l.TopLocationSortIndex, nullStr(l.IconDefault),
		jsonOrNil(l.Custom), jsonOrNil(l.Localization), l.Source,
		l.StartTime, l.EndTime, l.IsSearchable,
	).Scan(&l.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.NewNotFound("location not found")
	}
	return err
}

func (r *locationRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return softDelete(ctx, r.pool, "locations", id.String())
}

func (r *locationRepo) SetCategories(ctx context.Context, locationID uuid.UUID, categoryIDs []uuid.UUID) error {
	_, err := r.pool.Exec(ctx,
		`DELETE FROM location_category_links WHERE location_id = $1`, locationID)
	if err != nil {
		return err
	}
	for _, catID := range categoryIDs {
		_, err := r.pool.Exec(ctx,
			`INSERT INTO location_category_links (location_id, category_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`,
			locationID, catID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *locationRepo) SetTopLocation(ctx context.Context, id uuid.UUID, isTop bool, sortIndex *int) error {
	tag, err := r.pool.Exec(ctx,
		`UPDATE locations SET is_top_location=$2, top_location_sort_index=$3 WHERE id=$1 AND deleted_at IS NULL`,
		id, isTop, sortIndex,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.NewNotFound("location not found")
	}
	return nil
}

func (r *locationRepo) ListTopLocations(ctx context.Context, venueID uuid.UUID) ([]*domain.Location, error) {
	q := `SELECT ` + locationSelectCols + `
		FROM locations
		WHERE venue_id = $1
		  AND is_top_location = TRUE
		  AND common_location_type != $2
		  AND deleted_at IS NULL
		ORDER BY updated_at DESC`

	rows, err := r.pool.Query(ctx, q, venueID, domain.LocationTypeMemo)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var locations []*domain.Location
	for rows.Next() {
		loc, err := scanLocation(rows)
		if err != nil {
			return nil, err
		}
		locations = append(locations, loc)
	}
	return locations, rows.Err()
}

// ── Images ────────────────────────────────────────────────────────────────────

func (r *locationRepo) ListImages(ctx context.Context, locationID uuid.UUID) ([]*domain.LocationImage, error) {
	const q = `
		SELECT id, location_id, original, small, medium, large, created_at, updated_at, deleted_at
		FROM location_images WHERE location_id = $1 AND deleted_at IS NULL`

	rows, err := r.pool.Query(ctx, q, locationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var imgs []*domain.LocationImage
	for rows.Next() {
		img, err := scanLocationImage(rows)
		if err != nil {
			return nil, err
		}
		imgs = append(imgs, img)
	}
	return imgs, rows.Err()
}

func (r *locationRepo) CreateImage(ctx context.Context, img *domain.LocationImage) error {
	if img.ID == uuid.Nil {
		img.ID = newID()
	}
	const q = `
		INSERT INTO location_images (id, location_id, original, small, medium, large)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING created_at, updated_at`

	return r.pool.QueryRow(ctx, q,
		img.ID, img.LocationID,
		nullStr(img.Original), nullStr(img.Small), nullStr(img.Medium), nullStr(img.Large),
	).Scan(&img.CreatedAt, &img.UpdatedAt)
}

func (r *locationRepo) DeleteImage(ctx context.Context, id uuid.UUID) error {
	return softDelete(ctx, r.pool, "location_images", id.String())
}

func (r *locationRepo) GeoSearch(ctx context.Context, lat, lng, radiusKm float64, venueID *uuid.UUID) ([]*domain.Location, error) {
	const q = `
		WITH haversine AS (
			SELECT *,
				6371 * 2 * ASIN(SQRT(
					POWER(SIN(RADIANS($1 - common_latitude) / 2), 2) +
					COS(RADIANS($1)) * COS(RADIANS(common_latitude)) *
					POWER(SIN(RADIANS($2 - common_longitude) / 2), 2)
				)) AS distance_km
			FROM locations
			WHERE deleted_at IS NULL
			  AND common_latitude IS NOT NULL AND common_longitude IS NOT NULL
			  AND ($4::uuid IS NULL OR venue_id = $4)
		)
		SELECT ` + locationSelectCols + `
		FROM haversine
		WHERE distance_km <= $3
		ORDER BY distance_km
		LIMIT 100`

	rows, err := r.pool.Query(ctx, q, lat, lng, radiusKm, venueID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*domain.Location
	for rows.Next() {
		loc, err := scanLocation(rows)
		if err != nil {
			return nil, err
		}
		results = append(results, loc)
	}
	return results, rows.Err()
}

func (r *locationRepo) FindByExternalID(ctx context.Context, venueID uuid.UUID, externalID string, locationType int) (*domain.Location, error) {
	q := `SELECT ` + locationSelectCols + `
		FROM locations
		WHERE venue_id = $1 AND external_id = $2 AND common_location_type = $3 AND deleted_at IS NULL
		LIMIT 1`

	l, err := scanLocation(r.pool.QueryRow(ctx, q, venueID, externalID, locationType))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.NewNotFound("location not found")
	}
	return l, err
}

func (r *locationRepo) SearchByName(ctx context.Context, venueID uuid.UUID, q string, limit int) ([]*domain.Location, error) {
	const query = `SELECT ` + locationSelectCols + `
		FROM locations
		WHERE venue_id = $1 AND common_name ILIKE $2 AND deleted_at IS NULL
		ORDER BY common_name LIMIT $3`
	rows, err := r.pool.Query(ctx, query, venueID, "%"+q+"%", limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*domain.Location
	for rows.Next() {
		l, err := scanLocation(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, l)
	}
	return out, rows.Err()
}

// ── helpers ───────────────────────────────────────────────────────────────────

func (r *locationRepo) loadCategories(ctx context.Context, locationID uuid.UUID) ([]*domain.LocationCategory, error) {
	const q = `
		SELECT lc.id, lc.venue_id, lc.external_id, lc.name, lc.short_name, lc.color,
		       lc.icon, lc.icon_default, lc.sort_index, lc.visible, lc.description,
		       lc.type, lc.image, lc.localization, lc.source,
		       lc.created_at, lc.updated_at, lc.deleted_at
		FROM location_categories lc
		JOIN location_category_links lcl ON lcl.category_id = lc.id
		WHERE lcl.location_id = $1 AND lc.deleted_at IS NULL`

	rows, err := r.pool.Query(ctx, q, locationID)
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
	return cats, rows.Err()
}

func scanLocation(row pgx.Row) (*domain.Location, error) {
	var l domain.Location
	var (
		extID, shortName, desc, color                                    *string
		addr, logo, largeLogo, medLogo, smallLogo                        *string
		socWeb, socTw, socTk, socFb, socIg, contactEmail, contactPhone   *string
		topLogo, topLogoType, boothNum, boothSize, boothSvcs, boothProds *string
		personName, personTitle, roomNum, roomDept, roomEquip            *string
		iconDefault, source                                              *string
		lat, lng                                                         *float64
		workHours, custom, localization                                  []byte
		deletedAt                                                        *time.Time
	)
	err := row.Scan(
		&l.ID, &l.VenueID, &l.LevelID, &l.MainCategoryID, &extID,
		&l.CommonHidden, &l.CommonName, &shortName, &desc, &color,
		&l.CommonLocationType, &l.CommonLocationSubType,
		&lat, &lng, &addr,
		&l.CommonLocationState, &l.CommonLocationStateStartDate, &l.CommonLocationStateEndDate,
		&logo, &largeLogo, &medLogo, &smallLogo,
		&socWeb, &socTw, &socTk, &socFb, &socIg,
		&contactEmail, &contactPhone, &l.CommonShowShortName,
		&topLogo, &topLogoType, &workHours,
		&boothNum, &l.BoothEventDate, &boothSize, &boothSvcs, &boothProds,
		&personName, &personTitle,
		&roomNum, &roomDept, &l.RoomBedCount, &roomEquip,
		&l.IsTopLocation, &l.TopLocationSortIndex, &iconDefault,
		&custom, &localization, &source, &l.StartTime, &l.EndTime, &l.IsSearchable,
		&l.CreatedAt, &l.UpdatedAt, &deletedAt,
	)
	if err != nil {
		return nil, err
	}
	if lat != nil {
		l.CommonLatitude = *lat
	}
	if lng != nil {
		l.CommonLongitude = *lng
	}
	derefStr(&l.ExternalID, extID)
	derefStr(&l.CommonShortName, shortName)
	derefStr(&l.CommonDescription, desc)
	derefStr(&l.CommonColor, color)
	derefStr(&l.CommonAddress, addr)
	derefStr(&l.CommonLogo, logo)
	derefStr(&l.CommonLargeLogo, largeLogo)
	derefStr(&l.CommonMediumLogo, medLogo)
	derefStr(&l.CommonSmallLogo, smallLogo)
	derefStr(&l.CommonSocialWebsite, socWeb)
	derefStr(&l.CommonSocialTwitter, socTw)
	derefStr(&l.CommonSocialTiktok, socTk)
	derefStr(&l.CommonSocialFacebook, socFb)
	derefStr(&l.CommonSocialInstagram, socIg)
	derefStr(&l.CommonContactEmail, contactEmail)
	derefStr(&l.CommonContactPhone, contactPhone)
	derefStr(&l.TopLogo, topLogo)
	derefStr(&l.TopLogoType, topLogoType)
	derefStr(&l.BoothNumber, boothNum)
	derefStr(&l.BoothSize, boothSize)
	derefStr(&l.BoothServicesOffered, boothSvcs)
	derefStr(&l.BoothProductsShowcased, boothProds)
	derefStr(&l.PersonFullName, personName)
	derefStr(&l.PersonJobTitle, personTitle)
	derefStr(&l.RoomNumber, roomNum)
	derefStr(&l.RoomDepartment, roomDept)
	derefStr(&l.RoomEquipmentDetails, roomEquip)
	derefStr(&l.IconDefault, iconDefault)
	derefStr(&l.Source, source)
	l.PlaceWorkHours = json.RawMessage(workHours)
	l.Custom = json.RawMessage(custom)
	l.Localization = json.RawMessage(localization)
	l.DeletedAt = deletedAt
	return &l, nil
}

// ── Memo helpers ──────────────────────────────────────────────────────────────

func (r *locationRepo) ListMemos(ctx context.Context, venueID uuid.UUID) ([]*domain.Location, error) {
	rows, err := r.pool.Query(ctx, `SELECT `+locationSelectCols+`
		FROM locations
		WHERE venue_id = $1
		  AND common_location_type = $2
		  AND deleted_at IS NULL
		ORDER BY updated_at DESC`,
		venueID, domain.LocationTypeMemo)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var memos []*domain.Location
	for rows.Next() {
		m, err := scanLocation(rows)
		if err != nil {
			return nil, err
		}
		memos = append(memos, m)
	}
	return memos, rows.Err()
}

func (r *locationRepo) FindMemoByID(ctx context.Context, id uuid.UUID) (*domain.Location, error) {
	row := r.pool.QueryRow(ctx, `SELECT `+locationSelectCols+`
		FROM locations
		WHERE id = $1 AND common_location_type = $2 AND deleted_at IS NULL`,
		id, domain.LocationTypeMemo)
	return scanLocation(row)
}

func scanLocationImage(row pgx.Row) (*domain.LocationImage, error) {
	var img domain.LocationImage
	var orig, small, medium, large *string
	var deletedAt *time.Time
	err := row.Scan(
		&img.ID, &img.LocationID, &orig, &small, &medium, &large,
		&img.CreatedAt, &img.UpdatedAt, &deletedAt,
	)
	if err != nil {
		return nil, err
	}
	derefStr(&img.Original, orig)
	derefStr(&img.Small, small)
	derefStr(&img.Medium, medium)
	derefStr(&img.Large, large)
	img.DeletedAt = deletedAt
	return &img, nil
}

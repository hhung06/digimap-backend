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

type amenityRepo struct {
	pool *pgxpool.Pool
}

func NewAmenityRepository(pool *pgxpool.Pool) repository.AmenityRepository {
	return &amenityRepo{pool: pool}
}

const amenitySelectCols = `
	id, common_name, common_short_name, common_description, common_color,
	common_location_type, common_latitude, common_longitude, common_address,
	common_location_state, common_location_state_start_date, common_location_state_end_date,
	common_logo, common_social_website, common_social_twitter, common_social_tiktok,
	common_social_facebook, common_social_instagram, common_contact_email, common_contact_phone,
	place_work_hours, booth_number, booth_event_date, booth_size,
	booth_services_offered, booth_products_showcased,
	person_full_name, person_job_title,
	room_number, room_department, room_bed_count, room_equipment_details,
	localization, created_at, updated_at, deleted_at`

func (r *amenityRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.Amenity, error) {
	q := `SELECT ` + amenitySelectCols + ` FROM amenities WHERE id = $1 AND deleted_at IS NULL`
	a, err := scanAmenity(r.pool.QueryRow(ctx, q, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.NewNotFound("amenity not found")
	}
	return a, err
}

func (r *amenityRepo) List(ctx context.Context, p domain.Pagination) ([]*domain.Amenity, int64, error) {
	const countQ = `SELECT COUNT(*) FROM amenities WHERE deleted_at IS NULL`
	q := `SELECT ` + amenitySelectCols + ` FROM amenities WHERE deleted_at IS NULL ORDER BY created_at DESC LIMIT $1 OFFSET $2`

	var total int64
	if err := r.pool.QueryRow(ctx, countQ).Scan(&total); err != nil {
		return nil, 0, err
	}

	rows, err := r.pool.Query(ctx, q, p.PageSize, p.Offset())
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var amenities []*domain.Amenity
	for rows.Next() {
		a, err := scanAmenity(rows)
		if err != nil {
			return nil, 0, err
		}
		amenities = append(amenities, a)
	}
	return amenities, total, rows.Err()
}

func (r *amenityRepo) ListByVenue(ctx context.Context, venueID uuid.UUID) ([]*domain.Amenity, error) {
	q := `SELECT ` + amenitySelectCols + `
		FROM amenities a
		JOIN venue_amenities va ON va.amenity_id = a.id
		WHERE va.venue_id = $1 AND va.deleted_at IS NULL AND a.deleted_at IS NULL
		ORDER BY a.common_name`

	rows, err := r.pool.Query(ctx, q, venueID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var amenities []*domain.Amenity
	for rows.Next() {
		a, err := scanAmenity(rows)
		if err != nil {
			return nil, err
		}
		amenities = append(amenities, a)
	}
	return amenities, rows.Err()
}

func (r *amenityRepo) Create(ctx context.Context, a *domain.Amenity) error {
	if a.ID == uuid.Nil {
		a.ID = newID()
	}
	const q = `
		INSERT INTO amenities (
			id, common_name, common_short_name, common_description, common_color,
			common_location_type, common_latitude, common_longitude, common_address,
			common_location_state, common_location_state_start_date, common_location_state_end_date,
			common_logo, common_social_website, common_social_twitter, common_social_tiktok,
			common_social_facebook, common_social_instagram, common_contact_email, common_contact_phone,
			place_work_hours, booth_number, booth_event_date, booth_size,
			booth_services_offered, booth_products_showcased,
			person_full_name, person_job_title,
			room_number, room_department, room_bed_count, room_equipment_details,
			localization
		) VALUES (
			$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,
			$17,$18,$19,$20,$21,$22,$23,$24,$25,$26,$27,$28,$29,$30,$31,$32,$33
		) RETURNING created_at, updated_at`

	return r.pool.QueryRow(ctx, q,
		a.ID, a.CommonName, nullStr(a.CommonShortName), nullStr(a.CommonDescription),
		nullStr(a.CommonColor), a.CommonLocationType,
		a.CommonLatitude, a.CommonLongitude, nullStr(a.CommonAddress),
		a.CommonLocationState, a.CommonLocationStateStartDate, a.CommonLocationStateEndDate,
		nullStr(a.CommonLogo),
		nullStr(a.CommonSocialWebsite), nullStr(a.CommonSocialTwitter),
		nullStr(a.CommonSocialTiktok), nullStr(a.CommonSocialFacebook),
		nullStr(a.CommonSocialInstagram),
		nullStr(a.CommonContactEmail), nullStr(a.CommonContactPhone),
		jsonOrNil(a.PlaceWorkHours),
		nullStr(a.BoothNumber), a.BoothEventDate, nullStr(a.BoothSize),
		nullStr(a.BoothServicesOffered), nullStr(a.BoothProductsShowcased),
		nullStr(a.PersonFullName), nullStr(a.PersonJobTitle),
		nullStr(a.RoomNumber), nullStr(a.RoomDepartment),
		a.RoomBedCount, nullStr(a.RoomEquipmentDetails),
		jsonOrNil(a.Localization),
	).Scan(&a.CreatedAt, &a.UpdatedAt)
}

func (r *amenityRepo) Update(ctx context.Context, a *domain.Amenity) error {
	const q = `
		UPDATE amenities SET
			common_name=$2, common_short_name=$3, common_description=$4, common_color=$5,
			common_location_type=$6, common_latitude=$7, common_longitude=$8, common_address=$9,
			common_location_state=$10, common_location_state_start_date=$11, common_location_state_end_date=$12,
			common_logo=$13, common_social_website=$14, common_social_twitter=$15,
			common_social_tiktok=$16, common_social_facebook=$17, common_social_instagram=$18,
			common_contact_email=$19, common_contact_phone=$20,
			place_work_hours=$21, booth_number=$22, booth_event_date=$23, booth_size=$24,
			booth_services_offered=$25, booth_products_showcased=$26,
			person_full_name=$27, person_job_title=$28,
			room_number=$29, room_department=$30, room_bed_count=$31, room_equipment_details=$32,
			localization=$33
		WHERE id=$1 AND deleted_at IS NULL
		RETURNING updated_at`

	err := r.pool.QueryRow(ctx, q,
		a.ID, a.CommonName, nullStr(a.CommonShortName), nullStr(a.CommonDescription),
		nullStr(a.CommonColor), a.CommonLocationType,
		a.CommonLatitude, a.CommonLongitude, nullStr(a.CommonAddress),
		a.CommonLocationState, a.CommonLocationStateStartDate, a.CommonLocationStateEndDate,
		nullStr(a.CommonLogo),
		nullStr(a.CommonSocialWebsite), nullStr(a.CommonSocialTwitter),
		nullStr(a.CommonSocialTiktok), nullStr(a.CommonSocialFacebook),
		nullStr(a.CommonSocialInstagram),
		nullStr(a.CommonContactEmail), nullStr(a.CommonContactPhone),
		jsonOrNil(a.PlaceWorkHours),
		nullStr(a.BoothNumber), a.BoothEventDate, nullStr(a.BoothSize),
		nullStr(a.BoothServicesOffered), nullStr(a.BoothProductsShowcased),
		nullStr(a.PersonFullName), nullStr(a.PersonJobTitle),
		nullStr(a.RoomNumber), nullStr(a.RoomDepartment),
		a.RoomBedCount, nullStr(a.RoomEquipmentDetails),
		jsonOrNil(a.Localization),
	).Scan(&a.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.NewNotFound("amenity not found")
	}
	return err
}

func (r *amenityRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return softDelete(ctx, r.pool, "amenities", id.String())
}

func (r *amenityRepo) LinkToVenue(ctx context.Context, va *domain.VenueAmenity) error {
	if va.ID == uuid.Nil {
		va.ID = newID()
	}
	const q = `
		INSERT INTO venue_amenities (id, venue_id, amenity_id)
		VALUES ($1, $2, $3)
		ON CONFLICT (venue_id, amenity_id) DO UPDATE SET deleted_at = NULL, updated_at = NOW()
		RETURNING created_at, updated_at`

	return r.pool.QueryRow(ctx, q, va.ID, va.VenueID, va.AmenityID).
		Scan(&va.CreatedAt, &va.UpdatedAt)
}

func (r *amenityRepo) UnlinkFromVenue(ctx context.Context, venueID, amenityID uuid.UUID) error {
	tag, err := r.pool.Exec(ctx,
		`UPDATE venue_amenities SET deleted_at = NOW(), updated_at = NOW()
		 WHERE venue_id = $1 AND amenity_id = $2 AND deleted_at IS NULL`,
		venueID, amenityID,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.NewNotFound("venue amenity link not found")
	}
	return nil
}

func scanAmenity(row pgx.Row) (*domain.Amenity, error) {
	var a domain.Amenity
	var (
		shortName, desc, color, logo                                             *string
		addr, socWebsite, socTwitter, socTiktok, socFacebook, socInstagram       *string
		contactEmail, contactPhone                                               *string
		boothNum, boothSize, boothSvcs, boothProds                              *string
		personName, personTitle                                                  *string
		roomNum, roomDept, roomEquip                                             *string
		workHours, localization                                                  []byte
		deletedAt                                                                *time.Time
	)
	err := row.Scan(
		&a.ID, &a.CommonName, &shortName, &desc, &color,
		&a.CommonLocationType, &a.CommonLatitude, &a.CommonLongitude, &addr,
		&a.CommonLocationState, &a.CommonLocationStateStartDate, &a.CommonLocationStateEndDate,
		&logo, &socWebsite, &socTwitter, &socTiktok, &socFacebook, &socInstagram,
		&contactEmail, &contactPhone,
		&workHours, &boothNum, &a.BoothEventDate, &boothSize, &boothSvcs, &boothProds,
		&personName, &personTitle,
		&roomNum, &roomDept, &a.RoomBedCount, &roomEquip,
		&localization, &a.CreatedAt, &a.UpdatedAt, &deletedAt,
	)
	if err != nil {
		return nil, err
	}
	derefStr(&a.CommonShortName, shortName)
	derefStr(&a.CommonDescription, desc)
	derefStr(&a.CommonColor, color)
	derefStr(&a.CommonLogo, logo)
	derefStr(&a.CommonAddress, addr)
	derefStr(&a.CommonSocialWebsite, socWebsite)
	derefStr(&a.CommonSocialTwitter, socTwitter)
	derefStr(&a.CommonSocialTiktok, socTiktok)
	derefStr(&a.CommonSocialFacebook, socFacebook)
	derefStr(&a.CommonSocialInstagram, socInstagram)
	derefStr(&a.CommonContactEmail, contactEmail)
	derefStr(&a.CommonContactPhone, contactPhone)
	derefStr(&a.BoothNumber, boothNum)
	derefStr(&a.BoothSize, boothSize)
	derefStr(&a.BoothServicesOffered, boothSvcs)
	derefStr(&a.BoothProductsShowcased, boothProds)
	derefStr(&a.PersonFullName, personName)
	derefStr(&a.PersonJobTitle, personTitle)
	derefStr(&a.RoomNumber, roomNum)
	derefStr(&a.RoomDepartment, roomDept)
	derefStr(&a.RoomEquipmentDetails, roomEquip)
	a.PlaceWorkHours = json.RawMessage(workHours)
	a.Localization = json.RawMessage(localization)
	a.DeletedAt = deletedAt
	return &a, nil
}

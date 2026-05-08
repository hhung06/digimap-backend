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

type appUserRepo struct {
	pool *pgxpool.Pool
}

func NewAppUserRepository(pool *pgxpool.Pool) repository.AppUserRepository {
	return &appUserRepo{pool: pool}
}

const appUserSelectCols = `
	id, venue_id, external_id, source, type,
	first_name, last_name, first_name_en, last_name_en,
	email, phone, company_name, company_name_en, department,
	position, position_en, section, app, token, token_expiration,
	survey_flag, staff_lead_flag, visitor_type, interests, other_interests,
	is_consented, ip_address, user_agent, business_name,
	created_at, updated_at, deleted_at`

func (r *appUserRepo) FindByToken(ctx context.Context, venueID uuid.UUID, token string) (*domain.AppUser, error) {
	q := `SELECT ` + appUserSelectCols + `
		FROM app_users
		WHERE venue_id = $1 AND token = $2 AND deleted_at IS NULL
		LIMIT 1`
	u, err := scanAppUser(r.pool.QueryRow(ctx, q, venueID, token))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.NewNotFound("app user not found")
	}
	return u, err
}

func (r *appUserRepo) Create(ctx context.Context, u *domain.AppUser) error {
	if u.ID == uuid.Nil {
		u.ID = newID()
	}
	interests := u.Interests
	if interests == nil {
		interests = json.RawMessage("[]")
	}
	const q = `
		INSERT INTO app_users (
			id, venue_id, external_id, source, type,
			first_name, last_name, email, phone, token, token_expiration,
			survey_flag, is_consented, interests, visitor_type, other_interests,
			ip_address, user_agent, business_name
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19)
		RETURNING created_at, updated_at`
	return r.pool.QueryRow(ctx, q,
		u.ID, u.VenueID, nullStr(u.ExternalID), u.Source, u.Type,
		nullStr(u.FirstName), nullStr(u.LastName), nullStr(u.Email), nullStr(u.Phone),
		nullStr(u.Token), u.TokenExpiration,
		u.SurveyFlag, u.IsConsented, interests, u.VisitorType, nullStr(u.OtherInterests),
		nullStr(u.IPAddress), nullStr(u.UserAgent), nullStr(u.BusinessName),
	).Scan(&u.CreatedAt, &u.UpdatedAt)
}

func scanAppUser(row pgx.Row) (*domain.AppUser, error) {
	var u domain.AppUser
	var (
		externalID, firstName, lastName, firstNameEn, lastNameEn *string
		email, phone, companyName, companyNameEn                 *string
		department, position, positionEn, app, token             *string
		otherInterests, ipAddress, userAgent, businessName       *string
		section, visitorType                                     *int
		tokenExpiration, deletedAt                               *time.Time
		interests                                                []byte
	)
	err := row.Scan(
		&u.ID, &u.VenueID, &externalID, &u.Source, &u.Type,
		&firstName, &lastName, &firstNameEn, &lastNameEn,
		&email, &phone, &companyName, &companyNameEn, &department,
		&position, &positionEn, &section, &app, &token, &tokenExpiration,
		&u.SurveyFlag, &u.StaffLeadFlag, &visitorType, &interests, &otherInterests,
		&u.IsConsented, &ipAddress, &userAgent, &businessName,
		&u.CreatedAt, &u.UpdatedAt, &deletedAt,
	)
	if err != nil {
		return nil, err
	}
	derefStr(&u.ExternalID, externalID)
	derefStr(&u.FirstName, firstName)
	derefStr(&u.LastName, lastName)
	derefStr(&u.FirstNameEn, firstNameEn)
	derefStr(&u.LastNameEn, lastNameEn)
	derefStr(&u.Email, email)
	derefStr(&u.Phone, phone)
	derefStr(&u.CompanyName, companyName)
	derefStr(&u.CompanyNameEn, companyNameEn)
	derefStr(&u.Department, department)
	derefStr(&u.Position, position)
	derefStr(&u.PositionEn, positionEn)
	derefStr(&u.App, app)
	derefStr(&u.Token, token)
	derefStr(&u.OtherInterests, otherInterests)
	derefStr(&u.IPAddress, ipAddress)
	derefStr(&u.UserAgent, userAgent)
	derefStr(&u.BusinessName, businessName)
	u.Section = section
	u.VisitorType = visitorType
	u.TokenExpiration = tokenExpiration
	u.Interests = json.RawMessage(interests)
	u.DeletedAt = deletedAt
	return &u, nil
}

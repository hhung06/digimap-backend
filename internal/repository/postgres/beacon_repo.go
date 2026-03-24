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

type beaconRepo struct {
	pool *pgxpool.Pool
}

// NewBeaconRepository returns a BeaconRepository backed by PostgreSQL.
func NewBeaconRepository(pool *pgxpool.Pool) repository.BeaconRepository {
	return &beaconRepo{pool: pool}
}

const beaconSelectCols = `id, venue_id, level_id, element_id, name, hw_id, vendor_key, lot_key,
	uuid_val, mac, radius, battery, position_x, position_y, is_enable,
	major, minor, voltage, tx_power, created_at, updated_at`

func (r *beaconRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.Beacon, error) {
	q := `SELECT ` + beaconSelectCols + ` FROM beacons WHERE id = $1 AND deleted_at IS NULL`
	b, err := scanBeacon(r.pool.QueryRow(ctx, q, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.NewNotFound("beacon not found")
	}
	return b, err
}

func (r *beaconRepo) List(ctx context.Context, venueID uuid.UUID, p domain.Pagination) ([]*domain.Beacon, int64, error) {
	var total int64
	if err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM beacons WHERE venue_id = $1 AND deleted_at IS NULL`, venueID,
	).Scan(&total); err != nil {
		return nil, 0, err
	}

	q := `SELECT ` + beaconSelectCols + ` FROM beacons WHERE venue_id = $1 AND deleted_at IS NULL ORDER BY created_at LIMIT $2 OFFSET $3`
	rows, err := r.pool.Query(ctx, q, venueID, p.PageSize, p.Offset())
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var beacons []*domain.Beacon
	for rows.Next() {
		b, err := scanBeacon(rows)
		if err != nil {
			return nil, 0, err
		}
		beacons = append(beacons, b)
	}
	return beacons, total, rows.Err()
}

func (r *beaconRepo) Create(ctx context.Context, b *domain.Beacon) error {
	if b.ID == uuid.Nil {
		b.ID = newID()
	}
	const q = `
		INSERT INTO beacons (id, venue_id, level_id, element_id, name, hw_id, vendor_key, lot_key,
		                     uuid_val, mac, radius, battery, position_x, position_y, is_enable,
		                     major, minor, voltage, tx_power)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19)
		RETURNING created_at, updated_at`

	return r.pool.QueryRow(ctx, q,
		b.ID, b.VenueID, b.LevelID, b.ElementID,
		nullStr(b.Name), nullStr(b.HwID), nullStr(b.VendorKey), nullStr(b.LotKey),
		nullStr(b.UUIDVal), nullStr(b.MAC),
		b.Radius, b.Battery, b.PositionX, b.PositionY, b.IsEnable,
		b.Major, b.Minor, b.Voltage, b.TxPower,
	).Scan(&b.CreatedAt, &b.UpdatedAt)
}

func (r *beaconRepo) Update(ctx context.Context, b *domain.Beacon) error {
	const q = `
		UPDATE beacons
		SET level_id=$2, element_id=$3, name=$4, hw_id=$5, vendor_key=$6, lot_key=$7,
		    uuid_val=$8, mac=$9, radius=$10, battery=$11, position_x=$12, position_y=$13,
		    is_enable=$14, major=$15, minor=$16, voltage=$17, tx_power=$18
		WHERE id=$1 AND deleted_at IS NULL
		RETURNING updated_at`

	err := r.pool.QueryRow(ctx, q,
		b.ID, b.LevelID, b.ElementID,
		nullStr(b.Name), nullStr(b.HwID), nullStr(b.VendorKey), nullStr(b.LotKey),
		nullStr(b.UUIDVal), nullStr(b.MAC),
		b.Radius, b.Battery, b.PositionX, b.PositionY, b.IsEnable,
		b.Major, b.Minor, b.Voltage, b.TxPower,
	).Scan(&b.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.NewNotFound("beacon not found")
	}
	return err
}

func (r *beaconRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return softDelete(ctx, r.pool, "beacons", id.String())
}

func scanBeacon(row scanner) (*domain.Beacon, error) {
	var b domain.Beacon
	var name, hwID, vendorKey, lotKey, uuidVal, mac *string
	err := row.Scan(
		&b.ID, &b.VenueID, &b.LevelID, &b.ElementID,
		&name, &hwID, &vendorKey, &lotKey, &uuidVal, &mac,
		&b.Radius, &b.Battery, &b.PositionX, &b.PositionY, &b.IsEnable,
		&b.Major, &b.Minor, &b.Voltage, &b.TxPower,
		&b.CreatedAt, &b.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	derefStr(&b.Name, name)
	derefStr(&b.HwID, hwID)
	derefStr(&b.VendorKey, vendorKey)
	derefStr(&b.LotKey, lotKey)
	derefStr(&b.UUIDVal, uuidVal)
	derefStr(&b.MAC, mac)
	return &b, nil
}

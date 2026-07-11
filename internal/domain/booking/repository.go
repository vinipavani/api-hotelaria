package booking

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(database *pgxpool.Pool) *Repository {
	return &Repository{
		db: database,
	}
}

func (r *Repository) Create(ctx context.Context, b *Booking) error {
	query := `
		INSERT INTO bookings (room_id, guest_name, guest_document, check_in_date, status)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, room_id, guest_name, guest_document, check_in_date, check_out_date, status, created_at;
	`

	row := r.db.QueryRow(ctx, query, b.RoomID, b.GuestName, b.GuestDocument, b.CheckInDate, b.Status)
	err := scanBookingRow(row, b)

	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) UpdateCheckOut(ctx context.Context, RoomID int64, checkOutDate string) (*Booking, error) {
	query := `
		UPDATE bookings 
		SET check_out_date = $1, status = 'finalizada'
		WHERE room_id = $2 AND status = 'em_estadia'
		RETURNING id, room_id, guest_name, guest_document, check_in_date, check_out_date, status, created_at;
	`

	var b Booking
	row := r.db.QueryRow(ctx, query, checkOutDate, RoomID)
	err := scanBookingRow(row, &b)

	if err != nil {
		return nil, err
	}

	return &b, nil
}

func (r *Repository) isBookingAvailable(ctx context.Context, RoomID int64) (bool, error) {
	query := `
		SELECT COUNT(*)
		FROM bookings
		WHERE room_id = $1 AND status = $2;	
	`

	var count int
	row := r.db.QueryRow(ctx, query, RoomID, BookingStatusInProgress)
	err := row.Scan(&count)
	if err != nil {
		return false, err
	}

	return count == 0, nil
}

func scanBookingRow(row pgx.Row, b *Booking) error {
	return row.Scan(
		&b.ID,
		&b.RoomID,
		&b.GuestName,
		&b.GuestDocument,
		&b.CheckInDate,
		&b.CheckOutDate,
		&b.Status,
		&b.CreatedAt,
	)
}

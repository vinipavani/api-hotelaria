package booking

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateTestBooking(ctx context.Context, db *pgxpool.Pool, roomID int64, name, doc string, status BookingStatus) (*Booking, error) {
	query := `
		INSERT INTO bookings (room_id, guest_name, guest_document, check_in_date, status) 
		VALUES ($1, $2, $3, $4, $5) 
		RETURNING id, room_id, guest_name, guest_document, check_in_date, check_out_date, status, created_at;
	`

	checkInTime := time.Now().Truncate(time.Second)
	var b *Booking

	row := db.QueryRow(ctx, query, roomID, name, doc, checkInTime, status)
	err := scanBookingRow(row, b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

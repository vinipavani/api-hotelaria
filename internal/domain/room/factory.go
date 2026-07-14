package room

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateTestRoom(ctx context.Context, db *pgxpool.Pool, hotelID int64, number string, roomType RoomType) (*Room, error) {
	query := `
		INSERT INTO rooms (hotel_id, number, type, capacity, per_night_value) 
		VALUES ($1, $2, $3, $4, $5) 
		RETURNING id, hotel_id, number, type, capacity, per_night_value, created_at;
	`
	capacity := 2
	price := 150.00

	if roomType == RoomTypeSuite {
		capacity = 4
		price = 500.00
	}

	var rm Room

	err := db.QueryRow(ctx, query, hotelID, number, roomType, capacity, price).Scan(
		&rm.ID,
		&rm.HotelID,
		&rm.Number,
		&rm.Type,
		&rm.Capacity,
		&rm.PerNightValue,
		&rm.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &rm, nil
}

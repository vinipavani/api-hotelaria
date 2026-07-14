package hotel

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateTestHotel(ctx context.Context, db *pgxpool.Pool, name, city string) (*Hotel, error) {
	query := `
		INSERT INTO hotels (name, city) 
		VALUES ($1, $2) 
		RETURNING id, name, city, created_at;
	`

	var h *Hotel
	row := db.QueryRow(ctx, query, name, city)
	err := scanHotelRow(row, h)
	if err != nil {
		return nil, err
	}

	return h, nil
}

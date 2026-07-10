package room

import (
	"context"
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

func (r *Repository) Insert(ctx context.Context, rm *Room) error {
	query := `
		INSERT INTO rooms (hotel_id, number, type, capacity, per_night_value)
		 
		VALUES (
			$1,
			(SELECT LPAD((COALESCE(MAX(number::INTEGER), 0) + 1)::TEXT, 4, '0') FROM rooms WHERE hotel_id = $1), 
			$2,
			$3, 
			$4
		)
		RETURNING id, number, created_at;
	`

	err := r.db.QueryRow(ctx, query, rm.HotelID, rm.Type, rm.Capacity, rm.PerNightValue).Scan(&rm.ID, &rm.Number, &rm.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}
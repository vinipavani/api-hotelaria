package room

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(database *pgxpool.Pool) *Repository {
	return &Repository{
		db: database,
	}
}

func (r *Repository) FindAll(ctx context.Context, HotelID string) ([]*Room, error) {
	query := `
		Select id, hotel_id, number, type, capacity, per_night_value, created_at
		From rooms
		Where hotel_id = $1
	`
	rows, err := r.db.Query(ctx, query, HotelID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var rooms []*Room

	rooms, err = buildList(rows)
	if err != nil {
		return nil, err
	}

	return rooms, nil
}

func buildList(dbRows pgx.Rows) ([]*Room, error) {
	rooms := []*Room{}

	for dbRows.Next() {
		var r Room

		err := dbRows.Scan(&r.ID, &r.HotelID, &r.Number, &r.Type, &r.Capacity, &r.PerNightValue, &r.CreatedAt)
		if err != nil {
			return nil, err
		}

		rooms = append(rooms, &r)
	}

	return rooms, nil
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
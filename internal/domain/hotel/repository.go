package hotel

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

func (r *Repository) FindAll(ctx context.Context) ([]*Hotel, error) {
	query := `
		SELECT id, name, city, created_at 
		FROM hotels 
		ORDER BY id ASC;
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var hotels []*Hotel
	hotels, err = buildList(rows)
	if err != nil {
		return nil, err
	}

	return hotels, nil
}

func (r *Repository) FindByID(ctx context.Context, HotelID int64) (*Hotel, error) {
	query := `
		SELECT name, city
		FROM hotels
		WHERE id = $1
	`

	var h Hotel
	row := r.db.QueryRow(ctx, query, HotelID)
	err := scanHotelRow(row, &h)
	if err != nil {
		return nil, err
	}

	return &h, nil
}

func (r *Repository) Insert(ctx context.Context, h *Hotel) error {
	query := `
		INSERT INTO hotels (name, city) 
		VALUES ($1, $2) 
		RETURNING id, created_at;
	`

	err := r.db.QueryRow(ctx, query, h.Name, h.City).Scan(&h.ID, &h.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func buildList(dbRows pgx.Rows) ([]*Hotel, error) {
	hotels := []*Hotel{}

	for dbRows.Next() {
		var h Hotel

		err := dbRows.Scan(&h.ID, &h.Name, &h.City, &h.CreatedAt)
		if err != nil {
			return nil, err
		}

		hotels = append(hotels, &h)
	}

	return hotels, nil
}

func scanHotelRow(row pgx.Row, h *Hotel) error {
	return row.Scan(
		&h.ID,
		&h.Name,
		&h.City,
		&h.CreatedAt,
	)
}

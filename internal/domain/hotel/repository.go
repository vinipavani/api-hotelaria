package hotel

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
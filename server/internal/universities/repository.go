package universities

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) List(
	ctx context.Context,
	limit, offset int,
) ([]University, error) {
	query := `
		SELECT id::text, name, country
		FROM universities
		ORDER BY name, id
		LIMIT $1 OFFSET $2
	`

	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query universities: %w", err)
	}
	defer rows.Close()

	universities := make([]University, 0, limit)

	for rows.Next() {
		var university University
		err := rows.Scan(
			&university.ID,
			&university.Name,
			&university.Country,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to read university: %w", err)
		}

		universities = append(universities, university)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate through universities: %w", err)
	}

	return universities, nil
}

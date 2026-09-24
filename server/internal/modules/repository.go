package modules

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
	universityID, search string,
	limit, offset int,
) ([]Module, error) {
	// strpos does substring search
	query := `
		SELECT id::text, university_id::text, code, title
		FROM modules
		WHERE university_id = $1
			AND (
				$2 = ''
				OR strpos(lower(code), lower($2)) > 0
				OR strpos(lower(title), lower($2)) > 0 
			)
			ORDER BY code, id
			LIMIT $3, OFFSET $4
	`

	rows, err := r.pool.Query(
		ctx, query, universityID, search, limit, offset,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to query modules: %w", err)
	}
	defer rows.Close()

	modules := make([]Module, 0, limit)

	for rows.Next() {
		var module Module
		if err := rows.Scan(
			&module.ID,
			&module.UniversityID,
			&module.Code,
			&module.Title,
		); err != nil {
			return nil, fmt.Errorf("failed to read module: %w", err)
		}

		modules = append(modules, module)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate modules: %w", err)
	}

	return modules, nil
}

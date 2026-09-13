package universities

import (
	"context"
	"errors"
)

var ErrInvalidPagination = errors.New(
	"limit must be between 1 to 100, and offset must be non-negative",
)

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) List(
	ctx context.Context,
	limit, offset int,
) ([]University, error) {
	if limit < 1 || limit > 100 || offset < 0 {
		return nil, ErrInvalidPagination
	}
	return s.repository.List(ctx, limit, offset)
}

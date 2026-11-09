package modules

import (
	"context"
	"modmapper/server/internal/universities"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
)

type Service struct {
	repository   *Repository
	universities *universities.Repository
}

func NewService(
	repository *Repository,
	universityRepository *universities.Repository,
) *Service {
	return &Service{
		repository:   repository,
		universities: universityRepository,
	}
}

func (s *Service) List(
	ctx context.Context,
	universityID, search string,
	limit, offset int,
) ([]Module, error) {
	var id pgtype.UUID
	if err := id.Scan(universityID); err != nil || !id.Valid {
		return nil, ErrInvalidUniversityID
	}

	if limit < 1 || limit > 100 || offset < 0 {
		return nil, ErrInvalidPagination
	}

	exists, err := s.universities.Exists(ctx, universityID)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, ErrUniversityNotFound
	}

	return s.repository.List(ctx, universityID, strings.TrimSpace(search), limit, offset)
}

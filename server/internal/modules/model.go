package modules

import "errors"

var (
	ErrInvalidUniversityID = errors.New("invalid university ID")
	ErrInvalidPagination   = errors.New("limit must be between 1 and 100, and offset must be non-negative")
	ErrUniversityNotFound  = errors.New("university not found")
)

type Module struct {
	ID           string `json:"id"`
	UniversityID string `json:"university_id"`
	Code         string `json:"code"`
	Title        string `json:"title"`
}

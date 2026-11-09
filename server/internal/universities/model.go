package universities

import "errors"

var ErrInvalidPagination = errors.New(
	"limit must be between 1 to 100, and offset must be non-negative",
)

type University struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Country string `json:"country"`
}

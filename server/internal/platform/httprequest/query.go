package httprequest

import (
	"net/http"
	"strconv"
)

func ReadIntQuery(
	r *http.Request,
	key string,
	defaultValue int,
) (int, error) {
	value := r.URL.Query().Get(key)
	if value == "" {
		return defaultValue, nil
	}

	return strconv.Atoi(value)
}

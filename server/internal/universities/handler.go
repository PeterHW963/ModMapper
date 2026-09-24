package universities

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /universities", h.List)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	limit, err := readIntQuery(r, "limit", 40)
	if err != nil { // error due to Atoi
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "limit must be an integer",
		})
		return
	}

	offset, err := readIntQuery(r, "offset", 0)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "offset must be an integer",
		})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	universities, err := h.service.List(ctx, limit, offset)

	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidPagination):
			writeJSON(w, http.StatusBadRequest, map[string]string{
				"error": err.Error(),
			})

		case errors.Is(err, context.DeadlineExceeded):
			log.Printf("university lookup timed out: %v", err)
			writeJSON(w, http.StatusServiceUnavailable, map[string]string{
				"error": "university lookup timed out",
			})

		case errors.Is(err, context.Canceled):
			// request cancelled (completed), stop handling
			return

		default:
			log.Printf("failed to list universities: %v", err)
			writeJSON(w, http.StatusInternalServerError, map[string]string{
				"error": "failed to retrieve universities",
			})
		}

		return
	}

	writeJSON(w, http.StatusOK, universities)
}

func readIntQuery(
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

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("failed to write JSON response: %v", err)
	}
}

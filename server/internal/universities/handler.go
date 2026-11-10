package universities

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"modmapper/server/internal/platform/httprequest"
	"modmapper/server/internal/platform/httpresponse"
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
	limit, err := httprequest.ReadIntQuery(r, "limit", 40)
	if err != nil { // error due to Atoi
		httpresponse.WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": "limit must be an integer",
		})
		return
	}

	offset, err := httprequest.ReadIntQuery(r, "offset", 0)
	if err != nil {
		httpresponse.WriteJSON(w, http.StatusBadRequest, map[string]string{
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
			httpresponse.WriteJSON(w, http.StatusBadRequest, map[string]string{
				"error": err.Error(),
			})

		case errors.Is(err, context.DeadlineExceeded):
			log.Printf("university lookup timed out: %v", err)
			httpresponse.WriteJSON(w, http.StatusServiceUnavailable, map[string]string{
				"error": "university lookup timed out",
			})

		case errors.Is(err, context.Canceled):
			// request cancelled (completed), stop handling
			return

		default:
			log.Printf("failed to list universities: %v", err)
			httpresponse.WriteJSON(w, http.StatusInternalServerError, map[string]string{
				"error": "failed to retrieve universities",
			})
		}

		return
	}

	httpresponse.WriteJSON(w, http.StatusOK, universities)
}

package modules

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
	mux.HandleFunc("GET /universities/{id}/modules", h.List)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	limit, err := httprequest.ReadIntQuery(r, "limit", 50)
	if err != nil {
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

	modules, err := h.service.List(
		ctx,
		r.PathValue("id"),
		r.URL.Query().Get("q"),
		limit,
		offset,
	)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidUniversityID), errors.Is(err, ErrInvalidPagination):
			httpresponse.WriteJSON(w, http.StatusBadRequest, map[string]string{
				"error": err.Error(),
			})

		case errors.Is(err, ErrUniversityNotFound):
			httpresponse.WriteJSON(w, http.StatusNotFound, map[string]string{
				"error": err.Error(),
			})

		case errors.Is(err, context.Canceled):
			return

		case errors.Is(err, context.DeadlineExceeded):
			log.Printf("module lookup timed out: %v", err)
			httpresponse.WriteJSON(w, http.StatusServiceUnavailable, map[string]string{
				"error": "module lookupi timed out",
			})

		default:
			log.Printf("failed to list modules: %v", err)
			httpresponse.WriteJSON(w, http.StatusInternalServerError, map[string]string{
				"error": "failed to retrieve modules",
			})
		}
		return
	}

	httpresponse.WriteJSON(w, http.StatusOK, modules)
}

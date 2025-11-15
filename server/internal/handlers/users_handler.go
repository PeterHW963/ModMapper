package handlers

import (
	"modmapper/server/internal/httpx"
	"modmapper/server/internal/store"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UsersHandler struct {
	store store.UsersStore
}

// constructor for user handler that wires dependencies
func NewUsersHandler(store store.UsersStore) *UsersHandler {
	return &UsersHandler{store: store}
}

func (h *UsersHandler) Register(router chi.Router) {
	// function connect the http endpoints to router
	router.Get("/", h.list)
}

func (h *UsersHandler) list(writer http.ResponseWriter, req *http.Request) {
	query := req.URL.Query() // get the query param
	search := query.Get("query")
	limit, _ := strconv.ParseInt(query.Get("limit"), 10, 64)
	skip, _ := strconv.ParseInt(query.Get("skip"), 10, 64)

	users, err := h.store.List(req.Context(), search, limit, skip)
	if err != nil {
		httpx.WriteErr(writer, http.StatusInternalServerError, err)
		return
	}
	httpx.WriteJSON(writer, http.StatusOK, map[string]any{"users": users})
}

func (h *UsersHandler) get(writer http.ResponseWriter, req *http.Request) {
	id, err := primitive.ObjectIDFromHex(chi.URLParam(req, "id")) // get the path
	if err != nil {
		// can't find the path var
		httpx.WriteErr(writer, http.StatusNotFound, err)
		return
	}
	user, err := h.store.GetByID(req.Context(), id)
	if err != nil {
		// can't find the user
		httpx.WriteErr(writer, http.StatusNotFound, err)
		return
	}
	httpx.WriteJSON(writer, http.StatusOK, user)
}

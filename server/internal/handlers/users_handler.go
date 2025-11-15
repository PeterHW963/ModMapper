package handlers

import (
	"encoding/json"
	"modmapper/server/internal/httpx"
	"modmapper/server/internal/models"
	"modmapper/server/internal/store"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson"
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
	router.Get("/{id}", h.get)
	router.Post("/", h.create)
	router.Patch("/{id}", h.update)
	router.Delete("/{id}", h.delete)
}

func (h *UsersHandler) list(writer http.ResponseWriter, req *http.Request) {
	query := req.URL.Query() // get the query param
	search := query.Get("q")
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

func (h *UsersHandler) create(writer http.ResponseWriter, req *http.Request) {
	var inUser = models.User{}
	// .Decode function converts a json string into a go value like a struct
	if err := json.NewDecoder(req.Body).Decode(&inUser); err != nil {
		httpx.WriteErr(writer, http.StatusBadRequest, err)
		return
	}
	outUser, err := h.store.Create(req.Context(), inUser)
	if err != nil {
		httpx.WriteErr(writer, http.StatusInternalServerError, err)
		return
	}
	httpx.WriteJSON(writer, http.StatusCreated, outUser)
}

func (h *UsersHandler) update(writer http.ResponseWriter, req *http.Request) {
	id, err := primitive.ObjectIDFromHex(chi.URLParam(req, "id"))
	if err != nil {
		httpx.WriteErr(writer, http.StatusNotFound, err)
		return
	}
	var patch map[string]any
	if err := json.NewDecoder(req.Body).Decode(&patch); err != nil {
		httpx.WriteErr(writer, http.StatusBadRequest, err)
		return
	}
	out, err := h.store.Update(req.Context(), id, bson.M(patch)) // make bson map
	if err != nil || out == nil {
		httpx.WriteErr(writer, http.StatusNotFound, err)
		return
	}
	httpx.WriteJSON(writer, http.StatusOK, out)
}

func (h *UsersHandler) delete(writer http.ResponseWriter, req *http.Request) {
	id, err := primitive.ObjectIDFromHex(chi.URLParam(req, "id"))
	if err != nil {
		httpx.WriteErr(writer, http.StatusNotFound, err)
		return
	}
	if err := h.store.Delete(req.Context(), id); err != nil {
		httpx.WriteErr(writer, http.StatusNotFound, err)
		return
	}
	writer.WriteHeader(http.StatusNoContent) // 204 completed delete req
}

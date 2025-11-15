package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"go.mongodb.org/mongo-driver/mongo"

	"modmapper/server/internal/httpx"
	"modmapper/server/internal/store"
)

/*
Builds a router, sets up middleware (CORS), creates database-backed stores, creates handlers, and then registers the routes.
*/

type Config struct {
	CORSOrigin string
}

func NewRouter(cfg Config, db *mongo.Database) *chi.Mux {
	router := chi.NewRouter()

	// tell chi, for every req, run this middleware first.
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{cfg.CORSOrigin},
		AllowedMethods: []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type"},
		MaxAge:         300, // preflight (browser-issued OPTIONS) reqs only kept for 300 secs max
	}))

	// Health
	router.Get("/api/healthz", func(writer http.ResponseWriter, _ *http.Request) {
		httpx.WriteJSON(writer, http.StatusOK, map[string]any{"status": "ok"})
	})

	// Stores
	usersStore := store.NewUsersStore(db)

	// Handlers
	usersH := NewUsersHandler(usersStore)

	// Routers
	router.Route("/api/users", usersH.Register)

	return router
}

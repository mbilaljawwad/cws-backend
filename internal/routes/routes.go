package routes

import (
	"cws-backend/internal/database"
	"cws-backend/internal/handlers"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Routes struct {
	AppRouter *chi.Mux
}

// setting up user routes
func setupUserRoutes(dbm *database.DBManager) *chi.Mux {
	router := chi.NewRouter()
	userHandler := handlers.NewUserHandler(dbm)
	router.Get("/", userHandler.GetUsers)
	return router
}

func setupRoutes(dbm *database.DBManager) *chi.Mux {
	router := chi.NewRouter()

	// healthcheck
	router.Get("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// setup user routes
	router.Mount("/users", setupUserRoutes(dbm))
	return router
}

func NewRoutes(dbm *database.DBManager) *Routes {
	return &Routes{
		AppRouter: setupRoutes(dbm),
	}
}

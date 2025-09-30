package routes

import (
	"cws-backend/internal/handlers"
	"cws-backend/pkg/database"
	"net/http"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
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

	// Swagger documentation
	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"), // The url pointing to API definition
	))

	// API routes
	router.Route("/api/v1", func(r chi.Router) {
		// setup user routes
		r.Mount("/users", setupUserRoutes(dbm))
	})

	return router
}

func NewRoutes(dbm *database.DBManager) *Routes {
	return &Routes{
		AppRouter: setupRoutes(dbm),
	}
}

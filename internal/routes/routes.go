package routes

import (
	"cws-backend/pkg/database"
	"net/http"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Routes struct {
	AppRouter *chi.Mux
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
	router.Route("/api/v1", func(apiRouter chi.Router) {
		setupSuperAdminRoutes(apiRouter, dbm)
		setupUserRoutes(apiRouter, dbm)

	})
	return router
}

func NewRoutes(dbm *database.DBManager) *Routes {
	return &Routes{
		AppRouter: setupRoutes(dbm),
	}
}

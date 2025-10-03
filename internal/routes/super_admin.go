package routes

import (
	auth "cws-backend/internal/handlers/super_admin"
	"cws-backend/internal/repository"
	"cws-backend/internal/services"
	"cws-backend/pkg/database"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func setupSuperAdminRoutes(apiRouter chi.Router, dbm *database.DBManager) {
	apiRouter.Route("/super-admin", func(superAdminRouter chi.Router) {
		// setting up super admin chi middleware
		superAdminRouter.Use(middleware.Logger)
		superAdminRouter.Use(middleware.Recoverer)
		superAdminRouter.Use(render.SetContentType(render.ContentTypeJSON))

		// setting up models, services, repository
		repo := repository.NewSuperAdminRepo(dbm.DB)
		service := services.NewSuperAdminService(repo)
		authHandler := auth.NewAuthHandler(service)

		// Setting up super admin public routes
		superAdminRouter.Post("/auth/login", authHandler.Login)

	})
}

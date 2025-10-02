package routes

import (
	"cws-backend/internal/handlers"
	"cws-backend/pkg/database"

	"github.com/go-chi/chi/v5"
)

func setupUserRoutes(apiRouter chi.Router, dbm *database.DBManager) {
	userHandler := handlers.NewUserHandler(dbm)
	apiRouter.Group(func(userRouter chi.Router) {
		userRouter.Get("/users", userHandler.GetUsers)
	})
}

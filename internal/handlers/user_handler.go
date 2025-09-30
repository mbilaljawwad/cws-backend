package handlers

import (
	"cws-backend/pkg/database"
	"net/http"
)

// User represents a user in the system
type User struct {
	ID    int    `json:"id" example:"1"`
	Name  string `json:"name" example:"John Doe"`
	Email string `json:"email" example:"john.doe@example.com"`
}

type UserHandler struct {
	dbm *database.DBManager
}

func NewUserHandler(dbm *database.DBManager) *UserHandler {
	return &UserHandler{
		dbm: dbm,
	}
}

// @Summary Get users
// @Description Get all userss
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {array} User
// @Router /users [get]
func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

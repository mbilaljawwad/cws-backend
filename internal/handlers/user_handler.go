package handlers

import (
	"cws-backend/internal/database"
	"net/http"
)

type UserHandler struct {
	dbm *database.DBManager
}

func NewUserHandler(dbm *database.DBManager) *UserHandler {
	return &UserHandler{
		dbm: dbm,
	}
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

package handlers

import "cws-backend/pkg/database"

type AuthHandler struct {
	dbm *database.DBManager
}

// Set super admin value in request context if login is from super admin
func NewAuthHandler(dbm *database.DBManager) *AuthHandler {
	return &AuthHandler{
		dbm: dbm,
	}
}

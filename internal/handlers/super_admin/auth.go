package auth

import (
	"cws-backend/internal/models"
	"cws-backend/internal/services"
	"cws-backend/internal/utils"
	"log"
	"net/http"

	"github.com/go-chi/render"
)

type AuthHandler struct {
	service *services.SuperAdminService
}

func NewAuthHandler(service *services.SuperAdminService) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

// @Summary Super Admin Login
// @Description super admin login
// @Tags super-admin
// @Param request body models.SuperAdminAuthRequest true "Super Admin User Credentials"
// @Accept json
// @Produce json
// @Success 200 {object} models.SuperAdminUser
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /super-admin/auth/login [POST]
func (a *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var requestBody models.SuperAdminAuthRequest
	if err := render.DecodeJSON(r.Body, &requestBody); err != nil {
		log.Printf("Error decoding JSON: %v", err)
		utils.SendJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	user, err := a.service.GetByEmail(r.Context(), requestBody.Email, requestBody.Password)
	if err != nil {
		log.Printf("Error getting user: %v", err)
		utils.SendJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, user)
}

package models

type SuperAdminAuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SuperAdminUser struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

package repository

import (
	"context"
	"cws-backend/internal/models"
	"database/sql"
	"errors"
	"log"

	"github.com/jmoiron/sqlx"
)

type SuperAdminRepo interface {
	GetByEmail(ctx context.Context, email, password string) (models.SuperAdminUser, error)
}

type superAdminRepo struct {
	db *sqlx.DB
}

func NewSuperAdminRepo(db *sqlx.DB) SuperAdminRepo {
	return &superAdminRepo{
		db: db,
	}
}

func (r *superAdminRepo) GetByEmail(ctx context.Context, email string, password string) (models.SuperAdminUser, error) {
	var user models.SuperAdminUser = models.SuperAdminUser{}
	query := `
		SELECT * FROM super_admin WHERE email = $1 and password = $2
	`
	err := r.db.GetContext(ctx, &user, query, email, password)
	log.Printf("User: %v", user)
	if errors.Is(err, sql.ErrNoRows) {
		return user, errors.New("user not found")
	}
	if err != nil {
		return user, err
	}

	return user, nil
}

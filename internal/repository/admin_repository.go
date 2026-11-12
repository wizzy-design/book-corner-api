package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wizzy-design/book-corner-api/internal/models"
)

func CreateAdmin(pool *pgxpool.Pool, email, passwordHash string) (models.Admin, error){
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		INSERT INTO admins (email, password_hash)
		VALUES (LOWER($1), $2)
		RETURNING id, email, password_hash, created_at
	`

	var admin models.Admin
	err := pool.QueryRow(ctx, query, email, passwordHash).Scan(
		&admin.ID,
		&admin.Email,
		&admin.PasswordHash,
		&admin.CreatedAt,
	)
	if err != nil {
		return models.Admin{}, err
	}

	return admin, nil
}
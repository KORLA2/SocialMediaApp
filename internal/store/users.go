package store

import (
	"context"
	"database/sql"

	"github.com/KORLA2/SocialMedia/models"
)

type UsersStore struct {
	db *sql.DB
}

func (u *UsersStore) Create(ctx context.Context, User *models.User) error {

	query := `INSERT INTO users (email,username,password) 
	VALUES ($1,$2,$3) RETURNING id,created_at `

	if err := u.db.QueryRowContext(
		ctx,
		query,
		User.Email,
		User.Username,
		User.Password,
	).Scan(&User.ID, &User.CreatedAt); err != nil {
		return err
	}

	return nil
}

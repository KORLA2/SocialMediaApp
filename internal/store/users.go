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

func (u *UsersStore) GetUserByID(ctx context.Context, userID int) (*models.User, error) {

	query := ` Select id, email,user_name,password,created_at from users
	where id=$1
	`
	var user models.User

	if err := u.db.QueryRowContext(ctx, query, userID).
		Scan(&user.ID, &user.Email, &user.Username, &user.Password, &user.CreatedAt); err != nil {
		return nil, err
	}

	return &user, nil

}

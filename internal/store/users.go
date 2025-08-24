package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/KORLA2/SocialMedia/models"
)

type UsersStore struct {
	db *sql.DB
}

func (u *UsersStore) Create(ctx context.Context, tx *sql.Tx, User *models.User) error {

	query := `INSERT INTO users (email,username,password) 
	VALUES ($1,$2,$3) RETURNING id,created_at `

	if err := tx.QueryRowContext(
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

func (u *UsersStore) CreateAndInvite(ctx context.Context, user *models.User,
	token string, invitationExp time.Duration) error {

	return withTx(u.db, ctx, user, func(tx *sql.Tx) error {

		if err := u.Create(ctx, tx, user); err != nil {

			return err // RollBack;
		}

		if err := u.CreateUserInvite(ctx, tx, user.ID, token, invitationExp); err != nil {
			return err
		}
		return nil
	})

}

func (u *UsersStore) CreateUserInvite(ctx context.Context, tx *sql.Tx, userID int, token string, invitationExp time.Duration) error {

	query := `Insert into user_invitations (user_id,token,expiry) values($1,$2,$3)`

	_, err := tx.ExecContext(ctx, query, userID, token, time.Now().Add(invitationExp))

	if err != nil {
		return err
	}

	return nil
}

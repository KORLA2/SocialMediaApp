package store

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"log"
	"time"

	"github.com/KORLA2/SocialMedia/models"
)

type UsersStore struct {
	db *sql.DB
}

func (u *UsersStore) create(ctx context.Context, tx *sql.Tx, User *models.User) error {

	query := `INSERT INTO users (email,user_name,password,role_id) 
	VALUES ($1,$2,$3,$4) RETURNING id,created_at `

	if err := tx.QueryRowContext(
		ctx,
		query,
		User.Email,
		User.Username,
		User.Password,
		User.Role.Level,
	).Scan(&User.ID, &User.CreatedAt); err != nil {
		return err
	}

	return nil
}

func (u *UsersStore) GetUserByID(ctx context.Context, userID int) (*models.User, error) {

	query := ` Select id, email,user_name,password,created_at , is_active,r.role_id,r.name,r.level,r.description from users
	join roles r on users.role_id=r.role_id
	where id=$1 and is_active=$2
	`
	var user models.User

	if err := u.db.QueryRowContext(ctx, query, userID, true).
		Scan(&user.ID,
			&user.Email,
			&user.Username,
			&user.Password,
			&user.CreatedAt,
			&user.IsActive,
			&user.Role.ID,
			&user.Role.Name,
			&user.Role.Level,
			&user.Role.Description,
		); err != nil {

		switch err.Error() {

		case "sql: no rows in result set":
			return nil, errors.New("either follower user id is not valid or follower user not verified")

		default:

			return nil, err
		}
	}

	return &user, nil

}

type LoginPayload struct {
	ID       int    `json:"id"`
	Username string `json:"user_name"`
	Password string `json:"password"`
}

func (u *UsersStore) GetUserByUserName(ctx context.Context, UserName string) (*LoginPayload, error) {

	query := ` select id,user_name,password
from users where user_name=$1 and is_active=$2
`
	var user LoginPayload
	if err := u.db.QueryRowContext(ctx, query, UserName, true).
		Scan(&user.ID, &user.Username, &user.Password); err != nil {

		switch err {

		case sql.ErrNoRows:
			return nil, errors.New("either your user id is not valid or you are not verified")

		default:

			return nil, err
		}

	}
	return &user, nil

}

func (u *UsersStore) CreateAndInvite(ctx context.Context, user *models.User,
	token string, invitationExp time.Duration) error {

	return withTx(u.db, ctx, func(tx *sql.Tx) error {

		if err := u.create(ctx, tx, user); err != nil {

			return err // RollBack;
		}

		if err := u.createUserInvite(ctx, tx, user.ID, token, invitationExp); err != nil {
			return err
		}
		return nil
	})

}

func (u *UsersStore) Activate(ctx context.Context, token string) error {

	return withTx(u.db, ctx, func(tx *sql.Tx) error {
		// mail send

		// user ative =1

		userID, err := u.getUserFromToken(ctx, tx, token)
		if err != nil {
			log.Print("Problem in UserToke to ID")
			return err
		}

		if err := u.updateUser(ctx, tx, userID); err != nil {
			return nil
		}
		if err := u.deleteUserInvitation(ctx, tx, userID); err != nil {
			return err
		}
		return nil

	})

}

func (u *UsersStore) deleteUserInvitation(ctx context.Context, tx *sql.Tx, userID int) error {

	query := ` Delete from user_invitations where user_id =$1`

	_, err := tx.ExecContext(ctx, query, userID)

	if err != nil {
		return err
	}
	return nil
}

func (u *UsersStore) updateUser(ctx context.Context, tx *sql.Tx, userID int) error {

	query := `Update users set is_active=$1 where id=$2`

	_, err := tx.ExecContext(ctx, query, true, userID)
	if err != nil {
		return err
	}
	return nil

}
func (u *UsersStore) getUserFromToken(ctx context.Context, tx *sql.Tx, token string) (int, error) {
	hash := sha256.Sum256([]byte(token))
	hashToken := hex.EncodeToString(hash[:])
	log.Print(hashToken)
	query := `select user_id 
		from  user_invitations 
		where token=$1 and expiry>$2
	
	`

	var err error
	var userID int
	if err := tx.QueryRowContext(ctx, query, hashToken, time.Now()).Scan(&userID); err != nil {
		return -1, err
	}

	return userID, err
}

func (u *UsersStore) createUserInvite(ctx context.Context, tx *sql.Tx, userID int, token string, invitationExp time.Duration) error {

	query := `Insert into user_invitations (user_id,token,expiry) values($1,$2,$3)`

	_, err := tx.ExecContext(ctx, query, userID, token, time.Now().Add(invitationExp))

	if err != nil {
		return err
	}

	return nil
}

func (u *UsersStore) Delete(ctx context.Context, userID int) error {

	return withTx(u.db, ctx, func(tx *sql.Tx) error {

		// Delete User and Invitation

		if err := u.deleteUser(ctx, tx, userID); err != nil {
			return err
		}

		if err := u.deleteUserInvitation(ctx, tx, userID); err != nil {
			return err
		}

		return nil
	})

}

func (u *UsersStore) deleteUser(ctx context.Context, tx *sql.Tx, userID int) error {

	query := ` delete users where user_id=-$1`

	if _, err := tx.ExecContext(ctx, query, userID); err != nil {
		return err
	}

	return nil
}

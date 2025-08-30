package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/KORLA2/SocialMedia/models"
)

type Storage struct {
	Posts interface {
		Create(context.Context, *models.Post) error
		GetPostByID(context.Context, int) (*models.Post, error)
		DeletePostByID(context.Context, int) error
		UpdatePostByID(context.Context, *models.Post) error
		Feed(context.Context, int, PaginatedQuery) ([]models.UserFeed, error)
	}
	Users interface {
		CreateAndInvite(context.Context, *models.User, string, time.Duration) error
		GetUserByID(context.Context, int) (*models.User, error)
		GetUserByUserName(context.Context, string, *models.User) error
		Activate(context.Context, string) error
		Delete(context.Context, int) error
	}
	Comments interface {
		GetCommentsByPostID(context.Context, int) ([]models.Comment, error)
		Create(context.Context, *models.Comment) error
	}
	Followers interface {
		Create(context.Context, int, int) error
		Delete(context.Context, int, int) error
	}
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{
		Posts:     &PostsStore{db},
		Users:     &UsersStore{db},
		Comments:  &CommentsStore{db},
		Followers: &FollowStore{db},
	}

}

func withTx(db *sql.DB, ctx context.Context, fn func(*sql.Tx) error) error {

	tx, err := db.BeginTx(ctx, nil)

	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		_ = tx.Rollback() // RollBack
		return err
	}
	tx.Commit()
	return nil
}

package store

import (
	"context"
	"database/sql"

	"github.com/KORLA2/SocialMedia/models"
)

type Storage struct {
	Posts interface {
		Create(context.Context, *models.Post) error
		GetPostByID(context.Context, int) (*models.Post, error)
		DeletePostByID(context.Context, int) error
		UpdatePostByID(context.Context, *models.Post) error
		Feed(context.Context, int) ([]models.UserFeed, error)
	}
	Users interface {
		Create(context.Context, *models.User) error
		GetUserByID(context.Context, int) (*models.User, error)
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

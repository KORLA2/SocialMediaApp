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
	}
	Users interface {
		Create(context.Context, *models.User) error
	}
	Comments interface {
		GetCommentsByPostID(context.Context, int) ([]models.Comment, error)
	}
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{
		Posts:    &PostsStore{db},
		Users:    &UsersStore{db},
		Comments: &CommentsStore{db},
	}
}

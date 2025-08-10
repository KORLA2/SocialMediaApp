package store

import (
	"context"
	"database/sql"

	"github.com/KORLA2/SocialMedia/models"
	"github.com/lib/pq"
)

type PostsStore struct {
	db *sql.DB
}

func (s *PostsStore) Create(ctx context.Context, post *models.Post) error {

	query := `INSERT INTO posts (content,user_id,title,tags) 
         VALUES ($1,$2,$3,$4) RETURNING id,created_at,updated_at

`

	if err := s.db.QueryRowContext(
		ctx,
		query,
		post.Content,
		post.User_ID,
		post.Title,
		pq.Array(post.Tags)).Scan(
		&post.ID,
		&post.CreatedAt,
		&post.UpdatedAt,
	); err != nil {
		return err
	}

	return nil
}

func (s *PostsStore) GetPostByID(ctx context.Context, postID int) (*models.Post, error) {

	query := `SELECT id,title,content,tags,created_at,updated_at from posts WHERE id=$1`

	var post models.Post
	if err := s.db.QueryRowContext(
		ctx,
		query,
		postID).Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		pq.Array(&post.Tags),
		&post.CreatedAt,
		&post.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return &post, nil

}

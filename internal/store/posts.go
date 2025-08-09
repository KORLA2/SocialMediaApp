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

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
         VALUES ($1,$2,$3,$4) RETURNING id,created_at,updated_at `

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

func (s *PostsStore) DeletePostByID(ctx context.Context, postID int) error {

	query := `
	DELETE FROM posts WHERE id=$1
	`
	res, err := s.db.ExecContext(ctx, query, postID)

	if err != nil {
		return err
	}
	if rowsAffected, _ := res.RowsAffected(); rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
func (s *PostsStore) UpdatePostByID(ctx context.Context, post *models.Post) error {

	query := `
	Update posts set content=$1, title=$2, updated_at=NOW() WHERE id=$3 RETURNING updated_at
	`

	err := s.db.QueryRowContext(ctx, query, post.Content, post.Title, post.ID).Scan(&post.UpdatedAt)
	if err != nil {

		return err
	}

	return nil
}

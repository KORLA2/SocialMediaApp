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

func (s *PostsStore) Feed(ctx context.Context, userID int, fq PaginatedQuery) ([]models.UserFeed, error) {

	query := `
		select p.id as post_id,u.user_name,p.user_id,p.title,p.content,p.tags,p.created_at, count(c.id) as comments_count,
		
            (CASE WHEN p.title ILIKE '%' || $4 || '%' THEN 1 ELSE 0 END) +
            (CASE WHEN p.content ILIKE '%' || $4 || '%' THEN 1 ELSE 0 END) +
            (SELECT COUNT(*) 
             FROM unnest(p.tags)  t
             WHERE t ILIKE '%' || $4 || '%')
       
                AS match_score

		from posts p 

		left join comments c on c.post_id=p.id 

		left join followers f on p.user_id=f.follower_id and  f.user_id=$1
		join users u on p.user_id=u.id
		where (p.user_id = $1 or f.follower_id is not null) 
		group by p.id,u.user_name
		order by match_score desc,p.created_at ` + fq.Sort + `
		LIMIT $2
		OFFSET $3
`

	rows, err := s.db.QueryContext(ctx, query, userID, fq.Limit, fq.Offset, fq.Search)

	if err != nil {
		return nil, err
	}
	var feed []models.UserFeed
	for rows.Next() {

		var u models.UserFeed
		u.User = models.User{}
		if err := rows.Scan(
			&u.ID,
			&u.User.Username,
			&u.User.ID,
			&u.Title,
			&u.Content,
			pq.Array(&u.Tags),
			&u.CreatedAt,
			&u.Comments_Count,
			&u.PostScore,
		); err != nil {
			return nil, err
		}
		feed = append(feed, u)

	}
	return feed, nil

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

	query := `SELECT id,title,content,user_id,tags,created_at,updated_at from posts WHERE id=$1`

	var post models.Post

	if err := s.db.QueryRowContext(
		ctx,
		query,
		postID).Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		&post.User_ID,
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

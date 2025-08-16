package store

import (
	"context"
	"database/sql"

	"github.com/KORLA2/SocialMedia/models"
)

type CommentsStore struct {
	db *sql.DB
}

func (c *CommentsStore) GetCommentsByPostID(ctx context.Context, postID int) ([]models.Comment, error) {
	query := `
	SELECT c.id,c.post_id,c.user_id,c.content,c.created_at,u.user_name,u.email FROM comments c
	JOIN users u ON c.user_id=u.id
	WHERE c.post_id=$1
	ORDER BY c.created_at DESC
	`
	rows, err := c.db.QueryContext(ctx, query, postID)

	if err != nil {

		return nil, err
	}
	defer rows.Close()
	var comments []models.Comment
	for rows.Next() {
		var comment models.Comment
		comment.User = models.User{}
		if err := rows.Scan(
			&comment.ID,
			&comment.PostID,
			&comment.UserID,
			&comment.Content,
			&comment.CreatedAt,
			&comment.User.Username,
			&comment.User.Email,
		); err != nil {
			return nil, err
		}

		comments = append(comments, comment)

	}
	return comments, nil

}


func (c* CommentsStore)Create(ctx context.Context,comment *models.Comment)error{

	query :=`INSERT INTO comments (post_id,user_id,content)
	 values($1,$2,$3) RETURNING id,created_at
	`
if err:= c.db.QueryRowContext(ctx,
	query,
    comment.PostID,
	comment.UserID,
	comment.Content,
).Scan(
	&comment.ID,
	&comment.CreatedAt,
);err!=nil{
	return err
}

return nil;
}
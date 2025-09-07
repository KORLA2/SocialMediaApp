package store

import (
	"context"
	"database/sql"
)

type FollowStore struct {
	db *sql.DB
}

func (f *FollowStore) Create(ctx context.Context, followerID, userID int) error {

	query := `Insert into followers (user_id,follower_id) 
	values($1,$2)	
	`
	if _, err := f.db.ExecContext(ctx, query, userID, followerID); err != nil {
		return err
	}

	return nil

}

func (f *FollowStore) Delete(ctx context.Context, followerID, userID int) error {

	query := `Delete from followers 
	where follower_id=$1 and user_id=$2
	`

	if _, err := f.db.ExecContext(ctx, query, followerID, userID); err != nil {
		return err
	}
	return nil

}

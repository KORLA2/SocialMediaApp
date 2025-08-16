package models

import "time"

type Post struct {
	ID        int       `josn:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	User_ID   int       `json:"user_id"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Comments  []Comment `json:"comments"`
}

type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

type Comment struct {
	ID        int       `json:"id"`
	PostID    int       `json:"post_id"`
	Content   string    `json:"content"`
	UserID    int       `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	User      User
}

type Follower struct {
	UserID     int       `json:"userid"`
	FollowerID int       `json:"followerid"`
	CreatedAt  time.Time `json:"created_at"`
}

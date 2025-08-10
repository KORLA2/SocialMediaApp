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
}

type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

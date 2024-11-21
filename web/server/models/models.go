package models

import (
	"time"
)

type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password,omitempty"` // Stored as hash
	CreatedAt time.Time `json:"created_at"`
}

type Post struct {
	ID         int64      `json:"id"`
	UserID     int64      `json:"user_id"`
	Title      string     `json:"title"`
	Content    string     `json:"content"`
	CreatedAt  time.Time  `json:"created_at"`
	Username   string     `json:"username"`
	Likes      int        `json:"likes"`
	Dislikes   int        `json:"dislikes"`
	Categories []Category `json:"categories,omitempty"`
}

type Comment struct {
	ID        int64     `json:"id"`
	PostID    int64     `json:"post_id"`
	UserID    int64     `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	Likes     int       `json:"likes"`
	Dislikes  int       `json:"dislikes"`
}

type Category struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}
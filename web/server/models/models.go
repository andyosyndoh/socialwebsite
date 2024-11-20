package models

import (
	"time"
)

type User struct {
	ID        int64
	Username  string
	Email     string
	Password  string // Stored as hash
	CreatedAt time.Time
}

type Post struct {
	ID         int64
	UserID     int64
	Title      string
	Content    string
	CreatedAt  time.Time
	Categories []Category
	Likes      int
	Dislikes   int
}

type Comment struct {
	ID        int64
	PostID    int64
	UserID    int64
	Content   string
	CreatedAt time.Time
	Likes     int
	Dislikes  int
}

type Category struct {
	ID   int64
	Name string
}
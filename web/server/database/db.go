package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	Conn *sql.DB
}

func NewDatabase() *Database {
	// Open SQLite database connection
	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	return &Database{Conn: db}
}

func (db *Database) InitializeTables() error {
	// Create users table
	_, err := db.Conn.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL UNIQUE,
			email TEXT NOT NULL UNIQUE,
			password_hash TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return err
	}

	// Create categories table
	_, err = db.Conn.Exec(`
		CREATE TABLE IF NOT EXISTS categories (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE
		)
	`)
	if err != nil {
		return err
	}

	// Create posts table
	_, err = db.Conn.Exec(`
		CREATE TABLE IF NOT EXISTS posts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			title TEXT NOT NULL,
			content TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY(user_id) REFERENCES users(id)
		)
	`)
	if err != nil {
		return err
	}

	// Create post_categories junction table
	_, err = db.Conn.Exec(`
		CREATE TABLE IF NOT EXISTS post_categories (
			post_id INTEGER,
			category_id INTEGER,
			PRIMARY KEY(post_id, category_id),
			FOREIGN KEY(post_id) REFERENCES posts(id),
			FOREIGN KEY(category_id) REFERENCES categories(id)
		)
	`)
	if err != nil {
		return err
	}

	// Create comments table
	_, err = db.Conn.Exec(`
		CREATE TABLE IF NOT EXISTS comments (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			post_id INTEGER NOT NULL,
			user_id INTEGER NOT NULL,
			content TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY(post_id) REFERENCES posts(id),
			FOREIGN KEY(user_id) REFERENCES users(id)
		)
	`)
	if err != nil {
		return err
	}

	// Create likes table
	_, err = db.Conn.Exec(`
		CREATE TABLE IF NOT EXISTS likes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			content_type TEXT NOT NULL, -- 'post' or 'comment'
			content_id INTEGER NOT NULL,
			is_like BOOLEAN NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY(user_id) REFERENCES users(id)
		)
	`)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) Close() {
	db.Conn.Close()
}

func (db *Database) PopulateMockData() error {
	// Insert mock categories
	_, err := db.Conn.Exec(`
        INSERT OR IGNORE INTO categories (name) VALUES 
        ('Technology'), 
        ('Sports'), 
        ('Music'), 
        ('Movies')
    `)
	if err != nil {
		return fmt.Errorf("error inserting categories: %v", err)
	}

	// Insert a mock user (you'd typically use password hashing in a real app)
	_, err = db.Conn.Exec(`
        INSERT OR IGNORE INTO users (username, email, password_hash) VALUES 
        ('testuser', 'test@example.com', 'hashed_password')
    `)
	if err != nil {
		return fmt.Errorf("error inserting user: %v", err)
	}

	// Insert mock posts
	_, err = db.Conn.Exec(`
        INSERT OR IGNORE INTO posts (user_id, title, content) VALUES 
        (1, 'First Post', 'This is the content of the first post'),
        (1, 'Another Interesting Post', 'Some more content here')
    `)
	if err != nil {
		return fmt.Errorf("error inserting posts: %v", err)
	}

	// Link posts to categories
	_, err = db.Conn.Exec(`
        INSERT OR IGNORE INTO post_categories (post_id, category_id) VALUES 
        (1, 1),  -- First post in Technology category
        (2, 2)   -- Second post in Sports category
    `)
	if err != nil {
		return fmt.Errorf("error linking posts to categories: %v", err)
	}

	return nil
}

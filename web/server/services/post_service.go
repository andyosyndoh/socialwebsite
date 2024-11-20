// services/post_service.go
package services

import (
	"database/sql"

	"forum/web/server/models"
)

type PostService struct {
	db *sql.DB
}

func NewPostService(db *sql.DB) *PostService {
	return &PostService{db: db}
}

func (s *PostService) CreatePost(post *models.Post) (int64, error) {
	// Begin transaction
	tx, err := s.db.Begin()
	if err != nil {
		return 0, err
	}

	// Insert post
	result, err := tx.Exec(
		"INSERT INTO posts (user_id, title, content) VALUES (?, ?, ?)",
		post.UserID, post.Title, post.Content,
	)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// Get post ID
	postID, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// Insert categories
	for _, category := range post.Categories {
		_, err = tx.Exec(
			"INSERT INTO post_categories (post_id, category_id) VALUES (?, ?)",
			postID, category.ID,
		)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return postID, nil
}

func (s *PostService) GetPostsByCategory(categoryID int64) ([]models.Post, error) {
	// Query to get posts by category with additional details
	query := `
		SELECT p.id, p.user_id, p.title, p.content, p.created_at, 
		       (SELECT COUNT(*) FROM likes WHERE content_type = 'post' AND content_id = p.id AND is_like = 1) as likes,
		       (SELECT COUNT(*) FROM likes WHERE content_type = 'post' AND content_id = p.id AND is_like = 0) as dislikes
		FROM posts p
		JOIN post_categories pc ON p.id = pc.post_id
		WHERE pc.category_id = ?
	`

	rows, err := s.db.Query(query, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		err := rows.Scan(&post.ID, &post.UserID, &post.Title, &post.Content,
			&post.CreatedAt, &post.Likes, &post.Dislikes)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

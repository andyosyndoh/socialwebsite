package services

import (
	"database/sql"
	"errors"

	"forum/web/server/models"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	db *sql.DB
}

func NewUserService(db *sql.DB) *UserService {
	return &UserService{db: db}
}

func (s *UserService) Register(user *models.User) error {
	// Check if email already exists
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", user.Email).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Insert user
	_, err = s.db.Exec(
		"INSERT INTO users (username, email, password_hash) VALUES (?, ?, ?)",
		user.Username, user.Email, string(hashedPassword),
	)
	return err
}

func (s *UserService) Login(email, password string) (*models.User, error) {
	var user models.User
	var storedPassword string

	// Fetch user by email
	err := s.db.QueryRow(
		"SELECT id, username, email, password_hash FROM users WHERE email = ?",
		email,
	).Scan(&user.ID, &user.Username, &user.Email, &storedPassword)

	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, err
	}

	// Compare passwords
	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password))
	if err != nil {
		return nil, errors.New("invalid password")
	}

	return &user, nil
}

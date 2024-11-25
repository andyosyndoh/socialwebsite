package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"forum/internals/renders"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"-"` // Don't expose the password hash
}

type Userlogin struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"` // Don't expose the password hash
}

// LoginRequest represents the login request payload
type LoginRequest struct {
	Username    string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse represents the login response payload
type LoginResponse struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}

type Post struct {
    ID        int       `json:"id"`
    Username  string    `json:"username"`
    Content   string    `json:"content"`
    Timestamp time.Time `json:"timestamp"`
    Likes     int       `json:"likes"`
    Dislikes  int       `json:"dislikes"`
}


// HomeHandler handles the homepage route '/'
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		renders.RenderTemplate(w, "home.page.html", nil)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// RegisterHandler handles the registration route '/register'
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	// Ensure the request method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	fmt.Println("here1")
	// Parse the JSON body
	var user Userlogin
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Basic input validation
	if user.Username == "" || user.Email == "" || user.Password == "" {
		log.Println("error")
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}
	fmt.Println("here2")
	// Check if the email is already registered
	var exists bool

	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	// Check if the username already exists
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = ?)", user.Username).Scan(&exists)
	if err != nil {
		fmt.Println("here3")
		log.Println("Database error:", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	if exists {
		fmt.Println("here4")
		http.Error(w, "Username already taken", http.StatusConflict)
		return
	}

	// Insert the new user into the database
	_, err = db.Exec("INSERT INTO users (username, email, password) VALUES (?, ?, ?)", user.Username, user.Email, user.Password)
	if err != nil {
		log.Println("Database error:", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	response := map[string]string{"message": "User registered successfully"}
	json.NewEncoder(w).Encode(response)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse the JSON body
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	// Query the user from the database
	var user User
	err = db.QueryRow("SELECT id, username, email, password FROM users WHERE username = ?", req.Username).
		Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	// Compare the provided password with the stored hashed password
	if user.Password != req.Password {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "user_session",
		Value:    user.Username, // You could store more info, like username or user ID
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour), // Cookie expires in 24 hours
		HttpOnly: true,                           // Prevent access via JavaScript (for security)
	})

	// Respond with success
	resp := LoginResponse{
		Username: user.Username,
		Message:  "Login successful",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve the user_session cookie
	cookie, err := r.Cookie("user_session")
	if err != nil || cookie == nil {
		http.Error(w, "Not logged in", http.StatusUnauthorized)
		return
	}

	// If cookie is found, return user data
	// In a real app, you could fetch more info from the database or the cookie value
	user := map[string]string{
		"username": cookie.Value, // Using the cookie value as the username for this example
	}

	// Respond with user details
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Clear the session cookie by setting MaxAge to -1
	http.SetCookie(w, &http.Cookie{
		Name:   "user_session",
		Value:  "",
		MaxAge: -1, // This invalidates the cookie
		Path:   "/",
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Logged out"})
}

func GetAllPosts(w http.ResponseWriter, r *http.Request) {
    // Set content type to JSON
    w.Header().Set("Content-Type", "application/json")

    // Connect to database
    db, err := sql.Open("sqlite3", "./forumposts.db")
    if err != nil {
        http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
        return
    }
    defer db.Close()

    // Query posts
    rows, err := db.Query(`
        SELECT posts.id, users.username, posts.content, posts.timestamp, 
               COUNT(CASE WHEN likes.type = 'like' THEN 1 END) AS likes,
               COUNT(CASE WHEN likes.type = 'dislike' THEN 1 END) AS dislikes
        FROM posts
        LEFT JOIN users ON posts.user_id = users.id
        LEFT JOIN likes ON posts.id = likes.post_id
        GROUP BY posts.id
        ORDER BY posts.timestamp DESC
    `)
    if err != nil {
        http.Error(w, "Failed to query posts", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    // Parse query results
    var posts []Post
    for rows.Next() {
        var post Post
        err := rows.Scan(&post.ID, &post.Username, &post.Content, &post.Timestamp, &post.Likes, &post.Dislikes)
        if err != nil {
            http.Error(w, "Failed to parse post data", http.StatusInternalServerError)
            return
        }
        posts = append(posts, post)
    }

    // Send response
    json.NewEncoder(w).Encode(posts)
}


// NotFoundHandler handles unknown routes; 404 status
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	renders.RenderTemplate(w, "notfound.page.html", nil)
}

// BadRequestHandler handles bad requests routes
func BadRequestHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	renders.RenderTemplate(w, "badrequest.page.html", nil)
}

// ServerErrorHandler handles server failures that result in status 500
func ServerErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	renders.RenderTemplate(w, "serverError.page.html", nil)
}

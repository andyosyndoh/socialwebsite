package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"text/template"
	"time"

	"forum/web/server/models"
	"forum/web/server/services"

	"github.com/google/uuid"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// Register handles user registration
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	// Only allow POST method
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse JSON request body
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate input
	if user.Email == "" || user.Username == "" || user.Password == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// Attempt registration
	err = h.userService.Register(&user)
	if err != nil {
		if strings.Contains(err.Error(), "email already exists") {
			http.Error(w, "Email already registered", http.StatusConflict)
		} else {
			http.Error(w, "Registration failed", http.StatusInternalServerError)
		}
		return
	}

	// Create session cookie
	sessionToken := uuid.New().String()
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Path:     "/",
	})

	// Respond with success
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
}

// Login handles user authentication
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	// Only allow POST method
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse login credentials
	var loginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&loginRequest)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Attempt login
	user, err := h.userService.Login(loginRequest.Email, loginRequest.Password)
	if err != nil {
		if err.Error() == "user not found" || err.Error() == "invalid password" {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		} else {
			http.Error(w, "Login failed", http.StatusInternalServerError)
		}
		return
	}

	// Create session cookie
	sessionToken := uuid.New().String()
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Path:     "/",
	})

	// Respond with user info (excluding sensitive data)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
	})
}

type PostHandler struct {
	postService *services.PostService
}

func NewPostHandler(postService *services.PostService) *PostHandler {
	return &PostHandler{postService: postService}
}

// CreatePost handles creating a new post
func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	// Only allow POST method
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse post data
	var post models.Post
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate input
	if post.Title == "" || post.Content == "" {
		http.Error(w, "Title and content are required", http.StatusBadRequest)
		return
	}

	// Get user ID from session (placeholder - you'll implement proper session management)
	post.UserID = getCurrentUserID(r)
	if post.UserID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Create post
	postID, err := h.postService.CreatePost(&post)
	if err != nil {
		http.Error(w, "Failed to create post", http.StatusInternalServerError)
		return
	}

	// Respond with created post ID
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int64{"post_id": postID})
}

// GetPostsByCategory retrieves posts for a specific category
func (h *PostHandler) GetPostsByCategory(w http.ResponseWriter, r *http.Request) {
	// Only allow GET method
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse category ID from query parameter
	categoryIDStr := r.URL.Query().Get("category_id")
	categoryID, err := strconv.ParseInt(categoryIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	// Retrieve posts
	posts, err := h.postService.GetPostsByCategory(categoryID)
	if err != nil {
		http.Error(w, "Failed to retrieve posts", http.StatusInternalServerError)
		return
	}

	// Respond with posts
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(posts)
}

func (h *PostHandler) GetAllPosts(w http.ResponseWriter, r *http.Request) {
	// Parse the home template
	tmpl, err := template.ParseFiles("web/templates/index.html")
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Retrieve all posts
	posts, err := h.postService.GetAllPosts()
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to retrieve posts", http.StatusInternalServerError)
		return
	}

	// Execute template with posts data
	err = tmpl.Execute(w, posts)
	if err != nil {
		fmt.Println("Error executing template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// Middleware for authentication
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check for session cookie
		cookie, err := r.Cookie("session_token")
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Validate session token (placeholder - you'll implement proper session management)
		if !isValidSession(cookie.Value) {
			http.Error(w, "Invalid session", http.StatusUnauthorized)
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	}
}

// Placeholder for session validation
func isValidSession(token string) bool {
	// TODO: Implement actual session validation
	// This would typically involve checking against a sessions table or cache
	return token != ""
}

// Placeholder for getting current user ID
func getCurrentUserID(r *http.Request) int64 {
	// TODO: Implement actual user ID retrieval from session
	// This is a placeholder that returns a mock user ID
	return 1
}

// Router setup (typically in main.go or a separate routing file)
func SetupRoutes(userService *services.UserService, postService *services.PostService) *http.ServeMux {
	// Create handlers
	userHandler := NewUserHandler(userService)
	postHandler := NewPostHandler(postService)

	// Create a new router
	mux := http.NewServeMux()

	// User routes
	// mux.HandleFunc("/", HomeRoute)
	staticHandler := http.FileServer(http.Dir("web/static"))
    mux.Handle("/static/", http.StripPrefix("/static/", staticHandler))
	mux.HandleFunc("/", postHandler.GetAllPosts)
	mux.HandleFunc("/register", userHandler.Register)
	mux.HandleFunc("/login", userHandler.Login)

	// Post routes
	mux.HandleFunc("/posts/create", AuthMiddleware(postHandler.CreatePost))
	mux.HandleFunc("/posts/category", postHandler.GetPostsByCategory)

	return mux
}

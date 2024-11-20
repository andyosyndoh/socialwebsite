package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"forum/web/server/database"
	"forum/web/server/routes"
	"forum/web/server/services"
)

func main() {
	db := database.NewDatabase()
	defer db.Close()
	userService := services.NewUserService(db.Conn)
	postService := services.NewPostService(db.Conn)
	// Define the path to the static files
	staticDir := "./web/static/" // Relative path from `server/main.go` to `web/static`

	// Get the absolute path for better error handling
	absStaticDir, err := filepath.Abs(staticDir)
	if err != nil {
		log.Fatalf("Failed to get absolute path of static directory: %v", err)
	}

	// Serve static files
	fs := http.FileServer(http.Dir(absStaticDir))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Printf("Serving static files from %s on http://localhost:8080/static/", absStaticDir)
	routes := routes.SetupRoutes(userService, postService)
	fmt.Println("Server running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", routes); err != nil {
		log.Fatal("ListenAndServe error:", err)
	}
}

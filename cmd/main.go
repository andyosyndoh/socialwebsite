package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"forum/internals/routes"

	"forum/utils"
)

func main() {
	// Initialize SQLite database
	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create users table if it doesn't exist
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL,
		email TEXT NOT NULL,
		password TEXT NOT NULL
	)`)
	if err != nil {
		log.Fatal("Error creating users table:", err)
	}

	mux := http.NewServeMux()
	routes.RegisterRoutes(mux)

	wrappedMux := routes.RouteChecker(mux)

	server := &http.Server{
		Addr:    ":8080",
		Handler: wrappedMux,
	}

	fmt.Println("server running @http://localhost:8080\n=====================================")
	err = server.ListenAndServe()
	if err != nil {
		utils.ErrorHandler("web")
	}
}

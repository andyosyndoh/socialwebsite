package routes

import (
	"net/http"
	"strings"

	"forum/internals/handlers"
	"forum/internals/renders"
)

// Allowed routes
var allowedRoutes = map[string]bool{
	"/": true,
	"/register" : true,
	"/login" : true,
	"/profile" : true,
	"/logout" : true,
	"/posts" : true,
}

// RouteChecker is a middleware that checkes allowed routes
func RouteChecker(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/static/") {
			// Static(w,r)
			next.ServeHTTP(w, r)
			return
		}

		if _, ok := allowedRoutes[r.URL.Path]; !ok {
			handlers.NotFoundHandler(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// RegisterRoutes manages the routes
func RegisterRoutes(mux *http.ServeMux) {
	staticDir := renders.GetProjectRoot("views", "static")
	fs := http.FileServer(http.Dir(staticDir))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handlers.HomeHandler(w, r)
	})

	mux.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		handlers.RegisterHandler(w, r)
	})

	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
        handlers.LoginHandler(w, r)
    })

	mux.HandleFunc("/profile", func(w http.ResponseWriter, r *http.Request) {
        handlers.ProfileHandler(w, r)
    })

	mux.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
        handlers.LogoutHandler(w, r)
    })

	mux.HandleFunc("/posts", func(w http.ResponseWriter, r *http.Request) {
        handlers.GetAllPosts(w, r)
    })
}

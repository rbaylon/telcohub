package main

import (
	"fmt"
	"log"
	"net/http"

	"telcohub/db"
	"telcohub/handlers"
	"telcohub/middleware"

	ghandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

func main() {
	var (
		app_ip       = db.GetEnvVariable("APP_IP")
		app_port     = db.GetEnvVariable("APP_PORT")
		cors_allowed = db.GetEnvCors("CORS_ALLOWED")
		secret       = db.GetEnvVariable("APP_SECRET")
	)

	credentials := ghandlers.AllowCredentials()
	methods := ghandlers.AllowedMethods([]string{"POST"})
	//ttl := ghandlers.MaxAge(3600)
	origins := ghandlers.AllowedOrigins(cors_allowed)

	// Initialize DB
	db.Init()

	// session
	store := sessions.NewCookieStore([]byte(secret))

	// Create router
	r := mux.NewRouter()

	// Auth routes
	r.HandleFunc("/register", handlers.Register).Methods("POST")
	r.HandleFunc("/register.html", handlers.RegisterUi).Methods("GET")
	r.HandleFunc("/login", handlers.Login).Methods("POST")
	r.HandleFunc("/login.html", handlers.LoginUi).Methods("GET")
	r.HandleFunc("/logout", handlers.Logout).Methods("GET")

	// Marker routes (protected)
	marker := r.PathPrefix("/marker").Subrouter()
	marker.Use(middleware.AuthMiddleware(store))
	marker.HandleFunc("/create", handlers.CreateMarker).Methods("POST")
	marker.HandleFunc("/edit/{id}", handlers.EditMarker).Methods("POST")
	marker.HandleFunc("/delete/{id}", handlers.DeleteMarker).Methods("POST")
	marker.HandleFunc("/list", handlers.ListMarkers).Methods("GET")

	// Admin routes (admin-only middleware)
	admin := r.PathPrefix("/admin").Subrouter()
	admin.Use(middleware.RequireRole("admin", store))
	admin.HandleFunc("/dashboard", handlers.AdminDashboard).Methods("GET")
	admin.HandleFunc("/ui", handlers.AdminDashboardUi(store)).Methods("GET")
	admin.HandleFunc("/user/{id}/role", handlers.UpdateUserRole).Methods("POST")
	admin.HandleFunc("/user/{id}/delete", handlers.DeleteUser).Methods("DELETE")
	admin.HandleFunc("/category/list", handlers.ListCategories).Methods("GET")
	admin.HandleFunc("/category/delete/{id}", handlers.DeleteCategory).Methods("POST")
	admin.HandleFunc("/category/create", handlers.CreateCategory).Methods("POST")
	admin.HandleFunc("/category/create.html", handlers.CreateCategoryUi).Methods("GET")

	group := r.PathPrefix("/groups").Subrouter()
	group.Use(middleware.AuthMiddleware(store))
	group.HandleFunc("/list", handlers.ShowGroups).Methods("GET")
	group.HandleFunc("/create", handlers.CreateGroup).Methods("POST")
	group.HandleFunc("/{id}/add", handlers.AddToGroup).Methods("POST")
	group.HandleFunc("/user/{id}/toggle", handlers.ToggleAdmin).Methods("POST")
	group.HandleFunc("/user/{id}/remove", handlers.RemoveMember).Methods("DELETE")

	// Serve static files and templates
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("assets"))))

	r.HandleFunc("/", handlers.Home).Methods("GET")
	r.HandleFunc("/profile.html", handlers.ShowProfile).Methods("GET")
	r.HandleFunc("/change-password.html", handlers.ShowChangePassword).Methods("GET")
	r.HandleFunc("/change-password", handlers.ChangePassword).Methods("POST")

	// Start server
	log.Println("Server running at http://localhost:8080")
	http.ListenAndServe(fmt.Sprintf("%s:%s", app_ip, app_port), ghandlers.CORS(credentials, methods, origins)(r))
}

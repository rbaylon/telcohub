package middleware

import (
	"net/http"
	// Ensure store reference is accessible
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

func AuthMiddleware(store *sessions.CookieStore) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, _ := store.Get(r, "session-id")
			if _, ok := session.Values["user_id"]; !ok {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// Role-based access (admin-only)
func RequireRole(role string, store *sessions.CookieStore) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, _ := store.Get(r, "session-id")
			if session.Values["role"] != role {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

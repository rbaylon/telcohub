package handlers

import (
	"net/http"
	"text/template"

	"github.com/gorilla/sessions"
)

func Home(store *sessions.CookieStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session-id")
		role, ok := session.Values["role"].(string)
		username, uok := session.Values["username"].(string)

		if !ok || !uok {
			http.Redirect(w, r, "/login.html", http.StatusSeeOther)
			return
		}

		tmpl := template.Must(template.ParseFiles("templates/index.html"))
		tmpl.Execute(w, map[string]interface{}{
			"IsAdmin":  role == "admin",
			"Username": username,
		})
	}
}

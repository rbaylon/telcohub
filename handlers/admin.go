package handlers

import (
	"encoding/json"
	"net/http"
	"text/template"

	"telcohub/db"
	"telcohub/models"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

func AdminDashboardUi(store *sessions.CookieStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session-id")
		if _, ok := session.Values["user_id"]; !ok {
			http.Redirect(w, r, "/login.html", http.StatusSeeOther)
			return
		}

		tmpl := template.Must(template.ParseFiles("templates/admin.html"))
		tmpl.Execute(w, nil)
	}
}

func AdminDashboard(w http.ResponseWriter, r *http.Request) {
	var users []models.User
	db.DB.Preload("Markers").Find(&users)

	response := []map[string]interface{}{}
	for _, u := range users {
		var history []map[string]interface{}
		for _, m := range u.Markers {
			history = append(history, map[string]interface{}{
				"id":          m.ID,
				"title":       m.Title,
				"description": m.Description,
				"latitude":    m.Latitude,
				"longitude":   m.Longitude,
				"created_at":  m.CreatedAt.Format("2006-01-02 15:04"),
			})
		}
		response = append(response, map[string]interface{}{
			"id":            u.ID,
			"username":      u.Username,
			"role":          u.Role,
			"markerCount":   len(u.Markers),
			"markerHistory": history,
		})
	}
	json.NewEncoder(w).Encode(response)
}

func UpdateUserRole(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	role := r.FormValue("role")

	if role != "admin" && role != "user" {
		http.Error(w, "Invalid role", http.StatusBadRequest)
		return
	}

	db.DB.Model(&models.User{}).Where("id = ?", id).Update("role", role)
	w.WriteHeader(http.StatusOK)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	db.DB.Delete(&models.User{}, id)
	w.WriteHeader(http.StatusNoContent)
}

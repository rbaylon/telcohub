package handlers

import (
	"log"
	"net/http"
	"telcohub/db"
	"telcohub/models"
	"text/template"
)

func Home(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-id")
	role, ok := session.Values["role"].(string)
	username, uok := session.Values["username"].(string)
	id, uid := session.Values["user_id"].(uint)
	if !uid {
		log.Println("session error: user_id not found")
	}

	if !ok || !uok {
		http.Redirect(w, r, "/login.html", http.StatusSeeOther)
		return
	}
	var categories []models.Category
	db.DB.Find(&categories)
	var groups []models.Group
	db.DB.Where("owner_id = ?", id).Find(&groups)
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, map[string]interface{}{
		"IsAdmin":    role == "admin",
		"Username":   username,
		"Categories": categories,
		"Groups":     groups,
	})
}

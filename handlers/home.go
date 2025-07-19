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
	var gus []models.GroupUser
	db.DB.Preload("Group").Where("user_id = ? AND is_admin = ?", id, true).Find(&gus)
	for _, gu := range gus {
		groups = append(groups, gu.Group)
	}
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, map[string]interface{}{
		"IsAdmin":    role == "admin",
		"Username":   username,
		"Categories": categories,
		"Groups":     groups,
	})
}

func ShowLandingPage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/landing.html"))
	tmpl.Execute(w, nil)
}

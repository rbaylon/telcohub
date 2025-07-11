package handlers

import (
	"net/http"
	"telcohub/db"
	"telcohub/models"
	"telcohub/utils"
	"text/template"
)

func ShowProfile(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-id")
	//id, _ := session.Values["user_id"].(uint)
	role, ok := session.Values["role"].(string)
	username, uok := session.Values["username"].(string)

	if !ok || !uok {
		http.Redirect(w, r, "/login.html", http.StatusSeeOther)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/profile.html"))
	tmpl.Execute(w, map[string]interface{}{
		"Username": username,
		"Role":     role,
	})
}

func ShowChangePassword(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-id")
	_, ok := session.Values["role"].(string)
	_, uok := session.Values["username"].(string)

	if !ok || !uok {
		http.Redirect(w, r, "/login.html", http.StatusSeeOther)
		return
	}
	tmpl := template.Must(template.ParseFiles("templates/change-password.html"))
	tmpl.Execute(w, nil)
}

func ChangePassword(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-id")
	id, ok := session.Values["user_id"].(uint)

	if !ok {
		http.Redirect(w, r, "/login.html", http.StatusSeeOther)
		return
	}

	current := r.FormValue("current")
	newPass := r.FormValue("new")

	var user models.User
	db.DB.First(&user, id)

	if !utils.CheckPasswordHash(current, user.Password) {
		http.Error(w, "Incorrect current password", http.StatusForbidden)
		return
	}

	hash, err := utils.HashPassword(newPass)
	if err != nil {
		http.Error(w, "Password encryption failed", http.StatusInternalServerError)
		return
	}

	user.Password = hash
	db.DB.Save(&user)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

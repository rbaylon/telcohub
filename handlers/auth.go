package handlers

import (
	"net/http"
	"telcohub/db"
	"telcohub/models"
	"telcohub/utils"
	"text/template"

	"github.com/gorilla/sessions"
)

var secret = db.GetEnvVariable("APP_SECRET")
var store = sessions.NewCookieStore([]byte(secret))

func Register(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	hash, err := utils.HashPassword(password)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	user := models.User{Username: username, Password: hash}
	if err := db.DB.Create(&user).Error; err != nil {
		http.Error(w, "Username taken", http.StatusConflict)
		return
	}

	http.Redirect(w, r, "/login.html", http.StatusSeeOther)
}

func RegisterUi(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/register.html"))
	tmpl.Execute(w, nil)
}

func Login(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	var user models.User
	if err := db.DB.Where("username = ?", username).First(&user).Error; err != nil {
		http.Error(w, "Invalid login", http.StatusUnauthorized)
		return
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		http.Error(w, "Incorrect password", http.StatusUnauthorized)
		return
	}

	session, _ := store.Get(r, "session-id")
	session.Values["user_id"] = user.ID
	session.Values["role"] = user.Role
	session.Values["username"] = user.Username
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func LoginUi(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/login.html"))
	tmpl.Execute(w, nil)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-id")
	session.Options.MaxAge = -1
	session.Save(r, w)
	http.Redirect(w, r, "/login.html", http.StatusSeeOther)
}

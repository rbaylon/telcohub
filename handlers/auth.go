package handlers

import (
	"net/http"
	"telcohub/controllers"
	"telcohub/db"
	"telcohub/models"
	"telcohub/utils"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
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

	utils.StartUserSession(w, r, user, store)
	//session, _ := store.Get(r, "session-id")
	//session.Values["user_id"] = user.ID
	//session.Values["role"] = user.Role
	//session.Values["username"] = user.Username
	//session.Save(r, w)

	http.Redirect(w, r, "/gis", http.StatusSeeOther)
}

func LoginUi(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/login.html"))
	tmpl.Execute(w, nil)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-id")
	session.Options.MaxAge = -1
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func BeginAuth(w http.ResponseWriter, r *http.Request) {
	providerName := mux.Vars(r)["provider"]
	provider, err := goth.GetProvider(providerName)
	if err != nil {
		http.Error(w, "Provider not found", http.StatusBadRequest)
		return
	}

	sess, err := provider.BeginAuth("state-token")
	if err != nil {
		http.Error(w, "Error starting auth", http.StatusInternalServerError)
		return
	}

	url, _ := sess.GetAuthURL()
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func CompleteAuth(w http.ResponseWriter, r *http.Request) {
	providerName := mux.Vars(r)["provider"]
	provider, _ := goth.GetProvider(providerName)

	value, _ := provider.UnmarshalSession(r.URL.Query().Encode())
	userData, err := provider.FetchUser(value)
	if err != nil {
		http.Error(w, "Auth failed", http.StatusUnauthorized)
		return
	}

	// 🔐 Auto-register user or find existing one
	user := controllers.FindOrCreateUserByEmail(userData.Email, providerName)

	// Start session
	utils.StartUserSession(w, r, user, store)

	http.Redirect(w, r, "/gis", http.StatusSeeOther)
}

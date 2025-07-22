package handlers

import (
	"fmt"
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
		http.Error(w, "Invalid provider", http.StatusBadRequest)
		return
	}

	authSession, err := provider.BeginAuth("state-token")
	if err != nil {
		http.Error(w, "Auth failed", http.StatusInternalServerError)
		return
	}

	// 🧠 Save session for callback
	session, _ := store.Get(r, "session-id")
	session.Values["provider"] = providerName
	session.Values["auth"] = authSession.Marshal()
	session.Save(r, w)

	// Redirect to provider's login URL
	url, _ := authSession.GetAuthURL()
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func CompleteAuth(w http.ResponseWriter, r *http.Request) {
	// Restore saved session
	session, _ := store.Get(r, "session-id")
	providerName := session.Values["provider"].(string)
	authSessionStr := session.Values["auth"].(string)

	provider, _ := goth.GetProvider(providerName)
	authSession, _ := provider.UnmarshalSession(authSessionStr)

	// 💡 Exchange code for token
	_, err := authSession.Authorize(provider, r.URL.Query())
	if err != nil {
		http.Error(w, fmt.Sprintf("Token exchange failed: %s", err), http.StatusUnauthorized)
		return
	}

	user, err := provider.FetchUser(authSession)
	if err != nil {
		http.Error(w, fmt.Sprintf("User info fetch failed: %s", err), http.StatusInternalServerError)
		return
	}

	// 🔐 Find or create account
	account := controllers.FindOrCreateUserByEmail(user.Email, providerName)
	utils.StartUserSession(w, r, account, store)

	// ✅ Redirect after login
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

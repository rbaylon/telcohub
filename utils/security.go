package utils

import (
	"net/http"
	"telcohub/db"
	"telcohub/models"

	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/google"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes a plain password using bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPasswordHash verifies a password against its hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// get user session
func GetUserFromSession(r *http.Request, store *sessions.CookieStore) (models.User, error) {
	session, err := store.Get(r, "session-id")
	if err != nil {
		return models.User{}, err
	}
	id, ok := session.Values["user_id"].(uint)
	if !ok {
		return models.User{}, err
	}
	var user models.User
	db.DB.First(&user, id)
	return user, nil
}

func OAuthInit(google_client_id string, google_client_secret string, callback_url string) {
	goth.UseProviders(
		google.New(google_client_id, google_client_secret, callback_url, "email", "profile"),
	)
}

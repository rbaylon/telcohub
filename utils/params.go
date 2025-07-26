package utils

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/gorilla/mux"
)

func GetParamID(r *http.Request) uint {
	vars := mux.Vars(r)
	idStr := vars["id"]

	idUint64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return 0 // handle safely if needed
	}
	return uint(idUint64)
}

func ValidateCredsInput(username string, password string, w http.ResponseWriter) error {
	if username == "" || password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return fmt.Errorf("Username and password are required")
	}

	if len(username) < 3 || !regexp.MustCompile(`^[a-zA-Z0-9_.-]+$`).MatchString(username) {
		http.Error(w, "Invalid username format", http.StatusBadRequest)
		return fmt.Errorf("Invalid username format: %s", username)
	}

	if len(password) < 6 {
		http.Error(w, "Password too short", http.StatusBadRequest)
		return fmt.Errorf("Password too short")
	}
	return nil
}

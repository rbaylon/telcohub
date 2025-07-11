package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"telcohub/db"
	"telcohub/models"

	"github.com/gorilla/mux"
)

func getUserFromSession(r *http.Request) (models.User, error) {
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

func CreateMarker(w http.ResponseWriter, r *http.Request) {
	user, err := getUserFromSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	lat, _ := strconv.ParseFloat(r.FormValue("latitude"), 64)
	lng, _ := strconv.ParseFloat(r.FormValue("longitude"), 64)
	cat_id, _ := strconv.Atoi(r.FormValue("category_id"))

	marker := models.Marker{
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
		Latitude:    lat,
		Longitude:   lng,
		UserID:      user.ID,
		CategoryID:  uint(cat_id),
	}
	db.DB.Create(&marker)
	w.WriteHeader(http.StatusCreated)
}

func EditMarker(w http.ResponseWriter, r *http.Request) {
	user, err := getUserFromSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	id := mux.Vars(r)["id"]

	var marker models.Marker
	db.DB.First(&marker, id)

	if marker.UserID != user.ID && user.Role != "admin" {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	cat_id, _ := strconv.Atoi(r.FormValue("category_id"))
	marker.Title = r.FormValue("title")
	marker.Description = r.FormValue("description")
	marker.CategoryID = uint(cat_id)
	db.DB.Save(&marker)
	w.WriteHeader(http.StatusOK)
}

func DeleteMarker(w http.ResponseWriter, r *http.Request) {
	user, err := getUserFromSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	id := mux.Vars(r)["id"]

	var marker models.Marker
	db.DB.First(&marker, id)

	if marker.UserID != user.ID && user.Role != "admin" {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	db.DB.Delete(&marker)
	w.WriteHeader(http.StatusNoContent)
}

func ListMarkers(w http.ResponseWriter, r *http.Request) {
	user, err := getUserFromSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var markers []models.Marker
	if user.Role == "admin" {
		db.DB.Preload("User").Preload("Category").Find(&markers)
	} else {
		db.DB.Preload("User").Preload("Category").Where("user_id = ?", user.ID).Find(&markers)
	}

	json.NewEncoder(w).Encode(markers)
}

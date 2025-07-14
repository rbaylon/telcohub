package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"telcohub/db"
	"telcohub/models"
	"telcohub/utils"

	"github.com/gorilla/mux"
)

func CreateMarker(w http.ResponseWriter, r *http.Request) {
	user, err := utils.GetUserFromSession(r, store)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	lat, _ := strconv.ParseFloat(r.FormValue("latitude"), 64)
	lng, _ := strconv.ParseFloat(r.FormValue("longitude"), 64)
	cat_id, _ := strconv.Atoi(r.FormValue("category_id"))
	group_id, _ := strconv.Atoi(r.FormValue("group_id"))
	marker := models.Marker{
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
		Latitude:    lat,
		Longitude:   lng,
		UserID:      user.ID,
		CategoryID:  uint(cat_id),
		GroupID:     uint(group_id),
	}

	// Check if user is group admin
	var gu models.GroupUser
	db.DB.Where("user_id = ? AND group_id = ?", user.ID, marker.GroupID).First(&gu)

	if !gu.IsAdmin {
		if user.Role != "admin" {
			http.Error(w, "Unauthorized", http.StatusForbidden)
			return
		}
	}

	db.DB.Create(&marker)
	w.WriteHeader(http.StatusCreated)
}

func EditMarker(w http.ResponseWriter, r *http.Request) {
	user, err := utils.GetUserFromSession(r, store)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	id := mux.Vars(r)["id"]

	var marker models.Marker
	db.DB.First(&marker, id)

	cat_id, _ := strconv.Atoi(r.FormValue("category_id"))
	group_id, _ := strconv.Atoi(r.FormValue("group_id"))
	marker.Title = r.FormValue("title")
	marker.Description = r.FormValue("description")
	marker.CategoryID = uint(cat_id)
	marker.GroupID = uint(group_id)

	// Check if user is group admin
	var gu models.GroupUser
	db.DB.Where("user_id = ? AND group_id = ?", user.ID, marker.GroupID).First(&gu)
	if !gu.IsAdmin {
		if user.Role != "admin" {
			http.Error(w, "Unauthorized", http.StatusForbidden)
			return
		}
	}
	db.DB.Save(&marker)
	w.WriteHeader(http.StatusOK)
}

func DeleteMarker(w http.ResponseWriter, r *http.Request) {
	user, err := utils.GetUserFromSession(r, store)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	id := mux.Vars(r)["id"]

	var marker models.Marker
	db.DB.First(&marker, id)

	// Check if user is group admin
	var gu models.GroupUser
	db.DB.Where("user_id = ? AND group_id = ?", user.ID, marker.GroupID).First(&gu)

	if !gu.IsAdmin {
		if user.Role != "admin" {
			http.Error(w, "Unauthorized", http.StatusForbidden)
			return
		}
	}

	db.DB.Delete(&marker)
	w.WriteHeader(http.StatusNoContent)
}

func ListMarkers(w http.ResponseWriter, r *http.Request) {
	user, err := utils.GetUserFromSession(r, store)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var markers []models.Marker
	var groupIDs []uint
	if user.Role == "admin" {
		db.DB.Preload("Group").Preload("User").Preload("Category").Find(&markers)
	} else {
		db.DB.Model(&models.GroupUser{}).Where("user_id = ?", user.ID).Pluck("group_id", &groupIDs)
		db.DB.Preload("Group").Preload("User").Preload("Category").Where("group_id IN ?", groupIDs).Find(&markers)
	}

	json.NewEncoder(w).Encode(markers)
}

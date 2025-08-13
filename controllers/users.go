package controllers

import (
	"strings"
	"telcohub/db"
	"telcohub/models"
)

func FindOrCreateUserByEmail(email string, provider string) models.User {
	var user models.User
	err := db.DB.Where("email = ?", email).First(&user).Error
	if err == nil {
		return user // already exists
	}

	// Create new user
	user = models.User{
		Email:    email,
		Username: email,    // derive username
		Provider: provider, // e.g. "google" or "facebook"
		Password: "",       // no local password
		Role:     "user",
	}
	db.DB.Create(&user)

	emailaddress := strings.Split(email, "@")
	group := models.Group{Name: emailaddress[0], OwnerID: user.ID}
	db.DB.Create(&group)

	db.DB.Create(&models.GroupUser{
		UserID:  user.ID,
		GroupID: group.ID,
		IsAdmin: true,
	})
	return user
}

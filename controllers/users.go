package controllers

import (
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
	return user
}

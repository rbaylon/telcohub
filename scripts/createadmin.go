package main

import (
	"fmt"
	"log"
	"strings"

	"telcohub/db"
	"telcohub/models"
	"telcohub/utils"
)

func main() {
	// 📦 Initialize DB
	var admin_user = strings.TrimSpace(db.GetEnvVariable("APP_ADMIN"))
	var admin_pw = strings.TrimSpace(db.GetEnvVariable("APP_ADMIN_PW"))
	db.Init()

	hash, err := utils.HashPassword(admin_pw)
	if err != nil {
		log.Fatalf("Password hash error: %v", err)
	}

	admin := models.User{
		Username: admin_user,
		Password: hash,
		Role:     "admin",
	}

	if err := db.DB.Create(&admin).Error; err != nil {
		log.Fatalf("Failed to create admin: %v", err)
	}

	fmt.Println("Admin user created successfully!")
}

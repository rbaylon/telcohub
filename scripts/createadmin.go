package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"telcohub/db"
	"telcohub/models"
	"telcohub/utils"
)

func main() {
	// 📦 Initialize DB
	db.Init()

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter admin username: ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)

	fmt.Print("Enter admin password: ")
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	hash, err := utils.HashPassword(password)
	if err != nil {
		log.Fatalf("Password hash error: %v", err)
	}

	admin := models.User{
		Username: username,
		Password: hash,
		Role:     "admin",
	}

	if err := db.DB.Create(&admin).Error; err != nil {
		log.Fatalf("Failed to create admin: %v", err)
	}

	fmt.Println("✅ Admin user created successfully!")
}

package db

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"telcohub/models" // Replace with actual module name
)

// DB is a globally accessible database instance
var DB *gorm.DB

// Init initializes the database and runs migrations
func Init() {
	var err error
	DB, err = gorm.Open(sqlite.Open("markers.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Database connected successfully")

	// Run migrations for all models
	if err := DB.AutoMigrate(&models.User{}, &models.Marker{}); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Println("Database migrated")
}

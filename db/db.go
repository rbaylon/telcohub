package db

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"telcohub/models" // Replace with actual module name
)

// DB is a globally accessible database instance
var DB *gorm.DB

// Init initializes the database and runs migrations
func Init() {
	var (
		dbname = GetEnvVariable("DB_NAME")
		dbhost = GetEnvVariable("DB_HOST")
		dbuser = GetEnvVariable("DB_USER")
		dbpass = GetEnvVariable("DB_PASS")
		err    error
	)
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 TimeZone=Asia/Singapore", dbhost, dbuser, dbpass, dbname)
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	//DB, err = gorm.Open(sqlite.Open("markers.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Database connected successfully")

	// Run migrations for all models
	if err := DB.AutoMigrate(
		&models.User{},
		&models.Marker{},
		&models.Category{},
		&models.Group{},
		&models.GroupUser{},
	); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Println("Database migrated")
}

func GetEnvVariable(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}
	return os.Getenv(key)
}

func GetEnvCors(key string) []string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}
	var cors []string
	if key == "CORS_ALLOWED" {
		cors_allowed := os.Getenv(key)
		// Parse the string into a slice
		cors = strings.Split(cors_allowed, ",")
	}
	return cors
}

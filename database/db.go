package database

import (
	"fmt"
	"log"
	"os"
	"pos-backend/models" // Make sure this matches your module name in go.mod

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := os.Getenv("DB_URL") // We will put the Nhost string in .env
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to database!", err)
	}

	// Auto-Migrate: This creates the tables automatically for you
	// We will add Product and Order models here later
	database.AutoMigrate(&models.User{})

	DB = database
	fmt.Println("🚀 Database connected successfully!")
}
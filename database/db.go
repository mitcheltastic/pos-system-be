package database

import (
	"fmt"
	"log"
	"os"

	"pos-backend/models" // Matches your go.mod

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	// 1. Try to load .env for local dev. 
	// We use `_` to ignore errors because on Vercel, there is no .env file!
	_ = godotenv.Load()

	// 2. Grab the connection string
	dsn := os.Getenv("DB_URL")
	if dsn == "" {
		log.Fatal("DB_URL is not set! Check your .env file or Vercel environment variables.")
	}

	// 3. Connect to Postgres
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database!", err)
	}

	// 4. Auto-Migrate tables
	database.AutoMigrate(&models.User{}, &models.Product{}, &models.Order{}, &models.OrderItem{})

	DB = database
	fmt.Println("🚀 Database connected successfully!")
}
package main

import (
	"log"
	"pos-backend/database"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// 1. Load .env variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// 2. Connect to Database
	database.ConnectDatabase()

	// 3. Initialize Router
	r := gin.Default()

	// 4. Test Route
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "POS Backend is running!",
		})
	})

	// 5. Run Server
	r.Run(":8080") // Runs on localhost:8080
}
package main

import (
	"log"
	"pos-backend/controllers" // Import your new controller
	"pos-backend/database"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// 1. Load .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// 2. Connect DB
	database.ConnectDatabase()

	// 3. Init Router
	r := gin.Default()

	// 4. Routes
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "POS Backend is running!"})
	})

	// Auth Routes Group
	auth := r.Group("/auth")
	{
		auth.POST("/register", controllers.Register)
		auth.POST("/login", controllers.Login)
	}

	// Product Routes Group
	products := r.Group("/products")
    {
        products.POST("/", controllers.CreateProduct)       // Create
        products.GET("/", controllers.GetProducts)          // List & Filter
        products.PUT("/:id", controllers.UpdateProduct)     // Update
        products.DELETE("/:id", controllers.DeleteProduct)  // Delete
    }

	orders := r.Group("/orders")
    {
        orders.POST("/", controllers.CreateOrder) // Cashier creates order
        orders.GET("/", controllers.GetOrders)    // Admin sees history
    }

	dashboard := r.Group("/dashboard")
    {
        dashboard.GET("/", controllers.GetDashboardStats)
    }

	// 5. Run
	r.Run(":8080")
}
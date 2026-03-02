package main

import (
	"log"
	"pos-backend/controllers" 
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
		auth.POST("/change-password", controllers.ChangePassword)
		
		auth.GET("/profile/:id", controllers.GetProfile)    // ✅ Add this line (Fetch data)
		auth.PUT("/profile/:id", controllers.UpdateProfile) // (Update data)
	}

	storeProfile := r.Group("/store")
	{
		storeProfile.POST("/", controllers.CreateStoreProfile)
		storeProfile.GET("/", controllers.GetStoreProfiles)
		storeProfile.PUT("/:id", controllers.UpdateStoreProfile)
		storeProfile.DELETE("/:id", controllers.DeleteStoreProfile)
	}

	// Product Routes Group
	products := r.Group("/products")
	{
		products.POST("/", controllers.CreateProduct)
		products.GET("/", controllers.GetProducts)
		products.GET("/summary", controllers.GetInventorySummary) // ✅ Add this for the 3 UI cards
		products.PUT("/:id", controllers.UpdateProduct)
		products.DELETE("/:id", controllers.DeleteProduct)
	}

	orders := r.Group("/orders")
	{
		orders.POST("/", controllers.CreateOrder)
		orders.GET("/", controllers.GetOrders)
		orders.PUT("/:id/status", controllers.UpdateOrderStatus) // ✅ Add this line!
	}

	paymentMethods := r.Group("/payment-methods")
	{
		paymentMethods.POST("/", controllers.CreatePaymentMethod)
		paymentMethods.GET("/", controllers.GetPaymentMethods)
		paymentMethods.PUT("/:id/status", controllers.UpdatePaymentMethodStatus)
	}

	analytics := r.Group("/analytics")
	{
		analytics.GET("/dashboard", controllers.GetAnalytics)
	}

	dashboard := r.Group("/dashboard")
    {
        dashboard.GET("/", controllers.GetDashboardStats)
    }

	// 5. Run
	r.Run(":8080")
}
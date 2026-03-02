package api

import (
	"net/http"
	"pos-backend/controllers"
	"pos-backend/database"

	"github.com/gin-gonic/gin"
)

var app *gin.Engine

// init() runs exactly once when the Vercel serverless function cold-starts
func init() {
	database.ConnectDatabase()

	app = gin.Default()

	// --- CORS MIDDLEWARE (Crucial for your Frontend Team) ---
	app.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*") // Change * to your frontend URL later for security
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// --- ROUTES ---
	app.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "POS Backend is running on Vercel!"})
	})

	auth := app.Group("/auth")
	{
		auth.POST("/register", controllers.Register)
		auth.POST("/login", controllers.Login)
	}

	products := app.Group("/products")
	{
		products.POST("/", controllers.CreateProduct)
		products.GET("/", controllers.GetProducts)
		products.PUT("/:id", controllers.UpdateProduct)
		products.DELETE("/:id", controllers.DeleteProduct)
	}

	orders := app.Group("/orders")
	{
		orders.POST("/", controllers.CreateOrder)
		orders.GET("/", controllers.GetOrders)
	}

	dashboard := app.Group("/dashboard")
	{
		dashboard.GET("/", controllers.GetDashboardStats)
	}
}

// Handler is the actual exported function that Vercel looks for
func Handler(w http.ResponseWriter, r *http.Request) {
	app.ServeHTTP(w, r)
}
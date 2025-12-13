package controllers

import (
	"net/http"
	"pos-backend/database"
	"pos-backend/models"

	"github.com/gin-gonic/gin"
)

// Input format for the checkout request
type OrderItemInput struct {
	ProductID uint `json:"product_id" binding:"required"`
	Quantity  int  `json:"quantity" binding:"required"`
}

type CreateOrderInput struct {
	CashierID     uint             `json:"cashier_id" binding:"required"` // In real app, get this from JWT
	PaymentMethod string           `json:"payment_method" binding:"required"`
	Items         []OrderItemInput `json:"items" binding:"required"`
}

func CreateOrder(c *gin.Context) {
	var input CreateOrderInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Start a Database Transaction (Safety first!)
	tx := database.DB.Begin()

	var totalAmount float64
	var orderItems []models.OrderItem

	// Loop through items to calculate total & verify they exist
	for _, itemInput := range input.Items {
		var product models.Product
		
		// Fetch fresh product data (especially Price)
		if err := tx.First(&product, itemInput.ProductID).Error; err != nil {
			tx.Rollback() // Cancel everything if product invalid
			c.JSON(http.StatusBadRequest, gin.H{"error": "Product not found: " + string(rune(itemInput.ProductID))})
			return
		}

		// Calculate cost
		subtotal := product.Price * float64(itemInput.Quantity)
		totalAmount += subtotal

		// Create the OrderItem object
		orderItems = append(orderItems, models.OrderItem{
			ProductID: itemInput.ProductID,
			Quantity:  itemInput.Quantity,
			Price:     product.Price, // Save the price at THIS moment
		})
	}

	// Create the Main Order
	order := models.Order{
		CashierID:     input.CashierID,
		TotalAmount:   totalAmount,
		PaymentMethod: input.PaymentMethod,
		Items:         orderItems, // GORM will automagically create the OrderItems too!
	}

	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}

	// Commit the transaction
	tx.Commit()

	c.JSON(http.StatusCreated, gin.H{
		"message": "Order created successfully",
		"order_id": order.ID,
		"total": totalAmount,
	})
}

// Get All Orders (For Admin Dashboard)
func GetOrders(c *gin.Context) {
	var orders []models.Order
	// Preload "Items" and "Items.Product" so we see the details in the JSON
	database.DB.Preload("Items.Product").Find(&orders)
	c.JSON(http.StatusOK, gin.H{"data": orders})
}
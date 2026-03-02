package controllers

import (
	"fmt"
	"net/http"
	"time"

	"pos-backend/database"
	"pos-backend/models"

	"github.com/gin-gonic/gin"
)

type OrderItemInput struct {
	ProductID uint `json:"product_id" binding:"required"`
	Quantity  int  `json:"quantity" binding:"required"`
}

type CreateOrderInput struct {
	CashierID     uint             `json:"cashier_id" binding:"required"` 
	CustomerName  string           `json:"customer_name"` // Optional, from UI
	PaymentMethod string           `json:"payment_method" binding:"required"`
	Status        string           `json:"status"` // Optional, defaults to Completed
	Items         []OrderItemInput `json:"items" binding:"required"`
}

type UpdateOrderStatusInput struct {
	Status string `json:"status" binding:"required"` // Completed, Pending, Cancelled
}

// 1. Create Order & Deduct Stock
func CreateOrder(c *gin.Context) {
	var input CreateOrderInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx := database.DB.Begin()

	var totalAmount float64
	var orderItems []models.OrderItem

	for _, itemInput := range input.Items {
		var product models.Product
		
		if err := tx.First(&product, itemInput.ProductID).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Product ID %d not found", itemInput.ProductID)})
			return
		}

		// ✅ NEW: Check if there is enough stock
		if product.Stock < itemInput.Quantity {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Not enough stock for %s. Only %d left.", product.Name, product.Stock)})
			return
		}

		// ✅ NEW: Deduct the stock
		product.Stock -= itemInput.Quantity
		if err := tx.Save(&product).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update inventory"})
			return
		}

		subtotal := product.Price * float64(itemInput.Quantity)
		totalAmount += subtotal

		orderItems = append(orderItems, models.OrderItem{
			ProductID: itemInput.ProductID,
			Quantity:  itemInput.Quantity,
			Price:     product.Price, 
		})
	}

	// Determine status (default to Completed if not provided)
	orderStatus := "Completed"
	if input.Status != "" {
		orderStatus = input.Status
	}

	// Generate a simple unique order number using Unix timestamp
	orderNumber := fmt.Sprintf("#ORD-%d", time.Now().Unix())

	order := models.Order{
		OrderNumber:   orderNumber,
		CustomerName:  input.CustomerName,
		CashierID:     input.CashierID,
		TotalAmount:   totalAmount,
		PaymentMethod: input.PaymentMethod,
		Status:        orderStatus,
		Items:         orderItems, 
	}

	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}

	tx.Commit()

	c.JSON(http.StatusCreated, gin.H{
		"message": "Order created successfully",
		"data": order,
	})
}

// 2. Get All Orders (Handles Search Bar and Status Tabs)
func GetOrders(c *gin.Context) {
	var orders []models.Order
	
	search := c.Query("search") // Catches "John Doe" or "#ORD-123"
	status := c.Query("status") // Catches "Completed", "Pending", "Cancelled"

	query := database.DB.Preload("Items.Product")

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if search != "" {
		// Search by Order Number or Customer Name
		query = query.Where("order_number ILIKE ? OR customer_name ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// Order by newest first
	query.Order("created_at desc").Find(&orders)

	c.JSON(http.StatusOK, gin.H{"data": orders})
}

// 3. Update Order Status (For moving Pending -> Completed or Cancelled)
func UpdateOrderStatus(c *gin.Context) {
	id := c.Param("id")
	var order models.Order

	if err := database.DB.First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	var input UpdateOrderStatusInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Note: If you cancel an order, you might want to add logic here to add the stock back!
	// For now, we just update the text status.
	database.DB.Model(&order).Update("status", input.Status)

	c.JSON(http.StatusOK, gin.H{"message": "Order status updated", "data": order})
}
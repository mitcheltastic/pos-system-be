package controllers

import (
	"net/http"
	"pos-backend/database" // Adjust to your module name if needed
	"pos-backend/models"

	"github.com/gin-gonic/gin"
)

type CreatePaymentMethodInput struct {
	Name   string `json:"payment_method" binding:"required"`
	Status string `json:"status"` // Optional when creating, will default to "active" if empty
}

type UpdatePaymentMethodStatusInput struct {
	Status string `json:"status" binding:"required"` // "active" or "inactive"
}

// 1. Create a new Payment Method (Admin)
func CreatePaymentMethod(c *gin.Context) {
	var input CreatePaymentMethodInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Default status if none provided
	status := "active"
	if input.Status != "" {
		status = input.Status
	}

	paymentMethod := models.PaymentMethod{
		Name:   input.Name,
		Status: status,
	}

	if err := database.DB.Create(&paymentMethod).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create payment method (might already exist)"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Payment method created", "data": paymentMethod})
}

// 2. Get Payment Methods (Cashier needs this for checkout)
func GetPaymentMethods(c *gin.Context) {
	var paymentMethods []models.PaymentMethod

	// Optional: Let the frontend filter by ?status=active
	statusFilter := c.Query("status")

	if statusFilter != "" {
		database.DB.Where("status = ?", statusFilter).Find(&paymentMethods)
	} else {
		database.DB.Find(&paymentMethods)
	}

	c.JSON(http.StatusOK, gin.H{"data": paymentMethods})
}

// 3. Update Status (Admin turning a method on/off)
func UpdatePaymentMethodStatus(c *gin.Context) {
	id := c.Param("id")
	var paymentMethod models.PaymentMethod

	// Verify it exists
	if err := database.DB.First(&paymentMethod, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payment method not found"})
		return
	}

	var input UpdatePaymentMethodStatusInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update only the status column
	if err := database.DB.Model(&paymentMethod).Update("status", input.Status).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Payment method status updated",
		"data":    paymentMethod,
	})
}
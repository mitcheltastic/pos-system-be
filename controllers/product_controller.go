package controllers

import (
	"net/http"
	"pos-backend/database"
	"pos-backend/models"

	"github.com/gin-gonic/gin"
)

// Input struct ensures we only get specific data from the user
type CreateProductInput struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required"`
	Category    string  `json:"category" binding:"required"`
	ImageURL    string  `json:"image_url"`
}

// 1. Create Product (Admin only usually)
func CreateProduct(c *gin.Context) {
	var input CreateProductInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product := models.Product{
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
		Category:    input.Category,
		ImageURL:    input.ImageURL,
	}

	database.DB.Create(&product)
	c.JSON(http.StatusCreated, gin.H{"data": product})
}

// 2. Get All Products (With optional Category Filter)
func GetProducts(c *gin.Context) {
	var products []models.Product
	
	// Check if user sent ?category=Coffee
	category := c.Query("category")

	if category != "" {
		// Filter by category
		database.DB.Where("category = ?", category).Find(&products)
	} else {
		// Return everything
		database.DB.Find(&products)
	}

	c.JSON(http.StatusOK, gin.H{"data": products})
}

// 3. Update Product
func UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	var product models.Product

	// Check if product exists
	if err := database.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	// Validate input
	var input CreateProductInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update fields
	database.DB.Model(&product).Updates(models.Product{
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
		Category:    input.Category,
		ImageURL:    input.ImageURL,
	})

	c.JSON(http.StatusOK, gin.H{"data": product})
}

// 4. Delete Product
func DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	var product models.Product

	if err := database.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	database.DB.Delete(&product)
	c.JSON(http.StatusOK, gin.H{"message": "Product deleted"})
}
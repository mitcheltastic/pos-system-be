package controllers

import (
	"context"
	"mime/multipart"
	"net/http"

	"pos-backend/database"
	"pos-backend/models"

	"github.com/gin-gonic/gin"
	"github.com/imagekit-developer/imagekit-go/v2"
)

type CreateProductInput struct {
	Name         string                `form:"name" binding:"required"`
	Description  string                `form:"description"`
	Price        float64               `form:"price" binding:"required"`
	Category     string                `form:"category" binding:"required"`
	Stock        int                   `form:"stock" binding:"required"`         // Added
	ReorderLevel int                   `form:"reorder_level" binding:"required"` // Added
	Image        *multipart.FileHeader `form:"image"`
}

// 1. Create Product
func CreateProduct(c *gin.Context) {
	var input CreateProductInput

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var finalImageURL string

	if input.Image != nil {
		file, err := input.Image.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
			return
		}
		defer file.Close()

		client := imagekit.NewClient()
		resp, err := client.Files.Upload(
			context.Background(),
			imagekit.FileUploadParams{
				File:     file,
				FileName: input.Image.Filename,
			},
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Upload failed: " + err.Error()})
			return
		}

		finalImageURL = resp.URL
	}

	product := models.Product{
		Name:         input.Name,
		Description:  input.Description,
		Price:        input.Price,
		Category:     input.Category,
		Stock:        input.Stock,
		ReorderLevel: input.ReorderLevel,
		ImageURL:     finalImageURL,
	}

	if err := database.DB.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Product created successfully", "data": product})
}

// 2. Get All Products (Handles the Search Bar and Category Filters)
func GetProducts(c *gin.Context) {
	var products []models.Product
	
	category := c.Query("category")
	search := c.Query("search") // Catches text from the "Search inventory..." bar

	query := database.DB.Model(&models.Product{})

	if category != "" {
		query = query.Where("category = ?", category)
	}
	
	if search != "" {
		// ILIKE is Postgres-specific for case-insensitive search
		query = query.Where("name ILIKE ?", "%"+search+"%")
	}

	query.Find(&products)

	c.JSON(http.StatusOK, gin.H{"data": products})
}

// 3. Update Product
func UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	var product models.Product

	if err := database.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	var input CreateProductInput
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	finalImageURL := product.ImageURL

	if input.Image != nil {
		file, err := input.Image.Open()
		if err == nil {
			defer file.Close()
			client := imagekit.NewClient()
			resp, err := client.Files.Upload(
				context.Background(),
				imagekit.FileUploadParams{
					File:     file,
					FileName: input.Image.Filename,
				},
			)
			if err == nil {
				finalImageURL = resp.URL 
			}
		}
	}

	database.DB.Model(&product).Updates(models.Product{
		Name:         input.Name,
		Description:  input.Description,
		Price:        input.Price,
		Category:     input.Category,
		Stock:        input.Stock,
		ReorderLevel: input.ReorderLevel,
		ImageURL:     finalImageURL,
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

// 5. Get Inventory Summary (Powers the top 3 cards in the UI)
func GetInventorySummary(c *gin.Context) {
	var totalItems int64
	var lowStockAlerts int64
	var outOfStock int64

	// Count total product types
	database.DB.Model(&models.Product{}).Count(&totalItems)

	// Count items where stock is running low (but not empty)
	database.DB.Model(&models.Product{}).Where("stock <= reorder_level AND stock > 0").Count(&lowStockAlerts)

	// Count items completely out of stock
	database.DB.Model(&models.Product{}).Where("stock = 0").Count(&outOfStock)

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"total_items":      totalItems,
			"low_stock_alerts": lowStockAlerts,
			"out_of_stock":     outOfStock,
		},
	})
}
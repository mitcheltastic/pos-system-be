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
	Name        string                `form:"name" binding:"required"`
	Description string                `form:"description"`
	Price       float64               `form:"price" binding:"required"`
	Category    string                `form:"category" binding:"required"`
	Image       *multipart.FileHeader `form:"image"`
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
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
		Category:    input.Category,
		ImageURL:    finalImageURL,
	}

	if err := database.DB.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Product created successfully", "data": product})
}

// 2. Get All Products (With optional Category Filter)
func GetProducts(c *gin.Context) {
	var products []models.Product
	
	category := c.Query("category")

	if category != "" {
		database.DB.Where("category = ?", category).Find(&products)
	} else {
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

	var input CreateProductInput
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Default to existing image URL
	finalImageURL := product.ImageURL

	// If a NEW image is uploaded, send it to ImageKit
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
				finalImageURL = resp.URL // Overwrite with new URL
			}
		}
	}

	// Update fields
	database.DB.Model(&product).Updates(models.Product{
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
		Category:    input.Category,
		ImageURL:    finalImageURL,
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
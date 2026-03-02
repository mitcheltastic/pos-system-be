package controllers

import (
	"net/http"
	"pos-backend/database" // Adjust to match your go.mod
	"pos-backend/models"

	"github.com/gin-gonic/gin"
)

type StoreProfileInput struct {
	Name    string `json:"name" binding:"required"`
	Address string `json:"address"`
	Phone   string `json:"phone"`
	TaxID   string `json:"tax_id"`
}

// 1. Create Store Profile
func CreateStoreProfile(c *gin.Context) {
	var input StoreProfileInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	store := models.StoreProfile{
		Name:    input.Name,
		Address: input.Address,
		Phone:   input.Phone,
		TaxID:   input.TaxID,
	}

	if err := database.DB.Create(&store).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create store profile"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Store profile created", "data": store})
}

// 2. Get All Store Profiles (Usually just fetches the 1 store you have)
func GetStoreProfiles(c *gin.Context) {
	var stores []models.StoreProfile
	database.DB.Find(&stores)
	c.JSON(http.StatusOK, gin.H{"data": stores})
}

// 3. Update Store Profile
func UpdateStoreProfile(c *gin.Context) {
	id := c.Param("id")
	var store models.StoreProfile

	// Find the store
	if err := database.DB.First(&store, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Store profile not found"})
		return
	}

	var input StoreProfileInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update the fields
	database.DB.Model(&store).Updates(models.StoreProfile{
		Name:    input.Name,
		Address: input.Address,
		Phone:   input.Phone,
		TaxID:   input.TaxID,
	})

	c.JSON(http.StatusOK, gin.H{"message": "Store profile updated", "data": store})
}

// 4. Delete Store Profile (Rarely used, but good for completeness)
func DeleteStoreProfile(c *gin.Context) {
	id := c.Param("id")
	var store models.StoreProfile

	if err := database.DB.First(&store, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Store profile not found"})
		return
	}

	database.DB.Delete(&store)
	c.JSON(http.StatusOK, gin.H{"message": "Store profile deleted"})
}
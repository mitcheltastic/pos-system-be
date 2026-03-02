package controllers

import (
	"net/http"
	"pos-backend/database" // Ensure this matches your go.mod module name
	"pos-backend/models"

	"github.com/gin-gonic/gin"
)

// --- INPUT STRUCTS ---

type RegisterInput struct {
	Name     string `json:"name" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role" binding:"required"` // "admin" or "cashier"
}

type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// ✅ NEW: Struct for changing password
type ChangePasswordInput struct {
	Username        string `json:"username" binding:"required"`
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required"`
}

// --- FUNCTIONS ---

// 1. Register Function
func Register(c *gin.Context) {
	var input RegisterInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := models.User{
		Name:     input.Name,
		Username: input.Username,
		Password: input.Password, 
		Role:     input.Role,
	}

	result := database.DB.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username already exists or database error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully", "data": user})
}

// 2. Login Function
func Login(c *gin.Context) {
	var input LoginInput
	var user models.User

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	if user.Password != input.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"role":    user.Role,
		"name":    user.Name,
	})
}

// ✅ 3. NEW: Change Password Function
func ChangePassword(c *gin.Context) {
	var input ChangePasswordInput
	var user models.User

	// Validate JSON input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find the user in the database
	if err := database.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Verify the current password matches
	if user.Password != input.CurrentPassword {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect current password"})
		return
	}

	// Update to the new password
	// GORM's .Update() only updates specific columns, which is perfect here
	if err := database.DB.Model(&user).Update("password", input.NewPassword).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password updated successfully"})
}
package controllers

import (
	"net/http"
	"pos-backend/database" // Ensure this matches your go.mod module name
	"pos-backend/models"

	"github.com/gin-gonic/gin"
)

// 1. Define what we expect from the user (The Input)
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

// 2. Register Function (Create User)
func Register(c *gin.Context) {
	var input RegisterInput

	// Validate JSON input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create User Model
	user := models.User{
		Name:     input.Name,
		Username: input.Username,
		Password: input.Password, // Saving as plain text as requested
		Role:     input.Role,
	}

	// Save to Database
	result := database.DB.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username already exists or database error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully", "data": user})
}

// 3. Login Function (Authenticate)
func Login(c *gin.Context) {
	var input LoginInput
	var user models.User

	// Validate JSON
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user exists
	if err := database.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	// Check Password (Plain text check)
	if user.Password != input.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect password"})
		return
	}

	// Return User Data (In a real app, we would return a JWT Token here)
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"role":    user.Role,
		"name":    user.Name,
	})
}
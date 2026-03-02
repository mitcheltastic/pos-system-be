package controllers

import (
	"context"
	"mime/multipart"
	"net/http"

	"pos-backend/database" // Ensure this matches your go.mod module name
	"pos-backend/models"

	"github.com/gin-gonic/gin"
	"github.com/imagekit-developer/imagekit-go/v2"
)

// --- INPUT STRUCTS ---

type RegisterInput struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name"`
	Username  string `json:"username" binding:"required"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Password  string `json:"password" binding:"required"`
	Role      string `json:"role" binding:"required"` // "admin" or "cashier"
}

type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ChangePasswordInput struct {
	Username        string `json:"username" binding:"required"`
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required"`
}

type UpdateProfileInput struct {
	FirstName      string                `form:"first_name" binding:"required"`
	LastName       string                `form:"last_name"`
	Email          string                `form:"email"`
	Phone          string                `form:"phone"`
	ProfilePicture *multipart.FileHeader `form:"profile_picture"` // Handled as a file
}

// --- FUNCTIONS ---

// 1. Register Function (Create User)
func Register(c *gin.Context) {
	var input RegisterInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := models.User{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Username:  input.Username,
		Email:     input.Email,
		Phone:     input.Phone,
		Password:  input.Password, // Saving as plain text as requested
		Role:      input.Role,
		// ProfilePicture remains empty until updated
	}

	result := database.DB.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username or Email already exists"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully", "data": user})
}

// 2. Login Function (Authenticate)
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

	// Returning the full profile data so the frontend can display the avatar and name immediately
	c.JSON(http.StatusOK, gin.H{
		"message":         "Login successful",
		"id":              user.ID,
		"role":            user.Role,
		"first_name":      user.FirstName,
		"last_name":       user.LastName,
		"profile_picture": user.ProfilePicture, 
	})
}

// 3. Change Password Function
func ChangePassword(c *gin.Context) {
	var input ChangePasswordInput
	var user models.User

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if user.Password != input.CurrentPassword {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect current password"})
		return
	}

	if err := database.DB.Model(&user).Update("password", input.NewPassword).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password updated successfully"})
}

// 4. Update Profile Function (Handles ImageKit Uploads)
func UpdateProfile(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	if err := database.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var input UpdateProfileInput
	// Use ShouldBind instead of ShouldBindJSON to handle multipart/form-data
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	finalImageURL := user.ProfilePicture

	// Process image upload if a new file is provided
	if input.ProfilePicture != nil {
		file, err := input.ProfilePicture.Open()
		if err == nil {
			defer file.Close()

			client := imagekit.NewClient()
			resp, err := client.Files.Upload(
				context.Background(),
				imagekit.FileUploadParams{
					File:     file,
					FileName: input.ProfilePicture.Filename,
				},
			)

			if err == nil {
				finalImageURL = resp.URL 
			}
		}
	}

	// Update user data in the database
	database.DB.Model(&user).Updates(models.User{
		FirstName:      input.FirstName,
		LastName:       input.LastName,
		Email:          input.Email,
		Phone:          input.Phone,
		ProfilePicture: finalImageURL,
	})

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile updated successfully",
		"data":    user,
	})
}

// 5. Get Profile Function (Fetch user data without password)
func GetProfile(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	// Find the user by ID
	if err := database.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Return the data, explicitly mapping it so we DON'T send the password back
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"id":              user.ID,
			"first_name":      user.FirstName,
			"last_name":       user.LastName,
			"username":        user.Username,
			"email":           user.Email,
			"phone":           user.Phone,
			"profile_picture": user.ProfilePicture,
			"role":            user.Role,
		},
	})
}
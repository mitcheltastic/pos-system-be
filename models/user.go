package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Username       string `json:"username" gorm:"unique;not null"`
	Email          string `json:"email" gorm:"unique"` // Added email
	Phone          string `json:"phone"`               // Added phone
	ProfilePicture string `json:"profile_picture"`     // Stores the ImageKit URL
	Password       string `json:"password"`            // Plain text for now
	Role           string `json:"role"`                // "admin" or "cashier"
}
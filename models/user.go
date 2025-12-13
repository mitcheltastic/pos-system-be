package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string `json:"name"`
	Username string `json:"username" gorm:"unique"`
	Password string `json:"password"` // Storing plain text for now as requested
	Role     string `json:"role"`     // "admin" or "cashier"
}
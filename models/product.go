package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Category    string  `json:"category"`  // e.g., "Coffee", "Dessert"
	ImageURL    string  `json:"image_url"` // Optional: link to an image
}
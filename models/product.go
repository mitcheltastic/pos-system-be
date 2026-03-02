package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	Price        float64 `json:"price"`
	Category     string  `json:"category"`
	ImageURL     string  `json:"image_url"`
	Stock        int     `json:"stock"`         // NEW: Current quantity
	ReorderLevel int     `json:"reorder_level"` // NEW: Threshold for low stock warning
}
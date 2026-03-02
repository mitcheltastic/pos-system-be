package models

import "gorm.io/gorm"

type PaymentMethod struct {
	gorm.Model
	Name   string `json:"payment_method" gorm:"unique;not null"` // e.g., "Cash", "BCA Transfer", "QRIS"
	Status string `json:"status" gorm:"default:'active'"`        // e.g., "active" or "inactive"
}
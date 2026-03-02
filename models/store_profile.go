package models

import "gorm.io/gorm"

type StoreProfile struct {
	gorm.Model
	Name    string `json:"name" gorm:"not null"`
	Address string `json:"address"`
	Phone   string `json:"phone"`
	TaxID   string `json:"tax_id"`
}
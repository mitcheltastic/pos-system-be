package models

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	CashierID     uint        `json:"cashier_id"`     // Who processed this?
	TotalAmount   float64     `json:"total_amount"`
	PaymentMethod string      `json:"payment_method"` // "Cash", "QRIS"
	Items         []OrderItem `json:"items" gorm:"foreignKey:OrderID"`
}

type OrderItem struct {
	gorm.Model
	OrderID   uint    `json:"order_id"`
	ProductID uint    `json:"product_id"`
	Product   Product `json:"product"` // Link to product info
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`   // Snapshot of price at time of sale
}
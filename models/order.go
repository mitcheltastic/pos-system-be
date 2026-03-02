package models

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	OrderNumber   string      `json:"order_number" gorm:"unique;not null"` // e.g., #ORD-1709400000
	CustomerName  string      `json:"customer_name"`                       // e.g., John Doe
	CashierID     uint        `json:"cashier_id"`
	TotalAmount   float64     `json:"total_amount"`
	PaymentMethod string      `json:"payment_method"`
	Status        string      `json:"status" gorm:"default:'Completed'"`   // Completed, Pending, Cancelled
	Items         []OrderItem `json:"items" gorm:"foreignKey:OrderID"`
}

type OrderItem struct {
	gorm.Model
	OrderID   uint    `json:"order_id"`
	ProductID uint    `json:"product_id"`
	Product   Product `json:"product" gorm:"foreignKey:ProductID"` // Loads product details
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"` // Price at the time of sale
}
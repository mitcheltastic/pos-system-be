package controllers

import (
	"net/http"
	"pos-backend/database"
	"pos-backend/models"
	"time"

	"github.com/gin-gonic/gin"
)

type DashboardData struct {
	TotalRevenue      float64        `json:"total_revenue"`
	TotalOrders       int64          `json:"total_orders"`
	TodayRevenue      float64        `json:"today_revenue"`
	PopularProducts   []ProductStats `json:"popular_products"`
}

type ProductStats struct {
	Name      string  `json:"name"`
	TotalSold int     `json:"total_sold"`
}

func GetDashboardStats(c *gin.Context) {
	var stats DashboardData

	// 1. Calculate All-Time Revenue
	// SELECT SUM(total_amount) FROM orders
	database.DB.Model(&models.Order{}).
		Select("COALESCE(SUM(total_amount), 0)").
		Scan(&stats.TotalRevenue)

	// 2. Count Total Orders
	// SELECT COUNT(*) FROM orders
	database.DB.Model(&models.Order{}).Count(&stats.TotalOrders)

	// 3. Calculate Revenue for TODAY only
	// SELECT SUM(total_amount) FROM orders WHERE created_at >= '2023-XX-XX 00:00:00'
	startOfDay := time.Now().Truncate(24 * time.Hour)
	database.DB.Model(&models.Order{}).
		Where("created_at >= ?", startOfDay).
		Select("COALESCE(SUM(total_amount), 0)").
		Scan(&stats.TodayRevenue)

	// 4. Find Top 5 Best Selling Products (Advanced Query)
	// Joins orders_items with products, groups by product, counts quantity
	database.DB.Table("order_items").
		Select("products.name, sum(order_items.quantity) as total_sold").
		Joins("join products on products.id = order_items.product_id").
		Group("products.name").
		Order("total_sold desc").
		Limit(5).
		Scan(&stats.PopularProducts)

	c.JSON(http.StatusOK, gin.H{"data": stats})
}
package controllers

import (
	"net/http"
	"time"

	"pos-backend/database"

	"github.com/gin-gonic/gin"
)

// --- Structs to format the JSON exactly how the frontend needs it ---

type DashboardSummary struct {
	TotalRevenue  float64 `json:"total_revenue"`
	TotalOrders   int64   `json:"total_orders"`
	TotalCustomers int64  `json:"customers"`
	AvgOrderValue float64 `json:"avg_order_value"`
}

type ChartData struct {
	Label string  `json:"label"` // e.g., Date or Category Name
	Value float64 `json:"value"` // Revenue or Percentage
}

type TopProductData struct {
	Rank        int     `json:"rank"`
	ProductName string  `json:"product_name"`
	UnitsSold   int     `json:"units_sold"`
	Revenue     float64 `json:"revenue"`
}

type AnalyticsResponse struct {
	Summary       DashboardSummary `json:"summary"`
	SalesTrend    []ChartData      `json:"sales_trend"`
	SalesCategory []ChartData      `json:"sales_by_category"`
	TopProducts   []TopProductData `json:"top_products"`
}

// --- The Main Endpoint ---

func GetAnalytics(c *gin.Context) {
	timeframe := c.Query("timeframe") // "today", "week", "month", "year"

	// 1. Determine the Start Date based on the timeframe filter
	now := time.Now()
	var startDate time.Time

	switch timeframe {
	case "today":
		startDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	case "week":
		// Go back 7 days
		startDate = now.AddDate(0, 0, -7)
	case "month":
		startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	case "year":
		startDate = time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
	default:
		startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()) // Default to this month
	}

	var response AnalyticsResponse

	// 2. Calculate Summary Cards
	database.DB.Table("orders").
		Where("status = ? AND created_at >= ?", "Completed", startDate).
		Select("COALESCE(SUM(total_amount), 0) as total_revenue, COUNT(id) as total_orders, COUNT(DISTINCT customer_name) as total_customers, COALESCE(AVG(total_amount), 0) as avg_order_value").
		Scan(&response.Summary)

	// 3. Sales Trend (Line & Bar Chart data)
	// Groups revenue by Date. The frontend can format these dates into "Mon, Tue" or "Jan, Feb".
	database.DB.Table("orders").
		Select("DATE(created_at) as label, SUM(total_amount) as value").
		Where("status = ? AND created_at >= ?", "Completed", startDate).
		Group("DATE(created_at)").
		Order("label ASC").
		Scan(&response.SalesTrend)

	// 4. Sales by Category (Pie Chart)
	// Joins order items with products to figure out which categories sold the most
	database.DB.Table("order_items").
		Select("products.category as label, SUM(order_items.quantity * order_items.price) as value").
		Joins("left join products on products.id = order_items.product_id").
		Joins("left join orders on orders.id = order_items.order_id").
		Where("orders.status = ? AND orders.created_at >= ?", "Completed", startDate).
		Group("products.category").
		Scan(&response.SalesCategory)

	// 5. Top Selling Products (Table)
	var topProducts []TopProductData
	database.DB.Table("order_items").
		Select("products.name as product_name, SUM(order_items.quantity) as units_sold, SUM(order_items.quantity * order_items.price) as revenue").
		Joins("left join products on products.id = order_items.product_id").
		Joins("left join orders on orders.id = order_items.order_id").
		Where("orders.status = ? AND orders.created_at >= ?", "Completed", startDate).
		Group("products.id, products.name").
		Order("units_sold DESC").
		Limit(5). // Top 5 products
		Scan(&topProducts)

	// Add the rank numbers 1, 2, 3...
	for i := range topProducts {
		topProducts[i].Rank = i + 1
	}
	response.TopProducts = topProducts

	// Return the massive payload!
	c.JSON(http.StatusOK, gin.H{"data": response})
}
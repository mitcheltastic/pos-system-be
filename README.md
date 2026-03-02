# 🛒 POS System Backend

A modern Point of Sale (POS) backend service built with Go and Gin, designed for small to medium retail operations. This system provides comprehensive inventory management, order processing, user authentication, and sales analytics with role-based access control.

## 📋 Table of Contents

- [Features](#features)
- [Tech Stack](#tech-stack)
- [Project Structure](#project-structure)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Environment Variables](#environment-variables)
- [API Documentation](#api-documentation)
- [Database Schema](#database-schema)
- [User Roles & Permissions](#user-roles--permissions)
- [Running the Application](#running-the-application)

## ✨ Features

- **User Authentication & Authorization**
  - User registration with role assignment (Admin/Cashier)
  - Login functionality with role-based responses
  - Role-based access control for different operations

- **Product Management**
  - Create, read, update, and delete products
  - Categorized product organization
  - Product descriptions and pricing
  - Image URL support for product photos

- **Order Processing**
  - Create orders with multiple items
  - Support for multiple payment methods (Cash, QRIS)
  - Real-time order total calculation
  - Order item tracking with price snapshots

- **Dashboard & Analytics**
  - View sales statistics and metrics
  - Order history and trends
  - Performance insights for admins

- **Database Deployment**
  - PostgreSQL integration
  - Automatic schema migration with GORM
  - Vercel-ready deployment configuration

## 🚀 Tech Stack

| Component | Technology | Version |
|-----------|-----------|---------|
| **Language** | Go | 1.24.4+ |
| **Web Framework** | Gin | v1.11.0 |
| **Database** | PostgreSQL | Latest |
| **ORM** | GORM | v1.31.1 |
| **Database Driver** | gorm.io/driver/postgres | v1.6.0 |
| **Environment Management** | godotenv | v1.5.1 |
| **Image Handling** | ImageKit Go SDK | v2.2.0 |

## 📂 Project Structure

```text
pos-system-be/
├── api/                   # Serverless API functions
│   └── index.go           # Vercel entrypoint
├── controllers/           # Request handlers & business logic
│   ├── auth_controller.go      # User registration & login
│   ├── product_controller.go   # Product CRUD operations
│   ├── order_controller.go     # Order creation & retrieval
│   └── dashboard_controller.go # Analytics & statistics
├── database/              # Database connection & initialization
│   └── db.go              # PostgreSQL connection setup
├── models/                # Data models & schemas
│   ├── user.go            # User model (Admin/Cashier)
│   ├── product.go         # Product model
│   └── order.go           # Order & OrderItem models
├── middleware/            # Authentication & authorization
├── main.go                # Application entry point
├── go.mod                 # Go module dependencies
├── go.sum                 # Dependency checksums
├── .env                   # Environment variables (local only)
├── .env.example           # Environment variables template
├── vercel.json            # Vercel deployment configuration
└── README.md              # This file
```

## 📋 Prerequisites

Before running the application, ensure you have:

- **Go** 1.24.4 or higher ([Download](https://golang.org/dl/))
- **PostgreSQL** database (local or cloud-hosted like Nhost)
- **Git** for version control

## 🔧 Installation

### 1. Clone the Repository

```bash
git clone <repository-url>
cd pos-system-be
```

### 2. Install Dependencies

```bash
go mod download
go mod verify
```

### 3. Set Up Environment Variables

Create a `.env` file in the root directory:

```bash
cp .env.example .env
```

Or manually create it with the required variables (see [Environment Variables](#environment-variables) section).

### 4. Verify Installation

```bash
go mod tidy
```

## 🔐 Environment Variables

Create a `.env` file in the root directory with the following variables:

```env
# Database Configuration
DB_URL=postgresql://user:password@localhost:5432/pos_system

# Application Port (optional, defaults to 8080)
PORT=8080

# Environment
ENV=development
```

### Database URL Format

PostgreSQL connection string:
```
postgresql://[username]:[password]@[host]:[port]/[database_name]
```

**Example for local PostgreSQL:**
```
postgresql://postgres:yourpassword@localhost:5432/pos_system
```

**Example for Nhost (cloud):**
```
postgresql://[project-ref]_[user]:[password]@[host].nhost.run:5432/postgres
```

### Important Notes
- **Never commit `.env` file** - it contains sensitive credentials
- Update `.gitignore` to exclude `.env` if not already done
- For Vercel deployment, set environment variables in Vercel dashboard

## 📡 API Documentation

### Base URL
```
http://localhost:8080
```

### Health Check
```
GET /
Response: { "message": "POS Backend is running!" }
```

### Authentication Endpoints

#### Register User
```http
POST /auth/register
Content-Type: application/json

{
  "name": "John Doe",
  "username": "johndoe",
  "password": "securepassword",
  "role": "admin"  // or "cashier"
}

Response (201 Created):
{
  "message": "User registered successfully",
  "data": {
    "id": 1,
    "name": "John Doe",
    "username": "johndoe",
    "role": "admin",
    "created_at": "2026-03-02T10:30:00Z"
  }
}
```

#### Login User
```http
POST /auth/login
Content-Type: application/json

{
  "username": "johndoe",
  "password": "securepassword"
}

Response (200 OK):
{
  "message": "Login successful",
  "role": "admin",
  "name": "John Doe"
}
```

### Product Endpoints

#### Create Product
```http
POST /products
Content-Type: application/json

{
  "name": "Espresso",
  "description": "Strong black coffee",
  "price": 3.50,
  "category": "Coffee",
  "image_url": "https://example.com/espresso.jpg"
}

Response (201 Created):
{
  "id": 1,
  "name": "Espresso",
  "description": "Strong black coffee",
  "price": 3.50,
  "category": "Coffee",
  "image_url": "https://example.com/espresso.jpg",
  "created_at": "2026-03-02T10:30:00Z"
}
```

#### Get All Products
```http
GET /products?category=Coffee&limit=10&page=1

Response (200 OK):
[
  {
    "id": 1,
    "name": "Espresso",
    "description": "Strong black coffee",
    "price": 3.50,
    "category": "Coffee",
    "image_url": "https://example.com/espresso.jpg"
  }
]
```

#### Update Product
```http
PUT /products/:id
Content-Type: application/json

{
  "name": "Premium Espresso",
  "price": 4.50
}

Response (200 OK):
{
  "message": "Product updated successfully",
  "data": { product_details }
}
```

#### Delete Product
```http
DELETE /products/:id

Response (200 OK):
{
  "message": "Product deleted successfully"
}
```

### Order Endpoints

#### Create Order
```http
POST /orders
Content-Type: application/json

{
  "cashier_id": 1,
  "payment_method": "Cash",  // or "QRIS"
  "items": [
    {
      "product_id": 1,
      "quantity": 2,
      "price": 3.50
    },
    {
      "product_id": 2,
      "quantity": 1,
      "price": 5.00
    }
  ]
}

Response (201 Created):
{
  "id": 1,
  "cashier_id": 1,
  "total_amount": 12.00,
  "payment_method": "Cash",
  "items": [ ... ],
  "created_at": "2026-03-02T10:30:00Z"
}
```

#### Get All Orders
```http
GET /orders?limit=10&page=1

Response (200 OK):
[
  {
    "id": 1,
    "cashier_id": 1,
    "total_amount": 12.00,
    "payment_method": "Cash",
    "items": [ ... ],
    "created_at": "2026-03-02T10:30:00Z"
  }
]
```

### Dashboard Endpoints

#### Get Dashboard Statistics
```http
GET /dashboard

Response (200 OK):
{
  "total_revenue": 15000.00,
  "total_orders": 45,
  "orders_today": 12,
  "top_products": [ ... ],
  "revenue_by_category": { ... }
}
```

## 💾 Database Schema

### Users Table
Stores user account information with role-based access control.

```sql
CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  username VARCHAR(255) UNIQUE NOT NULL,
  password VARCHAR(255) NOT NULL,
  role VARCHAR(50) NOT NULL CHECK (role IN ('admin', 'cashier')),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

**Fields:**
- `id` - Unique user identifier
- `name` - Full name of the user
- `username` - Unique username for login
- `password` - Password (currently stored as plain text - SECURITY NOTE: should be hashed)
- `role` - User role: 'admin' or 'cashier'

### Products Table
Maintains the product catalog with pricing and categorization.

```sql
CREATE TABLE products (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  description TEXT,
  price DECIMAL(10, 2) NOT NULL,
  category VARCHAR(100),
  image_url VARCHAR(255),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

**Fields:**
- `id` - Unique product identifier
- `name` - Product name
- `description` - Detailed product description
- `price` - Product price (max 99,999.99)
- `category` - Product category (e.g., "Coffee", "Dessert")
- `image_url` - URL to product image

### Orders Table
Records all transactions processed by the system.

```sql
CREATE TABLE orders (
  id SERIAL PRIMARY KEY,
  cashier_id INTEGER NOT NULL,
  total_amount DECIMAL(10, 2) NOT NULL,
  payment_method VARCHAR(50) NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (cashier_id) REFERENCES users(id)
);
```

**Fields:**
- `id` - Unique order identifier
- `cashier_id` - ID of the cashier who processed the order
- `total_amount` - Total amount of the order
- `payment_method` - Payment method used (Cash, QRIS, etc.)

### Order Items Table
Details individual items within each order.

```sql
CREATE TABLE order_items (
  id SERIAL PRIMARY KEY,
  order_id INTEGER NOT NULL,
  product_id INTEGER NOT NULL,
  quantity INTEGER NOT NULL,
  price DECIMAL(10, 2) NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (order_id) REFERENCES orders(id),
  FOREIGN KEY (product_id) REFERENCES products(id)
);
```

**Fields:**
- `id` - Unique order item identifier
- `order_id` - Reference to the parent order
- `product_id` - Reference to the ordered product
- `quantity` - Number of units ordered
- `price` - Price snapshot at time of sale

## 👥 User Roles & Permissions

### Admin Role
- **Capabilities:**
  - Create, read, update, and delete products
  - View all orders and sales history
  - Access dashboard and analytics
  - View detailed reports
  - Manage user accounts (planned)

### Cashier Role
- **Capabilities:**
  - Process new orders
  - View product catalog
  - Access limited order history
  - Cannot modify products
  - Cannot access admin dashboard

## 🚀 Running the Application

### Local Development

1. **Start the server:**
```bash
go run main.go
```

2. **Expected output:**
```
🚀 Database connected successfully!
[GIN-debug] Loaded HTML Templates (2): 
[GIN-debug] GET    /               --> main.main.func1 (3 handlers)
[GIN-debug] POST   /auth/register  --> controllers.Register (3 handlers)
[GIN-debug] POST   /auth/login     --> controllers.Login (3 handlers)
[GIN-debug] POST   /products       --> controllers.CreateProduct (3 handlers)
[GIN-debug] GET    /products       --> controllers.GetProducts (3 handlers)
[GIN-debug] PUT    /products/:id   --> controllers.UpdateProduct (3 handlers)
[GIN-debug] DELETE /products/:id   --> controllers.DeleteProduct (3 handlers)
[GIN-debug] POST   /orders         --> controllers.CreateOrder (3 handlers)
[GIN-debug] GET    /orders         --> controllers.GetOrders (3 handlers)
[GIN-debug] GET    /dashboard      --> controllers.GetDashboardStats (3 handlers)
[GIN-debug] Listening and serving HTTP on :8080
```

3. **Test the health endpoint:**
```bash
curl http://localhost:8080/
# Response: {"message":"POS Backend is running!"}
```

### Vercel Deployment

1. **Push to GitHub:**
```bash
git push origin main
```

2. **Deploy to Vercel:**
   - Connect your GitHub repository to Vercel
   - Set environment variables in Vercel dashboard
   - Deploy automatically on push

3. **Access your API:**
```
https://your-project-name.vercel.app
```

## 🔒 Security Considerations

⚠️ **Important Security Notes:**

1. **Password Hashing** - Passwords are currently stored as plain text. For production, implement proper hashing using bcrypt or argon2.
   
2. **JWT Implementation** - Replace response messages with JWT tokens for stateless authentication.

3. **Input Validation** - Implement comprehensive validation for all inputs.

4. **HTTPS** - Always use HTTPS in production.

5. **Environment Variables** - Never commit `.env` files. Use secure secret management.

6. **SQL Injection** - Current implementation uses GORM parameterized queries, which are safe.

7. **CORS** - Configure CORS appropriately for your frontend.

## 📦 Dependencies

All Go dependencies are managed in `go.mod`. To view them:

```bash
go mod graph
```

To update dependencies:

```bash
go get -u ./...
go mod tidy
```

## 🛠️ Development Tips

### Debug Mode
Enable Gin debug logging:
```go
gin.SetMode(gin.DebugMode)
```

### Database Debugging
Use Postgres CLI to inspect database:
```bash
psql -U user -d pos_system -h localhost
```

### Useful PostgreSQL Commands
```sql
-- View all tables
\dt

-- View specific table structure
\d orders

-- View all data
SELECT * FROM users;
```

## 📝 Future Enhancements

- [ ] JWT token-based authentication
- [ ] Password hashing with bcrypt
- [ ] Advanced reporting and exports
- [ ] Inventory tracking and low-stock alerts
- [ ] Multi-store support
- [ ] Customer loyalty program
- [ ] Real-time notifications
- [ ] API rate limiting
- [ ] Comprehensive API documentation (Swagger)
- [ ] Unit and integration tests

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit changes (`git commit -m 'Add AmazingFeature'`)
4. Push to branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License.

## 💬 Support

For questions or issues, please open an issue on GitHub or contact the development team.

---

**Last Updated:** March 2, 2026  
**Version:** 1.0.0

A backend service for a Point of Sale (POS) system designed for small retail operations. This system manages roles (Admin, Cashier), inventory, and transaction logging.

## 🚀 Tech Stack

* **Language:** Golang
* **Framework:** Gin (Web Framework)
* **Database:** PostgreSQL (via Nhost)
* **ORM:** GORM
* **Auth:** JWT (Planned) / Role-based Access Control

## 📂 Project Structure

```text
pos-system-be/
├── controllers/      # Route handlers (Logic)
├── database/         # Database connection & config
├── models/           # DB Schema definitions (Structs)
├── middleware/       # Auth & Role validation
├── main.go           # Entry point
└── .env              # Environment variables
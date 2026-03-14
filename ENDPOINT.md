# POS System Backend API Endpoints

This document lists all available endpoints in the POS (Point of Sale) System Backend, along with their explanations.

## Base URL
The server runs on `http://localhost:8080` by default.

## Root Endpoint

- **GET /**  
  Returns a simple message indicating the POS Backend is running.  
  Response: `{"message": "POS Backend is running!"}`

## Authentication Endpoints (`/auth`)

- **POST /auth/register**  
  Registers a new user (admin or cashier).  
  Body: `{"first_name": "string", "last_name": "string", "username": "string", "email": "string", "phone": "string", "password": "string", "role": "admin"|"cashier"}`  
  Response: User data on success.

- **POST /auth/login**  
  Authenticates a user and returns profile data.  
  Body: `{"username": "string", "password": "string"}`  
  Response: User ID, role, name, and profile picture.

- **POST /auth/change-password**  
  Changes the password for a user.  
  Body: `{"username": "string", "current_password": "string", "new_password": "string"}`  
  Response: Success message.

- **GET /auth/profile/:id**  
  Fetches user profile data by user ID (excluding password).  
  Response: User details including ID, name, username, email, phone, profile picture, and role.

- **PUT /auth/profile/:id**  
  Updates user profile, including optional profile picture upload via ImageKit.  
  Body (form-data): `first_name`, `last_name`, `email`, `phone`, `profile_picture` (file).  
  Response: Updated user data.

## Store Profile Endpoints (`/store`)

- **POST /store/**  
  Creates a new store profile.  
  Body: `{"name": "string", "address": "string", "phone": "string", "tax_id": "string"}`  
  Response: Created store profile data.

- **GET /store/**  
  Retrieves all store profiles (typically one store).  
  Response: Array of store profiles.

- **PUT /store/:id**  
  Updates an existing store profile by ID.  
  Body: `{"name": "string", "address": "string", "phone": "string", "tax_id": "string"}`  
  Response: Updated store profile data.

- **DELETE /store/:id**  
  Deletes a store profile by ID.  
  Response: Success message.

## Product Endpoints (`/products`)

- **POST /products/**  
  Creates a new product, with optional image upload via ImageKit.  
  Body (form-data): `name`, `description`, `price`, `category`, `stock`, `reorder_level`, `image` (file).  
  Response: Created product data.

- **GET /products/**  
  Retrieves all products, with optional filtering by category and search query.  
  Query params: `?category=string&search=string`  
  Response: Array of products.

- **GET /products/summary**  
  Gets inventory summary for dashboard cards: total items, low stock alerts, out of stock.  
  Response: Summary statistics.

- **PUT /products/:id**  
  Updates an existing product by ID, with optional image upload.  
  Body (form-data): Same as create.  
  Response: Updated product data.

- **DELETE /products/:id**  
  Deletes a product by ID.  
  Response: Success message.

## Order Endpoints (`/orders`)

- **POST /orders/**  
  Creates a new order, deducts stock from products, and calculates total.  
  Body: `{"cashier_id": number, "customer_name": "string", "payment_method": "string", "status": "string", "items": [{"product_id": number, "quantity": number}]}`  
  Response: Created order data with order number.

- **GET /orders/**  
  Retrieves all orders, with optional filtering by status and search (order number or customer name).  
  Query params: `?status=string&search=string`  
  Response: Array of orders with preloaded items and products.

- **PUT /orders/:id/status**  
  Updates the status of an order (e.g., Pending to Completed or Cancelled).  
  Body: `{"status": "Completed"|"Pending"|"Cancelled"}`  
  Response: Updated order data.

## Payment Method Endpoints (`/payment-methods`)

- **POST /payment-methods/**  
  Creates a new payment method (e.g., Cash, Card).  
  Body: `{"payment_method": "string", "status": "active"|"inactive"}`  
  Response: Created payment method data.

- **GET /payment-methods/**  
  Retrieves all payment methods, with optional status filter.  
  Query params: `?status=active`  
  Response: Array of payment methods.

- **PUT /payment-methods/:id/status**  
  Updates the status of a payment method (active/inactive).  
  Body: `{"status": "active"|"inactive"}`  
  Response: Updated payment method data.

## Analytics Endpoints (`/analytics`)

- **GET /analytics/dashboard**  
  Retrieves comprehensive analytics data for the dashboard, including summary, sales trends, sales by category, and top products.  
  Query params: `?timeframe=today|week|month|year` (defaults to month).  
  Response: Detailed analytics data with charts and tables.

## Dashboard Endpoints (`/dashboard`)

- **GET /dashboard/**  
  Retrieves dashboard statistics: total revenue, total orders, today's revenue, and popular products.  
  Response: Dashboard data summary.</content>
<parameter name="filePath">d:\1. COLLEGE\1. HERE WE GO\0.1 CPS\12. POS RESEARCH BE\pos-system-be\ENDPOINT.md
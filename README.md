# 🛒 Simple POS Backend (Go + Gin)

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
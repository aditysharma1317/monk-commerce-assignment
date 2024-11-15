# Coupons Management API

A RESTful API for managing discount coupons on an e-commerce platform. This API provides support for creating, retrieving, deleting, and applying different types of coupons like cart-wise, product-wise, and BxGy offers.

---

## Features

### Coupon Types:

- **Cart-wise Coupons**: Discounts applied to the entire cart when the total value exceeds a threshold.
- **Product-wise Coupons**: Discounts applied to specific products in the cart.
- **BxGy Coupons**: "Buy X, Get Y" deals with configurable repetition limits.

### Key Endpoints:
- `POST /coupons`: Create a new coupon.
- `GET /coupons`: Retrieve all available coupons.
- `GET /coupons/{id}`: Retrieve a specific coupon by its ID.
- `DELETE /coupons/{id}`: Delete a coupon by its ID.
- `POST /applicable-coupons`: Fetch applicable coupons for a given cart.
- `POST /apply-coupon/{id}`: Apply a specific coupon to the cart and return the updated cart.

### Designed for Extensibility:
- Easily add new coupon types in the future with minimal code changes.

### Error Handling:
- Logs errors using `zap`.
- Rolls back transactions on failure to ensure data consistency.

---

## Tech Stack

- **Language**: Go (Golang)
- **Database**: PostgreSQL [CREATE DATABASE monk-commerce]
- **Framework**: Gin
- **ORM**: GORM
- **Migrations**: [golang-migrate](https://github.com/golang-migrate/migrate)
- **Logging**: Zap

---

## Setup Instructions

### 1. Clone the Repository

```bash
git clone https://github.com/your-repo/monk-commerce-assignment.git
cd monk-commerce-assignment
migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/monk-commerce?sslmode=disable" -verbose up
go run main.go
```
## Things to implement

- **Unit Testing**: Test core functionality of service and DAO layers.
- **Performance Enhancements**: Use caching for frequently accessed coupons.
- **Enhanced Validations**: Add stricter input validation to prevent invalid data.



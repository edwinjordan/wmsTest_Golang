# WMS (Warehouse Management System) Backend API

## Overview
RESTful API untuk sistem manajemen gudang (Warehouse Management System) yang dibangun menggunakan Golang dengan Repository Pattern. API ini mendukung manajemen produk, lokasi, dan pergerakan stok dengan business rules yang ketat.

## Tech Stack
- **Backend**: Golang 1.25+
- **Database**: PostgreSQL
- **Framework**: Gorilla Mux
- **Authentication**: JWT + API Key
- **Architecture**: Repository Pattern (No ORM)

## Features
- ✅ Authentication (JWT & API Key)
- ✅ Product Management
- ✅ Location Management
- ✅ Stock Movement Management
- ✅ Business Rules Validation
- ✅ RESTful API Design
- ✅ Database Migrations
- ✅ Comprehensive Error Handling

## Business Rules
1. **Stock OUT** tidak boleh melebihi stok tersedia
2. **Stock IN** tidak boleh melebihi kapasitas lokasi
3. Quantity produk auto-update saat ada pergerakan stok
4. Semua endpoint (kecuali login) wajib menggunakan authentication

## Installation & Setup

### Prerequisites
- Go 1.25 or higher
- PostgreSQL 12 or higher

### Environment Variables
Copy `.env.example` to `.env` dan sesuaikan konfigurasi:

```env
SERVICE_NAME=WMS_API
APP_ENVIRONMENT=local
APP_HOST=127.0.0.1
APP_PORT=8000
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASS=your_password
DB_NAME=wms_db
DATABASE_URL=postgres://postgres:your_password@localhost:5432/wms_db?sslmode=disable
JWT_SECRET=your-super-secret-jwt-key
```

### Database Setup
```bash
# Create database
createdb wms_db

# Run migrations
go run cmd/main.go migrate up
```

Seed Database
1. Clear and re-seed users table (useful for development)
```bash
go run cmd/main.go seed refresh
```

### Run the Application
```bash
# Install dependencies
go mod tidy

# Run the server
go run cmd/main.go
```

Server akan berjalan di `http://127.0.0.1:8000`

## API Documentation

### Base URL
```
http://127.0.0.1:8000/api/v1
```

### Authentication
API mendukung dua jenis autentikasi:

#### 1. JWT Token (Recommended)
```bash
# Login untuk mendapatkan token
POST /api/v1/auth/login
{
    "username": "admin",
    "password": "admin123"
}

# Gunakan token di header
Authorization: Bearer <jwt_token>
```

#### 2. API Key
```bash
# Gunakan API key di header
X-API-Key: <your_api_key>
```

### Health Check
```bash
GET /health
```

### Authentication Endpoints

#### Register User
```bash
POST /api/v1/auth/register
Content-Type: application/json

{
    "username": "newuser",
    "email": "user@example.com", 
    "password": "password123"
}
```

#### Login
```bash
POST /api/v1/auth/login
Content-Type: application/json

{
    "username": "admin",
    "password": "admin123"
}
```

#### Get Current User
```bash
GET /api/v1/auth/me
Authorization: Bearer <jwt_token>
```

### Product Endpoints

#### Create Product
```bash
POST /api/v1/products
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
    "sku": "PROD-001",
    "name": "Sample Product",
    "description": "Product description",
    "price": 99.99,
    "weight": 1.5,
    "dimensions": "10x20x5",
    "category": "Electronics"
}
```

#### Get All Products
```bash
GET /api/v1/products?limit=20&offset=0&search=laptop
Authorization: Bearer <jwt_token>
```

#### Get Product by ID
```bash
GET /api/v1/products/{id}
Authorization: Bearer <jwt_token>
```

#### Get Product by SKU
```bash
GET /api/v1/products/sku/{sku}
Authorization: Bearer <jwt_token>
```

#### Update Product
```bash
PUT /api/v1/products/{id}
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
    "name": "Updated Product Name",
    "price": 149.99
}
```

#### Delete Product
```bash
DELETE /api/v1/products/{id}
Authorization: Bearer <jwt_token>
```

### Location Endpoints

#### Create Location
```bash
POST /api/v1/locations
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
    "code": "A-01-01-01",
    "name": "Zone A Aisle 1 Rack 1 Shelf 1",
    "zone": "A",
    "aisle": "01", 
    "rack": "01",
    "shelf": "01",
    "capacity": 100,
    "temperature": 20.5
}
```

#### Get All Locations
```bash
GET /api/v1/locations?limit=20&offset=0&zone=A
Authorization: Bearer <jwt_token>
```

#### Get Location by ID
```bash
GET /api/v1/locations/{id}
Authorization: Bearer <jwt_token>
```

#### Get Location by Code
```bash
GET /api/v1/locations/code/{code}
Authorization: Bearer <jwt_token>
```

#### Update Location
```bash
PUT /api/v1/locations/{id}
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
    "name": "Updated Location Name",
    "capacity": 150
}
```

#### Delete Location
```bash
DELETE /api/v1/locations/{id}
Authorization: Bearer <jwt_token>
```

### Stock Management Endpoints

#### Create Stock Movement
```bash
POST /api/v1/stock-movements
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
    "product_id": 1,
    "location_id": 1,
    "type": "IN",
    "quantity": 50,
    "reference": "PO-2024-001",
    "notes": "Initial stock from supplier"
}
```

#### Get Stock Movements
```bash
GET /api/v1/stock-movements?limit=20&offset=0&product_id=1&type=IN&date_from=2024-01-01&date_to=2024-12-31
Authorization: Bearer <jwt_token>
```

#### Get Stock Movement by ID
```bash
GET /api/v1/stock-movements/{id}
Authorization: Bearer <jwt_token>
```


## API Response Format

### Success Response
```json
{
    "success": true,
    "data": {
        // response data
    },
    "meta": {
        "page": 1,
        "limit": 20,
        "total": 100,
        "total_pages": 5
    }
}
```

### Error Response
```json
{
    "success": false,
    "error": {
        "code": 400,
        "message": "Error message",
        "details": "Detailed error information"
    }
}
```

## Error Codes
- `400` - Bad Request (Invalid input)
- `401` - Unauthorized (Authentication required)
- `404` - Not Found (Resource not found)
- `409` - Conflict (Duplicate entry, insufficient stock, capacity exceeded)
- `500` - Internal Server Error

## Sample Data
Sistem dilengkapi dengan sample data:
- **Default Admin User**: 
  - Username: `admin`
  - Password: `admin123`
  - API Key: `wms_admin_default_api_key_change_in_production`
- **Sample Products**: Laptop, Mouse, Book, Chair
- **Sample Locations**: Multiple zones (A, B) with aisles, racks, and shelves
- **Sample Stock Movements**: Initial inventory transactions

## Architecture

### Database Schema
- **users**: User authentication and authorization
- **products**: Product catalog management  
- **locations**: Warehouse location hierarchy
- **stock_movements**: Historical stock transactions


## Production Deployment
1. Set proper environment variables
2. Change default JWT secret and API keys
3. Use production database credentials
4. Enable HTTPS
5. Set up proper logging and monitoring
6. Configure CORS for your frontend domain

## Menggunakan Docker
```bash
wsl -d Ubuntu

docker build -t wms .
```  

## Menjalankan dengan Docker (jika tersedia Dockerfile):
```bash
docker run --rm -p 8000:8000 \
  -e APP_HOST=0.0.0.0 \
  -e APP_PORT=8000 \
  -e DATABASE_URL="postgres://postgres:aero1996@host.docker.internal:5432/wms_db" \
  wms
```
## License
MIT License
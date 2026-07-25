# Product CRUD — Full-Stack Application

A simple full-stack CRUD application for managing products, built with **React + TypeScript + Vite** (frontend) and **Go + Fiber** (backend).

## Project Structure

```
product-crud/
├── backend/
│   ├── main.go                  # Application entry point
│   ├── go.mod                   # Go module dependencies
│   ├── models/
│   │   └── product.go           # Product model and DTOs
│   ├── store/
│   │   ├── product_store.go     # In-memory data store
│   │   └── product_store_test.go
│   ├── handlers/
│   │   ├── product_handler.go   # HTTP request handlers
│   │   └── product_handler_test.go
│   └── routes/
│       └── product_route.go     # API route definitions
├── frontend/
│   ├── index.html
│   ├── package.json
│   ├── vite.config.ts
│   ├── tsconfig.json
│   └── src/
│       ├── main.tsx
│       ├── App.tsx
│       ├── App.css
│       ├── index.css
│       ├── types/
│       │   └── product.ts       # TypeScript interfaces
│       ├── services/
│       │   └── productService.ts # API client (fetch)
│       └── components/
│           ├── ProductForm.tsx  # Add/Edit form
│           └── ProductList.tsx  # Product table
├── README.md
├── docs/
│   └── MONITORING.md            # Prometheus + Loki + Grafana guide
├── docker-compose.yml           # Docker stack (app + monitoring)
├── docker/
│   ├── backend/Dockerfile
│   ├── frontend/Dockerfile
│   └── nginx/nginx.conf
├── monitoring/
│   ├── prometheus/prometheus.yml
│   ├── loki/loki.yml
│   ├── promtail/promtail.yml
│   └── grafana/
└── .gitignore
```

## Prerequisites

- **Go** 1.22 or later
- **Node.js** 18 or later
- **npm**

## Installation

### Backend

```bash
cd backend
go mod tidy
```

### Frontend

```bash
cd frontend
npm install
```

## Running the Application

Start the backend and frontend in separate terminals.

### Backend (port 8080)

```bash
cd backend
go run .
```

The API will be available at `http://localhost:8080`.

### Frontend (port 5173)

```bash
cd frontend
npm run dev
```

Open `http://localhost:5173` in your browser.

## Running Tests (Backend)

```bash
cd backend
go test ./... -v
```

## Running with Docker (App + Monitoring)

Run the full stack including Prometheus, Loki, and Grafana:

```bash
docker compose up -d --build
```

| URL | Description |
|-----|-------------|
| http://localhost:5173 | Frontend |
| http://localhost:8080/api | Backend API |
| http://localhost:3000 | Grafana (admin / admin) |
| http://localhost:9099 | Prometheus |

See [docs/MONITORING.md](docs/MONITORING.md) for detailed monitoring setup, log queries, and troubleshooting.

Stop the stack:

```bash
docker compose down
```

## API Endpoints

Base URL: `http://localhost:8080/api`

| Method | Endpoint              | Description              |
|--------|-----------------------|--------------------------|
| GET    | `/products`           | Get all products         |
| GET    | `/products/:id`       | Get product by ID        |
| POST   | `/products`           | Create a new product     |
| PUT    | `/products/:id`       | Update product by ID     |
| DELETE | `/products/:id`       | Delete product by ID     |

### Product Model

```json
{
  "id": 1,
  "name": "Wireless Mouse",
  "description": "Ergonomic wireless mouse",
  "price": 29.99,
  "stock": 50
}
```

### Example Requests

#### Create Product (POST /api/products)

**Request Body:**

```json
{
  "name": "Laptop Stand",
  "description": "Adjustable aluminum laptop stand",
  "price": 39.99,
  "stock": 15
}
```

**Response (201 Created):**

```json
{
  "id": 4,
  "name": "Laptop Stand",
  "description": "Adjustable aluminum laptop stand",
  "price": 39.99,
  "stock": 15
}
```

#### Update Product (PUT /api/products/:id)

**Request Body:**

```json
{
  "name": "Laptop Stand Pro",
  "description": "Premium adjustable laptop stand",
  "price": 49.99,
  "stock": 10
}
```

**Response (200 OK):**

```json
{
  "id": 4,
  "name": "Laptop Stand Pro",
  "description": "Premium adjustable laptop stand",
  "price": 49.99,
  "stock": 10
}
```

#### Delete Product (DELETE /api/products/:id)

**Response (200 OK):**

```json
{
  "message": "Product deleted successfully"
}
```

### Validation Rules

- `name` — required, cannot be empty
- `price` — must be >= 0
- `stock` — must be >= 0

### Error Responses

**404 Not Found:**

```json
{
  "error": "Product not found"
}
```

**400 Bad Request:**

```json
{
  "error": "Name is required and cannot be empty"
}
```

## Features

### Backend

- In-memory mock data (3 initial products)
- Auto-incrementing unique product IDs
- Input validation with appropriate HTTP status codes
- CORS enabled for `http://localhost:5173`
- Prometheus metrics at `/metrics` (request rate, latency, Go runtime)

### Frontend

- Product table with Name, Description, Price, Stock, and Actions columns
- Shared form for Add and Edit operations
- Delete with confirmation dialog
- Loading and error states
- Auto-refresh after create, update, or delete
- Responsive layout

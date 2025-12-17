# GoShop - E-commerce Microservices

A production-ready E-commerce backend built with Go, gRPC, and Docker following Clean Architecture principles.

[![CI](https://github.com/herman-xphp/go-microservices-ecommerce/actions/workflows/ci.yml/badge.svg)](https://github.com/herman-xphp/go-microservices-ecommerce/actions/workflows/ci.yml)

## 🏗️ Architecture

```
┌──────────────────────────────────────────────────────────────┐
│                        API Gateway                            │
│                     (Coming Soon)                             │
└────────────────────────────┬─────────────────────────────────┘
                             │
        ┌────────────────────┼────────────────────┐
        │                    │                    │
        ▼                    ▼                    ▼
┌───────────────┐    ┌───────────────┐    ┌───────────────┐
│ Auth Service  │    │Product Service│◄───│ Order Service │
│   :8081/:9091 │    │   :8082/:9092 │gRPC│    :8083      │
└───────┬───────┘    └───────┬───────┘    └───────┬───────┘
        │                    │                    │
        └────────────────────┼────────────────────┘
                             │
                             ▼
                    ┌─────────────────┐
                    │   PostgreSQL    │
                    │     :5432       │
                    └─────────────────┘
```

## 🚀 Services

| Service | HTTP Port | gRPC Port | Description                              |
| ------- | --------- | --------- | ---------------------------------------- |
| Auth    | 8081      | 9091      | User registration, login, JWT validation |
| Product | 8082      | 9092      | Product catalog, inventory management    |
| Order   | 8083      | -         | Order creation, status management        |

## 📋 Features

- ✅ **Clean Architecture** (Handler → Service → Repository)
- ✅ **gRPC Inter-service Communication**
- ✅ **JWT Authentication**
- ✅ **Unit Tests** (16+ tests)
- ✅ **Docker & Docker Compose**
- ✅ **GitHub Actions CI/CD**
- ✅ **Environment-based Configuration**

## 🛠️ Tech Stack

- **Language**: Go 1.22+
- **Framework**: Gin (HTTP), gRPC
- **Database**: PostgreSQL 15
- **ORM**: GORM
- **Auth**: JWT (golang-jwt/jwt/v5)
- **Config**: Environment Variables
- **Container**: Docker, Docker Compose
- **CI/CD**: GitHub Actions

## 🏃 Quick Start

### Prerequisites

- Go 1.22+
- Docker & Docker Compose
- Make

### Setup

1. **Clone the repository**

   ```bash
   git clone https://github.com/herman-xphp/go-microservices-ecommerce.git
   cd go-microservices-ecommerce
   ```

2. **Copy environment file**

   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

3. **Start infrastructure**

   ```bash
   make up
   ```

4. **Run services locally (development)**

   ```bash
   # Terminal 1 - Auth Service
   go run ./cmd/auth-service

   # Terminal 2 - Product Service
   go run ./cmd/product-service

   # Terminal 3 - Order Service
   go run ./cmd/order-service
   ```

## 🧪 Testing

```bash
# Run all tests
go test ./... -v

# Run with coverage
go test ./... -cover -coverprofile=coverage.out

# View coverage report
go tool cover -html=coverage.out
```

## 📡 API Endpoints

### Auth Service (:8081)

| Method | Endpoint              | Description                  |
| ------ | --------------------- | ---------------------------- |
| POST   | /api/v1/auth/register | Register new user            |
| POST   | /api/v1/auth/login    | Login, returns JWT           |
| GET    | /api/v1/auth/profile  | Get user profile (protected) |

### Product Service (:8082)

| Method | Endpoint             | Description               |
| ------ | -------------------- | ------------------------- |
| GET    | /api/v1/products     | List products (paginated) |
| GET    | /api/v1/products/:id | Get product by ID         |
| POST   | /api/v1/products     | Create product            |
| PUT    | /api/v1/products/:id | Update product            |
| DELETE | /api/v1/products/:id | Delete product            |
| GET    | /api/v1/categories   | List categories           |
| POST   | /api/v1/categories   | Create category           |

### Order Service (:8083)

| Method | Endpoint                  | Description         |
| ------ | ------------------------- | ------------------- |
| POST   | /api/v1/orders            | Create order        |
| GET    | /api/v1/orders            | List user orders    |
| GET    | /api/v1/orders/:id        | Get order by ID     |
| PUT    | /api/v1/orders/:id/status | Update order status |
| POST   | /api/v1/orders/:id/cancel | Cancel order        |

## 🔧 Makefile Commands

```bash
# Docker
make up          # Start containers
make down        # Stop containers
make logs        # View logs

# Git (GitFlow)
make git-feature name="xxx"    # Create feature branch from develop
make git-merge-develop         # Merge to develop
make git-release version="x.x.x"  # Create release branch
make git-finish-release        # Finish release (merge to main + develop, tag)
```

## 📁 Project Structure

```
.
├── cmd/                    # Service entry points
│   ├── auth-service/
│   ├── product-service/
│   └── order-service/
├── pkg/                    # Shared packages
│   ├── config/
│   ├── database/
│   └── utils/
├── proto/                  # Protocol Buffer definitions
│   ├── auth/
│   └── product/
├── services/               # Service implementations
│   ├── auth/
│   │   ├── domain/
│   │   ├── dto/
│   │   ├── handler/
│   │   ├── repository/
│   │   ├── service/
│   │   └── grpc/
│   ├── product/
│   └── order/
├── .github/workflows/      # CI/CD pipelines
├── docker-compose.yml
├── Makefile
└── README.md
```

## 📜 License

MIT License

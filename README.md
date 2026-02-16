# List Manager API

[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8E?style=flat&logo=go)](https://go.dev/)
[![MongoDB](https://img.shields.io/badge/MongoDB-8.0-47A248?style=flat&logo=mongodb)](https://www.mongodb.com/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

A production-ready RESTful API for managing lists of items (products/users) with bulk operations support. Built with **Go** following **Clean Architecture** principles, **MongoDB** for persistence, and **Docker** for containerization.

## Table of Contents

- [Features](#features)
- [Tech Stack](#tech-stack)
- [Quick Start](#quick-start)
- [API Endpoints](#api-endpoints)
- [Architecture](#architecture)
- [Documentation](#documentation)
- [Development](#development)
- [Testing](#testing)
- [Deployment](#deployment)
- [Contributing](#contributing)

---

## Features

- **Complete CRUD Operations** for items (products, users, etc.)
- **Bulk Updates** - Update active status for all items at once
- **Health Checks** with MongoDB connectivity verification
- **Clean Architecture** - Separation of concerns with interface-driven design
- **MongoDB Integration** with transaction support
- **Docker Compose** for local development
- **Structured Logging** with Zap
- **OpenAPI Specification** for API documentation

---

## Tech Stack

| Component | Technology | Version |
|-----------|-------------|----------|
| **Language** | Go | 1.24+ |
| **HTTP Router** | Gorilla Mux | Latest |
| **Database** | MongoDB | 8.0 |
| **Logging** | Zap | Latest |
| **Testing** | Testify | Latest |
| **Containerization** | Docker Compose | Latest |

---

## Quick Start

### Prerequisites

- Go 1.24 or higher
- Docker and Docker Compose

### 1. Clone and Run

```bash
# Clone repository
git clone https://github.com/lucaspereirasilva0/list-manager-api.git
cd list-manager-api

# Start MongoDB with Docker Compose
docker-compose up -d mongodb

# Set environment variables
export MONGO_URI=mongodb://localhost:27017/
export MONGO_DB_NAME=listmanager
export PORT=8085

# Run application
go run cmd/api/main.go
```

### 2. Verify Installation

```bash
curl http://localhost:8085/healthz
```

Expected response:
```json
{
  "status": "up",
  "server": "up",
  "database": { "status": "connected" },
  "timestamp": "2026-02-13T10:30:00Z"
}
```

---

## API Endpoints

### Items Management

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/item` | Create a new item |
| `GET` | `/item?id={id}` | Get item by ID |
| `PUT` | `/item?id={id}` | Update an existing item |
| `DELETE` | `/item?id={id}` | Delete an item |
| `GET` | `/items` | List all items |
| `PUT` | `/items/active` | Bulk update active status |

### Health & Application

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/healthz` | Health check with MongoDB verification |
| `GET` | `/_app/version.json` | Application version (PWA auto-update) |

### Example Request

```bash
curl -X POST http://localhost:8085/item \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Product A",
    "active": true,
    "observation": "Premium product"
  }'
```

### Example Response

```json
{
  "id": "a1b2c3d4e5f6",
  "name": "Product A",
  "active": true,
  "observation": "Premium product",
  "createdAt": "2026-02-13T10:30:00Z",
  "updatedAt": "2026-02-13T10:30:00Z"
}
```

---

## Architecture

This project follows **Clean Architecture** with four distinct layers:

```
┌─────────────────────────────────────────┐
│         Handlers (HTTP Layer)            │  ← Request/Response, CORS, Logging
└──────────────────┬──────────────────────┘
                   │ depends on
┌──────────────────▼──────────────────────┐
│        Services (Business Logic)          │  ← Use cases, Orchestration
└──────────────────┬──────────────────────┘
                   │ depends on
┌──────────────────▼──────────────────────┐
│      Repositories (Data Access)          │  ← MongoDB, Local (in-memory)
└──────────────────┬──────────────────────┘
                   │ depends on
┌──────────────────▼──────────────────────┐
│         Domain (Business Entities)        │  ← Pure entities, No dependencies
└─────────────────────────────────────────┘
```

### Project Structure

```
list-manager-api/
├── cmd/                         # Application entry points
│   └── api/
│       ├── main.go                 # Bootstrap and dependency injection
│       ├── handlers/              # HTTP handlers and middleware
│       └── server/               # Server configuration
├── internal/                    # Private application code
│   ├── domain/                  # Business entities
│   ├── repository/              # Data access interfaces
│   │   ├── local/              # In-memory implementation
│   │   └── mongodb/            # MongoDB implementation
│   ├── service/                # Business logic layer
│   └── database/mongodb/       # MongoDB client
├── docs/                       # Documentation
│   ├── adr/                    # Architecture Decision Records
│   ├── diagrams/               # Visual architecture diagrams
│   ├── prd.md                  # Product Requirements Document
│   └── openapi.yaml            # OpenAPI specification
├── docker-compose.yml          # Service orchestration
├── Makefile                    # Build automation
├── CONTRIBUTING.md             # Contributing guidelines
├── CHANGELOG.md               # Version history
└── DEPLOYMENT.md              # Deployment guide
```

---

## Documentation

| Document | Description |
|----------|-------------|
| [Product Requirements Document (PRD)](docs/prd.md) | Functional and non-functional requirements |
| [Architecture Documentation](docs/architecture.md) | Detailed architecture overview |
| [Architecture Decision Records (ADRs)](docs/adr/) | 12 major architectural decisions |
| [OpenAPI Specification](docs/openapi.yaml) | Complete API specification |
| [Architecture Diagrams](docs/diagrams/architecture.md) | Visual architecture representations |
| [Request Flow Diagrams](docs/diagrams/request-flow.md) | Request/response lifecycle |
| [Deployment Diagrams](docs/diagrams/deployment.md) | Deployment architecture |
| [Deployment Guide](DEPLOYMENT.md) | Local, Docker, and Render deployment |
| [Contributing Guidelines](CONTRIBUTING.md) | Development workflow and standards |

---

## Development

### Prerequisites

- Go 1.24+
- Docker and Docker Compose
- Make (for automation commands)

### Setup

```bash
# Install dependencies
go mod download

# Start MongoDB
docker-compose up -d mongodb

# Run tests
make test

# Run with coverage
make coverage

# Run linter
make lint

# Format code
make fmt

# Build binary
make build
```

### Environment Variables

| Variable | Description | Required |
|-----------|-------------|-----------|
| `MONGO_URI` | MongoDB connection string | Yes |
| `MONGO_DB_NAME` | Database name | Yes |
| `PORT` | HTTP server port | No (default: 8080) |

---

## Testing

The project uses a testing strategy prioritizing **unit tests with mocks**:

```bash
# Run all tests
make test

# Run tests with coverage report
make coverage
```

### Testing Approach

- **Unit Tests (80%)**: Fast, isolated tests with mocks
- **Integration Tests (15%)**: Tests with real MongoDB container
- **Manual Tests (5%)**: Smoke tests before deployment

---

## Deployment

### Local Development

```bash
docker-compose up
```

### Docker

```bash
docker build -t list-manager-api .
docker run -p 8085:8085 -e MONGO_URI=mongodb://host.docker.internal:27017/ list-manager-api
```

### Render (Production)

See [DEPLOYMENT.md](DEPLOYMENT.md) for detailed Render deployment instructions.

---

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for:

- Development setup
- Code style and standards
- Testing guidelines
- Commit conventions
- Pull request process

---

## License

This project is licensed under the MIT License.

---

## Links

- [Repository](https://github.com/lucaspereirasilva0/list-manager-api)
- [Issues](https://github.com/lucaspereirasilva0/list-manager-api/issues)
- [Releases](https://github.com/lucaspereirasilva0/list-manager-api/releases)

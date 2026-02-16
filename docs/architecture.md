# Project Architecture Document: List Manager API

## 1. Introduction

This document describes the architecture of the "List Manager API", a robust backend service for managing lists of items, specifically products and users. The project focuses on integration with MongoDB as the primary data storage, implementing CRUD operations and adhering to Clean Architecture principles.

## 2. Goals and Requirements

The main objectives and requirements of this project include:

- Implement a list management API in Go.
- Integrate MongoDB for data persistence.
- Define data models for `Products` and `Users` with fields such as `_id`, `name`, `active` for products and `_id`, `created_by` for users.
- Implement transactional operations when applicable.
- Adhere to Clean Architecture: separation of handlers/controllers, services/use cases, repositories/data access and domain models.
- Follow idiomatic Go practices, modular design, testability, and best backend development practices.
- Ensure proper error handling, dependency injection, and context propagation.
- Develop comprehensive unit and integration tests.

## 3. Architecture Overview

The project architecture follows Clean Architecture principles, promoting separation of concerns and testability. The system is divided into distinct layers:

- **Handlers (or Controllers)**: Responsible for receiving HTTP requests, parsing inputs, and calling the Service layer.
- **Services (or Use Cases)**: Contain the main business logic of the application, orchestrating operations and interacting with the Repository layer.
- **Repositories (or Data Access)**: Abstract the complexity of data persistence, providing interfaces for CRUD operations.
- **Domain**: Defines business entities and their rules.
- **Database (MongoDB)**: The persistent data storage.

Communication between layers is done through interfaces, ensuring low coupling.

## 4. Key Components

- **`cmd/api/handlers`**: Handles HTTP requests, routing, and input validation. Contains `item.go` for item-related operations, `errors.go` for HTTP error handling, `middleware.go` for middlewares, and `model.go` for request/response models.
- **`internal/service`**: Contains business logic and use cases. `item.go` defines services for item operations. `parser.go` for parsing and validating input data.
- **`internal/repository`**: Defines interfaces for abstracting data storage.
  - **`internal/repository/mongodb`**: Concrete implementation of repository interfaces using MongoDB. Includes `repository.go` for CRUD operations of `Product` and `User`.
  - **`internal/repository/local`**: A local repository implementation (in-memory) that can be replaced or augmented by the MongoDB implementation.
- **`internal/domain`**: Contains data model definitions, such as `item.go` for `Product` and `User` structures.
- **`internal/database/mongodb`**: Manages connection and low-level operations with the MongoDB client, including `client.go` and `interfaces.go`.

## 4.1 Key Components Overview and Package Structure

The project structure is organized to follow Clean Architecture principles, with a clear separation of responsibilities. Below, we detail the main components and the organization of their respective packages.

### 4.1.1 Package Structure and Summary

The project is organized in the following directories and packages, each with a specific function:

- **`cmd/`**: Contains the main application entries.
  - **`cmd/api/`**: Where the main API application is defined.
    - **`cmd/api/handlers/`**: Package responsible for handling HTTP requests, routing, and input validation.
      - `cors.go`: Cross-Origin Resource Sharing (CORS) configurations.
      - `errors.go`: Definitions for HTTP error handling.
      - `handlers.go`: General request handler definitions.
      - `handlers_test.go`: Unit tests for HTTP handlers.
      - `item.go`: Specific handlers for item operations (products and users).
      - `middleware.go`: Implements HTTP middlewares for features such as authentication and logging.
      - `model.go`: Data model definitions for HTTP requests and responses.
      - `parser.go`: Utilities for parsing request data.
      - `version.go`: API version information.
    - **`cmd/api/main.go`**: The main entry point of the API application.
    - **`cmd/api/server/`**: Contains HTTP server configuration and initialization.
      - `server.go`: Responsible for configuring and starting the server.
- **`docs/`**: Contains project documentation.
  - `architecture.md`: This architecture document.
- **`internal/`**: Internal application code, not intended to be exposed publicly.
  - **`internal/database/`**: Contains abstraction for database interaction.
    - **`internal/database/mongodb/`**: Specific implementations for MongoDB.
      - `client.go`: Logic for MongoDB client connection and management.
      - `interfaces.go`: Interfaces for MongoDB client operations.
      - `wrappers.go`: "Wrapper" functions for low-level MongoDB operations.
  - **`internal/domain/`**: Defines domain entities and pure business rules.
    - `item.go`: Definition of `Product` and `User` data structures. Includes the `generateID()` function to create IDs compatible with MongoDB's `ObjectID`, and business logic methods such as `IsEmpty()` and `IsActive()`.
    - `item_test.go`: Unit tests for domain entities.
  - **`internal/repository/`**: Data persistence abstraction layer.
    - `errors.go`: Repository-specific errors.
    - `mock.go`: Mock implementations to facilitate repository unit tests.
    - `model.go`: Data models used internally by the repository layer.
    - `repository.go`: Interfaces that define contracts for data persistence operations.
    - **`internal/repository/local/`**: In-memory repository implementation for development/quick tests.
      - `service.go`: Local repository service.
    - **`internal/repository/mongodb/`**: Concrete implementation of repository interfaces for MongoDB.
      - `repository.go`: Logic for persistence of `Product` and `User` in MongoDB. Includes MongoDB transaction implementation for multi-document/collection operations, such as `CreateItemWithUser`, ensuring atomicity.
      - `repository_test.go`: Unit tests for the MongoDB repository.
  - **`internal/service/`**: Contains the main business logic (use cases).
    - `errors.go`: Service-specific errors.
    - `item.go`: Business logic for item operations (products and users).
    - `mock.go`: Mock implementations for service unit tests.
    - `parser.go`: Utilities for parsing and validating data within the service layer.
    - `service.go`: General service definitions.
    - `service_test.go`: Unit tests for services.
- **`memory-bank/`**: Directory containing context documents and project information.
  - `activeContext.md`, `productContext.md`, `progress.md`, `projectbrief.md`, `systemPatterns.md`, `techContext.md`: Various context documents.
- **`docker-compose.yml`**: Docker Compose configuration file for service orchestration (e.g., MongoDB).
- **`go.mod`**: Go module, defines project dependencies.
- **`go.sum`**: Go module dependency checksums.
- **`Makefile`**: Makefile for automating build, test, and deployment tasks.
- **`README.md`**: Project introduction document.

## 5. Data Persistence

MongoDB is the primary database, configured via `docker-compose.yml`. The `Product` and `User` entities are persisted with `bson` tags for correct mapping. Transactional operations are implemented when necessary to ensure data integrity, as seen in the `CreateItemWithUser` function in the MongoDB repository.

## 6. Testing Strategy

- **Unit Tests**: Prioritize extensive unit tests for the repository layer using mocks for MongoDB dependencies. This eliminates the need for a Docker instance for unit tests and ensures isolation.
- **Integration Tests**: Will be considered after unit tests are stable, with potential use of `testcontainers-go` to simulate a real database environment.
- **Manual Tests**: Initial functionality validation through manual tests to confirm basic CRUD operations.

## 7. Error Handling

Robust error handling is ensured, with error encapsulation and propagation to facilitate traceability and debugging. Domain and system-specific errors are handled appropriately.

## 8. Dependency Injection

Dependency injection is performed through constructor functions, ensuring that dependencies are passed explicitly and controlled, which improves testability and code modularity.

## 9. Observability

The application uses `go.uber.org/zap` for structured logging, which facilitates log analysis and debugging in production environments. Logs are configured to provide detailed information about the application flow.

**Future OpenTelemetry Integration**: Plans to integrate OpenTelemetry for distributed tracing and metrics. This integration will allow deeper visibility into performance and request flow through services, complementing existing logging with end-to-end tracing and standardized metrics collection.

## 10. Middleware Order

In `cmd/api/server/server.go`, middlewares are applied in the following order to ensure correct behavior:

1. **`CORSMiddleware`**: Applied first to handle CORS preflight requests before any other middleware or routing logic.
2. **`LoggingMiddleware`**: Applied after CORS to log requests that have passed the CORS check.
3. **Router (`mux.NewRouter()`)**: The router is applied last, ensuring that requests are logged and CORS is handled before routing is performed.

---

## 11. API Endpoints

The application exposes the following endpoints:

### Health & Application

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/healthz` | Health check with MongoDB connectivity verification |
| GET | `/_app/version.json` | Application version for PWA auto-update |

### Items Management

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/item` | Create a new item with name, active status, and optional observation |
| GET | `/item?id={id}` | Get a specific item by ID |
| PUT | `/item?id={id}` | Update an existing item by ID |
| DELETE | `/item?id={id}` | Delete an item by ID |
| GET | `/items` | List all items |
| PUT | `/items/active` | Bulk update active status for all items |

### Endpoint Details

**Create Item (`POST /item`):**
- Request body: `{ "name": "string", "active": boolean, "observation": "string?" }`
- Response: Item with generated ID and timestamps (201 Created)
- Validation: Name cannot be empty

**Get Item (`GET /item?id={id}`):**
- Query parameter: `id` (required)
- Response: Item data (200 OK) or 404 Not Found

**Update Item (`PUT /item?id={id}`):**
- Query parameter: `id` (required)
- Request body: `{ "name": "string", "active": boolean, "observation": "string?" }`
- Response: Updated item (200 OK) or 404 Not Found

**Delete Item (`DELETE /item?id={id}`):**
- Query parameter: `id` (required)
- Response: 204 No Content on success, 404 Not Found

**List Items (`GET /items`):**
- No parameters required
- Response: Array of all items (200 OK)

**Bulk Update Active (`PUT /items/active`):**
- Request body: `{ "active": boolean }`
- Response: `{ "matchedCount": number, "modifiedCount": number }` (200 OK)

**Health Check (`GET /healthz`):**
- No parameters required
- Response: Health status with server, database, and checks (200 OK)
- Verifies MongoDB connectivity before returning "up" status

---

## 12. Data Models

### Item Entity

```go
type Item struct {
    ID          string     // Unique identifier (hex-encoded random bytes)
    Name        string     // Item name (required, non-empty)
    Active      bool       // Active status flag
    Observation *string    // Optional observation/notes (nullable)
    CreatedAt   time.Time  // Creation timestamp
    UpdatedAt   time.Time  // Last update timestamp
}
```

### Health Check Response

```go
type HealthCheckResponse struct {
    Status    HealthStatus           // Overall status: up/degraded/down
    Server    ComponentStatus       // Server status: up/down
    Database  ComponentStatus       // Database connection status: connected/disconnected
    Timestamp string               // ISO 8601 timestamp
    Checks    map[string]Check     // Additional component checks
}
```

### Bulk Update Response

```go
type BulkActiveResponse struct {
    MatchedCount  int64  // Number of documents matched
    ModifiedCount int64  // Number of documents actually modified
}
```

---

## 13. Error Handling Strategy

The application implements comprehensive error handling across all layers:

### Handler Layer
- HTTP error responses with appropriate status codes
- Error middleware wraps all handlers for panic recovery
- JSON error responses: `{ "message": "error description" }`

### Service Layer
- Business logic validation errors
- Repository error wrapping with context
- Domain-specific error types

### Repository Layer
- Database operation errors
- Connection failure handling
- Query result errors (not found, etc.)

### HTTP Status Code Mapping

| Error Type | Status Code | Example |
|------------|-------------|----------|
| Validation Error | 400 Bad Request | Empty item name |
| Not Found | 404 Not Found | Item ID does not exist |
| Internal Error | 500 Internal Server Error | Database connection failure |

---
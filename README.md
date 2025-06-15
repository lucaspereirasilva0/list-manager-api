# List Backend API

A Go microservice that implements basic CRUD (Create, Read, Update, Delete) operations following Clean Architecture principles.

## Project Structure

```
.
├── cmd/                    # Application entry points
├── internal/              # Private application code
│   ├── domain/           # Domain models
│   ├── repository/       # Repository interface
│   └── service/         # Business logic
└── pkg/                  # Public code that can be used by other projects
```

## Prerequisites

- Go 1.24 or higher

## Database Options

The project is structured to support different database implementations through the `ItemRepository` interface. You can implement this interface for any database of your choice:

- PostgreSQL
- MongoDB
- MySQL
- SQLite
- etc.

To add a database:

1. Create a new package in `internal/repository/`
2. Implement the `ItemRepository` interface
3. Create necessary migrations (if applicable)

## Development

To implement a new database:

1. Create a new folder in `internal/repository/` for your implementation
2. Implement the `ItemRepository` interface
3. Add tests for your implementation
4. Update documentation

## Repository Interface

The `ItemRepository` interface defines the following methods that need to be implemented:

```go
type ItemRepository interface {
    Create(ctx context.Context, item *domain.Item) error
    Update(ctx context.Context, item *domain.Item) error
    Delete(ctx context.Context, id uuid.UUID) error
    GetByID(ctx context.Context, id uuid.UUID) (*domain.Item, error)
    List(ctx context.Context) ([]*domain.Item, error)
}
```

## Tests

Run tests with:

```bash
make test
```

## Contributing

1. Fork the project
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request 
# Contributing to List Manager API

Thank you for your interest in contributing to the List Manager API! This document provides guidelines and instructions for contributing to the project.

## Table of Contents

- [Development Setup](#development-setup)
- [Code Style and Standards](#code-style-and-standards)
- [Testing Guidelines](#testing-guidelines)
- [Commit Conventions](#commit-conventions)
- [Pull Request Process](#pull-request-process)
- [Code Review Guidelines](#code-review-guidelines)

---

## Development Setup

### Prerequisites

- **Go:** Version 1.24 or higher
- **Docker:** For running MongoDB locally
- **Make:** For running build automation commands
- **Git:** For version control

### Initial Setup

1. **Clone the repository:**
   ```bash
   git clone https://github.com/lucaspereirasilva0/list-manager-api.git
   cd list-manager-api
   ```

2. **Install dependencies:**
   ```bash
   go mod download
   ```

3. **Start MongoDB using Docker Compose:**
   ```bash
   docker-compose up -d mongodb
   ```

4. **Set environment variables:**
   ```bash
   export MONGO_URI=mongodb://localhost:27017/
   export MONGO_DB_NAME=listmanager
   export PORT=8085
   ```

5. **Run the application:**
   ```bash
   go run cmd/api/main.go
   ```

6. **Verify the application is running:**
   ```bash
   curl http://localhost:8085/healthz
   ```

### Development Workflow

```bash
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

# Run all quality checks
make
```

---

## Code Style and Standards

### Clean Architecture Principles

This project follows **Clean Architecture** with explicit layer separation:

1. **Domain Layer** (`internal/domain/`): Pure business entities without external dependencies
2. **Repository Layer** (`internal/repository/`): Data access abstractions
3. **Service Layer** (`internal/service/`): Business logic and use cases
4. **Handler Layer** (`cmd/api/handlers/`): HTTP request/response handling

**Key Principles:**
- **Dependency Inversion:** Outer layers depend on interfaces defined in inner layers
- **Interface-Driven:** Define small, purpose-specific interfaces
- **Composition Over Inheritance:** Favor small interfaces and composition
- **Single Responsibility:** Each component has one reason to change

### Naming Conventions

| Type | Convention | Example |
|------|-------------|----------|
| Package | `lowercase`, single word | `service`, `repository` |
| Interface | `PascalCase`, often ends with type | `ItemService`, `ItemRepository` |
| Struct | `PascalCase` | `Item`, `ItemHandler` |
| Method | `PascalCase` (exported), `camelCase` (private) | `CreateItem`, `validateInput` |
| Constant | `PascalCase` or `UPPER_SNAKE_CASE` | `HealthStatusUp`, `MaxItems` |
| Variable | `camelCase` | `itemService`, `mongoClient` |

### File Organization

```
internal/service/
├── service.go          # Interface definitions
├── item.go            # Item service implementation
├── errors.go          # Service-specific errors
├── parser.go          # Data parsing utilities
├── mock.go            # Mock implementations for testing
└── service_test.go     # Service tests
```

### Code Formatting

- Use `gofmt` for consistent formatting
- Run `make fmt` before committing
- Maximum line length: 120 characters (soft limit)
- Use `goimports` for import organization

### Error Handling

```go
// ✅ Good: Wrap errors with context
if err := repo.Create(ctx, item); err != nil {
    return fmt.Errorf("failed to create item: %w", err)
}

// ❌ Bad: Discard error context
if err := repo.Create(ctx, item); err != nil {
    return err
}
```

### Comments and Documentation

```go
// ItemService defines the contract for item business operations.
// All methods accept context for timeout/cancellation control.
type ItemService interface {
    // CreateItem creates a new item with the given data.
    // Returns the created item with ID and timestamps generated.
    CreateItem(ctx context.Context, item domain.Item) (domain.Item, error)
}
```

**GoDoc Comments:**
- Every exported package must have a package comment
- Every exported function must have a comment
- Comments should be complete sentences
- Focus on "what" and "why", not "how"

---

## Testing Guidelines

### Test Structure

```go
func TestItemService_CreateItem(t *testing.T) {
    // Setup
    mockRepo := new(MockItemRepository)
    service := NewItemService(mockRepo)

    // Given
    input := domain.Item{Name: "Test Item", Active: true}
    mockRepo.On("Create", mock.Anything, input).Return(
        domain.Item{ID: "123", Name: "Test Item", Active: true}, nil)

    // When
    result, err := service.CreateItem(context.Background(), input)

    // Then
    assert.NoError(t, err)
    assert.Equal(t, "123", result.ID)
    assert.Equal(t, "Test Item", result.Name)
    mockRepo.AssertExpectations(t)
}
```

### Table-Driven Tests

```go
func TestItem_IsActive(t *testing.T) {
    tests := []struct {
        name     string
        item     domain.Item
        expected bool
    }{
        {
            name:     "active item returns true",
            item:     domain.Item{Active: true},
            expected: true,
        },
        {
            name:     "inactive item returns false",
            item:     domain.Item{Active: false},
            expected: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            assert.Equal(t, tt.expected, tt.item.IsActive())
        })
    }
}
```

### BDD Naming Convention

Test names should follow the pattern: `Test[Function]_[Scenario]_[ExpectedResult]`

```go
TestItemService_CreateItem_WithValidData_ReturnsCreatedItem
TestItemService_CreateItem_WithEmptyName_ReturnsError
TestItemService_CreateItem_WhenRepositoryFails_ReturnsError
```

### Coverage Requirements

- Minimum coverage: **80%**
- Aim for: **90%+** for business-critical code
- Use `make coverage` to generate coverage reports

---

## Commit Conventions

This project follows [Conventional Commits](https://www.conventionalcommits.org/).

### Format

```
<type>(<scope>): <subject>

<body>

<footer>
```

### Types

| Type | Description | Examples |
|-------|-------------|-----------|
| `feat` | New feature | `feat: add bulk update endpoint` |
| `fix` | Bug fix | `fix: handle empty observation field` |
| `docs` | Documentation only | `docs: update ADR for MongoDB choice` |
| `style` | Style changes (formatting) | `style: run gofmt on handlers` |
| `refactor` | Code change without functionality change | `refactor: extract logger initialization` |
| `perf` | Performance improvement | `perf: optimize MongoDB query with index` |
| `test` | Adding or updating tests | `test: add integration tests for health check` |
| `chore` | Maintenance tasks | `chore: update Go to 1.24` |

### Examples

```bash
# Feature
git commit -m "feat(items): add bulk update active status endpoint"

# Bug fix
git commit -m "fix(repository): handle duplicate key error on create"

# Documentation
git commit -m "docs(adr): add OpenTelemetry integration decision record"

# Refactor
git commit -m "refactor(service): extract validation logic to separate function"
```

---

## Pull Request Process

### Before Opening a PR

1. **Ensure tests pass:** `make test`
2. **Run linter:** `make lint`
3. **Format code:** `make fmt`
4. **Update documentation:** If applicable
5. **Add tests:** For new features or bug fixes

### PR Title Format

Follow the same format as commits:

```
feat(items): add bulk update active status endpoint
fix(handlers): handle missing query parameter correctly
docs(readme): update deployment instructions
```

### PR Description Template

```markdown
## Summary
Brief description of changes.

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Testing
- [ ] Unit tests added/updated
- [ ] All tests passing
- [ ] Manual testing completed

## Checklist
- [ ] Code follows project style guidelines
- [ ] Self-review completed
- [ ] Comments added to complex code
- [ ] Documentation updated
- [ ] No new warnings generated
```

### Review Process

1. **Automated Checks:** CI runs tests, linting, and coverage
2. **Code Review:** At least one approval required
3. **Changes:** Address review comments
4. **Approval:** PR approved and ready to merge
5. **Merge:** Squash and merge to main branch

---

## Code Review Guidelines

### For Reviewers

**Focus Areas:**
1. Correctness: Does the code work as intended?
2. Architecture: Does it follow Clean Architecture principles?
3. Testing: Are tests comprehensive (edge cases, errors)?
4. Documentation: Is code self-documenting with necessary comments?
5. Style: Does it follow project conventions?

**Review Comments:**
- Be constructive and specific
- Explain the "why" behind suggestions
- Ask questions to understand intent
- Approve once all concerns are addressed

### For Authors

**Before Requesting Review:**
- Self-review your changes
- Ensure CI checks pass
- Update PR description if scope changes
- Keep PRs focused (one feature/concern per PR)

**During Review:**
- Respond to all comments
- Explain your reasoning if you disagree
- Update code based on feedback
- Mark conversations as resolved when done

---

## Versioning

This project follows [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

### Version Format

```
MAJOR.MINOR.PATCH

Example: 1.2.3
```

### Version Components

| Component | When to Bump | Example |
|-----------|--------------|---------|
| **MAJOR** | Incompatible API changes | 1.0.0 → 2.0.0 |
| **MINOR** | Backwards-compatible functionality additions | 1.0.0 → 1.1.0 |
| **PATCH** | Backwards-compatible bug fixes | 1.1.0 → 1.1.1 |

### Release Process

1. **Update version in code:**
   ```bash
   # Update version in cmd/api/handlers/version.go
   ```

2. **Update CHANGELOG.md:**
   - Move changes from `[Unreleased]` to new version section
   - Add release date
   - Include summary of changes

3. **Create and push tag:**
   ```bash
   git tag v1.2.3
   git push origin v1.2.3
   ```

4. **GitHub Actions will:**
   - Create GitHub release automatically
   - Deploy to production (Render)

### Commit-Based Versioning

The CI/CD pipeline uses `mathieudutour/github-tag-action` which analyzes commit messages to determine version bump:

| Commit Message | Version Bump |
|----------------|--------------|
| `feat: ...` | MINOR |
| `fix: ...` | PATCH |
| `perf: ...` | PATCH |
| `feat!: ...` or `BREAKING CHANGE:` | MAJOR |

**Example:**
```bash
# This will trigger a MINOR version bump (1.0.0 → 1.1.0)
git commit -m "feat(items): add bulk update endpoint"

# This will trigger a PATCH version bump (1.1.0 → 1.1.1)
git commit -m "fix(handlers): handle missing query parameter"

# This will trigger a MAJOR version bump (1.1.1 → 2.0.0)
git commit -m "feat!: rename active field to status"
```

---

## Questions?

- Open an issue for bugs or feature requests
- Start a discussion for architectural questions
- Check existing ADRs in `docs/adr/` for context on past decisions

Thank you for contributing! 🚀

# Changelog

All notable changes to the List Manager API project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- OpenAPI specification for all endpoints (`docs/openapi.yaml`)
- Architecture Decision Records (ADRs) for 12 major architectural decisions
- Product Requirements Document (PRD) documenting all functional and non-functional requirements
- Visual architecture diagrams with Mermaid
- Contributing guidelines (`CONTRIBUTING.md`)
- Deployment documentation (`DEPLOYMENT.md`)
- Changelog file (`CHANGELOG.md`)

### Changed
- Updated README with comprehensive documentation structure
- Enhanced architecture documentation with endpoint details

## [1.1.0] - 2025-06-25

### Added
- Observation field to Item model (`observation` string, optional)
- Bulk update endpoint for active status (`PUT /items/active`)
- Enhanced health check with MongoDB connection verification
- Application version endpoint for PWA auto-update (`GET /_app/version.json`)
- Component status tracking in health check response
- Bulk update response with matched and modified counts

### Changed
- Health check now verifies MongoDB connectivity
- Update operations preserve `createdAt` timestamp correctly

### Fixed
- Empty observation field handling in create/update operations

## [1.0.0] - 2025-06-24

### Added
- Initial release of List Manager API
- Complete CRUD operations for items:
  - `POST /item` - Create item
  - `GET /item?id={id}` - Get item by ID
  - `PUT /item?id={id}` - Update item
  - `DELETE /item?id={id}` - Delete item
  - `GET /items` - List all items
- MongoDB integration with repository pattern
- Clean Architecture implementation (4 layers)
- Docker Compose for local development
- Gorilla/Mux HTTP router
- Zap structured logging
- Health check endpoint (`GET /healthz`)
- Environment variable configuration
- Unit tests with testify framework
- Makefile for build automation
- GitHub Actions CI configuration
- Render deployment configuration

---

## Links

- [Repository](https://github.com/lucaspereirasilva0/list-manager-api)
- [Issues](https://github.com/lucaspereirasilva0/list-manager-api/issues)
- [Releases](https://github.com/lucaspereirasilva0/list-manager-api/releases)

---

*Note: For versioning guidelines and release process, see [CONTRIBUTING.md](CONTRIBUTING.md#versioning).*

# Go Project Template

A production-ready Go web application template with clean architecture, dependency injection, and OpenAPI code generation.

## Architecture

This project follows Clean Architecture principles with the following structure:

- `cmd/` - Application entry points
- `internal/app/` - Application composition and dependency injection
- `internal/domain/` - Business logic (entities, repositories, services)
- `internal/infrastructure/` - External concerns (database, HTTP, config)
- `internal/interfaces/` - Interface adapters (API handlers, middleware)
- `pkg/` - Public library code
- `api/` - OpenAPI specifications
- `gen/` - Generated code from OpenAPI specs

## Features

- Clean Architecture with dependency injection (uber/dig)
- OpenAPI code generation
- PostgreSQL with GORM
- XSRF token authentication with header/cookie verification
- Session management
- Docker containerization
- Comprehensive testing setup
- Database migrations
- Structured logging

## Requirements

- [Docker](https://www.docker.com/get-started/)
- [Task](https://taskfile.dev/installation/) - Task runner
- [Node.js](https://nodejs.org/) - For Lefthook installation

## Quick Start

1. Setup the development environment:
   ```bash
   task setup
   ```

2. Start the development environment:
   ```bash
   task up
   ```

   Or run both setup and start in one command:
   ```bash
   task
   ```

## Development Commands

| Command | Description |
|---------|-------------|
| `task` | Setup and start development environment |
| `task setup` | Setup development environment only |
| `task up` | Start services in development mode |
| `task build` | Build the application |
| `task test` | Run all tests |
| `task test:unit` | Run unit tests only |
| `task test:integration` | Run integration tests only |
| `task coverage` | Generate test coverage report |
| `task lint` | Run code linting |
| `task fmt` | Format Go code |
| `task generate` | Generate code from OpenAPI specs |
| `task migrate` | Run database migrations |
| `task seed` | Seed database with test data |
| `task ci` | Run all CI checks locally |
| `task install-tools` | Install development tools |
| `task down` | Stop all services |
| `task destroy` | Remove all containers and data |

## Project Structure

```
├── api/                    # OpenAPI specifications
├── build/                  # Build artifacts and Docker files
├── cmd/                    # Application entry points
├── configs/                # Configuration files
├── deployments/            # Deployment configurations
├── gen/                    # Generated code
├── internal/               # Private application code
│   ├── app/               # Composition root and DI container
│   ├── domain/            # Business logic
│   └── infrastructure/    # External integrations
├── pkg/                   # Public library code
├── scripts/               # Build and utility scripts
└── test/                  # Test configurations and data
```

## Configuration

Set the following environment variables:

- `DATABASE_URL` - PostgreSQL connection string
- `CSRF_SECRET` - Secret key for XSRF token generation
- `SESSION_SECRET` - Secret key for session management
- `PORT` - Server port (default: 8080)

## API Documentation

OpenAPI specifications are located in the `api/` directory. Use `task generate` to regenerate API code after making changes to the specifications.

View API documentation with built-in Swagger UI:
- Auth API: http://localhost:8080/docs/index.html (when server is running)

## XSRF Token Authentication

This application uses XSRF tokens for security:

1. Get XSRF token: `GET /csrf-token`
2. Include token in header: `X-XSRF-TOKEN: <token>`
3. Server verifies header token matches cookie token
4. Required for all POST requests to `/auth/*` and `/api/v1/auth/*`

## CI/CD

This project uses GitHub Actions for continuous integration and deployment:

### Automated Checks
- **Linting**: golangci-lint with comprehensive rule set
- **Testing**: Unit and integration tests with PostgreSQL
- **Security**: Gosec security scanner
- **Code Quality**: Static analysis and formatting checks
- **Coverage**: Code coverage reporting with Codecov

### Workflows
- **CI Pipeline**: Runs on every push and pull request
- **Dependency Updates**: Automated via Dependabot (weekly)
- **Security Scanning**: Integrated into CI pipeline

### Local Development
Run the same checks locally:
```bash
# Install development tools
task install-tools

# Run all CI checks (recommended before committing)
task ci

# Individual checks
task lint          # golangci-lint
task test          # tests with database
task fmt           # format code
```
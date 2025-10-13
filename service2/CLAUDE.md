# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go microservice (service2) that implements a mortgage application management API using clean architecture patterns. It's part of a saga pattern implementation for distributed systems, working alongside service1 (customer service).

## Technology Stack

- **Language**: Go 1.24
- **Web Framework**: Echo v4 (Labstack)
- **Database**: PostgreSQL with pgx/v5 driver
- **Environment**: Docker Compose for local development
- **Testing**: Go standard testing package with database integration tests

## Development Commands

### Running the Service
```bash
# Start PostgreSQL database
docker-compose up -d

# Run the service (requires .env file with DATABASE_URL)
cd api && go run main.go
```

### Testing
```bash
# Run all tests (requires running PostgreSQL instance)
go test ./...

# Run tests for specific package
go test ./api/internal/mortgages

# Run tests with verbose output
go test -v ./api/internal/mortgages
```

### Database Management
```bash
# Start PostgreSQL via Docker Compose
docker-compose up -d postgres

# Connect to database directly
docker exec -it postgres_db_service2 psql -U postgres -d service2_db
```

## Architecture

### Clean Architecture Layers
The codebase follows clean architecture with clear separation of concerns:

1. **Domain Layer** (`api/internal/mortgages/mortgages.go`):
   - `MortgageApplication` struct defines the domain model
   - `Repository` and `Service` interfaces define contracts
   - Contains customer_id reference to link with service1 (customer service)

2. **Repository Layer** (`MortgageRepository`):
   - Database access implementation using pgx driver
   - CRUD operations for mortgage_applications table
   - Additional query to get all applications by customer_id
   - Direct SQL queries for data persistence

3. **Service Layer** (`MortgageService`):
   - Business logic implementation
   - Currently a thin pass-through to repository (ready for business rules)

4. **Handler Layer** (`api/internal/mortgages/handlers.go`):
   - HTTP request/response handling
   - JSON serialization/deserialization
   - Echo framework integration
   - Auto-sets status to "pending" if not provided on creation

5. **Routing** (`api/internal/mortgages/routes.go`):
   - REST API endpoint definitions
   - Route registration with Echo

### Database Schema
- **mortgage_applications** table: id (UUID), customer_id (UUID), loan_amount (numeric), property_value (numeric), interest_rate (numeric), term_years (int), status (varchar), created_at (timestamp), modified_at (timestamp)
- Schema defined in both `schema.sql` and `api/main.go:createMortgageApplicationTable()`

### API Endpoints
- `POST /applications` - Create mortgage application
- `GET /applications/:id` - Read mortgage application by ID
- `PUT /applications/:id` - Update mortgage application
- `DELETE /applications/:id` - Delete mortgage application
- `GET /customers/:customerId/applications` - Get all applications for a specific customer

## Environment Configuration

Required environment variables:
- `DATABASE_URL` - PostgreSQL connection string
- Default for local dev: `postgres://postgres:postgres@localhost:5433/service2_db?sslmode=disable`

## Testing Strategy

Integration tests connect to real PostgreSQL database:
- Database setup/teardown per test
- Schema recreation from `schema.sql`
- Tests cover Repository and Service layers
- Test database URL: configurable via `DATABASE_URL` env var
- Tests include validation of GetByCustomerId functionality

## Key Patterns

1. **Dependency Injection**: Constructor functions (`New*` pattern) for all components
2. **Interface Segregation**: Separate Repository and Service interfaces
3. **Database Connection**: Single pgx.Conn instance passed through dependency chain
4. **Error Handling**: Go idiomatic error returns throughout the stack
5. **UUID Primary Keys**: All entities use UUID for distributed system compatibility

## Development Notes

- Service runs on port 8082 (different from service1 which runs on 8081)
- Database runs on port 5433 (different from service1 which uses 5432)
- Database tables are created automatically on startup if they don't exist
- Tests require a running PostgreSQL instance on localhost:5433
- The customer_id field references customers in service1, but there's no foreign key constraint (microservices pattern)
- Application status can be: "pending", "approved", "rejected"

## Distributed System Integration

This service is part of a distributed microservices architecture:
- **Service1** (port 8081): Customer management service
- **Service2** (port 8082): Mortgage application service (this service)
- Both services have independent databases following the database-per-service pattern
- Services communicate through customer_id references (eventual consistency)

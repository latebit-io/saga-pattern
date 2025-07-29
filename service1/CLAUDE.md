# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go microservice (service1) that implements a customer management API using clean architecture patterns. It's part of a saga pattern implementation for distributed systems.

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
go run main.go
```

### Testing
```bash
# Run all tests (requires running PostgreSQL instance)
go test ./...

# Run tests for specific package
go test ./api/pkg/customers

# Run tests with verbose output
go test -v ./api/pkg/customers
```

### Database Management
```bash
# Start PostgreSQL via Docker Compose
docker-compose up -d postgres

# Connect to database directly
docker exec -it postgres_db psql -U postgres -d service1_db
```

## Architecture

### Clean Architecture Layers
The codebase follows clean architecture with clear separation of concerns:

1. **Domain Layer** (`api/pkg/customers/customers.go`):
   - `Customer` and `Address` structs define the domain models
   - `Repository` and `Service` interfaces define contracts

2. **Repository Layer** (`CustomersRepository`):
   - Database access implementation using pgx driver
   - CRUD operations for customers table
   - Direct SQL queries for data persistence

3. **Service Layer** (`CustomerService`):
   - Business logic implementation
   - Currently a thin pass-through to repository (ready for business rules)

4. **Handler Layer** (`api/pkg/customers/handlers.go`):
   - HTTP request/response handling
   - JSON serialization/deserialization
   - Echo framework integration

5. **Routing** (`api/pkg/customers/routes.go`):
   - REST API endpoint definitions
   - Route registration with Echo

### Database Schema
- **customers** table: id (UUID), name, email
- **addresses** table: id (UUID), customersId (FK), number, street, city, province, postalCode
- Schema defined in both `schema.sql` and `main.go:createCustomerTable()`

### API Endpoints
- `POST /customers` - Create customer
- `GET /customers/:id` - Read customer by ID
- `PUT /customers/:id` - Update customer
- `DELETE /customers/:id` - Delete customer

## Environment Configuration

Required environment variables:
- `DATABASE_URL` - PostgreSQL connection string
- Default for local dev: `postgres://postgres:postgres@localhost:5432/service1_db?sslmode=disable`

## Testing Strategy

Integration tests connect to real PostgreSQL database:
- Database setup/teardown per test
- Schema recreation from `schema.sql`
- Tests cover Repository and Service layers
- Test database URL: configurable via `DATABASE_URL` env var

## Key Patterns

1. **Dependency Injection**: Constructor functions (`New*` pattern) for all components
2. **Interface Segregation**: Separate Repository and Service interfaces
3. **Database Connection**: Single pgx.Conn instance passed through dependency chain
4. **Error Handling**: Go idiomatic error returns throughout the stack
5. **UUID Primary Keys**: All entities use UUID for distributed system compatibility

## Development Notes

- Service runs on port 8081
- Database tables are created automatically on startup if they don't exist
- Tests require a running PostgreSQL instance on localhost:5432
- Address functionality is partially implemented (struct exists but not fully integrated)
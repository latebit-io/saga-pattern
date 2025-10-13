# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go microservice (service3) that implements a loan servicing API using clean architecture patterns. It's part of a saga pattern implementation for distributed systems, working alongside service1 (customer service) and service2 (mortgage application service).

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
go test ./api/internal/loans

# Run tests with verbose output
go test -v ./api/internal/loans
```

### Database Management
```bash
# Start PostgreSQL via Docker Compose
docker-compose up -d postgres

# Connect to database directly
docker exec -it postgres_db_service3 psql -U postgres -d service3_db
```

## Architecture

### Clean Architecture Layers
The codebase follows clean architecture with clear separation of concerns:

1. **Domain Layer** (`api/internal/loans/loans.go` and `api/internal/payments/payments.go`):
   - `Loan` struct defines the loan domain model
   - `Payment` struct defines the payment domain model
   - `Repository` and `Service` interfaces define contracts
   - Loans reference both customer_id (from service1) and mortgage_id (from service2)
   - Payments reference loan_id and customer_id

2. **Repository Layer** (`LoanRepository` and `PaymentRepository`):
   - Database access implementation using pgx driver
   - CRUD operations for loans and payments tables
   - Additional queries:
     - Get loans by customer_id
     - Get loan by mortgage_id
     - Get payments by loan_id
     - Get payments by customer_id
   - Direct SQL queries for data persistence

3. **Service Layer** (`LoanService` and `PaymentService`):
   - Business logic implementation
   - Currently a thin pass-through to repository (ready for business rules)

4. **Handler Layer** (`api/internal/loans/handlers.go` and `api/internal/payments/handlers.go`):
   - HTTP request/response handling
   - JSON serialization/deserialization
   - Echo framework integration
   - Auto-sets status to "active" for new loans if not provided
   - Auto-sets payment_type to "regular" for new payments if not provided

5. **Routing** (`api/internal/loans/routes.go` and `api/internal/payments/routes.go`):
   - REST API endpoint definitions
   - Route registration with Echo

### Database Schema

**loans** table:
- id (UUID)
- customer_id (UUID) - references customer from service1
- mortgage_id (UUID) - references mortgage application from service2
- loan_amount (numeric)
- interest_rate (numeric)
- term_years (int)
- monthly_payment (numeric)
- outstanding_balance (numeric)
- status (varchar: "active", "paid_off", "defaulted")
- start_date (timestamp)
- maturity_date (timestamp)
- created_at (timestamp)
- modified_at (timestamp)

**payments** table:
- id (UUID)
- loan_id (UUID) - references loan
- customer_id (UUID) - references customer from service1
- payment_amount (numeric)
- principal_amount (numeric)
- interest_amount (numeric)
- payment_date (timestamp)
- payment_type (varchar: "regular", "extra", "payoff")
- created_at (timestamp)

Schema defined in both `schema.sql` and `api/main.go` table creation functions.

### API Endpoints

**Loan Endpoints:**
- `POST /loans` - Create loan
- `GET /loans/:id` - Read loan by ID
- `PUT /loans/:id` - Update loan
- `DELETE /loans/:id` - Delete loan
- `GET /customers/:customerId/loans` - Get all loans for a specific customer
- `GET /mortgages/:mortgageId/loan` - Get loan by mortgage application ID

**Payment Endpoints:**
- `POST /payments` - Create payment
- `GET /payments/:id` - Read payment by ID
- `GET /loans/:loanId/payments` - Get all payments for a specific loan
- `GET /customers/:customerId/payments` - Get all payments for a specific customer

## Environment Configuration

Required environment variables:
- `DATABASE_URL` - PostgreSQL connection string
- Default for local dev: `postgres://postgres:postgres@localhost:5434/service3_db?sslmode=disable`

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

- Service runs on port 8083 (different from service1:8081 and service2:8082)
- Database runs on port 5434 (different from service1:5432 and service2:5433)
- Database tables are created automatically on startup if they don't exist
- Tests require a running PostgreSQL instance on localhost:5434
- The customer_id field references customers in service1
- The mortgage_id field references mortgage applications in service2
- No foreign key constraints (microservices pattern for loose coupling)
- Loan status values: "active", "paid_off", "defaulted"
- Payment type values: "regular", "extra", "payoff"

## Distributed System Integration

This service is part of a distributed microservices architecture:
- **Service1** (port 8081): Customer management service
- **Service2** (port 8082): Mortgage application service
- **Service3** (port 8083): Loan servicing service (this service)
- All services have independent databases following the database-per-service pattern
- Services communicate through entity ID references (eventual consistency)
- Loan servicing receives approved mortgages and converts them into active loans
- Payments are recorded against loans to track loan servicing

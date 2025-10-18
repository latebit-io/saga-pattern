# CLAUDE.md

CLAUDE you comments and suggestions should always be from the perspective of a highly skilled staff engineer.

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a **saga-client**, a Go client library that implements the Saga pattern for distributed transactions. It orchestrates multi-step operations across three microservices (service1, service2, service3) and handles compensation/rollback when failures occur.

The key innovation is flexible **compensation strategies** that determine how to handle failures when services are temporarily unavailable during rollback.

## Architecture Overview

### Core Components

**Saga Orchestrator** (`saga.go`)
- Generic `Saga[T any]` type that orchestrates multi-step operations
- Takes a data context struct and a list of `SagaStep[T]` objects
- Each step has an `Execute` function (forward action) and `Compensate` function (rollback action)
- Uses a fluent API for building: `NewSaga(data).AddStep(...).AddStep(...).Execute(ctx)`
- Supports pluggable compensation strategies via `WithCompensationStrategy()`

**Compensation Strategies** (`compensation_strategy.go`)
- `FailFastStrategy`: No retries, stops at first failure (default, original behavior)
- `RetryStrategy`: Retries each compensation with exponential backoff, stops at first failure
- `ContinueAllStrategy`: Retries with exponential backoff, attempts ALL steps even if some fail, collects errors
- All strategies implement `CompensationStrategy[T]` interface
- Custom strategies can be added by implementing the interface

**Customer Saga Example** (`customers_saga.go`)
- Demonstrates the pattern with a 3-step workflow: CreateCustomer → CreateApplication → ExportToServicing
- Data context: `CustomerSagaData` holds shared state (IDs, application details) passed between steps
- Uses `ContinueAllStrategy` with custom retry configuration (3 retries, 2s initial backoff)
- Each step accesses external services via injected clients (customersClient, applicationsClient, servicingClient)

### Service Integration

- `main.go` demonstrates creating clients for three external services (localhost:8081, :8082, :8083)
- Clients are created from service1, service2, service3 modules (local replacements in go.mod)
- Services are called within saga steps to perform distributed operations

## Development Commands

### Build
```bash
go build ./...
```

### Run the example
```bash
go run .
```

### Run tests
```bash
go test ./... -v
```

### Run a single test
```bash
go test ./... -v -run TestNameOrPattern
```

### Run tests with coverage
```bash
go test ./... -v -cover
```

### Format code
```bash
go fmt ./...
```

### Lint code (requires golangci-lint)
```bash
golangci-lint run ./...
```

### Check dependencies
```bash
go mod tidy
go mod verify
```

## Key Patterns and Conventions

### Generic Type Pattern
The saga implementation uses Go generics (`[T any]`) extensively:
- `Saga[T]` - saga instance for data type T
- `SagaStep[T]` - step definition for data type T
- `CompensationStrategy[T]` - strategy implementation for data type T

This allows type-safe saga definitions without losing the saga data between steps.

### Data Context Pattern
All steps operate on a shared data struct (e.g., `CustomerSagaData`):
- Steps read from the context to find previously created IDs
- Steps write results back to the context (e.g., storing IDs)
- Compensation functions can access the same context to know what to clean up

Example: After `CreateCustomer` step sets `data.CustomerID`, the compensation can use it to delete the customer.

### Fluent API
Saga configuration uses method chaining:
```go
NewSaga(data).
    WithCompensationStrategy(strategy).
    AddStep("name", executeFunc, compensateFunc).
    AddStep("name2", executeFunc2, compensateFunc2).
    Execute(ctx)
```

### Error Handling
- **Execution errors**: When a step fails, compensation is triggered automatically
- **Compensation errors**: Depends on strategy:
  - `FailFastStrategy`: Returns first error
  - `ContinueAllStrategy`: Returns `CompensationError` with detailed failures (check with `IsCompensationError()`)
- Use `IsCompensationError(err)` helper to distinguish permanent rollback failures requiring manual intervention

## Testing

### Test Files
- `compensation_strategy_test.go` - Tests all three compensation strategies with various failure scenarios

### Testing Patterns
- Table-driven tests for strategy behavior
- Mock context cancellation scenarios
- Retry behavior verification (attempt counts, backoff timing)
- Multiple failure collection in ContinueAllStrategy

### Important Test Scenarios
When adding new strategies or modifying compensation logic, ensure you test:
1. Success path (all steps execute, no compensation needed)
2. First step fails (no compensation needed)
3. Middle step fails (compensation of prior steps)
4. Last step fails (full compensation chain)
5. Compensation failure (test each strategy's handling)
6. Context cancellation during retries
7. Exponential backoff timing

## Important Files

- `saga.go` - Core saga orchestrator (the main generic saga engine)
- `compensation_strategy.go` - All compensation strategy implementations
- `customers_saga.go` - Example saga implementation showing how to use the library
- `compensation_strategy_test.go` - Comprehensive strategy tests
- `COMPENSATION_STRATEGIES.md` - Documentation on strategies, usage, and implementation details

## Common Tasks

### Add a new saga step
In your saga implementation, chain another `.AddStep()` call:
```go
saga.AddStep("StepName",
    func(ctx context.Context, data *DataType) error { /* execute */ },
    func(ctx context.Context, data *DataType) error { /* compensate */ })
```

### Switch compensation strategy
Change the strategy passed to `WithCompensationStrategy()`. Most common choice for production is `ContinueAllStrategy` with retries.

### Add retry capability to existing saga
Replace `FailFastStrategy` (default) with `RetryStrategy` or `ContinueAllStrategy` and configure `RetryConfig`.

### Create a custom compensation strategy
Implement the `CompensationStrategy[T]` interface in a new file and pass it to `WithCompensationStrategy()`.

## Compensation Strategy Selection Guide

| Strategy | Best For | Behavior |
|----------|----------|----------|
| **FailFastStrategy** | Testing, demos | No retries, stops at first compensation failure |
| **RetryStrategy** | APIs with temporary failures | Retries with backoff, stops at first failure |
| **ContinueAllStrategy** | Production, resilience | Retries all steps, collects all failures for logging/manual intervention |

Reference: See `COMPENSATION_STRATEGIES.md` for detailed examples and behavior documentation.

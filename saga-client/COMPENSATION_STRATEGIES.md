# Compensation Strategies

This saga implementation now supports flexible compensation strategies to handle scenarios where services are down during rollback.

## The Problem

When a saga step fails and compensation is needed, but a service is unavailable during rollback, you can end up with:
- Partial rollback (inconsistent state)
- Lost data
- No way to retry failed compensations

## Available Strategies

### 1. ContinueAllStrategy (Recommended)

**Best for:** Production environments where you need maximum reliability

**Behavior:**
- Retries each compensation with exponential backoff
- Continues trying to compensate ALL steps even if some fail
- Collects all errors for detailed reporting
- Logs all failures for manual intervention

**Example:**
```go
retryConfig := DefaultRetryConfig()
retryConfig.MaxRetries = 3
retryConfig.InitialBackoff = 2 * time.Second

strategy := NewContinueAllStrategy(retryConfig)

saga := NewSaga(data).
    WithCompensationStrategy(strategy).
    AddStep("Step1", exec1, comp1).
    Execute(ctx)
```

**Error Handling:**
```go
if err := saga.Execute(ctx); err != nil {
    if compErr, ok := IsCompensationError(err); ok {
        // Some compensations failed - needs manual intervention
        for _, failure := range compErr.Failures {
            log.Printf("Failed: %s after %d attempts: %v",
                failure.StepName, failure.Attempts, failure.Error)
        }
        // Alert ops team, log to DLQ, etc.
    }
}
```

### 2. RetryStrategy

**Best for:** When you want retries but need to stop at first failure

**Behavior:**
- Retries each compensation with exponential backoff
- Stops at first compensation failure
- Returns immediately with error

**Example:**
```go
retryConfig := DefaultRetryConfig()
strategy := NewRetryStrategy(retryConfig)

saga := NewSaga(data).
    WithCompensationStrategy(strategy).
    AddStep("Step1", exec1, comp1).
    Execute(ctx)
```

### 3. FailFastStrategy (Default)

**Best for:** Testing or when you want original behavior

**Behavior:**
- No retries
- Stops at first compensation failure
- Same as original implementation

**Example:**
```go
saga := NewSaga(data). // Uses FailFastStrategy by default
    AddStep("Step1", exec1, comp1).
    Execute(ctx)

// Or explicitly:
saga := NewSaga(data).
    WithCompensationStrategy(NewFailFastStrategy()).
    AddStep("Step1", exec1, comp1).
    Execute(ctx)
```

## Retry Configuration

```go
type RetryConfig struct {
    MaxRetries      int           // Number of retry attempts
    InitialBackoff  time.Duration // Starting backoff duration
    MaxBackoff      time.Duration // Maximum backoff duration
    BackoffMultiple float64       // Exponential multiplier
}

// Default configuration
DefaultRetryConfig() // 3 retries, 1s initial, 30s max, 2x multiplier

// Custom configuration
custom := RetryConfig{
    MaxRetries:      5,
    InitialBackoff:  2 * time.Second,
    MaxBackoff:      1 * time.Minute,
    BackoffMultiple: 1.5,
}
```

## Usage in customers_saga.go

The customer saga now uses ContinueAllStrategy with custom retry configuration:

```go
func (s *CustomersSaga) CreateCustomer(ctx context.Context, name, email string) error {
    data := &CustomerSagaData{...}

    // Configure compensation strategy
    retryConfig := DefaultRetryConfig()
    retryConfig.MaxRetries = 3
    retryConfig.InitialBackoff = 2 * time.Second
    compensationStrategy := NewContinueAllStrategy(retryConfig)

    // Create and execute saga
    err := NewSaga(data).
        WithCompensationStrategy(compensationStrategy).
        AddStep("CreateCustomer", execFunc, compFunc).
        AddStep("CreateApplication", execFunc2, compFunc2).
        Execute(ctx)

    return err
}
```

## Example Retry Behavior

With MaxRetries=3 and InitialBackoff=2s:

```
Attempt 1: Immediate
Attempt 2: After 2s
Attempt 3: After 4s
Attempt 4: After 8s
Total time: ~14 seconds
```

## Logging Output

### Success:
```
Executed: CreateCustomer
Executed: CreateApplication
Executed: ExportToServicing
```

### Failure with retry:
```
Executed: CreateCustomer
Executed: CreateApplication
Step ExportToServicing failed: service unavailable
⚠️  Compensation failed for CreateApplication (attempt 1/4): connection refused. Retrying in 2s...
⚠️  Compensation failed for CreateApplication (attempt 2/4): connection refused. Retrying in 4s...
✓ Compensated: CreateApplication
✓ Compensated: CreateCustomer
```

### Failure with permanent compensation failure:
```
Executed: CreateCustomer
Executed: CreateApplication
Step ExportToServicing failed: service unavailable
⚠️  Compensation failed for CreateApplication (attempt 1/4): connection refused. Retrying in 2s...
⚠️  Compensation failed for CreateApplication (attempt 2/4): connection refused. Retrying in 4s...
⚠️  Compensation failed for CreateApplication (attempt 3/4): connection refused. Retrying in 8s...
❌ CRITICAL: Compensation failed for CreateApplication after all retries: connection refused
✓ Compensated: CreateCustomer
```

## Adding Your Own Strategy

Implement the `CompensationStrategy` interface:

```go
type CustomStrategy struct {
    // your config
}

func (c *CustomStrategy) Compensate(
    ctx context.Context,
    steps []untypedStep,
    failedStepIndex int,
    data any,
    logger *log.Logger,
) error {
    // Your compensation logic
    for i := failedStepIndex - 1; i >= 0; i-- {
        step := steps[i]
        err := step.Compensate(ctx)
        // Handle error your way
    }
    return nil
}

// Use it
saga := NewSaga(data).
    WithCompensationStrategy(&CustomStrategy{}).
    AddStep(...).
    Execute(ctx)
```

## Next Steps

For even more resilience, consider:
1. Persistent saga log to resume compensation after restarts
2. Background worker to retry permanently failed compensations
3. Dead letter queue for manual intervention
4. Async saga execution with polling (see examples for API timeouts)
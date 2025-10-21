package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

// Example 1: Using ContinueAllStrategy (recommended for most cases)
// This will retry each compensation with exponential backoff, and continue
// trying to compensate all steps even if some fail
func ExampleContinueAllStrategy() {
	data := &CustomerSagaData{
		Name:  "John Doe",
		Email: "john@example.com",
		Application: ApplicationSagaData{
			LoanAmount:     100000,
			PropertyAmount: 200000,
			InterestRate:   3.5,
			TermYears:      30,
		},
	}

	// Configure retry behavior
	retryConfig := DefaultRetryConfig()
	retryConfig.MaxRetries = 3 // Retry up to 3 times
	retryConfig.InitialBackoff = 2 * time.Second
	retryConfig.MaxBackoff = 30 * time.Second

	strategy := NewContinueAllStrategy[CustomerSagaData](retryConfig)

	saga := NewSaga(NewNoStateStore(), uuid.New().String(), data).
		WithCompensationStrategy(strategy).
		AddStep("Step1", executeFunc1, compensateFunc1).
		AddStep("Step2", executeFunc2, compensateFunc2)

	err := saga.Execute(context.Background())
	if err != nil {
		// Check if it's a compensation error with details
		if compErr, ok := IsCompensationError(err); ok {
			log.Printf("Compensation had failures:")
			for _, failure := range compErr.Failures {
				log.Printf("  - Step %s failed after %d attempts: %v",
					failure.StepName, failure.Attempts, failure.Error)
			}
			// At this point, you might want to:
			// 1. Log to a dead letter queue for manual intervention
			// 2. Send alerts to ops team
			// 3. Store in a failure table for background retry worker
		} else {
			log.Printf("Execution failed: %v", err)
		}
	}
}

// Example 2: Using RetryStrategy (stops at first compensation failure)
// This will retry compensations but stop if any compensation fails
func ExampleRetryStrategy() {
	data := &CustomerSagaData{
		Name:  "Jane Smith",
		Email: "jane@example.com",
	}

	retryConfig := DefaultRetryConfig()
	retryConfig.MaxRetries = 5
	retryConfig.InitialBackoff = 1 * time.Second

	strategy := NewRetryStrategy[CustomerSagaData](retryConfig)

	saga := NewSaga(NewNoStateStore(), uuid.New().String(), data).
		WithCompensationStrategy(strategy).
		AddStep("Step1", executeFunc1, compensateFunc1).
		AddStep("Step2", executeFunc2, compensateFunc2)

	err := saga.Execute(context.Background())
	if err != nil {
		log.Printf("Saga failed: %v", err)
		// If compensation failed, you know that at least one step
		// was not compensated and all subsequent steps were skipped
	}
}

// Example 3: Using FailFastStrategy (original behavior, no retries)
// This will fail immediately on the first compensation error
func ExampleFailFastStrategy() {
	data := &CustomerSagaData{
		Name:  "Bob Johnson",
		Email: "bob@example.com",
	}

	strategy := NewFailFastStrategy[CustomerSagaData]()

	saga := NewSaga(NewNoStateStore(), uuid.New().String(), data).
		WithCompensationStrategy(strategy).
		AddStep("Step1", executeFunc1, compensateFunc1).
		AddStep("Step2", executeFunc2, compensateFunc2)

	err := saga.Execute(context.Background())
	if err != nil {
		log.Printf("Saga failed: %v", err)
	}
}

// Example 4: Default behavior (no strategy specified = FailFast)
func ExampleDefaultStrategy() {
	data := &CustomerSagaData{
		Name:  "Alice Williams",
		Email: "alice@example.com",
	}

	// No WithCompensationStrategy() call = uses FailFastStrategy by default
	saga := NewSaga(NewNoStateStore(), uuid.New().String(), data).
		AddStep("Step1", executeFunc1, compensateFunc1).
		AddStep("Step2", executeFunc2, compensateFunc2)

	err := saga.Execute(context.Background())
	if err != nil {
		log.Printf("Saga failed: %v", err)
	}
}

// Example 5: Custom retry configuration for specific use cases
func ExampleCustomRetryConfig() {
	data := &CustomerSagaData{
		Name:  "Charlie Brown",
		Email: "charlie@example.com",
	}

	// Custom configuration for slow/unreliable external services
	retryConfig := RetryConfig{
		MaxRetries:      10,              // Very persistent
		InitialBackoff:  5 * time.Second, // Start with longer wait
		MaxBackoff:      2 * time.Minute, // Cap at 2 minutes
		BackoffMultiple: 1.5,             // Slower exponential growth
	}

	strategy := NewContinueAllStrategy[CustomerSagaData](retryConfig)

	saga := NewSaga(NewNoStateStore(), uuid.New().String(), data).
		WithCompensationStrategy(strategy).
		AddStep("Step1", executeFunc1, compensateFunc1)

	err := saga.Execute(context.Background())
	if err != nil {
		log.Printf("Saga failed: %v", err)
	}

}

// Example 6: Handling different error types in your API
func HandleSagaError(err error) (statusCode int, message string) {
	if err == nil {
		return 200, "Success"
	}

	// Check if it's a compensation error
	if compErr, ok := IsCompensationError(err); ok {
		// Partial failure - some compensations failed
		// This is a critical error that needs manual intervention
		log.Printf("CRITICAL: Compensation failures detected")
		for _, failure := range compErr.Failures {
			log.Printf("  Failed to compensate %s: %v", failure.StepName, failure.Error)
		}

		return 500, fmt.Sprintf("Transaction failed with partial rollback. "+
			"%d step(s) could not be compensated. Please contact support.", len(compErr.Failures))
	}

	// Normal saga failure (rolled back successfully)
	return 400, fmt.Sprintf("Transaction failed: %v", err)
}

// Dummy functions for examples
func executeFunc1(ctx context.Context, data *CustomerSagaData) error {
	return nil
}

func compensateFunc1(ctx context.Context, data *CustomerSagaData) error {
	return nil
}

func executeFunc2(ctx context.Context, data *CustomerSagaData) error {
	return nil
}

func compensateFunc2(ctx context.Context, data *CustomerSagaData) error {
	return nil
}

package main

import (
	"context"
	"fmt"
	"log"
)

// SagaStep represents a single step in the saga with execute and compensate functions
type SagaStep[T any] struct {
	Name       string
	Execute    func(ctx context.Context, data *T) error
	Compensate func(ctx context.Context, data *T) error
}

// Saga represents the saga orchestrator
type Saga[T any] struct {
	Steps  []*SagaStep[T]
	Data   *T
	logger *log.Logger
}

// NewSaga creates a new saga instance
func NewSaga[T any](data *T) *Saga[T] {
	return &Saga[T]{
		Steps:  make([]*SagaStep[T], 0),
		Data:   data,
		logger: log.Default(),
	}
}

// NewSagaWithLogger creates a new saga instance with a custom logger
func NewSagaWithLogger[T any](data *T, logger *log.Logger) *Saga[T] {
	return &Saga[T]{
		Steps:  make([]*SagaStep[T], 0),
		Data:   data,
		logger: logger,
	}
}

// AddStep adds a step to the saga
func (s *Saga[T]) AddStep(name string, execute, compensate func(ctx context.Context, data *T) error) *Saga[T] {
	step := &SagaStep[T]{
		Name:       name,
		Execute:    execute,
		Compensate: compensate,
	}
	s.Steps = append(s.Steps, step)
	return s
}

// Execute runs the saga
func (s *Saga[T]) Execute(ctx context.Context) error {
	s.logger.Printf("Starting saga execution with %d steps", len(s.Steps))

	// Execute all steps
	for i, step := range s.Steps {
		s.logger.Printf("Executing step %d: %s", i+1, step.Name)

		if err := step.Execute(ctx, s.Data); err != nil {
			s.logger.Printf("Step %s failed: %v", step.Name, err)

			// Compensate all executed steps in reverse order
			if compErr := s.compensate(ctx, i); compErr != nil {
				s.logger.Printf("Compensation failed: %v", compErr)
				return fmt.Errorf("execution failed: %w, compensation failed: %w", err, compErr)
			}

			s.logger.Printf("Saga rolled back successfully")
			return fmt.Errorf("saga failed and rolled back: %w", err)
		}

		s.logger.Printf("Step %s completed successfully", step.Name)
	}

	s.logger.Printf("Saga completed successfully")
	return nil
}

// compensate runs compensation for executed steps in reverse order
func (s *Saga[T]) compensate(ctx context.Context, failedStepIndex int) error {
	s.logger.Printf("Starting compensation from step %d", failedStepIndex)

	// Compensate in reverse order, starting from the step before the failed one
	for i := failedStepIndex - 1; i >= 0; i-- {
		step := s.Steps[i]
		s.logger.Printf("Compensating step: %s", step.Name)

		if err := step.Compensate(ctx, s.Data); err != nil {
			return fmt.Errorf("compensation failed for step %s: %w", step.Name, err)
		}

		s.logger.Printf("Step %s compensated successfully", step.Name)
	}

	return nil
}
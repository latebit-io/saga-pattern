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
	for i, step := range s.Steps {
		if err := step.Execute(ctx, s.Data); err != nil {
			s.logger.Printf("Step %s failed: %v", step.Name, err)
			if compErr := s.compensate(ctx, i); compErr != nil {
				return fmt.Errorf("execution failed: %w, compensation failed: %w", err, compErr)
			}
			return fmt.Errorf("saga failed and rolled back: %w", err)
		}
		s.logger.Printf("Executed: %s", step.Name)
	}
	return nil
}

// compensate runs compensation for executed steps in reverse order
func (s *Saga[T]) compensate(ctx context.Context, failedStepIndex int) error {
	for i := failedStepIndex - 1; i >= 0; i-- {
		step := s.Steps[i]
		if err := step.Compensate(ctx, s.Data); err != nil {
			return fmt.Errorf("compensation failed for step %s: %w", step.Name, err)
		}
		s.logger.Printf("Compensated: %s", step.Name)
	}
	return nil
}
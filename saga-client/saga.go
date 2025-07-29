package main

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/google/uuid"
)

type SagaStep struct {
	Name        string
	Execute     func(ctx context.Context, data interface{}) error
	Compensate  func(ctx context.Context, data interface{}) error
	Executed    bool
	Compensated bool
}

// SagaStatus represents the current status of the saga
type SagaStatus string

const (
	StatusPending    SagaStatus = "pending"
	StatusExecuting  SagaStatus = "executing"
	StatusCompleted  SagaStatus = "completed"
	StatusFailed     SagaStatus = "failed"
	StatusRolledBack SagaStatus = "rolled_back"
)

// Saga represents the saga orchestrator
type Saga struct {
	ID     string
	Steps  []*SagaStep
	Status SagaStatus
	Data   interface{}
	mutex  sync.RWMutex
	logger *log.Logger
}

// NewSaga creates a new saga instance
func NewSaga(data interface{}) *Saga {
	return &Saga{
		ID:     uuid.New().String(),
		Steps:  make([]*SagaStep, 0),
		Status: StatusPending,
		Data:   data,
		logger: log.New(log.Writer(), fmt.Sprintf("[Saga-%s] ", uuid.New().String()[:8]), log.LstdFlags),
	}
}

// AddStep adds a step to the saga
func (s *Saga) AddStep(name string, execute, compensate func(ctx context.Context, data interface{}) error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	step := &SagaStep{
		Name:       name,
		Execute:    execute,
		Compensate: compensate,
	}
	s.Steps = append(s.Steps, step)
}

// Execute runs the saga
func (s *Saga) Execute(ctx context.Context) error {
	s.mutex.Lock()
	s.Status = StatusExecuting
	s.mutex.Unlock()

	s.logger.Printf("Starting saga execution with %d steps", len(s.Steps))

	// Execute all steps
	for i, step := range s.Steps {
		s.logger.Printf("Executing step %d: %s", i+1, step.Name)

		if err := step.Execute(ctx, s.Data); err != nil {
			s.logger.Printf("Step %s failed: %v", step.Name, err)
			s.mutex.Lock()
			s.Status = StatusFailed
			s.mutex.Unlock()

			// Compensate all executed steps in reverse order
			if compErr := s.compensate(ctx, i-1); compErr != nil {
				s.logger.Printf("Compensation failed: %v", compErr)
				return fmt.Errorf("execution failed: %v, compensation failed: %v", err, compErr)
			}

			s.mutex.Lock()
			s.Status = StatusRolledBack
			s.mutex.Unlock()
			return fmt.Errorf("saga failed and rolled back: %v", err)
		}

		step.Executed = true
		s.logger.Printf("Step %s completed successfully", step.Name)
	}

	s.mutex.Lock()
	s.Status = StatusCompleted
	s.mutex.Unlock()
	s.logger.Printf("Saga completed successfully")
	return nil
}

// compensate runs compensation for executed steps in reverse order
func (s *Saga) compensate(ctx context.Context, lastExecutedIndex int) error {
	s.logger.Printf("Starting compensation from step %d", lastExecutedIndex+1)

	for i := lastExecutedIndex; i >= 0; i-- {
		step := s.Steps[i]
		if step.Executed && !step.Compensated {
			s.logger.Printf("Compensating step: %s", step.Name)

			if err := step.Compensate(ctx, s.Data); err != nil {
				return fmt.Errorf("compensation failed for step %s: %v", step.Name, err)
			}

			step.Compensated = true
			s.logger.Printf("Step %s compensated successfully", step.Name)
		}
	}

	return nil
}

// GetStatus returns the current saga status
func (s *Saga) GetStatus() SagaStatus {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.Status
}

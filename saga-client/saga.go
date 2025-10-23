package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

type SagaStatus string

const (
	executing    SagaStatus = "EXECUTING"
	compensating SagaStatus = "COMPENSATING"
	complete     SagaStatus = "COMPLETE"
	failed       SagaStatus = "FAILED"
)

// SagaStep represents a single step in the saga with execute and compensate functions
type SagaStep[T any] struct {
	Name       string
	Execute    func(ctx context.Context, data *T) error
	Compensate func(ctx context.Context, data *T) error
}

// Saga represents the saga orchestrator
type Saga[T any] struct {
	SagaID               string
	Steps                []*SagaStep[T]
	Data                 *T
	logger               Logger
	compensationStrategy CompensationStrategy[T]
	stateStore           SagaStateStore
	metadata             map[string]string
	useState             bool
}

type SagaStateStore interface {
	SaveState(ctx context.Context, state *SagaState) error
	LoadState(ctx context.Context, sagaID string) (*SagaState, error)
	MarkComplete(ctx context.Context, sagaID string) error
}

type SagaState struct {
	SagaID            string
	TotalSteps        int
	CurrentStepIndex  int
	Status            SagaStatus
	Data              json.RawMessage
	ExecutedSteps     []string
	FailedStep        string
	CompensatedSteps  []string
	CompensatedStatus SagaStatus
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type Logger interface {
	Log(level string, msg string)
}

type DefaultLogger struct {
	logger *log.Logger
}

func NewDefaultLogger(logger *log.Logger) *DefaultLogger {
	return &DefaultLogger{logger: logger}
}

func (l *DefaultLogger) Log(level string, msg string) {
	l.logger.Printf("%s: %s", level, msg)
}

// NewSaga creates a new saga instance with default FailFast strategy
func NewSaga[T any](stateStore SagaStateStore, sagaID string, data *T) *Saga[T] {
	return &Saga[T]{
		SagaID:               sagaID,
		Steps:                make([]*SagaStep[T], 0),
		Data:                 data,
		stateStore:           stateStore,
		logger:               NewDefaultLogger(log.Default()),
		compensationStrategy: NewFailFastStrategy[T](),
	}
}

// NewSagaWithLogger creates a new saga instance with a custom logger and default FailFast strategy
func NewSagaWithLogger[T any](data *T, logger *log.Logger) *Saga[T] {
	return &Saga[T]{
		Steps:                make([]*SagaStep[T], 0),
		Data:                 data,
		logger:               NewDefaultLogger(log.Default()),
		compensationStrategy: NewFailFastStrategy[T](),
	}
}

// WithCompensationStrategy sets the compensation strategy for the saga (fluent API)
func (s *Saga[T]) WithCompensationStrategy(strategy CompensationStrategy[T]) *Saga[T] {
	s.compensationStrategy = strategy
	return s
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

// LoadState loads a saved state
func (s *Saga[T]) LoadState(sagaID string) *Saga[T] {
	s.useState = false
	// sagaState, err := s.loadState(ctx, s.SagaID)
	// if err != nil {
	// 	s.logger.Log("error", fmt.Sprintf("Failed to load state: %v", err))
	// }

	// if sagaState != nil {
	// 	err = json.Unmarshal(sagaState.Data, s.Data)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	useState = true
	// }
	// s.logger.Log("info", fmt.Sprintf("Using loaded state %t", useState))

	return s
}

// Execute runs the saga
func (s *Saga[T]) Execute(ctx context.Context) error {
	steps := []string{}
	compensated := []string{}
	for i, step := range s.Steps {
		if err := step.Execute(ctx, s.Data); err != nil {
			s.logger.Log("info", fmt.Sprintf("Step %s failed: %v", step.Name, err))
			if compErr := s.compensate(ctx, i); compErr != nil {
				return fmt.Errorf("execution failed: %w, compensation failed: %w", err, compErr)
			}
			// compensated = append(compensated, step.Name)
			// err := s.saveSagaFailed(ctx, steps, compensated, i, step.Name)
			// if err != nil {
			// 	s.logger.Log("info", fmt.Sprintf("Failed to update: %s", err))
			// }
			return fmt.Errorf("saga failed and rolled back: %w", err)
		}
		steps = append(steps, step.Name)
		err := s.saveStepComplete(ctx, steps, i)
		if err != nil {
			s.logger.Log("info", fmt.Sprintf("Failed to update: %s", err))
		}
		s.logger.Log("info", fmt.Sprintf("Executed: %s", step.Name))
	}

	err := s.saveComplete(ctx, steps, compensated, len(s.Steps))
	if err != nil {
		s.logger.Log("info", fmt.Sprintf("Failed to write: %s", err))
	}

	return nil
}

// compensate runs compensation for executed steps using the configured strategy
func (s *Saga[T]) compensate(ctx context.Context, failedStepIndex int) error {
	return s.compensationStrategy.Compensate(ctx, *s, failedStepIndex, s.logger)
}

func (s *Saga[T]) loadState(ctx context.Context, sagaID string) (*SagaState, error) {
	sagaState, err := s.stateStore.LoadState(ctx, sagaID)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(sagaState.Data, s.Data)
	if err != nil {
		return nil, err
	}

	return sagaState, nil
}

func (s *Saga[T]) saveComplete(ctx context.Context, executed []string, compensated []string, index int) error {
	marshaledData, err := json.Marshal(*s.Data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}
	state := SagaState{
		SagaID:           s.SagaID,
		TotalSteps:       len(s.Steps),
		CurrentStepIndex: index,
		Status:           complete,
		Data:             json.RawMessage(marshaledData),
		ExecutedSteps:    executed,
		CompensatedSteps: compensated,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	err = s.stateStore.SaveState(ctx, &state)
	if err != nil {
		return err
	}

	return nil
}

func (s *Saga[T]) saveStepComplete(ctx context.Context, executed []string, index int) error {
	marshaledData, err := json.Marshal(*s.Data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}
	state := SagaState{
		SagaID:           s.SagaID,
		TotalSteps:       len(s.Steps),
		CurrentStepIndex: index,
		Status:           executing,
		Data:             json.RawMessage(marshaledData),
		ExecutedSteps:    executed,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	err = s.stateStore.SaveState(ctx, &state)
	if err != nil {
		return err
	}

	return nil
}

func (s *Saga[T]) saveSagaFailed(ctx context.Context, executed []string, compensated []string, index int, failedStep string) error {
	marshaledData, err := json.Marshal(*s.Data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}
	state := SagaState{
		SagaID:           s.SagaID,
		TotalSteps:       len(s.Steps),
		CurrentStepIndex: index,
		Status:           failed,
		Data:             json.RawMessage(marshaledData),
		ExecutedSteps:    executed,
		FailedStep:       failedStep,
		CompensatedSteps: compensated,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	err = s.stateStore.SaveState(ctx, &state)
	if err != nil {
		return err
	}

	return nil
}

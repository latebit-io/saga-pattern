package main

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jackc/pgx/v5"
)

type SagaStateRecord struct {
	SagaID           string          `db:"saga_id" primaryKey:"true"`
	CorrelationID    string          `db:"correlation_id" index:"true"`
	Status           string          `db:"status" index:"true"`
	CurrentStep      int             `db:"current_step"`
	CompensatedSteps []string        `db:"compensated_steps"`
	ExecutedSteps    []string        `db:"executed_steps"`
	DataJSON         json.RawMessage `db:"data"`
	CreatedAt        time.Time       `db:"created_at"`
	UpdatedAt        time.Time       `db:"updated_at"`
}

type PostgresSagaStore struct {
	pool *pgx.Conn
}

func NewPostgresSagaStore(pool *pgx.Conn) *PostgresSagaStore {
	return &PostgresSagaStore{
		pool: pool,
	}
}

func (s *PostgresSagaStore) SaveState(ctx context.Context, state *SagaState) error {
	query := `
        INSERT INTO saga_states (saga_id, correlation_id, status, current_step, executed_steps, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        ON CONFLICT (saga_id) DO UPDATE
        SET status = $3, current_step = $4, executed_steps = $5, updated_at = $7
    `
	_, err := s.pool.Exec(ctx, query,
		state.SagaID,
		//state.CorrelationID,
		state.Status,
		state.CurrentStepIndex,
		state.ExecutedSteps,
		state.CreatedAt,
		time.Now(),
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresSagaStore) LoadState(ctx context.Context, sagaID string) (*SagaState, error) {
	query := `
        SELECT saga_id, correlation_id, status, current_step, executed_steps, created_at, updated_at
        FROM saga_states
        WHERE saga_id = $1
    `

	state := &SagaState{}
	var executedSteps []string

	err := s.pool.QueryRow(ctx, query, sagaID).Scan(
		&state.SagaID,
		//&state.CorrelationID,
		&state.Status,
		&state.CurrentStepIndex,
		&executedSteps,
		&state.CreatedAt,
		&state.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return state, nil
}

func (s *PostgresSagaStore) MarkComplete(ctx context.Context, sagaID string) error {
	return nil
}

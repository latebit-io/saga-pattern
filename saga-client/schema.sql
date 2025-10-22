-- psql -U postgres -h localhost < schema.sql
CREATE DATABASE saga_db;

-- Connect to the new database and create tables
\c saga_db;


CREATE TABLE saga_states (
    saga_id VARCHAR(36) PRIMARY KEY,
    status VARCHAR(50) NOT NULL,
    total_steps INT NOT NULL DEFAULT 0,
    current_step INT NOT NULL DEFAULT 0,
    executed_steps TEXT[] NOT NULL DEFAULT '{}',
    data JSONB NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Indexes for querying
CREATE INDEX idx_saga_states_status ON saga_states(status);
CREATE INDEX idx_saga_states_updated_at ON saga_states(updated_at);

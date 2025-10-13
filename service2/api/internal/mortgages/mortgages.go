package mortgages

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type MortgageApplication struct {
	Id            uuid.UUID `json:"id"`
	CustomerId    uuid.UUID `json:"customer_id"`
	LoanAmount    float64   `json:"loan_amount"`
	PropertyValue float64   `json:"property_value"`
	InterestRate  float64   `json:"interest_rate"`
	TermYears     int       `json:"term_years"`
	Status        string    `json:"status"` // pending, approved, rejected
	CreatedAt     time.Time `json:"created_at"`
	ModifiedAt    time.Time `json:"modified_at"`
}

type Repository interface {
	Create(ctx context.Context, application MortgageApplication) error
	Read(ctx context.Context, id uuid.UUID) (MortgageApplication, error)
	Update(ctx context.Context, application MortgageApplication) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByCustomerId(ctx context.Context, customerId uuid.UUID) ([]MortgageApplication, error)
}

type Service interface {
	Create(ctx context.Context, application MortgageApplication) error
	Read(ctx context.Context, id uuid.UUID) (MortgageApplication, error)
	Update(ctx context.Context, application MortgageApplication) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByCustomerId(ctx context.Context, customerId uuid.UUID) ([]MortgageApplication, error)
}

type MortgageRepository struct {
	conn *pgx.Conn
}

func NewMortgageRepository(conn *pgx.Conn) *MortgageRepository {
	return &MortgageRepository{conn}
}

func (m *MortgageRepository) Create(ctx context.Context, application MortgageApplication) error {
	sql := `INSERT INTO mortgage_applications
		(id, customer_id, loan_amount, property_value, interest_rate, term_years, status, created_at, modified_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())`

	_, err := m.conn.Exec(ctx, sql,
		application.Id,
		application.CustomerId,
		application.LoanAmount,
		application.PropertyValue,
		application.InterestRate,
		application.TermYears,
		application.Status,
	)
	if err != nil {
		return err
	}
	return nil
}

func (m *MortgageRepository) Read(ctx context.Context, id uuid.UUID) (MortgageApplication, error) {
	sql := `SELECT id, customer_id, loan_amount, property_value, interest_rate, term_years, status, created_at, modified_at
		FROM mortgage_applications WHERE id = $1`
	row := m.conn.QueryRow(ctx, sql, id)
	var application MortgageApplication
	err := row.Scan(
		&application.Id,
		&application.CustomerId,
		&application.LoanAmount,
		&application.PropertyValue,
		&application.InterestRate,
		&application.TermYears,
		&application.Status,
		&application.CreatedAt,
		&application.ModifiedAt,
	)
	if err != nil {
		return MortgageApplication{}, err
	}
	return application, nil
}

func (m *MortgageRepository) Update(ctx context.Context, application MortgageApplication) error {
	sql := `UPDATE mortgage_applications
		SET customer_id = $1, loan_amount = $2, property_value = $3, interest_rate = $4,
			term_years = $5, status = $6, modified_at = NOW()
		WHERE id = $7`
	_, err := m.conn.Exec(ctx, sql,
		application.CustomerId,
		application.LoanAmount,
		application.PropertyValue,
		application.InterestRate,
		application.TermYears,
		application.Status,
		application.Id,
	)
	if err != nil {
		return err
	}
	return nil
}

func (m *MortgageRepository) Delete(ctx context.Context, id uuid.UUID) error {
	sql := "DELETE FROM mortgage_applications WHERE id = $1"
	_, err := m.conn.Exec(ctx, sql, id)
	if err != nil {
		return err
	}
	return nil
}

func (m *MortgageRepository) GetByCustomerId(ctx context.Context, customerId uuid.UUID) ([]MortgageApplication, error) {
	sql := `SELECT id, customer_id, loan_amount, property_value, interest_rate, term_years, status, created_at, modified_at
		FROM mortgage_applications WHERE customer_id = $1 ORDER BY created_at DESC`
	rows, err := m.conn.Query(ctx, sql, customerId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var applications []MortgageApplication
	for rows.Next() {
		var app MortgageApplication
		err := rows.Scan(
			&app.Id,
			&app.CustomerId,
			&app.LoanAmount,
			&app.PropertyValue,
			&app.InterestRate,
			&app.TermYears,
			&app.Status,
			&app.CreatedAt,
			&app.ModifiedAt,
		)
		if err != nil {
			return nil, err
		}
		applications = append(applications, app)
	}
	return applications, nil
}

type MortgageService struct {
	repo Repository
}

func NewMortgageService(repo Repository) *MortgageService {
	return &MortgageService{repo}
}

func (m *MortgageService) Create(ctx context.Context, application MortgageApplication) error {
	return m.repo.Create(ctx, application)
}

func (m *MortgageService) Read(ctx context.Context, id uuid.UUID) (MortgageApplication, error) {
	return m.repo.Read(ctx, id)
}

func (m *MortgageService) Update(ctx context.Context, application MortgageApplication) error {
	return m.repo.Update(ctx, application)
}

func (m *MortgageService) Delete(ctx context.Context, id uuid.UUID) error {
	return m.repo.Delete(ctx, id)
}

func (m *MortgageService) GetByCustomerId(ctx context.Context, customerId uuid.UUID) ([]MortgageApplication, error) {
	return m.repo.GetByCustomerId(ctx, customerId)
}
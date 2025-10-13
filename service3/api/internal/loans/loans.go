package loans

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Loan struct {
	Id                 uuid.UUID `json:"id"`
	CustomerId         uuid.UUID `json:"customer_id"`
	MortgageId         uuid.UUID `json:"mortgage_id"`
	LoanAmount         float64   `json:"loan_amount"`
	InterestRate       float64   `json:"interest_rate"`
	TermYears          int       `json:"term_years"`
	MonthlyPayment     float64   `json:"monthly_payment"`
	OutstandingBalance float64   `json:"outstanding_balance"`
	Status             string    `json:"status"` // active, paid_off, defaulted
	StartDate          time.Time `json:"start_date"`
	MaturityDate       time.Time `json:"maturity_date"`
	CreatedAt          time.Time `json:"created_at"`
	ModifiedAt         time.Time `json:"modified_at"`
}

type Repository interface {
	Create(ctx context.Context, loan Loan) error
	Read(ctx context.Context, id uuid.UUID) (Loan, error)
	Update(ctx context.Context, loan Loan) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByCustomerId(ctx context.Context, customerId uuid.UUID) ([]Loan, error)
	GetByMortgageId(ctx context.Context, mortgageId uuid.UUID) (*Loan, error)
}

type Service interface {
	Create(ctx context.Context, loan Loan) error
	Read(ctx context.Context, id uuid.UUID) (Loan, error)
	Update(ctx context.Context, loan Loan) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByCustomerId(ctx context.Context, customerId uuid.UUID) ([]Loan, error)
	GetByMortgageId(ctx context.Context, mortgageId uuid.UUID) (*Loan, error)
}

type LoanRepository struct {
	conn *pgx.Conn
}

func NewLoanRepository(conn *pgx.Conn) *LoanRepository {
	return &LoanRepository{conn}
}

func (r *LoanRepository) Create(ctx context.Context, loan Loan) error {
	sql := `INSERT INTO loans
		(id, customer_id, mortgage_id, loan_amount, interest_rate, term_years,
		 monthly_payment, outstanding_balance, status, start_date, maturity_date,
		 created_at, modified_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, NOW(), NOW())`

	_, err := r.conn.Exec(ctx, sql,
		loan.Id,
		loan.CustomerId,
		loan.MortgageId,
		loan.LoanAmount,
		loan.InterestRate,
		loan.TermYears,
		loan.MonthlyPayment,
		loan.OutstandingBalance,
		loan.Status,
		loan.StartDate,
		loan.MaturityDate,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *LoanRepository) Read(ctx context.Context, id uuid.UUID) (Loan, error) {
	sql := `SELECT id, customer_id, mortgage_id, loan_amount, interest_rate, term_years,
		monthly_payment, outstanding_balance, status, start_date, maturity_date,
		created_at, modified_at
		FROM loans WHERE id = $1`
	row := r.conn.QueryRow(ctx, sql, id)
	var loan Loan
	err := row.Scan(
		&loan.Id,
		&loan.CustomerId,
		&loan.MortgageId,
		&loan.LoanAmount,
		&loan.InterestRate,
		&loan.TermYears,
		&loan.MonthlyPayment,
		&loan.OutstandingBalance,
		&loan.Status,
		&loan.StartDate,
		&loan.MaturityDate,
		&loan.CreatedAt,
		&loan.ModifiedAt,
	)
	if err != nil {
		return Loan{}, err
	}
	return loan, nil
}

func (r *LoanRepository) Update(ctx context.Context, loan Loan) error {
	sql := `UPDATE loans
		SET customer_id = $1, mortgage_id = $2, loan_amount = $3, interest_rate = $4,
			term_years = $5, monthly_payment = $6, outstanding_balance = $7, status = $8,
			start_date = $9, maturity_date = $10, modified_at = NOW()
		WHERE id = $11`
	_, err := r.conn.Exec(ctx, sql,
		loan.CustomerId,
		loan.MortgageId,
		loan.LoanAmount,
		loan.InterestRate,
		loan.TermYears,
		loan.MonthlyPayment,
		loan.OutstandingBalance,
		loan.Status,
		loan.StartDate,
		loan.MaturityDate,
		loan.Id,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *LoanRepository) Delete(ctx context.Context, id uuid.UUID) error {
	sql := "DELETE FROM loans WHERE id = $1"
	_, err := r.conn.Exec(ctx, sql, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *LoanRepository) GetByCustomerId(ctx context.Context, customerId uuid.UUID) ([]Loan, error) {
	sql := `SELECT id, customer_id, mortgage_id, loan_amount, interest_rate, term_years,
		monthly_payment, outstanding_balance, status, start_date, maturity_date,
		created_at, modified_at
		FROM loans WHERE customer_id = $1 ORDER BY created_at DESC`
	rows, err := r.conn.Query(ctx, sql, customerId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var loans []Loan
	for rows.Next() {
		var loan Loan
		err := rows.Scan(
			&loan.Id,
			&loan.CustomerId,
			&loan.MortgageId,
			&loan.LoanAmount,
			&loan.InterestRate,
			&loan.TermYears,
			&loan.MonthlyPayment,
			&loan.OutstandingBalance,
			&loan.Status,
			&loan.StartDate,
			&loan.MaturityDate,
			&loan.CreatedAt,
			&loan.ModifiedAt,
		)
		if err != nil {
			return nil, err
		}
		loans = append(loans, loan)
	}
	return loans, nil
}

func (r *LoanRepository) GetByMortgageId(ctx context.Context, mortgageId uuid.UUID) (*Loan, error) {
	sql := `SELECT id, customer_id, mortgage_id, loan_amount, interest_rate, term_years,
		monthly_payment, outstanding_balance, status, start_date, maturity_date,
		created_at, modified_at
		FROM loans WHERE mortgage_id = $1`
	row := r.conn.QueryRow(ctx, sql, mortgageId)
	var loan Loan
	err := row.Scan(
		&loan.Id,
		&loan.CustomerId,
		&loan.MortgageId,
		&loan.LoanAmount,
		&loan.InterestRate,
		&loan.TermYears,
		&loan.MonthlyPayment,
		&loan.OutstandingBalance,
		&loan.Status,
		&loan.StartDate,
		&loan.MaturityDate,
		&loan.CreatedAt,
		&loan.ModifiedAt,
	)
	if err != nil {
		return nil, err
	}
	return &loan, nil
}

type LoanService struct {
	repo Repository
}

func NewLoanService(repo Repository) *LoanService {
	return &LoanService{repo}
}

func (s *LoanService) Create(ctx context.Context, loan Loan) error {
	return s.repo.Create(ctx, loan)
}

func (s *LoanService) Read(ctx context.Context, id uuid.UUID) (Loan, error) {
	return s.repo.Read(ctx, id)
}

func (s *LoanService) Update(ctx context.Context, loan Loan) error {
	return s.repo.Update(ctx, loan)
}

func (s *LoanService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *LoanService) GetByCustomerId(ctx context.Context, customerId uuid.UUID) ([]Loan, error) {
	return s.repo.GetByCustomerId(ctx, customerId)
}

func (s *LoanService) GetByMortgageId(ctx context.Context, mortgageId uuid.UUID) (*Loan, error) {
	return s.repo.GetByMortgageId(ctx, mortgageId)
}
package payments

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Payment struct {
	Id              uuid.UUID `json:"id"`
	LoanId          uuid.UUID `json:"loan_id"`
	CustomerId      uuid.UUID `json:"customer_id"`
	PaymentAmount   float64   `json:"payment_amount"`
	PrincipalAmount float64   `json:"principal_amount"`
	InterestAmount  float64   `json:"interest_amount"`
	PaymentDate     time.Time `json:"payment_date"`
	PaymentType     string    `json:"payment_type"` // regular, extra, payoff
	CreatedAt       time.Time `json:"created_at"`
}

type Repository interface {
	Create(ctx context.Context, payment Payment) error
	Read(ctx context.Context, id uuid.UUID) (Payment, error)
	GetByLoanId(ctx context.Context, loanId uuid.UUID) ([]Payment, error)
	GetByCustomerId(ctx context.Context, customerId uuid.UUID) ([]Payment, error)
}

type Service interface {
	Create(ctx context.Context, payment Payment) error
	Read(ctx context.Context, id uuid.UUID) (Payment, error)
	GetByLoanId(ctx context.Context, loanId uuid.UUID) ([]Payment, error)
	GetByCustomerId(ctx context.Context, customerId uuid.UUID) ([]Payment, error)
}

type PaymentRepository struct {
	conn *pgx.Conn
}

func NewPaymentRepository(conn *pgx.Conn) *PaymentRepository {
	return &PaymentRepository{conn}
}

func (r *PaymentRepository) Create(ctx context.Context, payment Payment) error {
	sql := `INSERT INTO payments
		(id, loan_id, customer_id, payment_amount, principal_amount, interest_amount,
		 payment_date, payment_type, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW())`

	_, err := r.conn.Exec(ctx, sql,
		payment.Id,
		payment.LoanId,
		payment.CustomerId,
		payment.PaymentAmount,
		payment.PrincipalAmount,
		payment.InterestAmount,
		payment.PaymentDate,
		payment.PaymentType,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *PaymentRepository) Read(ctx context.Context, id uuid.UUID) (Payment, error) {
	sql := `SELECT id, loan_id, customer_id, payment_amount, principal_amount, interest_amount,
		payment_date, payment_type, created_at
		FROM payments WHERE id = $1`
	row := r.conn.QueryRow(ctx, sql, id)
	var payment Payment
	err := row.Scan(
		&payment.Id,
		&payment.LoanId,
		&payment.CustomerId,
		&payment.PaymentAmount,
		&payment.PrincipalAmount,
		&payment.InterestAmount,
		&payment.PaymentDate,
		&payment.PaymentType,
		&payment.CreatedAt,
	)
	if err != nil {
		return Payment{}, err
	}
	return payment, nil
}

func (r *PaymentRepository) GetByLoanId(ctx context.Context, loanId uuid.UUID) ([]Payment, error) {
	sql := `SELECT id, loan_id, customer_id, payment_amount, principal_amount, interest_amount,
		payment_date, payment_type, created_at
		FROM payments WHERE loan_id = $1 ORDER BY payment_date DESC`
	rows, err := r.conn.Query(ctx, sql, loanId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []Payment
	for rows.Next() {
		var payment Payment
		err := rows.Scan(
			&payment.Id,
			&payment.LoanId,
			&payment.CustomerId,
			&payment.PaymentAmount,
			&payment.PrincipalAmount,
			&payment.InterestAmount,
			&payment.PaymentDate,
			&payment.PaymentType,
			&payment.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		payments = append(payments, payment)
	}
	return payments, nil
}

func (r *PaymentRepository) GetByCustomerId(ctx context.Context, customerId uuid.UUID) ([]Payment, error) {
	sql := `SELECT id, loan_id, customer_id, payment_amount, principal_amount, interest_amount,
		payment_date, payment_type, created_at
		FROM payments WHERE customer_id = $1 ORDER BY payment_date DESC`
	rows, err := r.conn.Query(ctx, sql, customerId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []Payment
	for rows.Next() {
		var payment Payment
		err := rows.Scan(
			&payment.Id,
			&payment.LoanId,
			&payment.CustomerId,
			&payment.PaymentAmount,
			&payment.PrincipalAmount,
			&payment.InterestAmount,
			&payment.PaymentDate,
			&payment.PaymentType,
			&payment.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		payments = append(payments, payment)
	}
	return payments, nil
}

type PaymentService struct {
	repo Repository
}

func NewPaymentService(repo Repository) *PaymentService {
	return &PaymentService{repo}
}

func (s *PaymentService) Create(ctx context.Context, payment Payment) error {
	return s.repo.Create(ctx, payment)
}

func (s *PaymentService) Read(ctx context.Context, id uuid.UUID) (Payment, error) {
	return s.repo.Read(ctx, id)
}

func (s *PaymentService) GetByLoanId(ctx context.Context, loanId uuid.UUID) ([]Payment, error) {
	return s.repo.GetByLoanId(ctx, loanId)
}

func (s *PaymentService) GetByCustomerId(ctx context.Context, customerId uuid.UUID) ([]Payment, error) {
	return s.repo.GetByCustomerId(ctx, customerId)
}

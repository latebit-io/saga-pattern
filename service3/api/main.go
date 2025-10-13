package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"service3/api/internal/loans"
	"service3/api/internal/payments"
)

func main() {
	// Load .env file if it exists (optional - environment variables can also be set via docker-compose)
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
	}
	defer conn.Close(context.Background())

	err = createLoansTable(ctx, conn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create loans table: %v\n", err)
	}

	err = createPaymentsTable(ctx, conn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create payments table: %v\n", err)
	}

	e := echo.New()

	// Loans setup
	loanRepository := loans.NewLoanRepository(conn)
	loanService := loans.NewLoanService(loanRepository)
	loanHandler := loans.NewLoanHandler(loanService)
	loans.Routes(e, loanHandler)

	// Payments setup
	paymentRepository := payments.NewPaymentRepository(conn)
	paymentService := payments.NewPaymentService(paymentRepository)
	paymentHandler := payments.NewPaymentHandler(paymentService)
	payments.Routes(e, paymentHandler)

	e.Logger.Fatal(e.Start(":8083"))
}

func createLoansTable(ctx context.Context, conn *pgx.Conn) error {
	loansTable := `CREATE TABLE IF NOT EXISTS loans(
		id uuid PRIMARY KEY,
		customer_id uuid NOT NULL,
		mortgage_id uuid NOT NULL,
		loan_amount numeric NOT NULL,
		interest_rate numeric NOT NULL,
		term_years int NOT NULL,
		monthly_payment numeric NOT NULL,
		outstanding_balance numeric NOT NULL,
		status varchar NOT NULL,
		start_date timestamp NOT NULL,
		maturity_date timestamp NOT NULL,
		created_at timestamp NOT NULL,
		modified_at timestamp NOT NULL
	)`
	_, err := conn.Exec(ctx, loansTable)
	if err != nil {
		return err
	}

	return nil
}

func createPaymentsTable(ctx context.Context, conn *pgx.Conn) error {
	paymentsTable := `CREATE TABLE IF NOT EXISTS payments(
		id uuid PRIMARY KEY,
		loan_id uuid NOT NULL,
		customer_id uuid NOT NULL,
		payment_amount numeric NOT NULL,
		principal_amount numeric NOT NULL,
		interest_amount numeric NOT NULL,
		payment_date timestamp NOT NULL,
		payment_type varchar NOT NULL,
		created_at timestamp NOT NULL
	)`
	_, err := conn.Exec(ctx, paymentsTable)
	if err != nil {
		return err
	}

	return nil
}

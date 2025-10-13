package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"service2/api/internal/mortgages"
)

func main() {
	// Load .env file if it exists (optional - environment variables can also be set via docker-compose)
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
	}
	defer conn.Close(context.Background())

	err = createMortgageApplicationTable(ctx, conn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create mortgage_applications table: %v\n", err)
	}

	e := echo.New()

	mortgageRepository := mortgages.NewMortgageRepository(conn)
	mortgageService := mortgages.NewMortgageService(mortgageRepository)
	mortgageHandler := mortgages.NewMortgageHandler(mortgageService)
	mortgages.Routes(e, mortgageHandler)

	e.Logger.Fatal(e.Start(":8082"))
}

func createMortgageApplicationTable(ctx context.Context, conn *pgx.Conn) error {
	mortgageApplicationsTable := `CREATE TABLE IF NOT EXISTS mortgage_applications(
		id uuid PRIMARY KEY,
		customer_id uuid NOT NULL,
		loan_amount numeric NOT NULL,
		property_value numeric NOT NULL,
		interest_rate numeric NOT NULL,
		term_years int NOT NULL,
		status varchar NOT NULL,
		created_at timestamp NOT NULL,
		modified_at timestamp NOT NULL
	)`
	_, err := conn.Exec(ctx, mortgageApplicationsTable)
	if err != nil {
		return err
	}

	return nil
}

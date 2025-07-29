package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"service1/api/internal/customers"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
	}
	defer conn.Close(context.Background())

	err = createCustomerTable(ctx, conn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create customer table: %v\n", err)
	}

	e := echo.New()

	customersRepository := customers.NewCustomersRepository(conn)
	customersService := customers.NewCustomerService(customersRepository)
	customersHandler := customers.NewCustomersHandler(customersService)
	customers.Routes(e, customersHandler)

	e.Logger.Fatal(e.Start(":8081"))
}

func createCustomerTable(ctx context.Context, conn *pgx.Conn) error {
	customersTable := `CREATE TABLE IF NOT EXISTS customers(id uuid PRIMARY KEY, name varchar, email varchar)`
	_, err := conn.Exec(ctx, customersTable)
	if err != nil {
		return err
	}

	addressTable := `CREATE TABLE IF NOT EXISTS addresses(id uuid PRIMARY KEY, customersId uuid, number int, street varchar, city varchar, province varchar, postalCode varchar)`
	_, err = conn.Exec(ctx, addressTable)
	if err != nil {
		return err
	}

	return nil
}

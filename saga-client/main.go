package main

import (
	"context"
	"log"
	"os"

	customers "service1/api/pkg/client"
	applictions "service2/api/pkg/client"
	servicing "service3/api/pkg/client"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	pool, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	customersClient := customers.NewClient("http://localhost:8081")
	applicationsClient := applictions.NewClient("http://localhost:8082")
	servicingClient := servicing.NewClient("http://localhost:8083")
	stateStore := NewPostgresSagaStore(pool)
	saga := NewCustomersSaga(stateStore, customersClient, applicationsClient, servicingClient)
	err = saga.CreateCustomer(context.Background(), "John", "john@makes.beats")

	if err != nil {
		panic(err)
	}
}

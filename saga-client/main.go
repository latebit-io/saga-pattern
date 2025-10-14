package main

import (
	"context"

	customers "service1/api/pkg/client"
	applictions "service2/api/pkg/client"
	servicing "service3/api/pkg/client"
)

func main() {
	customersClient := customers.NewClient("http://localhost:8081")
	applicationsClient := applictions.NewClient("http://localhost:8082")
	servicingClient := servicing.NewClient("http://localhost:8083")

	saga := NewCustomersSaga(customersClient, applicationsClient, servicingClient)

	err := saga.CreateCustomer(
		context.Background(),
		"John",
		"john@makes.beats",
	)

	if err != nil {
		panic(err)
	}
}

package main

import (
	"context"

	service1 "service1/api/pkg/client"
)

func main() {
	service1Client := service1.NewClient("http://localhost:8081")
	saga := NewCustomersSaga(service1Client)

	err := saga.CreateCustomer(
		context.Background(),
		"John",
		"john@makes.beats",
	)

	if err != nil {
		panic(err)
	}
}
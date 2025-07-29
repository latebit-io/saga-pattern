package main

import (
	"context"

	"github.com/google/uuid"
	service1 "service1/api/pkg/client"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {
	service1Client := service1.NewClient("http://localhost:8081")
	saga := NewCustomersSaga(service1Client)
	err := saga.CreateCustomer(context.Background(),
		&service1.Customer{
			Id:    uuid.New(),
			Name:  "John",
			Email: "john@makes.beats",
		})
	if err != nil {
		panic(err)
	}
}

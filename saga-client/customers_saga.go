package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	service1 "service1/api/pkg/client"
)

// CustomerSagaData holds the shared data context for the customer saga
// Steps can read from and write to this struct to pass data between steps
type CustomerSagaData struct {
	// Input fields
	Name  string
	Email string

	// Populated by steps during execution
	CustomerID *uuid.UUID // Set by CreateCustomer step
	OrderID    *string    // Set by CreateOrder step (example)
}

type CustomersSaga struct {
	service1Client *service1.Client
}

func NewCustomersSaga(s1 *service1.Client) *CustomersSaga {
	return &CustomersSaga{
		service1Client: s1,
	}
}

func (s *CustomersSaga) CreateCustomer(ctx context.Context, name, email string) error {
	// Initialize the saga data context
	data := &CustomerSagaData{
		Name:  name,
		Email: email,
	}

	// Create and execute the saga
	err := NewSaga(data).
		AddStep(
			"CreateCustomer",
			func(ctx context.Context, data *CustomerSagaData) error {
				// Create customer and store the ID in the saga data
				customer, err := s.service1Client.Create(ctx, data.Name, data.Email)
				if err != nil {
					return fmt.Errorf("failed to create customer: %w", err)
				}
				data.CustomerID = &customer.Id
				return nil
			},
			func(ctx context.Context, data *CustomerSagaData) error {
				// Compensation: delete the customer using the ID from saga data
				if data.CustomerID == nil {
					return nil // Nothing to compensate
				}
				return s.service1Client.Delete(ctx, *data.CustomerID)
			},
		).
		AddStep(
			"CreateOrder",
			func(ctx context.Context, data *CustomerSagaData) error {
				// This step can access the CustomerID from the previous step
				if data.CustomerID == nil {
					return errors.New("customer ID not available")
				}

				// Simulate creating an order (intentionally failing for demo)
				return errors.New("order creation failed - this will trigger rollback")
			},
			func(ctx context.Context, data *CustomerSagaData) error {
				// Compensation: clean up order if it was created
				if data.OrderID != nil {
					// Delete order logic here
					return nil
				}
				return nil
			},
		).
		Execute(ctx)

	return err
}
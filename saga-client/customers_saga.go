package main

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	customers "service1/api/pkg/client"
	applictions "service2/api/pkg/client"
	servicing "service3/api/pkg/client"
)

// CustomerSagaData holds the shared data context for the customer saga
// Steps can read from and write to this struct to pass data between steps
type CustomerSagaData struct {
	// Input fields
	Name  string
	Email string

	// Populated by steps during execution
	CustomerID    *uuid.UUID // Set by CreateCustomer step
	ApplicationID *uuid.UUID

	Application ApplicationSagaData
}

type ApplicationSagaData struct {
	LoanAmount     float64
	PropertyAmount float64
	InterestRate   float64
	TermYears      int
}

type CustomersSaga struct {
	customersClient    *customers.Client
	applicationsClient *applictions.Client
	servicingClient    *servicing.Client
}

func NewCustomersSaga(customers *customers.Client,
	applications *applictions.Client, servicing *servicing.Client) *CustomersSaga {
	return &CustomersSaga{
		customersClient:    customers,
		applicationsClient: applications,
		servicingClient:    servicing,
	}
}

func (s *CustomersSaga) CreateCustomer(ctx context.Context, name, email string) error {
	// Initialize the saga data context
	data := &CustomerSagaData{
		Name:  name,
		Email: email,
		Application: ApplicationSagaData{
			LoanAmount:     1,
			PropertyAmount: 1,
			InterestRate:   1,
			TermYears:      1,
		},
	}

	// Create and execute the saga
	err := NewSaga(data).
		AddStep(
			"CreateCustomer",
			func(ctx context.Context, data *CustomerSagaData) error {
				// Create customer and store the ID in the saga data
				customer, err := s.customersClient.Create(ctx, data.Name, data.Email)
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
				return s.customersClient.Delete(ctx, *data.CustomerID)
			},
		).
		AddStep(
			"CreateApplication",
			func(ctx context.Context, data *CustomerSagaData) error {
				application, err := s.applicationsClient.Create(ctx, *data.CustomerID, data.Application.LoanAmount, data.Application.PropertyAmount, data.Application.InterestRate, data.Application.TermYears)
				if err != nil {
					return fmt.Errorf("failed to create application: %w", err)
				}
				data.ApplicationID = &application.Id
				return nil
			},
			func(ctx context.Context, data *CustomerSagaData) error {
				// Compensation: clean up order if it was created
				if data.ApplicationID != nil {
					return nil
				}
				return s.applicationsClient.Delete(ctx, *data.CustomerID)
			},
		).
		Execute(ctx)

	return err
}

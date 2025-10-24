package main

import (
	"context"
	"fmt"
	"time"

	customers "service1/api/pkg/client"
	applications "service2/api/pkg/client"
	servicing "service3/api/pkg/client"

	"github.com/google/uuid"
)

// CustomerSagaData holds the shared data context for the customer saga
// Steps can read from and write to this struct to pass data between steps
type CustomerSagaData struct {
	Name  string
	Email string

	CustomerID    *uuid.UUID
	ApplicationID *uuid.UUID
	LoanID        *uuid.UUID

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
	applicationsClient *applications.Client
	servicingClient    *servicing.Client
	stateStore         SagaStateStore
}

func NewCustomersSaga(stateStore SagaStateStore, customers *customers.Client,
	applications *applications.Client, servicing *servicing.Client) *CustomersSaga {
	return &CustomersSaga{
		customersClient:    customers,
		applicationsClient: applications,
		servicingClient:    servicing,
		stateStore:         stateStore,
	}
}

func (s *CustomersSaga) CreateCustomer(ctx context.Context, name, email string) error {
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

	// Configure compensation strategy with retry and continue-all behavior
	retryConfig := DefaultRetryConfig()
	retryConfig.MaxRetries = 3
	retryConfig.InitialBackoff = 2 * time.Second

	compensationStrategy := NewContinueAllStrategy[CustomerSagaData](retryConfig)

	// Create and execute the saga
	customerSaga := NewSaga(s.stateStore, uuid.New().String(), data).
		WithCompensationStrategy(compensationStrategy).
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
				if data.ApplicationID == nil {
					return nil
				}
				return s.applicationsClient.Delete(ctx, *data.ApplicationID)
			},
		).
		AddStep(
			"ExportToServicing",
			func(ctx context.Context, data *CustomerSagaData) error {
				return fmt.Errorf("failed to export loan: %w", "error")
				loan, err := s.servicingClient.CreateLoan(ctx, *data.CustomerID, *data.ApplicationID,
					data.Application.LoanAmount, data.Application.InterestRate, data.Application.TermYears,
					float64(100), data.Application.LoanAmount, time.Now(), time.Now().AddDate(1, 0, 0))
				if err != nil {
					return fmt.Errorf("failed to export loan: %w", err)
				}
				data.LoanID = &loan.Id
				return nil
			},
			func(ctx context.Context, data *CustomerSagaData) error {
				if data.LoanID != nil {
					return nil
				}
				return s.servicingClient.DeleteLoan(ctx, *data.LoanID)
			},
		)

	err := customerSaga.Execute(ctx)
	if err != nil {
		return customerSaga.Compensate(ctx)
	}

	return err
}

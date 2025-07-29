package main

import (
	"context"
	"errors"

	service1 "service1/api/pkg/client"
)

type CustomersSaga struct {
	service1Client *service1.Client
}

func NewCustomersSaga(s1 *service1.Client) *CustomersSaga {
	return &CustomersSaga{
		service1Client: s1,
	}
}

func (s *CustomersSaga) CreateCustomer(ctx context.Context, customer *service1.Customer) error {
	saga := NewSaga(customer)

	saga.AddStep(
		"CreateCustomer",
		func(ctx context.Context, data interface{}) error {
			input := data.(*service1.Customer)
			_, err := s.service1Client.Create(ctx, input.Name, input.Email)
			if err != nil {
				return err
			}
			return errors.New("whoops")
		},
		func(ctx context.Context, data interface{}) error {
			input := data.(*service1.Customer)
			return s.service1Client.Delete(ctx, input.Id)
		},
	)

	err := saga.Execute(ctx)
	if err != nil {
		return err
	}

	return nil
}

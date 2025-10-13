package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/google/uuid"
	"service2/api/internal/mortgages"
)

const path = "/applications"

type MortgageApplication = mortgages.MortgageApplication

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL:    baseURL,
		httpClient: &http.Client{},
	}
}

func (c *Client) Create(ctx context.Context, customerId uuid.UUID, loanAmount, propertyValue, interestRate float64, termYears int) (MortgageApplication, error) {
	payload := struct {
		CustomerId    uuid.UUID `json:"customer_id"`
		LoanAmount    float64   `json:"loan_amount"`
		PropertyValue float64   `json:"property_value"`
		InterestRate  float64   `json:"interest_rate"`
		TermYears     int       `json:"term_years"`
	}{
		CustomerId:    customerId,
		LoanAmount:    loanAmount,
		PropertyValue: propertyValue,
		InterestRate:  interestRate,
		TermYears:     termYears,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return MortgageApplication{}, err
	}

	req, err := http.NewRequest(http.MethodPost, c.baseURL+path, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return MortgageApplication{}, err
	}
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return MortgageApplication{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		return MortgageApplication{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	var application MortgageApplication
	err = json.NewDecoder(resp.Body).Decode(&application)
	if err != nil {
		return MortgageApplication{}, err
	}

	return application, nil
}

func (c *Client) Read(ctx context.Context, id uuid.UUID) (MortgageApplication, error) {
	fullURL, err := url.JoinPath(c.baseURL, path, id.String())
	if err != nil {
		return MortgageApplication{}, err
	}

	req, err := http.NewRequest(http.MethodGet, fullURL, nil)
	if err != nil {
		return MortgageApplication{}, err
	}
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return MortgageApplication{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return MortgageApplication{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	var application MortgageApplication
	err = json.NewDecoder(resp.Body).Decode(&application)
	if err != nil {
		return MortgageApplication{}, err
	}
	return application, nil
}

func (c *Client) Update(ctx context.Context, id uuid.UUID, customerId uuid.UUID, loanAmount, propertyValue, interestRate float64, termYears int, status string) (MortgageApplication, error) {
	payload := struct {
		CustomerId    uuid.UUID `json:"customer_id"`
		LoanAmount    float64   `json:"loan_amount"`
		PropertyValue float64   `json:"property_value"`
		InterestRate  float64   `json:"interest_rate"`
		TermYears     int       `json:"term_years"`
		Status        string    `json:"status"`
	}{
		CustomerId:    customerId,
		LoanAmount:    loanAmount,
		PropertyValue: propertyValue,
		InterestRate:  interestRate,
		TermYears:     termYears,
		Status:        status,
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return MortgageApplication{}, err
	}

	fullURL, err := url.JoinPath(c.baseURL, path, id.String())
	if err != nil {
		return MortgageApplication{}, err
	}

	req, err := http.NewRequest(http.MethodPut, fullURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return MortgageApplication{}, err
	}
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)

	if err != nil {
		return MortgageApplication{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return MortgageApplication{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	var application MortgageApplication
	err = json.NewDecoder(resp.Body).Decode(&application)
	if err != nil {
		return MortgageApplication{}, err
	}
	return application, nil
}

func (c *Client) Delete(ctx context.Context, id uuid.UUID) error {
	fullURL, err := url.JoinPath(c.baseURL, path, id.String())
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodDelete, fullURL, nil)
	if err != nil {
		return err
	}
	req = req.WithContext(ctx)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return nil
}

func (c *Client) GetByCustomerId(ctx context.Context, customerId uuid.UUID) ([]MortgageApplication, error) {
	fullURL, err := url.JoinPath(c.baseURL, "/customers", customerId.String(), "applications")
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, fullURL, nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	var applications []MortgageApplication
	err = json.NewDecoder(resp.Body).Decode(&applications)
	if err != nil {
		return nil, err
	}
	return applications, nil
}

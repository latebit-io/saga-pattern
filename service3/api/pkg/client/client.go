package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/google/uuid"
	"service3/api/internal/loans"
	"service3/api/internal/payments"
)

type Loan = loans.Loan
type Payment = payments.Payment

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

// Loan operations

func (c *Client) CreateLoan(ctx context.Context, customerId, mortgageId uuid.UUID, loanAmount, interestRate float64, termYears int, monthlyPayment, outstandingBalance float64, startDate, maturityDate time.Time) (Loan, error) {
	payload := struct {
		CustomerId         uuid.UUID `json:"customer_id"`
		MortgageId         uuid.UUID `json:"mortgage_id"`
		LoanAmount         float64   `json:"loan_amount"`
		InterestRate       float64   `json:"interest_rate"`
		TermYears          int       `json:"term_years"`
		MonthlyPayment     float64   `json:"monthly_payment"`
		OutstandingBalance float64   `json:"outstanding_balance"`
		StartDate          time.Time `json:"start_date"`
		MaturityDate       time.Time `json:"maturity_date"`
	}{
		CustomerId:         customerId,
		MortgageId:         mortgageId,
		LoanAmount:         loanAmount,
		InterestRate:       interestRate,
		TermYears:          termYears,
		MonthlyPayment:     monthlyPayment,
		OutstandingBalance: outstandingBalance,
		StartDate:          startDate,
		MaturityDate:       maturityDate,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return Loan{}, err
	}

	req, err := http.NewRequest(http.MethodPost, c.baseURL+"/loans", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return Loan{}, err
	}
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return Loan{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		return Loan{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	var loan Loan
	err = json.NewDecoder(resp.Body).Decode(&loan)
	if err != nil {
		return Loan{}, err
	}

	return loan, nil
}

func (c *Client) GetLoan(ctx context.Context, id uuid.UUID) (Loan, error) {
	fullURL, err := url.JoinPath(c.baseURL, "/loans", id.String())
	if err != nil {
		return Loan{}, err
	}

	req, err := http.NewRequest(http.MethodGet, fullURL, nil)
	if err != nil {
		return Loan{}, err
	}
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return Loan{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return Loan{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	var loan Loan
	err = json.NewDecoder(resp.Body).Decode(&loan)
	if err != nil {
		return Loan{}, err
	}
	return loan, nil
}

func (c *Client) UpdateLoan(ctx context.Context, id, customerId, mortgageId uuid.UUID, loanAmount, interestRate float64, termYears int, monthlyPayment, outstandingBalance float64, status string, startDate, maturityDate time.Time) (Loan, error) {
	payload := struct {
		CustomerId         uuid.UUID `json:"customer_id"`
		MortgageId         uuid.UUID `json:"mortgage_id"`
		LoanAmount         float64   `json:"loan_amount"`
		InterestRate       float64   `json:"interest_rate"`
		TermYears          int       `json:"term_years"`
		MonthlyPayment     float64   `json:"monthly_payment"`
		OutstandingBalance float64   `json:"outstanding_balance"`
		Status             string    `json:"status"`
		StartDate          time.Time `json:"start_date"`
		MaturityDate       time.Time `json:"maturity_date"`
	}{
		CustomerId:         customerId,
		MortgageId:         mortgageId,
		LoanAmount:         loanAmount,
		InterestRate:       interestRate,
		TermYears:          termYears,
		MonthlyPayment:     monthlyPayment,
		OutstandingBalance: outstandingBalance,
		Status:             status,
		StartDate:          startDate,
		MaturityDate:       maturityDate,
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return Loan{}, err
	}

	fullURL, err := url.JoinPath(c.baseURL, "/loans", id.String())
	if err != nil {
		return Loan{}, err
	}

	req, err := http.NewRequest(http.MethodPut, fullURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return Loan{}, err
	}
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)

	if err != nil {
		return Loan{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return Loan{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	var loan Loan
	err = json.NewDecoder(resp.Body).Decode(&loan)
	if err != nil {
		return Loan{}, err
	}
	return loan, nil
}

func (c *Client) DeleteLoan(ctx context.Context, id uuid.UUID) error {
	fullURL, err := url.JoinPath(c.baseURL, "/loans", id.String())
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

func (c *Client) GetLoansByCustomerId(ctx context.Context, customerId uuid.UUID) ([]Loan, error) {
	fullURL, err := url.JoinPath(c.baseURL, "/customers", customerId.String(), "loans")
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
	var loanList []Loan
	err = json.NewDecoder(resp.Body).Decode(&loanList)
	if err != nil {
		return nil, err
	}
	return loanList, nil
}

func (c *Client) GetLoanByMortgageId(ctx context.Context, mortgageId uuid.UUID) (Loan, error) {
	fullURL, err := url.JoinPath(c.baseURL, "/mortgages", mortgageId.String(), "loan")
	if err != nil {
		return Loan{}, err
	}

	req, err := http.NewRequest(http.MethodGet, fullURL, nil)
	if err != nil {
		return Loan{}, err
	}
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return Loan{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return Loan{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	var loan Loan
	err = json.NewDecoder(resp.Body).Decode(&loan)
	if err != nil {
		return Loan{}, err
	}
	return loan, nil
}

// Payment operations

func (c *Client) CreatePayment(ctx context.Context, loanId, customerId uuid.UUID, paymentAmount, principalAmount, interestAmount float64, paymentDate time.Time, paymentType string) (Payment, error) {
	payload := struct {
		LoanId          uuid.UUID `json:"loan_id"`
		CustomerId      uuid.UUID `json:"customer_id"`
		PaymentAmount   float64   `json:"payment_amount"`
		PrincipalAmount float64   `json:"principal_amount"`
		InterestAmount  float64   `json:"interest_amount"`
		PaymentDate     time.Time `json:"payment_date"`
		PaymentType     string    `json:"payment_type"`
	}{
		LoanId:          loanId,
		CustomerId:      customerId,
		PaymentAmount:   paymentAmount,
		PrincipalAmount: principalAmount,
		InterestAmount:  interestAmount,
		PaymentDate:     paymentDate,
		PaymentType:     paymentType,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return Payment{}, err
	}

	req, err := http.NewRequest(http.MethodPost, c.baseURL+"/payments", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return Payment{}, err
	}
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return Payment{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		return Payment{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	var payment Payment
	err = json.NewDecoder(resp.Body).Decode(&payment)
	if err != nil {
		return Payment{}, err
	}

	return payment, nil
}

func (c *Client) GetPayment(ctx context.Context, id uuid.UUID) (Payment, error) {
	fullURL, err := url.JoinPath(c.baseURL, "/payments", id.String())
	if err != nil {
		return Payment{}, err
	}

	req, err := http.NewRequest(http.MethodGet, fullURL, nil)
	if err != nil {
		return Payment{}, err
	}
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return Payment{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return Payment{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	var payment Payment
	err = json.NewDecoder(resp.Body).Decode(&payment)
	if err != nil {
		return Payment{}, err
	}
	return payment, nil
}

func (c *Client) GetPaymentsByLoanId(ctx context.Context, loanId uuid.UUID) ([]Payment, error) {
	fullURL, err := url.JoinPath(c.baseURL, "/loans", loanId.String(), "payments")
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
	var paymentList []Payment
	err = json.NewDecoder(resp.Body).Decode(&paymentList)
	if err != nil {
		return nil, err
	}
	return paymentList, nil
}

func (c *Client) GetPaymentsByCustomerId(ctx context.Context, customerId uuid.UUID) ([]Payment, error) {
	fullURL, err := url.JoinPath(c.baseURL, "/customers", customerId.String(), "payments")
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
	var paymentList []Payment
	err = json.NewDecoder(resp.Body).Decode(&paymentList)
	if err != nil {
		return nil, err
	}
	return paymentList, nil
}

package client

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/google/uuid"
	"service1/api/internal/customers"
)

const path = "/customers"

type Customer = customers.Customer

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

func (c *Client) Create(ctx context.Context, name, email string) (Customer, error) {
	payload := struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}{
		Name:  name,
		Email: email,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return Customer{}, err
	}

	req, err := http.NewRequest(http.MethodPost, c.baseURL+path, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return Customer{}, err
	}
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return Customer{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		return Customer{}, err
	}
	var customer Customer
	err = json.NewDecoder(resp.Body).Decode(&customer)
	if err != nil {
		return Customer{}, err
	}

	return customer, nil
}

func (c *Client) Read(ctx context.Context, id uuid.UUID) (Customer, error) {
	fullURL, err := url.JoinPath(c.baseURL, path, id.String())
	if err != nil {
		return Customer{}, err
	}

	req, err := http.NewRequest(http.MethodGet, fullURL, nil)
	if err != nil {
		return Customer{}, err
	}
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return Customer{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return Customer{}, err
	}
	var customer Customer
	err = json.NewDecoder(resp.Body).Decode(&customer)
	if err != nil {
		return Customer{}, err
	}
	return customer, nil
}

func (c *Client) Update(ctx context.Context, id uuid.UUID, name, email string) (Customer, error) {
	payload := struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}{
		Name:  name,
		Email: email,
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return Customer{}, err
	}

	fullURL, err := url.JoinPath(c.baseURL, path, id.String())
	if err != nil {
		return Customer{}, err
	}

	req, err := http.NewRequest(http.MethodPut, fullURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return Customer{}, err
	}
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)

	if err != nil {
		return Customer{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return Customer{}, err
	}
	var customer Customer
	err = json.NewDecoder(resp.Body).Decode(&customer)
	if err != nil {
		return Customer{}, err
	}
	return customer, nil
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
		return err
	}
	return nil
}

package customers

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Customer struct {
	Id         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	CreatedAt  time.Time `json:"created_at"`
	ModifiedAt time.Time `json:"modified_at"`
}

type Repository interface {
	Create(ctx context.Context, customer Customer) error
	Read(ctx context.Context, id uuid.UUID) (Customer, error)
	Update(ctx context.Context, customer Customer) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type Service interface {
	Create(ctx context.Context, customer Customer) error
	Read(ctx context.Context, id uuid.UUID) (Customer, error)
	Update(ctx context.Context, customer Customer) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type CustomersRepository struct {
	conn *pgx.Conn
}

func NewCustomersRepository(conn *pgx.Conn) *CustomersRepository {
	return &CustomersRepository{conn}
}

func (c *CustomersRepository) Create(ctx context.Context, customer Customer) error {
	sql := "INSERT INTO customers (id, name, email, created_at, modified_at) VALUES ($1, $2, $3, NOW(), NOW())"

	_, err := c.conn.Exec(ctx, sql, customer.Id, customer.Name, customer.Email)
	if err != nil {
		return err
	}
	return nil
}

func (c *CustomersRepository) Read(ctx context.Context, id uuid.UUID) (Customer, error) {
	sql := "SELECT id, name, email, created_at, modified_at FROM customers WHERE id = $1"
	row := c.conn.QueryRow(ctx, sql, id)
	var customer Customer
	err := row.Scan(&customer.Id, &customer.Name, &customer.Email, &customer.CreatedAt, &customer.ModifiedAt)
	if err != nil {
		return Customer{}, err
	}
	return customer, nil
}

func (c *CustomersRepository) Update(ctx context.Context, customer Customer) error {
	sql := "UPDATE customers SET name = $1, email = $2, modified_at = NOW() WHERE id = $3"
	_, err := c.conn.Exec(ctx, sql, customer.Name, customer.Email, customer.Id)
	if err != nil {
		return err
	}
	return nil
}

func (c *CustomersRepository) Delete(ctx context.Context, id uuid.UUID) error {
	sql := "DELETE FROM customers WHERE id = $1"
	_, err := c.conn.Exec(ctx, sql, id)
	if err != nil {
		return err
	}
	return nil
}

type CustomerService struct {
	repo Repository
}

func NewCustomerService(repo Repository) *CustomerService {
	return &CustomerService{repo}
}

func (c *CustomerService) Create(ctx context.Context, customer Customer) error {
	return c.repo.Create(ctx, customer)
}

func (c *CustomerService) Read(ctx context.Context, id uuid.UUID) (Customer, error) {
	return c.repo.Read(ctx, id)
}

func (c *CustomerService) Update(ctx context.Context, customer Customer) error {
	return c.repo.Update(ctx, customer)
}

func (c *CustomerService) Delete(ctx context.Context, id uuid.UUID) error {
	return c.repo.Delete(ctx, id)
}

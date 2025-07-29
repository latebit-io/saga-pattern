package customers

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func setupTestDB(t *testing.T) *pgx.Conn {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/service1_db?sslmode=disable"
	}

	conn, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	_, err = conn.Exec(context.Background(), "DROP TABLE IF EXISTS customers")
	if err != nil {
		t.Fatalf("Failed to drop existing customers table: %v", err)
	}

	schemaPath := filepath.Join("..", "..", "..", "schema.sql")
	schemaFile, err := os.Open(schemaPath)
	if err != nil {
		t.Fatalf("Failed to open schema.sql: %v", err)
	}
	defer schemaFile.Close()

	schemaSQL, err := io.ReadAll(schemaFile)
	if err != nil {
		t.Fatalf("Failed to read schema.sql: %v", err)
	}

	_, err = conn.Exec(context.Background(), string(schemaSQL))
	if err != nil {
		t.Fatalf("Failed to execute schema.sql: %v", err)
	}

	return conn
}

func teardownTestDB(t *testing.T, conn *pgx.Conn) {
	_, err := conn.Exec(context.Background(), "DELETE FROM customers")
	if err != nil {
		t.Errorf("Failed to clean up test data: %v", err)
	}
	conn.Close(context.Background())
}

func TestCustomersRepository_Create(t *testing.T) {
	conn := setupTestDB(t)
	defer teardownTestDB(t, conn)

	repo := NewCustomersRepository(conn)
	customer := Customer{
		Id:    uuid.New(),
		Name:  "John Doe",
		Email: "john@example.com",
	}

	err := repo.Create(context.Background(), customer)
	if err != nil {
		t.Errorf("Create failed: %v", err)
	}

	retrievedCustomer, err := repo.Read(context.Background(), customer.Id)
	if err != nil {
		t.Errorf("Read failed: %v", err)
	}

	if retrievedCustomer.Id != customer.Id {
		t.Errorf("Expected ID %v, got %v", customer.Id, retrievedCustomer.Id)
	}
	if retrievedCustomer.Name != customer.Name {
		t.Errorf("Expected Name %v, got %v", customer.Name, retrievedCustomer.Name)
	}
	if retrievedCustomer.Email != customer.Email {
		t.Errorf("Expected Email %v, got %v", customer.Email, retrievedCustomer.Email)
	}
}

func TestCustomersRepository_Read_NotFound(t *testing.T) {
	conn := setupTestDB(t)
	defer teardownTestDB(t, conn)

	repo := NewCustomersRepository(conn)
	nonExistentID := uuid.New()

	_, err := repo.Read(context.Background(), nonExistentID)
	if err == nil {
		t.Error("Expected error when reading non-existent customer, got nil")
	}
}

func TestCustomersRepository_Update(t *testing.T) {
	conn := setupTestDB(t)
	defer teardownTestDB(t, conn)

	repo := NewCustomersRepository(conn)
	customer := Customer{
		Id:    uuid.New(),
		Name:  "Jane Doe",
		Email: "jane@example.com",
	}

	err := repo.Create(context.Background(), customer)
	if err != nil {
		t.Errorf("Create failed: %v", err)
	}

	customer.Name = "Jane Smith"
	customer.Email = "jane.smith@example.com"

	err = repo.Update(context.Background(), customer)
	if err != nil {
		t.Errorf("Update failed: %v", err)
	}

	updatedCustomer, err := repo.Read(context.Background(), customer.Id)
	if err != nil {
		t.Errorf("Read failed: %v", err)
	}

	if updatedCustomer.Name != "Jane Smith" {
		t.Errorf("Expected Name 'Jane Smith', got %v", updatedCustomer.Name)
	}
	if updatedCustomer.Email != "jane.smith@example.com" {
		t.Errorf("Expected Email 'jane.smith@example.com', got %v", updatedCustomer.Email)
	}
}

func TestCustomersRepository_Delete(t *testing.T) {
	conn := setupTestDB(t)
	defer teardownTestDB(t, conn)

	repo := NewCustomersRepository(conn)
	customer := Customer{
		Id:    uuid.New(),
		Name:  "Bob Wilson",
		Email: "bob@example.com",
	}

	err := repo.Create(context.Background(), customer)
	if err != nil {
		t.Errorf("Create failed: %v", err)
	}

	err = repo.Delete(context.Background(), customer.Id)
	if err != nil {
		t.Errorf("Delete failed: %v", err)
	}

	_, err = repo.Read(context.Background(), customer.Id)
	if err == nil {
		t.Error("Expected error when reading deleted customer, got nil")
	}
}

func TestCustomerService_CRUD(t *testing.T) {
	conn := setupTestDB(t)
	defer teardownTestDB(t, conn)

	repo := NewCustomersRepository(conn)
	service := NewCustomerService(repo)

	customer := Customer{
		Id:    uuid.New(),
		Name:  "Alice Johnson",
		Email: "alice@example.com",
	}

	err := service.Create(context.Background(), customer)
	if err != nil {
		t.Errorf("Service Create failed: %v", err)
	}

	retrievedCustomer, err := service.Read(context.Background(), customer.Id)
	if err != nil {
		t.Errorf("Service Read failed: %v", err)
	}

	if retrievedCustomer.Name != customer.Name {
		t.Errorf("Expected Name %v, got %v", customer.Name, retrievedCustomer.Name)
	}

	customer.Name = "Alice Brown"
	err = service.Update(context.Background(), customer)
	if err != nil {
		t.Errorf("Service Update failed: %v", err)
	}

	updatedCustomer, err := service.Read(context.Background(), customer.Id)
	if err != nil {
		t.Errorf("Service Read after update failed: %v", err)
	}

	if updatedCustomer.Name != "Alice Brown" {
		t.Errorf("Expected updated Name 'Alice Brown', got %v", updatedCustomer.Name)
	}

	err = service.Delete(context.Background(), customer.Id)
	if err != nil {
		t.Errorf("Service Delete failed: %v", err)
	}

	_, err = service.Read(context.Background(), customer.Id)
	if err == nil {
		t.Error("Expected error when reading deleted customer via service, got nil")
	}
}

func TestCustomersRepository_MultipleOperations(t *testing.T) {
	conn := setupTestDB(t)
	defer teardownTestDB(t, conn)

	repo := NewCustomersRepository(conn)

	customers := []Customer{
		{Id: uuid.New(), Name: "Customer 1", Email: "customer1@example.com"},
		{Id: uuid.New(), Name: "Customer 2", Email: "customer2@example.com"},
		{Id: uuid.New(), Name: "Customer 3", Email: "customer3@example.com"},
	}

	for _, customer := range customers {
		err := repo.Create(context.Background(), customer)
		if err != nil {
			t.Errorf("Failed to create customer %v: %v", customer.Name, err)
		}
	}

	for _, customer := range customers {
		retrievedCustomer, err := repo.Read(context.Background(), customer.Id)
		if err != nil {
			t.Errorf("Failed to read customer %v: %v", customer.Name, err)
		}
		if retrievedCustomer.Name != customer.Name {
			t.Errorf("Expected Name %v, got %v", customer.Name, retrievedCustomer.Name)
		}
	}

	for _, customer := range customers {
		err := repo.Delete(context.Background(), customer.Id)
		if err != nil {
			t.Errorf("Failed to delete customer %v: %v", customer.Name, err)
		}
	}
}

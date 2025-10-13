package mortgages

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
		dbURL = "postgres://postgres:postgres@localhost:5433/service2_db?sslmode=disable"
	}

	conn, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	_, err = conn.Exec(context.Background(), "DROP TABLE IF EXISTS mortgage_applications")
	if err != nil {
		t.Fatalf("Failed to drop existing mortgage_applications table: %v", err)
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
	_, err := conn.Exec(context.Background(), "DELETE FROM mortgage_applications")
	if err != nil {
		t.Errorf("Failed to clean up test data: %v", err)
	}
	conn.Close(context.Background())
}

func TestMortgageRepository_Create(t *testing.T) {
	conn := setupTestDB(t)
	defer teardownTestDB(t, conn)

	repo := NewMortgageRepository(conn)
	application := MortgageApplication{
		Id:            uuid.New(),
		CustomerId:    uuid.New(),
		LoanAmount:    500000.00,
		PropertyValue: 650000.00,
		InterestRate:  3.5,
		TermYears:     30,
		Status:        "pending",
	}

	err := repo.Create(context.Background(), application)
	if err != nil {
		t.Errorf("Create failed: %v", err)
	}

	retrievedApp, err := repo.Read(context.Background(), application.Id)
	if err != nil {
		t.Errorf("Read failed: %v", err)
	}

	if retrievedApp.Id != application.Id {
		t.Errorf("Expected ID %v, got %v", application.Id, retrievedApp.Id)
	}
	if retrievedApp.CustomerId != application.CustomerId {
		t.Errorf("Expected CustomerId %v, got %v", application.CustomerId, retrievedApp.CustomerId)
	}
	if retrievedApp.LoanAmount != application.LoanAmount {
		t.Errorf("Expected LoanAmount %v, got %v", application.LoanAmount, retrievedApp.LoanAmount)
	}
	if retrievedApp.Status != application.Status {
		t.Errorf("Expected Status %v, got %v", application.Status, retrievedApp.Status)
	}
}

func TestMortgageRepository_Read_NotFound(t *testing.T) {
	conn := setupTestDB(t)
	defer teardownTestDB(t, conn)

	repo := NewMortgageRepository(conn)
	nonExistentID := uuid.New()

	_, err := repo.Read(context.Background(), nonExistentID)
	if err == nil {
		t.Error("Expected error when reading non-existent application, got nil")
	}
}

func TestMortgageRepository_Update(t *testing.T) {
	conn := setupTestDB(t)
	defer teardownTestDB(t, conn)

	repo := NewMortgageRepository(conn)
	application := MortgageApplication{
		Id:            uuid.New(),
		CustomerId:    uuid.New(),
		LoanAmount:    400000.00,
		PropertyValue: 550000.00,
		InterestRate:  4.0,
		TermYears:     25,
		Status:        "pending",
	}

	err := repo.Create(context.Background(), application)
	if err != nil {
		t.Errorf("Create failed: %v", err)
	}

	application.Status = "approved"
	application.InterestRate = 3.75

	err = repo.Update(context.Background(), application)
	if err != nil {
		t.Errorf("Update failed: %v", err)
	}

	updatedApp, err := repo.Read(context.Background(), application.Id)
	if err != nil {
		t.Errorf("Read failed: %v", err)
	}

	if updatedApp.Status != "approved" {
		t.Errorf("Expected Status 'approved', got %v", updatedApp.Status)
	}
	if updatedApp.InterestRate != 3.75 {
		t.Errorf("Expected InterestRate 3.75, got %v", updatedApp.InterestRate)
	}
}

func TestMortgageRepository_Delete(t *testing.T) {
	conn := setupTestDB(t)
	defer teardownTestDB(t, conn)

	repo := NewMortgageRepository(conn)
	application := MortgageApplication{
		Id:            uuid.New(),
		CustomerId:    uuid.New(),
		LoanAmount:    300000.00,
		PropertyValue: 400000.00,
		InterestRate:  3.25,
		TermYears:     20,
		Status:        "pending",
	}

	err := repo.Create(context.Background(), application)
	if err != nil {
		t.Errorf("Create failed: %v", err)
	}

	err = repo.Delete(context.Background(), application.Id)
	if err != nil {
		t.Errorf("Delete failed: %v", err)
	}

	_, err = repo.Read(context.Background(), application.Id)
	if err == nil {
		t.Error("Expected error when reading deleted application, got nil")
	}
}

func TestMortgageRepository_GetByCustomerId(t *testing.T) {
	conn := setupTestDB(t)
	defer teardownTestDB(t, conn)

	repo := NewMortgageRepository(conn)
	customerId := uuid.New()

	applications := []MortgageApplication{
		{Id: uuid.New(), CustomerId: customerId, LoanAmount: 500000, PropertyValue: 650000, InterestRate: 3.5, TermYears: 30, Status: "pending"},
		{Id: uuid.New(), CustomerId: customerId, LoanAmount: 400000, PropertyValue: 550000, InterestRate: 4.0, TermYears: 25, Status: "approved"},
		{Id: uuid.New(), CustomerId: uuid.New(), LoanAmount: 300000, PropertyValue: 400000, InterestRate: 3.25, TermYears: 20, Status: "pending"},
	}

	for _, app := range applications {
		err := repo.Create(context.Background(), app)
		if err != nil {
			t.Errorf("Failed to create application: %v", err)
		}
	}

	customerApps, err := repo.GetByCustomerId(context.Background(), customerId)
	if err != nil {
		t.Errorf("GetByCustomerId failed: %v", err)
	}

	if len(customerApps) != 2 {
		t.Errorf("Expected 2 applications for customer, got %d", len(customerApps))
	}

	for _, app := range customerApps {
		if app.CustomerId != customerId {
			t.Errorf("Expected CustomerId %v, got %v", customerId, app.CustomerId)
		}
	}
}

func TestMortgageService_CRUD(t *testing.T) {
	conn := setupTestDB(t)
	defer teardownTestDB(t, conn)

	repo := NewMortgageRepository(conn)
	service := NewMortgageService(repo)

	application := MortgageApplication{
		Id:            uuid.New(),
		CustomerId:    uuid.New(),
		LoanAmount:    450000.00,
		PropertyValue: 600000.00,
		InterestRate:  3.8,
		TermYears:     30,
		Status:        "pending",
	}

	err := service.Create(context.Background(), application)
	if err != nil {
		t.Errorf("Service Create failed: %v", err)
	}

	retrievedApp, err := service.Read(context.Background(), application.Id)
	if err != nil {
		t.Errorf("Service Read failed: %v", err)
	}

	if retrievedApp.LoanAmount != application.LoanAmount {
		t.Errorf("Expected LoanAmount %v, got %v", application.LoanAmount, retrievedApp.LoanAmount)
	}

	application.Status = "approved"
	err = service.Update(context.Background(), application)
	if err != nil {
		t.Errorf("Service Update failed: %v", err)
	}

	updatedApp, err := service.Read(context.Background(), application.Id)
	if err != nil {
		t.Errorf("Service Read after update failed: %v", err)
	}

	if updatedApp.Status != "approved" {
		t.Errorf("Expected updated Status 'approved', got %v", updatedApp.Status)
	}

	err = service.Delete(context.Background(), application.Id)
	if err != nil {
		t.Errorf("Service Delete failed: %v", err)
	}

	_, err = service.Read(context.Background(), application.Id)
	if err == nil {
		t.Error("Expected error when reading deleted application via service, got nil")
	}
}

func TestMortgageRepository_MultipleOperations(t *testing.T) {
	conn := setupTestDB(t)
	defer teardownTestDB(t, conn)

	repo := NewMortgageRepository(conn)

	applications := []MortgageApplication{
		{Id: uuid.New(), CustomerId: uuid.New(), LoanAmount: 500000, PropertyValue: 650000, InterestRate: 3.5, TermYears: 30, Status: "pending"},
		{Id: uuid.New(), CustomerId: uuid.New(), LoanAmount: 400000, PropertyValue: 550000, InterestRate: 4.0, TermYears: 25, Status: "approved"},
		{Id: uuid.New(), CustomerId: uuid.New(), LoanAmount: 300000, PropertyValue: 400000, InterestRate: 3.25, TermYears: 20, Status: "rejected"},
	}

	for _, app := range applications {
		err := repo.Create(context.Background(), app)
		if err != nil {
			t.Errorf("Failed to create application: %v", err)
		}
	}

	for _, app := range applications {
		retrievedApp, err := repo.Read(context.Background(), app.Id)
		if err != nil {
			t.Errorf("Failed to read application: %v", err)
		}
		if retrievedApp.LoanAmount != app.LoanAmount {
			t.Errorf("Expected LoanAmount %v, got %v", app.LoanAmount, retrievedApp.LoanAmount)
		}
	}

	for _, app := range applications {
		err := repo.Delete(context.Background(), app.Id)
		if err != nil {
			t.Errorf("Failed to delete application: %v", err)
		}
	}
}

# saga-pattern
An example implementation of the saga pattern in Go using microservices architecture.

## Architecture

This project consists of three microservices that implement a distributed saga pattern for a mortgage loan system:

- **Service 1** (port 8081): Customer Service - manages customer information
- **Service 2** (port 8082): Mortgage Application Service - handles mortgage applications
- **Service 3** (port 8083): Loan Servicing Service - manages active loans and payments

All services share a single PostgreSQL database server but use separate databases (schemas):
- `service1_db` - Customer database
- `service2_db` - Mortgage application database
- `service3_db` - Loan servicing database

## Quick Start

### Prerequisites
- Docker and Docker Compose installed
- Ports 5432, 8081, 8082, 8083 available

### Running All Services

From the project root directory:

```bash
# Build and start all services
docker-compose up --build

# Or run in detached mode
docker-compose up -d --build
```

This will:
1. Start a PostgreSQL server on port 5432
2. Create three separate databases (service1_db, service2_db, service3_db)
3. Build and start all three Go services
4. Set up networking between services

### Stopping Services

```bash
# Stop all services
docker-compose down

# Stop and remove volumes (clears database data)
docker-compose down -v
```

### Viewing Logs

```bash
# View all logs
docker-compose logs

# View logs for a specific service
docker-compose logs service1
docker-compose logs service2
docker-compose logs service3
docker-compose logs postgres

# Follow logs in real-time
docker-compose logs -f
```

## API Endpoints

### Service 1 - Customer Service (port 8081)
- `POST /customers` - Create customer
- `GET /customers/:id` - Get customer by ID
- `PUT /customers/:id` - Update customer
- `DELETE /customers/:id` - Delete customer

### Service 2 - Mortgage Application Service (port 8082)
- `POST /applications` - Create mortgage application
- `GET /applications/:id` - Get application by ID
- `PUT /applications/:id` - Update application (approve/reject)
- `DELETE /applications/:id` - Delete application
- `GET /customers/:customerId/applications` - Get all applications for a customer

### Service 3 - Loan Servicing Service (port 8083)
- `POST /loans` - Create loan
- `GET /loans/:id` - Get loan by ID
- `PUT /loans/:id` - Update loan
- `DELETE /loans/:id` - Delete loan
- `GET /customers/:customerId/loans` - Get all loans for a customer
- `GET /mortgages/:mortgageId/loan` - Get loan by mortgage ID
- `POST /payments` - Create payment
- `GET /payments/:id` - Get payment by ID
- `GET /loans/:loanId/payments` - Get all payments for a loan
- `GET /customers/:customerId/payments` - Get all payments for a customer

## Testing

Use the test-client.http files in each service directory to test the APIs with your HTTP client.

## Database Access

Connect to the PostgreSQL database:

```bash
# Using docker exec
docker exec -it saga_postgres psql -U postgres -d service1_db
docker exec -it saga_postgres psql -U postgres -d service2_db
docker exec -it saga_postgres psql -U postgres -d service3_db

# From host (if psql is installed)
psql -h localhost -U postgres -d service1_db
```

Default credentials:
- Username: `postgres`
- Password: `postgres`

## Development

### Running Individual Services Locally

Each service can be run independently for development:

```bash
# Start only the database
docker-compose up postgres

# Run service locally (example for service1)
cd service1/api
go run main.go
```

Make sure to update the `DATABASE_URL` in each service's `.env` file if running locally.

## Project Structure

```
saga-pattern/
├── docker-compose.yml          # Main orchestration file
├── init-db.sql                 # Database initialization script
├── service1/                   # Customer service
│   ├── Dockerfile
│   ├── docker-compose.yml      # Individual service compose file
│   ├── api/
│   └── ...
├── service2/                   # Mortgage application service
│   ├── Dockerfile
│   ├── docker-compose.yml
│   ├── api/
│   └── ...
└── service3/                   # Loan servicing service
    ├── Dockerfile
    ├── docker-compose.yml
    ├── api/
    └── ...
``` 

# Backend Assignment - Go Service

A Go-based service using Gin framework that handles limited stock product ordering and background settlement job processing.

## Features

### 1. Product Ordering System
- Limited stock product management
- Concurrent-safe order processing
- Stock validation and reduction
- Order tracking and retrieval

### 2. Settlement Job Processing
- Background job queue using Go channels
- Batch processing of transactions
- Job progress tracking and cancellation
- Worker pool configuration

## Tech Stack

- **Language**: Go 1.20+
- **Web Framework**: Gin
- **Database**: PostgreSQL
- **Containerization**: Docker Compose
- **Architecture**: Clean architecture with handlers, services, repositories

## Quick Start

### Prerequisites
- Docker and Docker Compose
- Go 1.20+ (for local development)

### Setup and Run

1. **Clone the repository**
   ```bash
   git clone https://github.com/nurfiansyah06/backend-assignment-indico.git
   cd backend-assignment-indico
   ```

2. **Start services with Docker Compose**
   ```bash
   docker-compose up -d
   ```

3. **Run database migrations using [https://github.com/golang-migrate/migrate]**

4. **Seed test data (optional)**
   ```bash
   # Generate ~1,000,000 transaction records
   docker-compose exec app go run cmd/seed/main.go
   ```

5. **Start the application**
   ```bash
   # Application starts automatically with docker-compose
   # Or run locally for development:
   go run main.go
   ```

The service will be available at `http://localhost:8080`

## API Endpoints

### Product Ordering

#### Create Order
```http
POST /orders
Content-Type: application/json

{
  "product_id": 1,
  "quantity": 1,
  "buyer_id": "user-123"
}
```

**Response:**
- `200 OK`: Order created successfully
- `400 Bad Request`: Invalid request data
- `409 Conflict`: OUT_OF_STOCK error

#### Get Order Details
```http
GET /orders/:id
```

**Response:**
```json
{
  "id": "order-123",
  "product_id": 1,
  "quantity": 1,
  "buyer_id": "user-123",
  "status": "completed",
  "created_at": "2025-01-01T10:00:00Z"
}
```

### Settlement Jobs

#### Create Settlement Job
```http
POST /jobs/settlement
Content-Type: application/json

{
  "from": "2025-01-01",
  "to": "2025-01-31"
}
```

**Response:**
```json
{
  "job_id": "job_123",
  "status": "QUEUED"
}
```
*HTTP Status: 202 Accepted*

#### Get Job Status
```http
GET /jobs/:id
```

**Response (Running):**
```json
{
  "job_id": "job_123",
  "status": "RUNNING",
  "progress": 63,
  "processed": 630000,
  "total": 1000000
}
```

**Response (Completed):**
```json
{
  "job_id": "job_123",
  "status": "COMPLETED",
  "progress": 100,
  "processed": 1000000,
  "total": 1000000,
  "download_url": "/downloads/job_123.csv"
}
```

#### Cancel Job
```http
POST /jobs/:id/cancel
```

#### Download Settlement File
```http
GET /downloads/:job_id.csv
```

## Configuration

### Environment Variables

```bash
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=backendtest

# Server
PORT=8080

# Job Processing
WORKERS=8                    # Number of worker goroutines (4-16 recommended)
BATCH_SIZE=10000            # Transaction batch size (5k-20k recommended)

# File Storage
SETTLEMENT_DIR=/tmp/settlements
```

### Docker Compose Configuration

The `docker-compose.yml` includes:
- PostgreSQL database with persistent storage
- Application service with environment variables
- Volume mounts for CSV file storage
- Health checks and dependency management

## Database Schema

### Products
```sql
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    stock INTEGER NOT NULL DEFAULT 0,
    price_cents INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);
```

### Orders
```sql
CREATE TABLE orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id INTEGER REFERENCES products(id),
    buyer_id VARCHAR(255) NOT NULL,
    quantity INTEGER NOT NULL,
    total_amount_cents INTEGER NOT NULL,
    status VARCHAR(50) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT NOW()
);
```

### Transactions
```sql
CREATE TABLE transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    merchant_id INTEGER NOT NULL,
    amount_cents INTEGER NOT NULL,
    fee_cents INTEGER NOT NULL,
    status VARCHAR(50) NOT NULL,
    paid_at TIMESTAMP NOT NULL
);
```

### Settlements
```sql
CREATE TABLE settlements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    merchant_id INTEGER NOT NULL,
    date DATE NOT NULL,
    gross_cents BIGINT NOT NULL,
    fee_cents BIGINT NOT NULL,
    net_cents BIGINT NOT NULL,
    txn_count INTEGER NOT NULL,
    generated_at TIMESTAMP DEFAULT NOW(),
    unique_run_id UUID NOT NULL,
    UNIQUE(merchant_id, date)
);
```

### Jobs
```sql
CREATE TABLE jobs (
    job_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    from_date DATE NOT NULL,
    to_date DATE NOT NULL,
    status VARCHAR(50) DEFAULT 'QUEUED',
    progress DECIMAL(5,2) DEFAULT 0.00,
    processed_count INTEGER DEFAULT 0,
    total_count INTEGER DEFAULT 0,
    result_path TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
```

## Architecture

```
├── cmd/
│   ├── migrate/           # Database migrations
│   └── seed/             # Test data seeding
├── internal/
│   ├── handlers/         # HTTP request handlers
│   ├── services/         # Business logic
│   ├── repositories/     # Data access layer
│   ├── models/          # Data models
│   └── jobs/            # Background job processing
├── migrations/          # SQL migration files
├── docker-compose.yml   # Docker services configuration
└── main.go             # Application entry point
```

## Key Features Implementation

### Concurrent Order Processing
- Database-level stock locking with `SELECT FOR UPDATE`
- Atomic stock reduction in single transaction
- Proper error handling for concurrent access

### Background Job System
- Channel-based job queue: `jobQueue := make(chan Job)`
- Worker pool with configurable size via `WORKERS` env var
- Batch processing to handle large datasets efficiently
- Context-based cancellation support

### Settlement Processing
- Aggregates transactions per merchant per day
- Upsert logic for idempotent processing
- CSV file generation with structured output
- Progress tracking and status updates

## Testing

### Run All Tests
```bash
go test ./...
```

### Concurrent Order Test
```bash
# Test that demonstrates no overselling with 500 concurrent orders
go test -v ./internal/handlers -run TestConcurrentOrders
```

This test verifies that when 500 concurrent orders are placed for a product with limited stock (e.g., 100), only 100 orders succeed and stock doesn't go below zero.

## Performance Guidelines

The system is designed to:
- Process ~1M transaction rows end-to-end in reasonable time
- Maintain stable memory usage through batching
- Demonstrate effective use of parallelism
- Handle high concurrent load for ordering

## Development

### Local Development Setup
```bash
# Start only database
docker-compose up -d postgres

# Run application locally
export DB_HOST=localhost
export DB_PORT=5432
export WORKERS=4
go run main.go
```

### Adding New Migrations
```bash
# Create new migration file
touch migrations/003_your_migration.sql
```

### Monitoring and Logs
```bash
# View application logs
docker-compose logs -f app

# View database logs
docker-compose logs -f postgres
```

## Troubleshooting

### Common Issues

1. **Job stuck in QUEUED status**
   - Check if background workers are running
   - Verify `WORKERS` environment variable is set
   - Check application logs for errors

2. **OUT_OF_STOCK errors**
   - Verify product stock in database
   - Check for proper transaction isolation

3. **Database connection issues**
   - Ensure PostgreSQL container is running
   - Verify connection parameters
   - Check if migrations have been applied

### Health Checks
```bash
# Check if services are running
docker-compose ps

# Test database connection
docker-compose exec postgres psql -U postgres -d backendtest -c "SELECT 1;"

# Test API endpoints
curl http://localhost:8080/health
```

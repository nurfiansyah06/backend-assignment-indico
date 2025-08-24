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
   
2. **Run database migrations using [Golang Migrate](https://github.com/golang-migrate/migrate)**


3. **Start the application**
   ```bash
   # Application starts automatically with docker-compose
   # Or run locally for development:
   go run cmd/main.go
   ```

The service will be available at `http://localhost:8080`

## API Endpoints
Documentation API [Postman](https://documenter.getpostman.com/view/11932880/2sB3BLkTkd)

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


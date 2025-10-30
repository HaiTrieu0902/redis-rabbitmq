# Quick Start Guide - Streaming Service

## Overview

The **Streaming Service** is a Go-based microservice for managing financial transactions with event streaming capabilities. It follows Clean Architecture principles and integrates with authentication and user services.

## What's Included

✅ **Clean Architecture** - Domain, Use Case, Infrastructure, and Delivery layers  
✅ **Transaction Management** - CRUD operations for transactions  
✅ **JWT Authentication** - Secure API with token validation  
✅ **RabbitMQ Integration** - Event publishing for transaction events  
✅ **PostgreSQL** - Persistent data storage  
✅ **RESTful API** - Well-structured HTTP endpoints  
✅ **Pagination** - Efficient data retrieval

## Prerequisites

- Go 1.23+
- PostgreSQL running on port 5432
- RabbitMQ running on port 5672
- Authentication service running on port 3001 (for JWT tokens)

## Quick Start

### 1. Install Dependencies

```bash
cd streaming
go mod download
```

### 2. Configure Environment

Edit `.env` file:

```env
PORT=3003
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=redis_mq
JWT_SECRET=your-secret-key-here
RABBITMQ_URL=amqp://guest:guest@localhost:5672/
```

⚠️ **Important:** Use the same `JWT_SECRET` as your authentication service!

### 3. Run the Service

```bash
# Option 1: Direct run
go run main.go

# Option 2: Build and run
go build -o bin/streaming.exe main.go
./bin/streaming.exe

# Option 3: Using Make
make run
```

The service will start on: **http://localhost:3003**

### 4. Verify Service is Running

```bash
curl http://localhost:3003/health
```

Expected response:

```json
{
  "status": "ok",
  "service": "streaming"
}
```

## Testing the API

### Step 1: Get JWT Token

First, login via the authentication service to get a JWT token:

```bash
# Open in browser
http://localhost:3001/auth/google
```

Save the `access_token` from the response.

### Step 2: Create a Transaction

```bash
curl -X POST http://localhost:3003/api/v1/transactions \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "amount": 100.50,
    "currency": "VND",
    "transaction_type": "deposit",
    "description": "My first transaction"
  }'
```

### Step 3: Get Your Transactions

```bash
curl http://localhost:3003/api/v1/transactions/my \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Step 4: Update Transaction Status

```bash
curl -X PATCH http://localhost:3003/api/v1/transactions/TRANSACTION_ID/status \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"status": "success"}'
```

## Project Structure

```
streaming/
├── domain/                     # Business entities and interfaces
│   ├── entity/                 # Domain models
│   │   └── transaction.go      # Transaction entity
│   └── repository/             # Repository interfaces
│       └── transaction_repository.go
│
├── infrastructure/             # External dependencies
│   ├── config/                 # Configuration management
│   │   └── config.go
│   ├── database/               # Database connection
│   │   └── postgres.go
│   ├── messaging/              # Message queue
│   │   └── rabbitmq.go
│   └── repository/             # Repository implementations
│       └── transaction_repository_impl.go
│
├── usecase/                    # Business logic
│   └── transaction_usecase.go
│
├── delivery/                   # HTTP layer
│   └── http/
│       ├── handler/            # Request handlers
│       │   └── transaction_handler.go
│       ├── middleware/         # Authentication middleware
│       │   └── auth_middleware.go
│       └── router/             # Route configuration
│           └── router.go
│
├── main.go                     # Application entry point
├── go.mod                      # Go dependencies
├── .env                        # Environment variables
├── README.md                   # Full documentation
└── API_DOCUMENTATION.md        # API reference
```

## Available Endpoints

| Method | Endpoint                          | Description             | Auth Required |
| ------ | --------------------------------- | ----------------------- | ------------- |
| GET    | `/health`                         | Health check            | ❌            |
| POST   | `/api/v1/transactions`            | Create transaction      | ✅            |
| GET    | `/api/v1/transactions/:id`        | Get transaction by ID   | ✅            |
| GET    | `/api/v1/transactions/my`         | Get user's transactions | ✅            |
| PATCH  | `/api/v1/transactions/:id/status` | Update status           | ✅            |
| GET    | `/api/v1/transactions`            | Get all transactions    | ✅            |
| GET    | `/api/v1/transactions/status`     | Filter by status        | ✅            |

## Transaction Types

- `deposit` - Deposit money
- `withdraw` - Withdraw money
- `purchase` - Purchase transaction

## Transaction Statuses

- `pending` - Transaction is pending
- `success` - Transaction completed successfully
- `failed` - Transaction failed

## Development Commands

```bash
# Build the application
make build

# Run the application
make run

# Run tests
make test

# Clean build artifacts
make clean

# Install/update dependencies
make deps

# Build for production
make build-prod
```

## RabbitMQ Events

The service publishes events to the `transaction_events` queue:

### Event: transaction.created

```json
{
  "event_type": "transaction.created",
  "transaction": { ... },
  "timestamp": "2025-10-14T10:30:00Z"
}
```

### Event: transaction.updated

```json
{
  "event_type": "transaction.updated",
  "transaction": { ... },
  "timestamp": "2025-10-14T10:35:00Z"
}
```

## Troubleshooting

### Service won't start

1. Check PostgreSQL is running:

   ```bash
   psql -U postgres -c "SELECT 1"
   ```

2. Check RabbitMQ is running:

   ```bash
   curl http://localhost:15672
   ```

3. Verify environment variables in `.env`

### Authentication fails

1. Ensure JWT_SECRET matches authentication service
2. Check token format: `Bearer <token>`
3. Verify token is not expired

### Database errors

1. Ensure database `redis_mq` exists
2. Run the SQL schema: `psql -U postgres -d redis_mq -f ../sql.sql`
3. Check database credentials

## Integration with Other Services

### Authentication Service (Port 3001)

- Provides JWT tokens via OAuth login
- Shared JWT secret for token validation

### User Service (Port 3002)

- Manages user profiles
- Uses same JWT authentication

### Shared Database

All services use the same PostgreSQL database for data consistency.

## Next Steps

1. ✅ Service is running
2. ✅ Test with curl or Postman
3. 📝 Integrate with your frontend
4. 📊 Monitor RabbitMQ events
5. 🔒 Secure with production JWT secrets
6. 📈 Add monitoring and logging

## Documentation

- **README.md** - Complete service documentation
- **API_DOCUMENTATION.md** - Detailed API reference with examples
- **INTEGRATION_GUIDE.md** - How to integrate with other services

## Support

If you encounter issues:

1. Check the logs in the terminal
2. Verify all prerequisites are running
3. Review the `.env` configuration
4. Check the API_DOCUMENTATION.md for examples

---

**Service Status:**

- ✅ Built successfully
- ✅ Dependencies installed
- ✅ Ready to run

**Run:** `go run main.go` or `make run`

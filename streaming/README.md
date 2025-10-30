# Streaming Service - Transaction Management

A Go-based microservice for managing transactions using Clean Architecture, PostgreSQL, and RabbitMQ.

## Architecture

This service follows Clean Architecture principles:

```
streaming/
├── domain/                 # Business entities and repository interfaces
│   ├── entity/
│   │   └── transaction.go
│   └── repository/
│       └── transaction_repository.go
├── infrastructure/         # External dependencies implementation
│   ├── config/
│   │   └── config.go
│   ├── database/
│   │   └── postgres.go
│   ├── messaging/
│   │   └── rabbitmq.go
│   └── repository/
│       └── transaction_repository_impl.go
├── usecase/               # Business logic
│   └── transaction_usecase.go
├── delivery/              # HTTP layer
│   └── http/
│       ├── handler/
│       │   └── transaction_handler.go
│       ├── middleware/
│       │   └── auth_middleware.go
│       └── router/
│           └── router.go
└── main.go               # Application entry point
```

## Features

- ✅ Transaction CRUD operations
- ✅ JWT-based authentication
- ✅ RabbitMQ event publishing
- ✅ PostgreSQL database
- ✅ RESTful API
- ✅ Clean Architecture
- ✅ Pagination support

## Prerequisites

- Go 1.23+
- PostgreSQL
- RabbitMQ
- Authentication service (for JWT tokens)

## Setup

1. **Install dependencies:**

```bash
go mod download
```

2. **Configure environment:**
   Copy `.env.example` to `.env` and update the values:

```bash
cp .env.example .env
```

3. **Run the service:**

```bash
go run main.go
```

## API Endpoints

All endpoints require JWT authentication (Bearer token from authentication service).

### Transactions

#### Create Transaction

```http
POST /api/v1/transactions
Authorization: Bearer <token>
Content-Type: application/json

{
  "amount": 100.50,
  "currency": "VND",
  "transaction_type": "deposit",
  "description": "Monthly deposit"
}
```

#### Get Transaction by ID

```http
GET /api/v1/transactions/:id
Authorization: Bearer <token>
```

#### Get User's Transactions

```http
GET /api/v1/transactions/my?page=1&page_size=10
Authorization: Bearer <token>
```

#### Update Transaction Status

```http
PATCH /api/v1/transactions/:id/status
Authorization: Bearer <token>
Content-Type: application/json

{
  "status": "success"
}
```

#### Get All Transactions

```http
GET /api/v1/transactions?page=1&page_size=10
Authorization: Bearer <token>
```

#### Get Transactions by Status

```http
GET /api/v1/transactions/status?status=pending&page=1&page_size=10
Authorization: Bearer <token>
```

### Health Check

```http
GET /health
```

## Transaction Types

- `deposit` - Deposit money
- `withdraw` - Withdraw money
- `purchase` - Purchase transaction

## Transaction Statuses

- `pending` - Transaction is pending
- `success` - Transaction completed successfully
- `failed` - Transaction failed

## RabbitMQ Events

The service publishes events to the `transaction_events` queue:

```json
{
  "event_type": "transaction.created",
  "transaction": { ... },
  "timestamp": "2025-10-14T00:00:00Z"
}
```

Event types:

- `transaction.created`
- `transaction.updated`

## Integration with Authentication Service

This service requires JWT tokens from the authentication service. The token should contain:

```json
{
  "user_id": "uuid",
  "email": "user@example.com"
}
```

## Development

Build:

```bash
go build -o streaming main.go
```

Run:

```bash
./streaming
```

## Environment Variables

| Variable     | Description              | Default                            |
| ------------ | ------------------------ | ---------------------------------- |
| PORT         | Server port              | 3003                               |
| GIN_MODE     | Gin mode (debug/release) | debug                              |
| DB_HOST      | PostgreSQL host          | localhost                          |
| DB_PORT      | PostgreSQL port          | 5432                               |
| DB_USER      | Database user            | postgres                           |
| DB_PASSWORD  | Database password        | postgres                           |
| DB_NAME      | Database name            | redis_mq                           |
| DB_SSLMODE   | SSL mode                 | disable                            |
| JWT_SECRET   | JWT secret key           | your-secret-key-here               |
| RABBITMQ_URL | RabbitMQ connection URL  | amqp://guest:guest@localhost:5672/ |

## License

MIT

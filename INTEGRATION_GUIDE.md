# Integration Guide: Streaming Service with Authentication & User Services

This guide explains how the three microservices work together.

## Architecture Overview

```
┌─────────────────────┐
│  Authentication     │  Port 3001
│  Service (NestJS)   │
│  - OAuth Login      │
│  - JWT Generation   │
└──────────┬──────────┘
           │
           │ JWT Token
           ▼
┌─────────────────────┐
│  User Service       │  Port 3002
│  (NestJS)           │
│  - User CRUD        │
│  - Profile Mgmt     │
└─────────────────────┘

           │ JWT Token
           ▼
┌─────────────────────┐
│  Streaming Service  │  Port 3003
│  (Go)               │
│  - Transactions     │
│  - Event Publishing │
└──────────┬──────────┘
           │
           │ Events
           ▼
┌─────────────────────┐
│  RabbitMQ           │  Port 5672
│  - Message Queue    │
└─────────────────────┘
           │
           ▼
┌─────────────────────┐
│  PostgreSQL         │  Port 5432
│  - Shared Database  │
└─────────────────────┘
```

## Services

### 1. Authentication Service (NestJS - Port 3001)

**Purpose:** User authentication via OAuth providers and JWT token generation

**Endpoints:**

- `GET /auth/google` - Login with Google
- `GET /auth/microsoft` - Login with Microsoft
- `GET /auth/github` - Login with GitHub
- `GET /auth/google/callback` - OAuth callback
- `GET /auth/microsoft/callback` - OAuth callback
- `GET /auth/github/callback` - OAuth callback

**Returns:** JWT token with user information

### 2. User Service (NestJS - Port 3002)

**Purpose:** User profile management

**Endpoints:**

- `GET /users/:id` - Get user by ID
- `GET /users/me` - Get current user profile
- `PATCH /users/:id` - Update user
- `DELETE /users/:id` - Delete user

**Authentication:** Requires JWT token from authentication service

### 3. Streaming Service (Go - Port 3003)

**Purpose:** Transaction management and event streaming

**Endpoints:**

- `POST /api/v1/transactions` - Create transaction
- `GET /api/v1/transactions/:id` - Get transaction
- `GET /api/v1/transactions/my` - Get user's transactions
- `PATCH /api/v1/transactions/:id/status` - Update status
- `GET /api/v1/transactions` - Get all transactions
- `GET /api/v1/transactions/status` - Filter by status

**Authentication:** Requires JWT token from authentication service

## Complete Flow

### 1. User Authentication Flow

```
1. User → Authentication Service: Login via OAuth
   GET http://localhost:3001/auth/google

2. User redirected to Google for authentication

3. Google → Authentication Service: OAuth callback

4. Authentication Service → Database: Create/update user

5. Authentication Service → User: Return JWT token
   {
     "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
     "user": { ... }
   }
```

### 2. Transaction Creation Flow

```
1. User → Streaming Service: Create transaction
   POST http://localhost:3003/api/v1/transactions
   Authorization: Bearer <jwt_token>
   {
     "amount": 100,
     "currency": "VND",
     "transaction_type": "deposit"
   }

2. Streaming Service: Validate JWT token

3. Streaming Service → Database: Save transaction

4. Streaming Service → RabbitMQ: Publish event
   Queue: transaction_events
   {
     "event_type": "transaction.created",
     "transaction": { ... }
   }

5. Streaming Service → User: Return transaction
   {
     "message": "transaction created successfully",
     "data": { ... }
   }
```

## Setup Instructions

### Prerequisites

- Node.js 16+ (for authentication and user services)
- Go 1.23+ (for streaming service)
- PostgreSQL
- RabbitMQ
- Redis (optional)

### 1. Start Infrastructure

```bash
# Start PostgreSQL, RabbitMQ, Redis
docker-compose up -d
```

### 2. Setup Database

```bash
# Run SQL script to create tables
psql -U postgres -d redis_mq -f sql.sql
```

### 3. Start Authentication Service

```bash
cd authenication
npm install
cp .env.example .env
# Edit .env with your OAuth credentials
npm run start:dev
```

Service will run on: http://localhost:3001

### 4. Start User Service

```bash
cd user
npm install
cp .env.example .env
# Edit .env
npm run start:dev
```

Service will run on: http://localhost:3002

### 5. Start Streaming Service

```bash
cd streaming
cp .env.example .env
# Edit .env
go mod tidy
go run main.go
```

Service will run on: http://localhost:3003

## Environment Configuration

### Authentication Service (.env)

```env
PORT=3001
DATABASE_URL=postgresql://postgres:postgres@localhost:5432/redis_mq
JWT_SECRET=your-secret-key-here

# OAuth Credentials
GOOGLE_CLIENT_ID=your-google-client-id
GOOGLE_CLIENT_SECRET=your-google-client-secret
GOOGLE_CALLBACK_URL=http://localhost:3001/auth/google/callback

MICROSOFT_CLIENT_ID=your-microsoft-client-id
MICROSOFT_CLIENT_SECRET=your-microsoft-client-secret
MICROSOFT_CALLBACK_URL=http://localhost:3001/auth/microsoft/callback

GITHUB_CLIENT_ID=your-github-client-id
GITHUB_CLIENT_SECRET=your-github-client-secret
GITHUB_CALLBACK_URL=http://localhost:3001/auth/github/callback

RABBITMQ_URL=amqp://guest:guest@localhost:5672/
REDIS_HOST=localhost
REDIS_PORT=6379
```

### User Service (.env)

```env
PORT=3002
DATABASE_URL=postgresql://postgres:postgres@localhost:5432/redis_mq
JWT_SECRET=your-secret-key-here
RABBITMQ_URL=amqp://guest:guest@localhost:5672/
REDIS_HOST=localhost
REDIS_PORT=6379
```

### Streaming Service (.env)

```env
PORT=3003
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=redis_mq
DB_SSLMODE=disable
JWT_SECRET=your-secret-key-here
RABBITMQ_URL=amqp://guest:guest@localhost:5672/
```

**Important:** All services must use the **same JWT_SECRET** to verify tokens!

## Testing the Integration

### Step 1: Get JWT Token

```bash
# Open in browser to login with Google
http://localhost:3001/auth/google

# Save the returned access_token
```

### Step 2: Create Transaction

```bash
curl -X POST http://localhost:3003/api/v1/transactions \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "amount": 100.50,
    "currency": "VND",
    "transaction_type": "deposit",
    "description": "Test transaction"
  }'
```

### Step 3: Get User's Transactions

```bash
curl -X GET http://localhost:3003/api/v1/transactions/my \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Step 4: Update Transaction Status

```bash
curl -X PATCH http://localhost:3003/api/v1/transactions/TRANSACTION_ID/status \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"status": "success"}'
```

## JWT Token Structure

The JWT token generated by the authentication service contains:

```json
{
  "user_id": "123e4567-e89b-12d3-a456-426614174000",
  "email": "user@example.com",
  "iat": 1729000000,
  "exp": 1729086400
}
```

The streaming service's JWT middleware validates this token and extracts the `user_id` to associate transactions with users.

## RabbitMQ Event Flow

### Transaction Events

Queue: `transaction_events`

Events published by the streaming service:

1. **Transaction Created**

```json
{
  "event_type": "transaction.created",
  "transaction": {
    "id": "uuid",
    "user_id": "uuid",
    "amount": 100.5,
    "currency": "VND",
    "transaction_type": "deposit",
    "status": "pending",
    "created_at": "2025-10-14T10:30:00Z"
  },
  "timestamp": "2025-10-14T10:30:00Z"
}
```

2. **Transaction Updated**

```json
{
  "event_type": "transaction.updated",
  "transaction": {
    "id": "uuid",
    "status": "success",
    "updated_at": "2025-10-14T10:35:00Z"
  },
  "timestamp": "2025-10-14T10:35:00Z"
}
```

### Consuming Events

Other services can consume these events:

```typescript
// NestJS example
@RabbitSubscribe({
  exchange: '',
  routingKey: 'transaction_events',
  queue: 'transaction_events',
})
async handleTransactionEvent(message: any) {
  console.log('Received transaction event:', message);
  // Process the event
}
```

## Database Schema

All services share the same PostgreSQL database with these relevant tables:

- `users` - User profiles
- `user_identity` - OAuth provider identities
- `transactions` - Transaction records
- `conversation` - Chat conversations (if applicable)
- `message` - Chat messages (if applicable)

## Troubleshooting

### JWT Token Invalid

- Ensure all services use the same `JWT_SECRET`
- Check token expiration
- Verify token format: `Bearer <token>`

### Database Connection Failed

- Check PostgreSQL is running
- Verify database credentials
- Ensure database `redis_mq` exists

### RabbitMQ Connection Failed

- Check RabbitMQ is running
- Verify connection URL
- Check firewall settings

### Transaction Not Found

- Verify transaction ID is valid UUID
- Check user has permission to view transaction
- Ensure database has the record

## Next Steps

1. Add Redis caching for frequently accessed data
2. Implement transaction rollback mechanisms
3. Add monitoring and logging
4. Implement rate limiting
5. Add WebSocket support for real-time updates
6. Create admin dashboard
7. Add transaction reporting and analytics

## Support

For issues or questions:

- Check service logs
- Verify environment variables
- Test database connectivity
- Check RabbitMQ management UI: http://localhost:15672

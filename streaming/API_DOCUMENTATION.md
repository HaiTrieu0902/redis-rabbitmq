# Streaming Service API Documentation

Base URL: `http://localhost:3003/api/v1`

All endpoints require JWT authentication via Bearer token obtained from the authentication service.

## Authentication

Add the JWT token to the Authorization header:

```
Authorization: Bearer <your_jwt_token>
```

The JWT token is obtained from the authentication service (port 3001) after logging in via Google, Microsoft, or GitHub.

## Endpoints

### 1. Create Transaction

Create a new transaction for the authenticated user.

**Endpoint:** `POST /transactions`

**Headers:**

```
Authorization: Bearer <token>
Content-Type: application/json
```

**Request Body:**

```json
{
  "amount": 100.5,
  "currency": "VND",
  "transaction_type": "deposit",
  "description": "Monthly deposit"
}
```

**Fields:**

- `amount` (required): Transaction amount (must be > 0)
- `currency` (required): Currency code (e.g., "VND", "USD")
- `transaction_type` (required): Type of transaction
  - `deposit` - Deposit money
  - `withdraw` - Withdraw money
  - `purchase` - Purchase transaction
- `description` (optional): Transaction description

**Response:** `201 Created`

```json
{
  "message": "transaction created successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "user_id": "123e4567-e89b-12d3-a456-426614174000",
    "amount": 100.5,
    "currency": "VND",
    "transaction_type": "deposit",
    "status": "pending",
    "description": "Monthly deposit",
    "created_at": "2025-10-14T10:30:00Z",
    "updated_at": "2025-10-14T10:30:00Z"
  }
}
```

---

### 2. Get Transaction by ID

Retrieve a specific transaction by its ID.

**Endpoint:** `GET /transactions/:id`

**Headers:**

```
Authorization: Bearer <token>
```

**Response:** `200 OK`

```json
{
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "user_id": "123e4567-e89b-12d3-a456-426614174000",
    "amount": 100.5,
    "currency": "VND",
    "transaction_type": "deposit",
    "status": "pending",
    "description": "Monthly deposit",
    "created_at": "2025-10-14T10:30:00Z",
    "updated_at": "2025-10-14T10:30:00Z"
  }
}
```

---

### 3. Get User's Transactions

Get all transactions for the authenticated user with pagination.

**Endpoint:** `GET /transactions/my`

**Headers:**

```
Authorization: Bearer <token>
```

**Query Parameters:**

- `page` (optional): Page number (default: 1)
- `page_size` (optional): Items per page (default: 10, max: 100)

**Example:** `GET /transactions/my?page=1&page_size=10`

**Response:** `200 OK`

```json
{
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "user_id": "123e4567-e89b-12d3-a456-426614174000",
      "amount": 100.5,
      "currency": "VND",
      "transaction_type": "deposit",
      "status": "pending",
      "description": "Monthly deposit",
      "created_at": "2025-10-14T10:30:00Z",
      "updated_at": "2025-10-14T10:30:00Z"
    }
  ],
  "page": 1,
  "page_size": 10
}
```

---

### 4. Update Transaction Status

Update the status of a transaction.

**Endpoint:** `PATCH /transactions/:id/status`

**Headers:**

```
Authorization: Bearer <token>
Content-Type: application/json
```

**Request Body:**

```json
{
  "status": "success"
}
```

**Valid Statuses:**

- `pending` - Transaction is pending
- `success` - Transaction completed successfully
- `failed` - Transaction failed

**Response:** `200 OK`

```json
{
  "message": "transaction status updated successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "user_id": "123e4567-e89b-12d3-a456-426614174000",
    "amount": 100.5,
    "currency": "VND",
    "transaction_type": "deposit",
    "status": "success",
    "description": "Monthly deposit",
    "created_at": "2025-10-14T10:30:00Z",
    "updated_at": "2025-10-14T10:35:00Z"
  }
}
```

---

### 5. Get All Transactions

Get all transactions in the system with pagination (admin function).

**Endpoint:** `GET /transactions`

**Headers:**

```
Authorization: Bearer <token>
```

**Query Parameters:**

- `page` (optional): Page number (default: 1)
- `page_size` (optional): Items per page (default: 10, max: 100)

**Example:** `GET /transactions?page=1&page_size=20`

**Response:** `200 OK`

```json
{
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "user_id": "123e4567-e89b-12d3-a456-426614174000",
      "amount": 100.5,
      "currency": "VND",
      "transaction_type": "deposit",
      "status": "success",
      "description": "Monthly deposit",
      "created_at": "2025-10-14T10:30:00Z",
      "updated_at": "2025-10-14T10:35:00Z"
    }
  ],
  "page": 1,
  "page_size": 20
}
```

---

### 6. Get Transactions by Status

Filter transactions by status with pagination.

**Endpoint:** `GET /transactions/status`

**Headers:**

```
Authorization: Bearer <token>
```

**Query Parameters:**

- `status` (required): Transaction status to filter
  - `pending`
  - `success`
  - `failed`
- `page` (optional): Page number (default: 1)
- `page_size` (optional): Items per page (default: 10, max: 100)

**Example:** `GET /transactions/status?status=pending&page=1&page_size=10`

**Response:** `200 OK`

```json
{
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "user_id": "123e4567-e89b-12d3-a456-426614174000",
      "amount": 100.5,
      "currency": "VND",
      "transaction_type": "deposit",
      "status": "pending",
      "description": "Monthly deposit",
      "created_at": "2025-10-14T10:30:00Z",
      "updated_at": "2025-10-14T10:30:00Z"
    }
  ],
  "status": "pending",
  "page": 1,
  "page_size": 10
}
```

---

### 7. Health Check

Check if the service is running.

**Endpoint:** `GET /health`

**No authentication required**

**Response:** `200 OK`

```json
{
  "status": "ok",
  "service": "streaming"
}
```

---

## Error Responses

All endpoints may return the following error responses:

**400 Bad Request**

```json
{
  "error": "invalid transaction ID"
}
```

**401 Unauthorized**

```json
{
  "error": "authorization header is required"
}
```

**404 Not Found**

```json
{
  "error": "transaction not found"
}
```

**500 Internal Server Error**

```json
{
  "error": "failed to create transaction: database error"
}
```

---

## RabbitMQ Events

When transactions are created or updated, events are published to the `transaction_events` queue:

**Event Structure:**

```json
{
  "event_type": "transaction.created",
  "transaction": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "user_id": "123e4567-e89b-12d3-a456-426614174000",
    "amount": 100.5,
    "currency": "VND",
    "transaction_type": "deposit",
    "status": "pending",
    "description": "Monthly deposit",
    "created_at": "2025-10-14T10:30:00Z",
    "updated_at": "2025-10-14T10:30:00Z"
  },
  "timestamp": "2025-10-14T10:30:00Z"
}
```

**Event Types:**

- `transaction.created` - New transaction created
- `transaction.updated` - Transaction status or details updated

---

## Integration with Authentication Service

### Getting a JWT Token

1. Start the authentication service (port 3001)
2. Navigate to one of the OAuth login endpoints:

   - Google: `http://localhost:3001/auth/google`
   - Microsoft: `http://localhost:3001/auth/microsoft`
   - GitHub: `http://localhost:3001/auth/github`

3. After successful authentication, you'll receive a JWT token in the response
4. Use this token in the `Authorization: Bearer <token>` header for all streaming service requests

### JWT Token Contents

The JWT token contains:

```json
{
  "user_id": "123e4567-e89b-12d3-a456-426614174000",
  "email": "user@example.com",
  "exp": 1729000000
}
```

---

## Example Usage with curl

### 1. Get JWT Token (from auth service)

```bash
# Login via Google OAuth (opens browser)
curl http://localhost:3001/auth/google
```

### 2. Create Transaction

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

### 3. Get My Transactions

```bash
curl -X GET "http://localhost:3003/api/v1/transactions/my?page=1&page_size=10" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### 4. Update Transaction Status

```bash
curl -X PATCH http://localhost:3003/api/v1/transactions/TRANSACTION_ID/status \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "status": "success"
  }'
```

---

## Testing

You can import the Postman collection from the root directory to test all endpoints.

The collection includes:

- Pre-configured authentication
- All API endpoints
- Example requests and responses

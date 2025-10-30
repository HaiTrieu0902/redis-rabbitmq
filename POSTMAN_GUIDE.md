# Postman Collection Guide

This guide explains how to use the Postman collection for testing all three microservices.

## Import Collection

1. Open Postman
2. Click **Import** button
3. Select `postman-collection.json`
4. Collection will be imported with all endpoints

## Collection Structure

The collection contains three services:

### 1. Authentication Service (Port 3001)

- Health Check
- Google OAuth Login
- GitHub OAuth Login
- Microsoft MSAL Login
- Refresh Token
- Get Current User
- Logout

### 2. User Service (Port 3002)

- Health Check
- Get My Profile
- Update My Profile
- Get User By ID
- Get All Users
- Update User By ID
- Delete My Account
- Delete User By ID

### 3. Streaming Service (Port 3003) - **NEW! 🎉**

- Health Check
- Create Transaction
- Get Transaction By ID
- Get My Transactions
- Update Transaction Status
- Get All Transactions
- Get Transactions By Status (Pending/Success/Failed)
- Create Withdraw Transaction
- Create Purchase Transaction

## Environment Variables

The collection uses these variables (auto-managed):

| Variable         | Description                 | Auto-Set                     |
| ---------------- | --------------------------- | ---------------------------- |
| `auth_url`       | http://localhost:3001       | ✅                           |
| `user_url`       | http://localhost:3002       | ✅                           |
| `streaming_url`  | http://localhost:3003       | ✅                           |
| `access_token`   | JWT token                   | ✅ (from OAuth login)        |
| `refresh_token`  | Refresh token               | ✅ (from OAuth login)        |
| `transaction_id` | Last created transaction ID | ✅ (from create transaction) |

## Quick Start Guide

### Step 1: Start All Services

Make sure all services are running:

```bash
# Authentication Service
cd authenication
npm run start:dev

# User Service
cd user
npm run start:dev

# Streaming Service
cd streaming
go run main.go
```

### Step 2: Authenticate

1. Open **Authentication Service** → **Google OAuth Login** (or GitHub/Microsoft)
2. Click **Send** (this will open in browser)
3. Complete OAuth login
4. Copy the `access_token` from response
5. The token is automatically saved to `{{access_token}}` variable

**Alternative:** Paste token manually:

- Click collection **Variables** tab
- Set `access_token` value
- Click **Save**

### Step 3: Test Streaming Service

Now you can use any Streaming Service endpoint:

#### 3a. Create a Transaction

1. Open **Streaming Service** → **Create Transaction**
2. Click **Send**
3. The `transaction_id` is automatically saved
4. Example response:

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
    "created_at": "2025-10-14T10:30:00Z"
  }
}
```

#### 3b. Get My Transactions

1. Open **Streaming Service** → **Get My Transactions**
2. Adjust pagination if needed: `page=1&page_size=10`
3. Click **Send**

#### 3c. Update Transaction Status

1. Open **Streaming Service** → **Update Transaction Status**
2. The `{{transaction_id}}` is auto-filled from step 3a
3. Change status in body: `pending`, `success`, or `failed`
4. Click **Send**

#### 3d. Filter By Status

1. Open **Streaming Service** → **Get Transactions By Status - Pending**
2. Click **Send**
3. Try other status endpoints: Success, Failed

## Complete Test Workflow

### Full Transaction Flow

```
1. Login → Get JWT Token
   Authentication Service → Google OAuth Login

2. Create Deposit
   Streaming Service → Create Transaction
   Body: {"amount": 100, "currency": "VND", "transaction_type": "deposit"}

3. View My Transactions
   Streaming Service → Get My Transactions

4. Update to Success
   Streaming Service → Update Transaction Status
   Body: {"status": "success"}

5. Create Withdrawal
   Streaming Service → Create Withdraw Transaction
   Body: {"amount": 50, "currency": "USD", "transaction_type": "withdraw"}

6. View All Transactions
   Streaming Service → Get All Transactions

7. Filter Pending
   Streaming Service → Get Transactions By Status - Pending
```

## Transaction Types

### Deposit Transaction

```json
{
  "amount": 100.5,
  "currency": "VND",
  "transaction_type": "deposit",
  "description": "Monthly salary deposit"
}
```

### Withdraw Transaction

```json
{
  "amount": 50.0,
  "currency": "USD",
  "transaction_type": "withdraw",
  "description": "ATM withdrawal"
}
```

### Purchase Transaction

```json
{
  "amount": 25.99,
  "currency": "USD",
  "transaction_type": "purchase",
  "description": "Online shopping"
}
```

## Transaction Statuses

- **pending** - Transaction is being processed
- **success** - Transaction completed successfully
- **failed** - Transaction failed

## Request Examples

### Create Transaction with Authorization

```
POST http://localhost:3003/api/v1/transactions
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json

{
  "amount": 100.50,
  "currency": "VND",
  "transaction_type": "deposit",
  "description": "Test transaction"
}
```

### Update Transaction Status

```
PATCH http://localhost:3003/api/v1/transactions/{transaction_id}/status
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json

{
  "status": "success"
}
```

## Pre-Request Scripts

The collection includes automatic scripts:

### Auto-Save Transaction ID

When you create a transaction, the response ID is automatically saved:

```javascript
if (pm.response.code === 201) {
  const response = pm.response.json();
  if (response.data && response.data.id) {
    pm.collectionVariables.set('transaction_id', response.data.id);
  }
}
```

This means you can immediately use `{{transaction_id}}` in other requests!

## Response Examples

### Success Response (201 Created)

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

### Paginated Response (200 OK)

```json
{
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "amount": 100.5,
      "currency": "VND",
      "transaction_type": "deposit",
      "status": "success"
    }
  ],
  "page": 1,
  "page_size": 10
}
```

### Error Response (401 Unauthorized)

```json
{
  "error": "authorization header is required"
}
```

### Error Response (400 Bad Request)

```json
{
  "error": "amount must be greater than 0"
}
```

## Troubleshooting

### "authorization header is required"

**Problem:** No JWT token or invalid format

**Solution:**

1. Login again: **Authentication Service** → **Google OAuth Login**
2. Check `access_token` variable has a value
3. Ensure header format: `Bearer {{access_token}}`

### "invalid token"

**Problem:** JWT token expired or invalid

**Solution:**

1. Get new token from authentication service
2. Ensure all services use same `JWT_SECRET`

### "transaction not found"

**Problem:** Invalid transaction ID

**Solution:**

1. Create a transaction first
2. Check `{{transaction_id}}` variable is set
3. Use the correct UUID format

### Service Not Responding

**Problem:** Service is not running

**Solution:**

```bash
# Check services
curl http://localhost:3001/health  # Auth
curl http://localhost:3002/health  # User
curl http://localhost:3003/health  # Streaming
```

## Testing Checklist

Use this checklist to test all functionality:

- [ ] Health checks (all services)
- [ ] OAuth login (Google/GitHub/Microsoft)
- [ ] Get current user profile
- [ ] Create deposit transaction
- [ ] Create withdraw transaction
- [ ] Create purchase transaction
- [ ] Get transaction by ID
- [ ] Get my transactions (with pagination)
- [ ] Update transaction status to success
- [ ] Update transaction status to failed
- [ ] Get all transactions
- [ ] Filter by pending status
- [ ] Filter by success status
- [ ] Filter by failed status

## Advanced Usage

### Test Multiple Currencies

```json
{"amount": 100, "currency": "VND", "transaction_type": "deposit"}
{"amount": 50, "currency": "USD", "transaction_type": "deposit"}
{"amount": 75, "currency": "EUR", "transaction_type": "deposit"}
```

### Test Pagination

```
page=1&page_size=5   // First 5 items
page=2&page_size=5   // Next 5 items
page=1&page_size=20  // First 20 items
```

### Test Invalid Scenarios

Try these to test error handling:

- Negative amount: `{"amount": -100, ...}`
- Invalid type: `{"transaction_type": "invalid", ...}`
- Invalid status: `{"status": "invalid"}`
- Missing required fields

## Export/Share Collection

To share with your team:

1. Click collection **...** menu
2. Select **Export**
3. Choose format: **Collection v2.1 (recommended)**
4. Save file
5. Share `postman-collection.json`

## Tips & Tricks

1. **Use Console** - View `console.log()` output from test scripts
2. **Save Responses** - Click **Save Response** to compare results
3. **Environments** - Create dev/staging/prod environments
4. **Tests Tab** - Add custom assertions
5. **Keyboard Shortcuts** - `Ctrl+Enter` to send request

---

**Collection Ready! 🚀**

You now have a complete Postman collection for testing all three microservices with automatic variable management and pre-configured requests.

**Quick Test:**

1. Send: **Authentication Service** → **Google OAuth Login**
2. Send: **Streaming Service** → **Create Transaction**
3. Send: **Streaming Service** → **Get My Transactions**

Done! ✅

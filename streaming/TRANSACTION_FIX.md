# Transaction Foreign Key Fix

## Problem

When creating a transaction, you encountered a foreign key constraint violation error:

```
failed to create transaction: pq: insert or update on table "transactions" violates foreign key constraint "transactions_user_id_fkey"
```

## Root Cause

The `transactions` table has a foreign key constraint on the `user_id` column that references the `users` table. This means you cannot create a transaction with a `user_id` that doesn't exist in the `users` table.

## Solution Implemented

Added user validation before creating a transaction to provide a clearer error message:

1. **Added `UserExists` method to repository interface** (`domain/repository/transaction_repository.go`)

   - New method: `UserExists(ctx context.Context, userID uuid.UUID) (bool, error)`

2. **Implemented `UserExists` in repository** (`infrastructure/repository/transaction_repository_impl.go`)

   - Checks if a user exists in the `users` table before attempting to create a transaction

3. **Added validation in usecase** (`usecase/transaction_usecase.go`)
   - Validates user existence before creating transaction
   - Returns clear error message: "user with ID {uuid} does not exist, please ensure the user is created before creating a transaction"

## How to Use Properly

### Step 1: Create a User First

Before creating a transaction, ensure the user exists in the `users` table. You need to call the authentication service or user service to create a user:

```bash
# Example: Create user via authentication service
POST http://localhost:3001/auth/register
{
  "email": "user@example.com",
  "name": "John Doe"
}
```

### Step 2: Get the User ID

After authentication, you should receive a JWT token. Decode this token to get the `user_id`, or use the authentication endpoint to get user details.

### Step 3: Create Transaction

Now you can create a transaction with the valid `user_id`:

```bash
POST http://localhost:8080/transactions
Authorization: Bearer <your-jwt-token>
{
  "amount": 100.50,
  "currency": "USD",
  "transaction_type": "deposit",
  "description": "Initial deposit"
}
```

## Testing the Fix

### Test 1: Create Transaction Without User (Should Fail)

```bash
POST http://localhost:8080/transactions
Authorization: Bearer <valid-jwt-token-but-user-not-in-db>
{
  "amount": 100.50,
  "currency": "USD",
  "transaction_type": "deposit"
}
```

**Expected Response (400 Bad Request):**

```json
{
  "error": "user with ID {uuid} does not exist, please ensure the user is created before creating a transaction"
}
```

### Test 2: Create Transaction With Valid User (Should Succeed)

```bash
# First create user
POST http://localhost:3001/auth/register
{
  "email": "test@example.com",
  "name": "Test User"
}

# Then create transaction
POST http://localhost:8080/transactions
Authorization: Bearer <jwt-token-from-above>
{
  "amount": 100.50,
  "currency": "USD",
  "transaction_type": "deposit",
  "description": "Test transaction"
}
```

**Expected Response (201 Created):**

```json
{
  "message": "transaction created successfully",
  "data": {
    "id": "...",
    "user_id": "...",
    "amount": 100.5,
    "currency": "USD",
    "transaction_type": "deposit",
    "status": "pending",
    "description": "Test transaction",
    "created_at": "...",
    "updated_at": "..."
  }
}
```

## Additional Notes

### Database Schema

The `transactions` table has this constraint:

```sql
user_id uuid REFERENCES users(id) ON DELETE CASCADE
```

This means:

- A transaction MUST have a valid `user_id` that exists in the `users` table
- If a user is deleted, all their transactions will be automatically deleted (CASCADE)

### Integration with Auth Service

Make sure your authentication service creates users in the same database that the streaming service uses. The services should share the same `users` table.

### Manual User Creation (For Testing)

If you need to create a user manually for testing:

```sql
INSERT INTO users (id, email, name, created_at, updated_at)
VALUES (
  'f47ac10b-58cc-4372-a567-0e02b2c3d479', -- Use this UUID in your transaction
  'test@example.com',
  'Test User',
  NOW(),
  NOW()
);
```

Then use this UUID (`f47ac10b-58cc-4372-a567-0e02b2c3d479`) when creating transactions.

## Error Handling Summary

| Error                                | Cause                                                 | Solution                                        |
| ------------------------------------ | ----------------------------------------------------- | ----------------------------------------------- |
| `user with ID {uuid} does not exist` | User not created before transaction                   | Create user first via auth service              |
| `foreign key constraint violation`   | (Should not happen now) User validation prevents this | -                                               |
| `user not authenticated`             | Missing or invalid JWT token                          | Include valid JWT token in Authorization header |

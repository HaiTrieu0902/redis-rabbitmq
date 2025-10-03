# API Examples & Testing

Tập hợp các ví dụ API calls để test hệ thống.

## 🔧 Setup

### Set Environment Variables (PowerShell)
```powershell
# After login, set your tokens
$access_token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
$refresh_token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
$user_id = "550e8400-e29b-41d4-a716-446655440000"
```

### Set Environment Variables (Bash)
```bash
export ACCESS_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
export REFRESH_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
export USER_ID="550e8400-e29b-41d4-a716-446655440000"
```

---

## 🔐 Authentication Service Examples

### 1. Health Check
```bash
curl -X GET http://localhost:3001/auth/health
```

**Expected Response:**
```json
{
  "status": "ok",
  "service": "authentication"
}
```

---

### 2. Google OAuth Login
```bash
# Open in browser
http://localhost:3001/auth/google
```

**Flow:**
1. Browser redirects to Google
2. User logs in
3. Google redirects back to callback
4. Receives tokens in URL

**Expected Redirect:**
```
http://localhost:3000/auth/callback?access_token=eyJ...&refresh_token=eyJ...
```

---

### 3. GitHub OAuth Login
```bash
# Open in browser
http://localhost:3001/auth/github
```

---

### 4. Microsoft MSAL Login
```bash
# Open in browser
http://localhost:3001/auth/msal
```

---

### 5. Refresh Access Token

**PowerShell:**
```powershell
curl -X POST http://localhost:3001/auth/refresh `
  -H "Content-Type: application/json" `
  -d "{\"refresh_token\": \"$refresh_token\"}"
```

**Bash:**
```bash
curl -X POST http://localhost:3001/auth/refresh \
  -H "Content-Type: application/json" \
  -d "{\"refresh_token\": \"$REFRESH_TOKEN\"}"
```

**Expected Response:**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
  "expires_in": 900,
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "user@example.com",
    "name": "John Doe",
    "avatar_url": "https://example.com/avatar.jpg"
  }
}
```

---

### 6. Get Current User (from Auth Service)

**PowerShell:**
```powershell
curl -X GET http://localhost:3001/auth/me `
  -H "Authorization: Bearer $access_token"
```

**Bash:**
```bash
curl -X GET http://localhost:3001/auth/me \
  -H "Authorization: Bearer $ACCESS_TOKEN"
```

**Expected Response:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "user@example.com",
  "name": "John Doe",
  "avatar_url": "https://example.com/avatar.jpg"
}
```

---

### 7. Logout

**PowerShell:**
```powershell
curl -X POST http://localhost:3001/auth/logout `
  -H "Authorization: Bearer $access_token"
```

**Bash:**
```bash
curl -X POST http://localhost:3001/auth/logout \
  -H "Authorization: Bearer $ACCESS_TOKEN"
```

**Expected Response:**
```json
{
  "message": "Logged out successfully"
}
```

---

## 👤 User Service Examples

### 1. Health Check
```bash
curl -X GET http://localhost:3002/users/health
```

**Expected Response:**
```json
{
  "status": "ok",
  "service": "user"
}
```

---

### 2. Get My Profile

**PowerShell:**
```powershell
curl -X GET http://localhost:3002/users/me `
  -H "Authorization: Bearer $access_token"
```

**Bash:**
```bash
curl -X GET http://localhost:3002/users/me \
  -H "Authorization: Bearer $ACCESS_TOKEN"
```

**Expected Response:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "user@example.com",
  "name": "John Doe",
  "avatar_url": "https://example.com/avatar.jpg",
  "created_at": "2025-01-01T00:00:00.000Z",
  "updated_at": "2025-01-01T00:00:00.000Z"
}
```

---

### 3. Update My Profile

**PowerShell:**
```powershell
curl -X PUT http://localhost:3002/users/me `
  -H "Authorization: Bearer $access_token" `
  -H "Content-Type: application/json" `
  -d '{"name": "Jane Doe", "avatar_url": "https://example.com/new-avatar.jpg"}'
```

**Bash:**
```bash
curl -X PUT http://localhost:3002/users/me \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name": "Jane Doe", "avatar_url": "https://example.com/new-avatar.jpg"}'
```

**Expected Response:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "user@example.com",
  "name": "Jane Doe",
  "avatar_url": "https://example.com/new-avatar.jpg",
  "created_at": "2025-01-01T00:00:00.000Z",
  "updated_at": "2025-01-01T12:00:00.000Z"
}
```

---

### 4. Get User By ID

**PowerShell:**
```powershell
curl -X GET "http://localhost:3002/users/$user_id" `
  -H "Authorization: Bearer $access_token"
```

**Bash:**
```bash
curl -X GET "http://localhost:3002/users/$USER_ID" \
  -H "Authorization: Bearer $ACCESS_TOKEN"
```

---

### 5. Get All Users (Paginated)

**PowerShell:**
```powershell
curl -X GET "http://localhost:3002/users?limit=10&offset=0" `
  -H "Authorization: Bearer $access_token"
```

**Bash:**
```bash
curl -X GET "http://localhost:3002/users?limit=10&offset=0" \
  -H "Authorization: Bearer $ACCESS_TOKEN"
```

**Expected Response:**
```json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "user1@example.com",
    "name": "User One",
    "avatar_url": "https://example.com/avatar1.jpg",
    "created_at": "2025-01-01T00:00:00.000Z",
    "updated_at": "2025-01-01T00:00:00.000Z"
  },
  {
    "id": "660e8400-e29b-41d4-a716-446655440001",
    "email": "user2@example.com",
    "name": "User Two",
    "avatar_url": "https://example.com/avatar2.jpg",
    "created_at": "2025-01-02T00:00:00.000Z",
    "updated_at": "2025-01-02T00:00:00.000Z"
  }
]
```

---

### 6. Update User By ID

**PowerShell:**
```powershell
curl -X PUT "http://localhost:3002/users/$user_id" `
  -H "Authorization: Bearer $access_token" `
  -H "Content-Type: application/json" `
  -d '{"name": "Updated Name"}'
```

**Bash:**
```bash
curl -X PUT "http://localhost:3002/users/$USER_ID" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name": "Updated Name"}'
```

---

### 7. Delete My Account

**PowerShell:**
```powershell
curl -X DELETE http://localhost:3002/users/me `
  -H "Authorization: Bearer $access_token"
```

**Bash:**
```bash
curl -X DELETE http://localhost:3002/users/me \
  -H "Authorization: Bearer $ACCESS_TOKEN"
```

**Expected Response:**
```json
{
  "message": "Account deleted successfully"
}
```

---

### 8. Delete User By ID

**PowerShell:**
```powershell
curl -X DELETE "http://localhost:3002/users/$user_id" `
  -H "Authorization: Bearer $access_token"
```

**Bash:**
```bash
curl -X DELETE "http://localhost:3002/users/$USER_ID" \
  -H "Authorization: Bearer $ACCESS_TOKEN"
```

---

## 🧪 Complete Test Flow

### Scenario: User Login → Get Profile → Update → Logout

**Step 1: Login with Google**
```bash
# Open in browser
http://localhost:3001/auth/google
```

**Step 2: Extract Tokens from Redirect URL**
```
http://localhost:3000/auth/callback?access_token=ABC...&refresh_token=XYZ...
```

**Step 3: Set Tokens**
```powershell
$access_token = "ABC..."
$refresh_token = "XYZ..."
```

**Step 4: Get My Profile**
```powershell
curl -X GET http://localhost:3002/users/me `
  -H "Authorization: Bearer $access_token"
```

**Step 5: Update Profile**
```powershell
curl -X PUT http://localhost:3002/users/me `
  -H "Authorization: Bearer $access_token" `
  -H "Content-Type: application/json" `
  -d '{"name": "My New Name"}'
```

**Step 6: Verify Update**
```powershell
curl -X GET http://localhost:3002/users/me `
  -H "Authorization: Bearer $access_token"
```

**Step 7: Logout**
```powershell
curl -X POST http://localhost:3001/auth/logout `
  -H "Authorization: Bearer $access_token"
```

**Step 8: Verify Token is Blacklisted**
```powershell
# This should fail with 401 Unauthorized
curl -X GET http://localhost:3002/users/me `
  -H "Authorization: Bearer $access_token"
```

---

## 🔄 Token Refresh Flow

### When Access Token Expires

**Step 1: Try API Call (will fail)**
```bash
curl -X GET http://localhost:3002/users/me \
  -H "Authorization: Bearer $EXPIRED_TOKEN"
```

**Response:**
```json
{
  "statusCode": 401,
  "message": "Unauthorized"
}
```

**Step 2: Refresh Token**
```bash
curl -X POST http://localhost:3001/auth/refresh \
  -H "Content-Type: application/json" \
  -d "{\"refresh_token\": \"$REFRESH_TOKEN\"}"
```

**Step 3: Update Access Token**
```powershell
# Extract new access_token from response
$access_token = "NEW_TOKEN_HERE"
```

**Step 4: Retry API Call**
```bash
curl -X GET http://localhost:3002/users/me \
  -H "Authorization: Bearer $ACCESS_TOKEN"
```

---

## 📊 Testing with Postman

### Import Collection
1. Open Postman
2. File > Import
3. Select `postman-collection.json`
4. Collection variables will auto-save tokens

### Test Flow in Postman
1. **Auth Service > Google OAuth Login**
   - Open in browser
   - Copy tokens from redirect URL
   
2. **Auth Service > Refresh Token**
   - Paste refresh_token in body
   - Tokens auto-saved to variables
   
3. **User Service > Get My Profile**
   - Uses saved access_token
   
4. **User Service > Update My Profile**
   - Edit body JSON
   - Submit
   
5. **Auth Service > Logout**
   - Token blacklisted

---

## 🐛 Error Examples

### 401 Unauthorized - Missing Token
```bash
curl -X GET http://localhost:3002/users/me
```

**Response:**
```json
{
  "statusCode": 401,
  "message": "Unauthorized"
}
```

---

### 401 Unauthorized - Invalid Token
```bash
curl -X GET http://localhost:3002/users/me \
  -H "Authorization: Bearer invalid_token"
```

**Response:**
```json
{
  "statusCode": 401,
  "message": "Unauthorized"
}
```

---

### 401 Unauthorized - Expired Token
```bash
curl -X GET http://localhost:3002/users/me \
  -H "Authorization: Bearer $EXPIRED_TOKEN"
```

**Response:**
```json
{
  "statusCode": 401,
  "message": "jwt expired"
}
```

**Solution:** Use refresh token to get new access token

---

### 404 Not Found - User Not Found
```bash
curl -X GET http://localhost:3002/users/invalid-uuid \
  -H "Authorization: Bearer $ACCESS_TOKEN"
```

**Response:**
```json
{
  "statusCode": 404,
  "message": "User not found"
}
```

---

### 400 Bad Request - Validation Error
```bash
curl -X PUT http://localhost:3002/users/me \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"email": "invalid-email"}'
```

**Response:**
```json
{
  "statusCode": 400,
  "message": [
    "email must be an email"
  ],
  "error": "Bad Request"
}
```

---

## 🔍 Debugging Tips

### Check Redis Cache
```bash
docker exec -it redis redis-cli

# See all keys
127.0.0.1:6379> KEYS *

# Get user cache
127.0.0.1:6379> GET "user:550e8400-e29b-41d4-a716-446655440000"

# Check token blacklist
127.0.0.1:6379> GET "blacklist:eyJhbG..."

# Exit
127.0.0.1:6379> exit
```

### Check Database
```sql
-- Connect to database
psql -U postgres -d vuihoi

-- See all users
SELECT * FROM users;

-- See user identities
SELECT u.email, ui.provider, ui.provider_user_id 
FROM users u 
JOIN user_identity ui ON u.id = ui.user_id;

-- Exit
\q
```

### Pretty Print JSON (PowerShell)
```powershell
curl http://localhost:3002/users/me `
  -H "Authorization: Bearer $access_token" | ConvertFrom-Json | ConvertTo-Json -Depth 10
```

### Pretty Print JSON (Bash with jq)
```bash
curl http://localhost:3002/users/me \
  -H "Authorization: Bearer $ACCESS_TOKEN" | jq
```

---

## 📝 Notes

- Access tokens expire in 15 minutes
- Refresh tokens expire in 7 days
- User cache in Redis expires in 1 hour
- Session cache expires in 7 days
- Token blacklist entries expire with token lifetime

---

Enjoy testing! 🚀

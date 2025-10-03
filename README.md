# Microservices Authentication & User Management

Hệ thống microservices với Authentication service và User service sử dụng NestJS, PostgreSQL, Redis, và OAuth2.0 (Google, GitHub, Microsoft MSAL).

## 🏗️ Kiến trúc

```
┌─────────────────────┐
│   Frontend (3000)   │
└──────────┬──────────┘
           │
           │ HTTP Requests
           ▼
┌──────────────────────────────────────────────┐
│                                              │
│  ┌────────────────────┐  ┌────────────────┐ │
│  │ Authentication     │  │ User Service   │ │
│  │ Service (3001)     │  │ (3002)         │ │
│  │                    │  │                │ │
│  │ - Google OAuth     │  │ - User CRUD    │ │
│  │ - GitHub OAuth     │  │ - Profile Mgmt │ │
│  │ - MSAL OAuth       │  │ - Caching      │ │
│  │ - JWT Generation   │  │                │ │
│  └────────┬───────────┘  └───────┬────────┘ │
│           │                      │          │
│           └──────────┬───────────┘          │
│                      │                      │
└──────────────────────┼──────────────────────┘
                       │
        ┌──────────────┴──────────────┐
        │                             │
        ▼                             ▼
┌───────────────┐            ┌────────────────┐
│ Redis (6379)  │            │ PostgreSQL     │
│               │            │ (5432)         │
│ - Session     │            │                │
│ - Cache       │            │ - users        │
│ - Pub/Sub     │            │ - user_identity│
└───────────────┘            └────────────────┘
```

## 🚀 Tính năng

### Authentication Service (Port 3001)
- ✅ Google OAuth 2.0 authentication
- ✅ GitHub OAuth authentication  
- ✅ Microsoft Azure AD (MSAL) authentication
- ✅ JWT token generation & validation
- ✅ Refresh token mechanism
- ✅ Token blacklist (logout)
- ✅ Redis session management
- ✅ User creation/update via OAuth

### User Service (Port 3002)
- ✅ Get user profile (with Redis caching)
- ✅ Update user information
- ✅ Delete user account
- ✅ List all users (paginated)
- ✅ Real-time sync với Auth service qua Redis Pub/Sub
- ✅ JWT authentication middleware

### Redis Integration
- ✅ Session management
- ✅ Token caching
- ✅ Token blacklist
- ✅ Pub/Sub messaging giữa services
- ✅ User data caching (1 hour TTL)

## 📦 Cài đặt

### 1. Prerequisites
```bash
# Đảm bảo đã cài đặt:
- Node.js >= 18
- PostgreSQL >= 14
- Docker (cho Redis)
```

### 2. Chạy Redis
```bash
cd d:\redis-mq
docker-compose up -d
```

### 3. Setup Database
```bash
# Import SQL schema
psql -U postgres -d vuihoi -f sql.sql
```

### 4. Install Dependencies

**Authentication Service:**
```bash
cd authenication
npm install
```

**User Service:**
```bash
cd user
npm install
```

## ⚙️ Configuration

### Authentication Service (.env)
```env
DATABASE_URL=postgresql://postgres:040202005173@localhost:5432/vuihoi
REDIS_HOST=localhost
REDIS_PORT=6379

JWT_SECRET=your-super-secret-jwt-key-change-in-production
JWT_EXPIRES_IN=15m
JWT_REFRESH_SECRET=your-super-secret-refresh-jwt-key-change-in-production
JWT_REFRESH_EXPIRES_IN=7d

# Google OAuth
GOOGLE_CLIENT_ID=your-google-client-id
GOOGLE_CLIENT_SECRET=your-google-client-secret
GOOGLE_REDIRECT_URI=http://localhost:3001/auth/callback/google

# GitHub OAuth
GITHUB_CLIENT_ID=your-github-client-id
GITHUB_CLIENT_SECRET=your-github-client-secret
GITHUB_REDIRECT_URI=http://localhost:3001/auth/callback/github

# Microsoft MSAL
APP_AZURE_CLIENT_ID=your-azure-client-id
APP_AZURE_TENANT_ID=your-azure-tenant-id
APP_AZURE_REDIRECT_URI=http://localhost:3001/auth/msal/callback
AZURE_CLIENT_SECRET=your-azure-client-secret

PORT=3001
```

### User Service (.env)
```env
DATABASE_URL=postgresql://postgres:040202005173@localhost:5432/vuihoi
REDIS_HOST=localhost
REDIS_PORT=6379

JWT_SECRET=your-super-secret-jwt-key-change-in-production
PORT=3002
```

## 🎯 Chạy Services

### Development Mode

**Terminal 1 - Authentication Service:**
```bash
cd authenication
npm run start:dev
```

**Terminal 2 - User Service:**
```bash
cd user
npm run start:dev
```

## 📝 API Documentation

### Authentication Service (http://localhost:3001)

#### 1. Google OAuth Login
```
GET /auth/google
```
Redirect user đến Google login page.

```
GET /auth/callback/google
```
Google callback URL, returns tokens.

#### 2. GitHub OAuth Login
```
GET /auth/github
```
Redirect user đến GitHub login page.

```
GET /auth/callback/github
```
GitHub callback URL, returns tokens.

#### 3. Microsoft MSAL Login
```
GET /auth/msal
```
Redirect user đến Microsoft login page.

```
GET /auth/msal/callback?code={code}
```
Microsoft callback URL, returns tokens.

#### 4. Refresh Token
```http
POST /auth/refresh
Content-Type: application/json

{
  "refresh_token": "your-refresh-token"
}
```

Response:
```json
{
  "access_token": "new-access-token",
  "refresh_token": "new-refresh-token",
  "expires_in": 900,
  "user": {
    "id": "uuid",
    "email": "user@example.com",
    "name": "User Name",
    "avatar_url": "https://..."
  }
}
```

#### 5. Logout
```http
POST /auth/logout
Authorization: Bearer {access_token}
```

#### 6. Get Current User
```http
GET /auth/me
Authorization: Bearer {access_token}
```

### User Service (http://localhost:3002)

#### 1. Get My Profile
```http
GET /users/me
Authorization: Bearer {access_token}
```

Response:
```json
{
  "id": "uuid",
  "email": "user@example.com",
  "name": "User Name",
  "avatar_url": "https://...",
  "created_at": "2025-01-01T00:00:00Z",
  "updated_at": "2025-01-01T00:00:00Z"
}
```

#### 2. Update My Profile
```http
PUT /users/me
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "name": "New Name",
  "avatar_url": "https://new-avatar.com/image.jpg"
}
```

#### 3. Get User By ID
```http
GET /users/{userId}
Authorization: Bearer {access_token}
```

#### 4. Get All Users (Paginated)
```http
GET /users?limit=20&offset=0
Authorization: Bearer {access_token}
```

#### 5. Delete My Account
```http
DELETE /users/me
Authorization: Bearer {access_token}
```

## 🔐 Authentication Flow

### OAuth Flow (Google/GitHub/MSAL)
```
1. Frontend → GET /auth/google
2. Redirect to Google OAuth consent screen
3. User approves
4. Google → GET /auth/callback/google?code=xxx
5. Auth Service:
   - Exchange code for Google tokens
   - Get user info from Google
   - Create/update user in database
   - Create user identity record
   - Generate JWT tokens
   - Cache in Redis
   - Publish event to Redis
6. Redirect to frontend with tokens
7. Frontend stores tokens
8. Frontend uses access_token for API calls
```

### JWT Flow
```
1. Frontend → API Request with Authorization: Bearer {token}
2. Service validates JWT signature
3. Service checks if token is blacklisted (Redis)
4. Service extracts userId from token
5. Service processes request
6. Return response
```

### Refresh Token Flow
```
1. Frontend → POST /auth/refresh with refresh_token
2. Auth Service validates refresh token
3. Check if token is blacklisted
4. Generate new access_token & refresh_token
5. Return new tokens
```

## 🔄 Microservices Communication

### Redis Pub/Sub Events

**Event: user:updated**
```json
{
  "userId": "uuid",
  "email": "user@example.com",
  "name": "User Name",
  "avatar_url": "https://...",
  "timestamp": "2025-01-01T00:00:00Z"
}
```

**Event: user:deleted**
```json
{
  "userId": "uuid",
  "timestamp": "2025-01-01T00:00:00Z"
}
```

## 🧪 Testing

### Test OAuth Login (Browser)
1. Mở browser
2. Navigate to `http://localhost:3001/auth/google`
3. Login với Google account
4. Sẽ được redirect về với tokens

### Test API với cURL

**Get user profile:**
```bash
curl -X GET http://localhost:3002/users/me \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

**Update profile:**
```bash
curl -X PUT http://localhost:3002/users/me \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name": "New Name"}'
```

## 📊 Database Schema

```sql
-- Users table
users (
  id UUID PRIMARY KEY,
  name VARCHAR(255),
  email VARCHAR(255) UNIQUE,
  avatar_url TEXT,
  created_at TIMESTAMP,
  updated_at TIMESTAMP
)

-- User identities (OAuth providers)
user_identity (
  id UUID PRIMARY KEY,
  user_id UUID REFERENCES users(id),
  provider VARCHAR(50), -- 'google', 'github', 'msal'
  provider_user_id VARCHAR(255),
  access_token TEXT,
  refresh_token TEXT,
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  UNIQUE(provider, provider_user_id)
)
```

## 🎨 Tech Stack

- **Framework:** NestJS
- **Database:** PostgreSQL with TypeORM
- **Cache/Queue:** Redis (ioredis)
- **Authentication:** Passport.js + JWT
- **OAuth Providers:** 
  - Google (passport-google-oauth20)
  - GitHub (passport-github2)
  - Microsoft Azure AD (@azure/msal-node)
- **Validation:** class-validator, class-transformer

## 🔒 Security Best Practices

1. ✅ JWT secrets được lưu trong environment variables
2. ✅ Passwords không được lưu (chỉ dùng OAuth)
3. ✅ Tokens có expiration time
4. ✅ Refresh tokens cho long-term sessions
5. ✅ Token blacklist khi logout
6. ✅ CORS configured properly
7. ✅ Input validation với class-validator
8. ✅ SQL injection protection với TypeORM

## 🐛 Troubleshooting

### Redis connection failed
```bash
# Check if Redis is running
docker ps

# Restart Redis
docker-compose restart redis
```

### Database connection error
```bash
# Check PostgreSQL is running
psql -U postgres -c "SELECT 1"

# Check DATABASE_URL in .env files
```

### OAuth redirect mismatch
- Đảm bảo redirect URIs trong .env match với Google/GitHub/Azure console
- Google Console: https://console.cloud.google.com/
- GitHub Apps: https://github.com/settings/developers
- Azure Portal: https://portal.azure.com/

## 📈 Monitoring

### Check Services Health
```bash
# Auth service
curl http://localhost:3001/auth/health

# User service
curl http://localhost:3002/users/health
```

### Redis CLI
```bash
docker exec -it redis redis-cli
> KEYS *
> GET user:{userId}
```

## 🚀 Production Deployment

1. Set `synchronize: false` trong TypeORM config
2. Use production database
3. Set strong JWT secrets
4. Configure Redis with password
5. Use HTTPS for all endpoints
6. Set proper CORS origins
7. Enable rate limiting
8. Add logging & monitoring (Winston, Sentry)
9. Use environment-specific .env files

## 📄 License

MIT

## 👨‍💻 Author

Microservices Authentication & User Management System

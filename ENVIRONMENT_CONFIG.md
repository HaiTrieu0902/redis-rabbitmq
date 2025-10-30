# Environment Configuration Summary

This document shows the correct environment configuration for all three microservices.

## ✅ Corrected Configuration

All services now use the same database, Redis, RabbitMQ, and JWT secret.

---

## 1. Authentication Service (Port 8000)

**File:** `authenication/.env`

```env
# Database
DATABASE_URL=postgresql://postgres:040202005173@localhost:5432/vuihoi

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379

# RabbitMQ
RABBITMQ_HOST=localhost
RABBITMQ_PORT=5672
RABBITMQ_USER=guest
RABBITMQ_PASSWORD=guest

# JWT
JWT_SECRET=your-super-secret-jwt-key-change-in-production
JWT_EXPIRES_IN=15m
JWT_REFRESH_SECRET=your-super-secret-refresh-jwt-key-change-in-production
JWT_REFRESH_EXPIRES_IN=7d

# OAuth (Google, GitHub, Microsoft)
# ... your OAuth credentials ...

# Service
PORT=8000
```

---

## 2. User Service (Port 3002)

**File:** `user/.env`

```env
# Database
DATABASE_URL=postgresql://postgres:040202005173@localhost:5432/vuihoi

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379

# RabbitMQ
RABBITMQ_HOST=localhost
RABBITMQ_PORT=5672
RABBITMQ_USER=guest
RABBITMQ_PASSWORD=guest

# JWT (must match authentication service)
JWT_SECRET=your-super-secret-jwt-key-change-in-production

# Service
PORT=3002
AUTH_SERVICE_HOST=localhost
AUTH_SERVICE_PORT=8000
```

---

## 3. Streaming Service (Port 3003)

**File:** `streaming/.env`

```env
# Server Configuration
PORT=3003
GIN_MODE=debug

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=040202005173
DB_NAME=vuihoi
DB_SSLMODE=disable

# JWT Configuration (must match authentication service)
JWT_SECRET=your-super-secret-jwt-key-change-in-production

# RabbitMQ Configuration
RABBITMQ_URL=amqp://guest:guest@localhost:5672/

# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
```

---

## Key Configuration Details

### Database

- **Host:** localhost
- **Port:** 5432
- **User:** postgres
- **Password:** `040202005173`
- **Database Name:** `vuihoi`

### Redis

- **Host:** localhost
- **Port:** 6379
- **Password:** (none)

### RabbitMQ

- **Host:** localhost
- **Port:** 5672
- **User:** guest
- **Password:** guest
- **URL:** `amqp://guest:guest@localhost:5672/`

### JWT

- **Secret:** `your-super-secret-jwt-key-change-in-production`
- ⚠️ **IMPORTANT:** All services must use the **same JWT_SECRET**!

---

## Running the Services

### 1. Start Infrastructure

Make sure PostgreSQL, Redis, and RabbitMQ are running:

```bash
# Check PostgreSQL
psql -U postgres -d vuihoi -c "SELECT 1"

# Check Redis
redis-cli ping

# Check RabbitMQ
curl http://localhost:15672
```

### 2. Start Authentication Service

```bash
cd authenication
npm install
npm run start:dev
```

Service runs on: http://localhost:8000

### 3. Start User Service

```bash
cd user
npm install
npm run start:dev
```

Service runs on: http://localhost:3002

### 4. Start Streaming Service

```bash
cd streaming
go mod tidy
go run main.go
```

Service runs on: http://localhost:3003

---

## Expected Output

### Authentication Service

```
✅ Database connection established
✅ RabbitMQ connection established
✅ Redis connection established
🚀 Authentication service is running on port 8000
```

### User Service

```
✅ Database connection established
✅ RabbitMQ connection established
✅ Redis connection established
🚀 User service is running on port 3002
```

### Streaming Service

```
✓ Database connection established
✓ RabbitMQ connection established
🚀 Streaming service is running on port 3003
```

---

## Troubleshooting

### Database Connection Error

**Error:** `password authentication failed for user "postgres"`

**Solution:**

1. Check password in `.env` matches your PostgreSQL password
2. Update all three `.env` files with correct password
3. Test connection: `psql -U postgres -d vuihoi`

### RabbitMQ Connection Error

**Error:** `Failed to connect to RabbitMQ`

**Solution:**

1. Ensure RabbitMQ is running
2. Check credentials: `guest:guest`
3. Verify URL: `amqp://guest:guest@localhost:5672/`
4. Test management UI: http://localhost:15672

### JWT Token Invalid

**Error:** `invalid token` or `unauthorized`

**Solution:**

1. Ensure all services use the **same** `JWT_SECRET`
2. Copy the JWT_SECRET from authentication service to other services
3. Restart all services after changing JWT_SECRET

### Redis Connection Error

**Error:** `Failed to connect to Redis`

**Solution:**

1. Ensure Redis is running: `redis-cli ping`
2. Check host and port in `.env`
3. If Redis has a password, add it to `REDIS_PASSWORD`

---

## Quick Fix Checklist

- [ ] PostgreSQL is running
- [ ] Redis is running
- [ ] RabbitMQ is running
- [ ] All `.env` files have correct database password
- [ ] All `.env` files have same JWT_SECRET
- [ ] Database name is `vuihoi` in all services
- [ ] RabbitMQ URL is correct
- [ ] All services can connect to infrastructure

---

## Test Connection

Run this command to test the streaming service:

```bash
cd streaming
go run main.go
```

If you see:

```
✓ Database connection established
✓ RabbitMQ connection established
🚀 Streaming service is running on port 3003
```

Then everything is configured correctly! ✅

---

## Next Steps

1. ✅ Fix all `.env` files
2. ✅ Test database connection
3. ✅ Start all services
4. 🧪 Test with Postman collection
5. 🎉 Start building!

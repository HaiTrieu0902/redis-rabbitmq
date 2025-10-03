# Architecture Documentation

## System Overview

```
┌───────────────────────────────────────────────────────────────┐
│                      Client Application                        │
│                     (Frontend - Port 3000)                     │
└───────────────────────────┬───────────────────────────────────┘
                            │
                            │ HTTP/HTTPS Requests
                            │
        ┌───────────────────┴────────────────────┐
        │                                        │
        ▼                                        ▼
┌────────────────────┐                  ┌────────────────────┐
│ Authentication     │                  │ User Service       │
│ Service            │                  │ (Port 3002)        │
│ (Port 3001)        │◄────────────────►│                    │
│                    │   Redis Pub/Sub  │                    │
│ - OAuth (Google)   │                  │ - User CRUD        │
│ - OAuth (GitHub)   │                  │ - Profile Mgmt     │
│ - OAuth (MSAL)     │                  │ - Redis Cache      │
│ - JWT Generation   │                  │ - JWT Validation   │
│ - Token Refresh    │                  │                    │
│ - Session Mgmt     │                  │                    │
└─────────┬──────────┘                  └──────────┬─────────┘
          │                                        │
          │                                        │
          └───────────┬────────────────────────────┘
                      │
          ┌───────────┴───────────────┐
          │                           │
          ▼                           ▼
┌──────────────────┐         ┌──────────────────┐
│ Redis            │         │ PostgreSQL       │
│ (Port 6379)      │         │ (Port 5432)      │
│                  │         │                  │
│ - Session Store  │         │ - users          │
│ - Token Cache    │         │ - user_identity  │
│ - Pub/Sub Queue  │         │ - conversations  │
│ - User Cache     │         │ - messages       │
│ - Token Blacklist│         │ - articles       │
└──────────────────┘         └──────────────────┘
```

## Authentication Flow

### OAuth Flow (Google/GitHub/Microsoft)

```
┌────────┐         ┌─────────┐         ┌──────────┐         ┌──────────┐
│ Client │         │  Auth   │         │  OAuth   │         │   DB     │
│        │         │ Service │         │ Provider │         │          │
└───┬────┘         └────┬────┘         └────┬─────┘         └────┬─────┘
    │                   │                   │                    │
    │ 1. GET /auth/google                   │                    │
    ├──────────────────►│                   │                    │
    │                   │                   │                    │
    │ 2. Redirect to OAuth                  │                    │
    │◄──────────────────┤                   │                    │
    │                   │                   │                    │
    │ 3. User Login & Consent               │                    │
    ├──────────────────────────────────────►│                    │
    │                   │                   │                    │
    │ 4. Callback with code                 │                    │
    │◄──────────────────────────────────────┤                    │
    │                   │                   │                    │
    │ 5. Forward code   │                   │                    │
    ├──────────────────►│                   │                    │
    │                   │ 6. Exchange token │                    │
    │                   ├──────────────────►│                    │
    │                   │                   │                    │
    │                   │ 7. Access Token   │                    │
    │                   │◄──────────────────┤                    │
    │                   │                   │                    │
    │                   │ 8. Get User Info  │                    │
    │                   ├──────────────────►│                    │
    │                   │◄──────────────────┤                    │
    │                   │                   │                    │
    │                   │ 9. Find/Create User                    │
    │                   ├───────────────────────────────────────►│
    │                   │◄───────────────────────────────────────┤
    │                   │                   │                    │
    │                   │ 10. Generate JWT  │                    │
    │                   │                   │                    │
    │ 11. JWT Tokens    │                   │                    │
    │◄──────────────────┤                   │                    │
    │                   │                   │                    │
```

### JWT Authentication Flow

```
┌────────┐         ┌─────────┐         ┌──────────┐         ┌──────────┐
│ Client │         │  User   │         │  Redis   │         │   DB     │
│        │         │ Service │         │          │         │          │
└───┬────┘         └────┬────┘         └────┬─────┘         └────┬─────┘
    │                   │                   │                    │
    │ 1. API Request + JWT                  │                    │
    ├──────────────────►│                   │                    │
    │                   │                   │                    │
    │                   │ 2. Verify JWT     │                    │
    │                   │                   │                    │
    │                   │ 3. Check Blacklist│                    │
    │                   ├──────────────────►│                    │
    │                   │◄──────────────────┤                    │
    │                   │                   │                    │
    │                   │ 4. Get Cached User│                    │
    │                   ├──────────────────►│                    │
    │                   │◄──────────────────┤                    │
    │                   │                   │                    │
    │                   │ 5. If not cached, query DB             │
    │                   ├───────────────────────────────────────►│
    │                   │◄───────────────────────────────────────┤
    │                   │                   │                    │
    │                   │ 6. Cache result   │                    │
    │                   ├──────────────────►│                    │
    │                   │                   │                    │
    │ 7. Response       │                   │                    │
    │◄──────────────────┤                   │                    │
    │                   │                   │                    │
```

## Data Flow

### User Update Flow with Redis Pub/Sub

```
┌────────┐    ┌─────────┐    ┌──────────┐    ┌─────────┐    ┌──────────┐
│ Client │    │  Auth   │    │  Redis   │    │  User   │    │   DB     │
│        │    │ Service │    │ Pub/Sub  │    │ Service │    │          │
└───┬────┘    └────┬────┘    └────┬─────┘    └────┬────┘    └────┬─────┘
    │              │              │              │              │
    │ 1. OAuth Login              │              │              │
    ├─────────────►│              │              │              │
    │              │              │              │              │
    │              │ 2. Update DB │              │              │
    │              ├─────────────────────────────────────────────►
    │              │◄────────────────────────────────────────────┤
    │              │              │              │              │
    │              │ 3. Publish Event            │              │
    │              ├─────────────►│              │              │
    │              │              │              │              │
    │              │              │ 4. Notify    │              │
    │              │              ├─────────────►│              │
    │              │              │              │              │
    │              │              │              │ 5. Invalidate Cache
    │              │              │              │              │
    │              │ 6. Return JWT               │              │
    │◄─────────────┤              │              │              │
    │              │              │              │              │
```

### Caching Strategy

```
┌─────────────────────────────────────────────────────────┐
│                  Redis Cache Layers                     │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  ┌───────────────────────────────────────────────────┐ │
│  │ Session Storage (TTL: 7 days)                     │ │
│  │ Key: session:{userId}                             │ │
│  │ Value: { email, name, loginTime }                 │ │
│  └───────────────────────────────────────────────────┘ │
│                                                         │
│  ┌───────────────────────────────────────────────────┐ │
│  │ Token Cache (TTL: 7 days)                         │ │
│  │ Key: token:{userId}                               │ │
│  │ Value: { accessToken, refreshToken }              │ │
│  └───────────────────────────────────────────────────┘ │
│                                                         │
│  ┌───────────────────────────────────────────────────┐ │
│  │ Token Blacklist (TTL: 15 minutes)                 │ │
│  │ Key: blacklist:{token}                            │ │
│  │ Value: true                                       │ │
│  └───────────────────────────────────────────────────┘ │
│                                                         │
│  ┌───────────────────────────────────────────────────┐ │
│  │ User Data Cache (TTL: 1 hour)                     │ │
│  │ Key: user:{userId}                                │ │
│  │ Value: { id, email, name, avatar_url, ... }       │ │
│  └───────────────────────────────────────────────────┘ │
│                                                         │
└─────────────────────────────────────────────────────────┘
```

## Database Schema

```sql
┌─────────────────────────────────────────┐
│              users                       │
├──────────────┬──────────────────────────┤
│ id           │ UUID PRIMARY KEY         │
│ name         │ VARCHAR(255)             │
│ email        │ VARCHAR(255) UNIQUE      │
│ avatar_url   │ TEXT                     │
│ created_at   │ TIMESTAMP                │
│ updated_at   │ TIMESTAMP                │
└──────────────┴──────────────────────────┘
                     │
                     │ 1:N
                     │
                     ▼
┌─────────────────────────────────────────┐
│          user_identity                   │
├──────────────┬──────────────────────────┤
│ id           │ UUID PRIMARY KEY         │
│ user_id      │ UUID FK → users(id)      │
│ provider     │ VARCHAR(50)              │
│              │ ('google','github',      │
│              │  'msal')                 │
│ provider_user│ VARCHAR(255)             │
│ _id          │                          │
│ access_token │ TEXT                     │
│ refresh_token│ TEXT                     │
│ created_at   │ TIMESTAMP                │
│ updated_at   │ TIMESTAMP                │
│              │                          │
│ UNIQUE(provider, provider_user_id)      │
└─────────────────────────────────────────┘
```

## Security Model

```
┌─────────────────────────────────────────────────────────┐
│                   Security Layers                        │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  Layer 1: OAuth 2.0                                     │
│  ├─ Google OAuth (OpenID Connect)                      │
│  ├─ GitHub OAuth                                        │
│  └─ Microsoft MSAL (Azure AD)                          │
│                                                         │
│  Layer 2: JWT Authentication                            │
│  ├─ Access Token (15 min expiry)                       │
│  ├─ Refresh Token (7 days expiry)                      │
│  ├─ Token Signature Verification                       │
│  └─ Token Blacklist (logout)                           │
│                                                         │
│  Layer 3: Session Management                            │
│  ├─ Redis Session Store                                │
│  ├─ Session Expiration                                 │
│  └─ Session Invalidation                               │
│                                                         │
│  Layer 4: Input Validation                              │
│  ├─ class-validator                                    │
│  ├─ class-transformer                                  │
│  └─ DTO validation                                     │
│                                                         │
│  Layer 5: Database Security                             │
│  ├─ TypeORM (SQL injection prevention)                │
│  ├─ Parameterized queries                              │
│  └─ Foreign key constraints                            │
│                                                         │
└─────────────────────────────────────────────────────────┘
```

## Deployment Architecture

```
┌─────────────────────────────────────────────────────────┐
│                    Production Setup                      │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  ┌───────────────────────────────────────────────────┐ │
│  │ Load Balancer (Nginx/AWS ALB)                     │ │
│  └──────────────────┬────────────────────────────────┘ │
│                     │                                   │
│         ┌───────────┴───────────┐                      │
│         │                       │                      │
│         ▼                       ▼                      │
│  ┌─────────────┐         ┌─────────────┐              │
│  │ Auth Service│         │ Auth Service│              │
│  │ Instance 1  │         │ Instance 2  │              │
│  └─────────────┘         └─────────────┘              │
│         │                       │                      │
│         ▼                       ▼                      │
│  ┌─────────────┐         ┌─────────────┐              │
│  │ User Service│         │ User Service│              │
│  │ Instance 1  │         │ Instance 2  │              │
│  └─────────────┘         └─────────────┘              │
│         │                       │                      │
│         └───────────┬───────────┘                      │
│                     │                                   │
│         ┌───────────┴───────────┐                      │
│         │                       │                      │
│         ▼                       ▼                      │
│  ┌─────────────┐         ┌─────────────┐              │
│  │ Redis       │         │ PostgreSQL  │              │
│  │ Cluster     │         │ (Master +   │              │
│  │ (Sentinel)  │         │  Replicas)  │              │
│  └─────────────┘         └─────────────┘              │
│                                                         │
└─────────────────────────────────────────────────────────┘
```

## Technology Stack

```
┌─────────────────────────────────────────────────────────┐
│                   Technology Stack                       │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  Backend Framework:                                     │
│  └─ NestJS 11.x (Node.js, TypeScript)                  │
│                                                         │
│  Database:                                              │
│  └─ PostgreSQL 14+ (with TypeORM)                      │
│                                                         │
│  Cache & Message Queue:                                 │
│  └─ Redis 7.x (with ioredis)                           │
│                                                         │
│  Authentication:                                        │
│  ├─ Passport.js                                        │
│  ├─ JWT (@nestjs/jwt)                                  │
│  ├─ passport-google-oauth20                            │
│  ├─ passport-github2                                   │
│  └─ @azure/msal-node                                   │
│                                                         │
│  Validation:                                            │
│  ├─ class-validator                                    │
│  └─ class-transformer                                  │
│                                                         │
│  Containerization:                                      │
│  └─ Docker & Docker Compose                            │
│                                                         │
└─────────────────────────────────────────────────────────┘
```

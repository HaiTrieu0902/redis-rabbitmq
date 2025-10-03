# RabbitMQ Integration Guide

## Overview

This project now uses **RabbitMQ** as a message broker for inter-service communication between the Authentication Service and User Service. RabbitMQ provides reliable, asynchronous messaging with features like message persistence, acknowledgments, and dead-letter queues.

## Architecture

```
┌─────────────────────────┐         ┌──────────────┐         ┌─────────────────────┐
│  Authentication Service │────────▶│   RabbitMQ   │────────▶│    User Service     │
│                         │         │   Exchange   │         │                     │
│  - Publishes Events:    │         │              │         │  - Consumes Events: │
│    • user.created       │         │  Topic-based │         │    • user.created   │
│    • user.updated       │         │   Routing    │         │    • user.updated   │
│                         │         │              │         │    • user.deleted   │
└─────────────────────────┘         └──────────────┘         └─────────────────────┘
```

## Message Flow

### 1. User Authentication (OAuth Login)

When a user logs in via OAuth (Google, GitHub, Microsoft):

1. **Authentication Service** creates or updates the user in the database
2. **Authentication Service** publishes an event to RabbitMQ:
   - `user.created` (for new users)
   - `user.updated` (for existing users)
3. **User Service** consumes the event from its queue
4. **User Service** invalidates cached user data in Redis

### 2. User Profile Update

When a user updates their profile in the User Service:

1. **User Service** updates the user in the database
2. **User Service** publishes `user.updated` event to RabbitMQ
3. Other services can consume this event to stay synchronized

### 3. User Deletion

When a user is deleted:

1. **User Service** deletes the user from the database
2. **User Service** publishes `user.deleted` event to RabbitMQ
3. Consumers handle cleanup (cache invalidation, etc.)

## Event Schema

### user.created / user.updated

```json
{
  "userId": "uuid-string",
  "email": "user@example.com",
  "name": "John Doe",
  "avatar_url": "https://example.com/avatar.jpg",
  "timestamp": "2025-10-03T12:34:56.789Z",
  "action": "created" // or "updated"
}
```

### user.deleted

```json
{
  "userId": "uuid-string",
  "action": "deleted",
  "timestamp": "2025-10-03T12:34:56.789Z"
}
```

## Configuration

### Environment Variables

Add these variables to your `.env` files in both services:

```bash
# RabbitMQ Configuration
RABBITMQ_HOST=localhost
RABBITMQ_PORT=5672
RABBITMQ_USER=admin
RABBITMQ_PASSWORD=admin
```

### Docker Compose

RabbitMQ is already configured in `docker-compose.yaml`:

```yaml
rabbitmq:
  image: rabbitmq:3-management
  container_name: rabbitmq
  restart: always
  ports:
    - '5672:5672' # AMQP port
    - '15672:15672' # Management UI
  environment:
    RABBITMQ_DEFAULT_USER: admin
    RABBITMQ_DEFAULT_PASS: admin
  volumes:
    - rabbitmq_data:/var/lib/rabbitmq
```

## Getting Started

### 1. Start RabbitMQ

```bash
docker-compose up -d rabbitmq
```

### 2. Access RabbitMQ Management UI

Open your browser and navigate to:

- URL: http://localhost:15672
- Username: `admin`
- Password: `admin`

### 3. Copy Environment Files

**Authentication Service:**

```bash
cd authenication
cp .env.example .env
```

**User Service:**

```bash
cd user
cp .env.example .env
```

Make sure to add the RabbitMQ configuration to your `.env` files.

### 4. Install Dependencies

Both services already have the RabbitMQ client installed:

```bash
npm install amqplib @types/amqplib
```

### 5. Start Services

**Authentication Service:**

```bash
cd authenication
npm run start:dev
```

**User Service:**

```bash
cd user
npm run start:dev
```

## Implementation Details

### RabbitMQ Service (`rabbitmq.service.ts`)

Both services include a `RabbitMQService` that provides:

- **Connection Management**: Automatic connection to RabbitMQ with error handling
- **Exchange Declaration**: Creates a topic exchange named `user_events`
- **Message Publishing**: Methods to publish user events
- **Message Consumption**: Subscribe to queues and handle messages
- **Graceful Shutdown**: Properly closes connections on application shutdown

### Key Features

#### 1. **Topic-based Routing**

Messages are routed using topic patterns:

- `user.created` - New user created
- `user.updated` - User information updated
- `user.deleted` - User deleted

#### 2. **Message Persistence**

Messages are marked as persistent, ensuring they survive RabbitMQ restarts.

#### 3. **Manual Acknowledgments**

Messages are manually acknowledged after successful processing, preventing message loss.

#### 4. **Error Handling**

Failed messages are rejected and requeued for retry (with `nack`).

#### 5. **Dual Publishing**

For backward compatibility, events are published to both:

- **RabbitMQ** (primary message broker)
- **Redis Pub/Sub** (legacy support)

## Code Examples

### Publishing an Event (Authentication Service)

```typescript
// In auth.service.ts
await this.rabbitMQService.publishUserCreated({
  userId: user.id,
  email: user.email,
  name: user.name,
  avatar_url: user.avatar_url,
  timestamp: new Date().toISOString(),
});
```

### Consuming Events (User Service)

```typescript
// In user.service.ts
await this.rabbitMQService.consume('user_service_queue', ['user.created', 'user.updated', 'user.deleted'], async (message: UserEventPayload) => {
  // Handle the message
  console.log('Received event:', message);
});
```

## Monitoring

### RabbitMQ Management UI

The management UI provides:

- **Connections**: View active connections from services
- **Channels**: Monitor channel activity
- **Exchanges**: See the `user_events` exchange and its bindings
- **Queues**: Monitor `user_service_queue` and message rates
- **Message Rates**: Publish/deliver rates, acknowledgments

### Logs

Both services log RabbitMQ activity:

**Authentication Service:**

```
[RabbitMQService] ✅ RabbitMQ connection established
[AuthService] 🐰 Published user.created event to RabbitMQ for user: abc-123
```

**User Service:**

```
[RabbitMQService] ✅ RabbitMQ connection established
[RabbitMQService] Queue "user_service_queue" bound to routing key "user.created"
[UserService] 🐰 Processing RabbitMQ event: {...}
```

## Benefits of RabbitMQ

### vs Redis Pub/Sub

| Feature                | RabbitMQ | Redis Pub/Sub |
| ---------------------- | -------- | ------------- |
| Message Persistence    | ✅ Yes   | ❌ No         |
| Message Acknowledgment | ✅ Yes   | ❌ No         |
| Message Queuing        | ✅ Yes   | ❌ No         |
| Guaranteed Delivery    | ✅ Yes   | ❌ No         |
| Complex Routing        | ✅ Yes   | ⚠️ Limited    |
| Dead Letter Queues     | ✅ Yes   | ❌ No         |
| Message Replay         | ✅ Yes   | ❌ No         |

### Why Use Both?

This implementation uses **both RabbitMQ and Redis**:

- **RabbitMQ**: For reliable, asynchronous inter-service communication
- **Redis**: For caching and fast in-memory data access

## Troubleshooting

### Connection Issues

If services can't connect to RabbitMQ:

1. Verify RabbitMQ is running:

   ```bash
   docker ps | grep rabbitmq
   ```

2. Check RabbitMQ logs:

   ```bash
   docker logs rabbitmq
   ```

3. Verify environment variables in `.env` files

### Message Not Being Consumed

1. Check if the queue exists in RabbitMQ Management UI
2. Verify queue bindings to the exchange
3. Check service logs for consumer setup
4. Ensure User Service is running

### Messages Stuck in Queue

1. Check for errors in User Service logs
2. View message details in RabbitMQ Management UI
3. Messages with errors are requeued automatically

## Best Practices

1. **Always acknowledge messages** after successful processing
2. **Use message persistence** for critical events
3. **Implement idempotent handlers** (handle duplicate messages gracefully)
4. **Monitor queue sizes** to detect processing bottlenecks
5. **Set up dead-letter queues** for failed messages (future enhancement)
6. **Use correlation IDs** for tracing messages across services (future enhancement)

## Future Enhancements

- [ ] Add dead-letter exchange for failed messages
- [ ] Implement message correlation IDs for tracing
- [ ] Add retry logic with exponential backoff
- [ ] Implement circuit breaker pattern
- [ ] Add message TTL (time-to-live)
- [ ] Set up monitoring and alerting
- [ ] Add more event types (password reset, email verification, etc.)

## References

- [RabbitMQ Documentation](https://www.rabbitmq.com/documentation.html)
- [amqplib (Node.js Client)](https://amqp-node.github.io/amqplib/)
- [RabbitMQ Best Practices](https://www.rabbitmq.com/best-practices.html)
- [Message Queue Patterns](https://www.enterpriseintegrationpatterns.com/patterns/messaging/)

---

**Note**: This implementation maintains backward compatibility with Redis Pub/Sub. You can gradually migrate to RabbitMQ-only messaging in the future.

# RabbitMQ Timing Issue Fix

## Problem

The User Service was crashing on startup with the error:

```
Cannot read properties of undefined (reading 'assertQueue')
at RabbitMQService.consume
```

## Root Cause

The issue occurred because of the initialization order in NestJS:

1. **UserService constructor** runs first
2. Constructor calls `subscribeToRabbitMQEvents()`
3. This tries to use `this.rabbitMQService.consume()`
4. But `RabbitMQService.onModuleInit()` hasn't completed yet
5. So `this.channel` is still `undefined`

**Timeline:**

```
1. UserService constructor ───┐
   │                           │
   ├─ subscribeToRedisEvents() │
   │                           │
   ├─ subscribeToRabbitMQ()────┼──> ❌ ERROR: channel is undefined
   │                           │
2. RabbitMQService.onModuleInit() ──> ✅ Connects and creates channel
```

## Solution

Implemented the `OnModuleInit` lifecycle hook in `UserService` to delay RabbitMQ subscription until after the module is fully initialized.

### Changes Made to `user/src/user/user.service.ts`

**Before:**

```typescript
@Injectable()
export class UserService {
  constructor(
    @InjectRepository(User)
    private userRepository: Repository<User>,
    private redisService: RedisService,
    private rabbitMQService: RabbitMQService,
  ) {
    this.subscribeToRedisEvents();
    this.subscribeToRabbitMQEvents(); // ❌ Too early!
  }
}
```

**After:**

```typescript
@Injectable()
export class UserService implements OnModuleInit {
  constructor(
    @InjectRepository(User)
    private userRepository: Repository<User>,
    private redisService: RedisService,
    private rabbitMQService: RabbitMQService,
  ) {
    this.subscribeToRedisEvents(); // ✅ Redis is fine (doesn't need async setup)
  }

  async onModuleInit() {
    // Wait for RabbitMQ connection to be established
    setTimeout(() => {
      void this.subscribeToRabbitMQEvents(); // ✅ Now runs after RabbitMQ is connected
    }, 1000);
  }
}
```

## New Timeline (Fixed)

```
1. UserService constructor ───┐
   │                           │
   ├─ subscribeToRedisEvents() │  ✅ Works fine
   │                           │
2. RabbitMQService.onModuleInit() ──> Connects to RabbitMQ
   │
3. UserService.onModuleInit()
   │
   ├─ setTimeout (1 second delay)
   │
   └─> subscribeToRabbitMQEvents() ──> ✅ Channel is now available!
```

## Why This Works

1. **NestJS Lifecycle Order:**

   - Constructors run first
   - Then `onModuleInit()` hooks run (in dependency order)
   - Global modules (like RabbitMQModule) initialize before dependent services

2. **setTimeout Delay:**

   - Gives RabbitMQ extra time to establish connection
   - Pushes subscription to next event loop tick
   - Ensures channel is fully initialized

3. **void Operator:**
   - Explicitly marks the async call as "fire and forget"
   - Satisfies TypeScript/ESLint rules
   - Promise errors are still caught by the try/catch inside `subscribeToRabbitMQEvents()`

## Why Authentication Service Doesn't Need This

The Authentication Service only **publishes** messages when users authenticate (on-demand action). It doesn't set up consumers in the constructor, so timing isn't an issue:

```typescript
// In auth.service.ts - only publishes when called
await this.rabbitMQService.publishUserCreated(payload);
```

## Testing

After this fix, the User Service should:

1. ✅ Start successfully without errors
2. ✅ Connect to RabbitMQ
3. ✅ Subscribe to the queue after connection is established
4. ✅ Receive and process messages from Authentication Service

## Logs to Expect

**Successful Startup:**

```
[RabbitMQService] ✅ RabbitMQ connection established
[RabbitMQService] Exchange "user_events" declared
[RabbitMQService] Queue "user_service_queue" bound to routing key "user.created"
[RabbitMQService] Queue "user_service_queue" bound to routing key "user.updated"
[RabbitMQService] Queue "user_service_queue" bound to routing key "user.deleted"
[RabbitMQService] 📨 Started consuming from queue "user_service_queue"
```

## Alternative Solutions (if needed)

If the 1-second timeout seems arbitrary, you could:

1. **Use a flag in RabbitMQService:**

   ```typescript
   private isReady = false;

   async onModuleInit() {
     await this.connect();
     this.isReady = true;
   }

   async consume() {
     if (!this.isReady) {
       throw new Error('RabbitMQ not ready');
     }
     // ...
   }
   ```

2. **Use an event emitter:**

   ```typescript
   // In RabbitMQService
   async onModuleInit() {
     await this.connect();
     this.eventEmitter.emit('rabbitmq:ready');
   }

   // In UserService
   this.eventEmitter.on('rabbitmq:ready', () => {
     this.subscribeToRabbitMQEvents();
   });
   ```

3. **Make consume() wait for connection:**
   ```typescript
   async consume() {
     while (!this.channel) {
       await new Promise(resolve => setTimeout(resolve, 100));
     }
     // proceed with consume logic
   }
   ```

For now, the `setTimeout` solution is simple and effective! ✅

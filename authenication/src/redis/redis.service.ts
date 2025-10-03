import { Injectable, Inject } from '@nestjs/common';
import Redis from 'ioredis';
import { REDIS_CLIENT } from './redis.module';

@Injectable()
export class RedisService {
  constructor(@Inject(REDIS_CLIENT) private readonly redisClient: Redis) {}

  async get(key: string): Promise<string | null> {
    return await this.redisClient.get(key);
  }

  async set(
    key: string,
    value: string,
    expirationSeconds?: number,
  ): Promise<void> {
    if (expirationSeconds) {
      await this.redisClient.setex(key, expirationSeconds, value);
    } else {
      await this.redisClient.set(key, value);
    }
  }

  async del(key: string): Promise<void> {
    await this.redisClient.del(key);
  }

  async exists(key: string): Promise<boolean> {
    const result = await this.redisClient.exists(key);
    return result === 1;
  }

  async setJson(
    key: string,
    value: any,
    expirationSeconds?: number,
  ): Promise<void> {
    console.log(
      `[RedisService] Setting JSON key: ${key}, TTL: ${expirationSeconds || 'none'}`,
    );
    const jsonString = JSON.stringify(value);
    await this.set(key, jsonString, expirationSeconds);
    console.log(`[RedisService] ✅ Successfully set key: ${key}`);
  }

  async getJson<T>(key: string): Promise<T | null> {
    const value = await this.get(key);
    if (!value) return null;
    try {
      return JSON.parse(value) as T;
    } catch (error) {
      return null;
    }
  }

  // For session management
  async setSession(
    userId: string,
    sessionData: any,
    expirationSeconds: number = 3600,
  ): Promise<void> {
    const key = `session:${userId}`;
    console.log(`[RedisService] Setting session key: ${key}`);
    await this.setJson(key, sessionData, expirationSeconds);
  }

  async getSession(userId: string): Promise<any> {
    const key = `session:${userId}`;
    return await this.getJson(key);
  }

  async deleteSession(userId: string): Promise<void> {
    const key = `session:${userId}`;
    await this.del(key);
  }

  // For token blacklist (logout)
  async blacklistToken(
    token: string,
    expirationSeconds: number,
  ): Promise<void> {
    const key = `blacklist:${token}`;
    await this.set(key, 'true', expirationSeconds);
  }

  async isTokenBlacklisted(token: string): Promise<boolean> {
    const key = `blacklist:${token}`;
    return await this.exists(key);
  }

  // Pub/Sub for microservices
  async publish(channel: string, message: any): Promise<void> {
    console.log(`[RedisService] Publishing to channel: ${channel}`, message);
    await this.redisClient.publish(channel, JSON.stringify(message));
    console.log(`[RedisService] ✅ Published to channel: ${channel}`);
  }

  subscribe(channel: string, callback: (message: any) => void): void {
    const subscriber = this.redisClient.duplicate();
    subscriber.subscribe(channel);
    subscriber.on('message', (ch, msg) => {
      if (ch === channel) {
        try {
          callback(JSON.parse(msg));
        } catch (error) {
          console.error('Error parsing message:', error);
        }
      }
    });
  }
}

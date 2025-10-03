import { Injectable, UnauthorizedException } from '@nestjs/common';
import { JwtService } from '@nestjs/jwt';
import { ConfigService } from '@nestjs/config';
import { InjectRepository } from '@nestjs/typeorm';
import { Repository } from 'typeorm';
import { User } from '../entities/user.entity';
import { UserIdentity } from '../entities/user-identity.entity';
import { RedisService } from '../redis/redis.service';
import { RabbitMQService } from '../rabbitmq/rabbitmq.service';
import { OAuthCallbackDto, TokenResponseDto } from '../dto/auth.dto';

@Injectable()
export class AuthService {
  constructor(
    @InjectRepository(User)
    private userRepository: Repository<User>,
    @InjectRepository(UserIdentity)
    private userIdentityRepository: Repository<UserIdentity>,
    private jwtService: JwtService,
    private configService: ConfigService,
    private redisService: RedisService,
    private rabbitMQService: RabbitMQService,
  ) {}

  async generateTokens(user: User): Promise<TokenResponseDto> {
    const payload = { sub: user.id, email: user.email };

    const accessToken = this.jwtService.sign(payload, {
      secret: this.configService.get('JWT_SECRET'),
      expiresIn: this.configService.get('JWT_EXPIRES_IN', '15m'),
    });

    const refreshToken = this.jwtService.sign(payload, {
      secret: this.configService.get('JWT_REFRESH_SECRET'),
      expiresIn: this.configService.get('JWT_REFRESH_EXPIRES_IN', '7d'),
    });

    console.log(`[AuthService] Generating tokens for user: ${user.id}`);
    console.log(`[AuthService] User email: ${user.email}, name: ${user.name}`);

    // Cache token in Redis
    try {
      await this.redisService.setJson(
        `token:${user.id}`,
        { accessToken, refreshToken },
        7 * 24 * 60 * 60, // 7 days
      );
      console.log(`[AuthService] ✅ Token cached in Redis: token:${user.id}`);
    } catch (error) {
      console.error(`[AuthService] ❌ Failed to cache token in Redis:`, error);
    }

    // Store session
    try {
      await this.redisService.setSession(
        user.id,
        { email: user.email, name: user.name },
        7 * 24 * 60 * 60,
      );
      console.log(
        `[AuthService] ✅ Session stored in Redis: session:${user.id}`,
      );
    } catch (error) {
      console.error(
        `[AuthService] ❌ Failed to store session in Redis:`,
        error,
      );
    }

    return {
      access_token: accessToken,
      refresh_token: refreshToken,
      expires_in: 900, // 15 minutes
      user: {
        id: user.id,
        email: user.email,
        name: user.name,
        avatar_url: user.avatar_url,
      },
    };
  }

  async handleOAuthCallback(
    oauthData: OAuthCallbackDto,
  ): Promise<TokenResponseDto> {
    // Validate required fields
    if (!oauthData.provider_user_id || !oauthData.email) {
      throw new UnauthorizedException('Missing required OAuth data');
    }

    // Find or create user identity
    let userIdentity = await this.userIdentityRepository.findOne({
      where: {
        provider: oauthData.provider,
        provider_user_id: oauthData.provider_user_id,
      },
      relations: ['user'],
    });

    let user: User;
    let isNewUser = false;

    if (userIdentity) {
      // User exists, update tokens
      userIdentity.access_token = oauthData.access_token || null;
      userIdentity.refresh_token = oauthData.refresh_token || null;
      await this.userIdentityRepository.save(userIdentity);
      user = userIdentity.user;

      // Update user info if changed
      const needsUpdate =
        (oauthData.email && user.email !== oauthData.email) ||
        (oauthData.name && user.name !== oauthData.name) ||
        (oauthData.avatar_url && user.avatar_url !== oauthData.avatar_url);

      if (needsUpdate) {
        if (oauthData.email) user.email = oauthData.email;
        if (oauthData.name) user.name = oauthData.name;
        if (oauthData.avatar_url) user.avatar_url = oauthData.avatar_url;
        await this.userRepository.save(user);
      }
    } else {
      // Check if user exists with this email
      user = await this.userRepository.findOne({
        where: { email: oauthData.email },
      });

      if (!user) {
        // Create new user
        user = this.userRepository.create({
          email: oauthData.email,
          name: oauthData.name || oauthData.email.split('@')[0],
          avatar_url: oauthData.avatar_url || null,
        });
        await this.userRepository.save(user);
        isNewUser = true;
      }

      // Create identity
      userIdentity = this.userIdentityRepository.create({
        user_id: user.id,
        provider: oauthData.provider,
        provider_user_id: oauthData.provider_user_id,
        access_token: oauthData.access_token || null,
        refresh_token: oauthData.refresh_token || null,
      });
      await this.userIdentityRepository.save(userIdentity);
    }

    // Publish user creation/update event to Redis (for backward compatibility)
    try {
      await this.redisService.publish('user:updated', {
        userId: user.id,
        email: user.email,
        name: user.name,
        avatar_url: user.avatar_url,
        timestamp: new Date().toISOString(),
      });
      console.log(
        `[AuthService] 📢 Published user:updated event to Redis for user: ${user.id}`,
      );
    } catch (error) {
      console.error(
        `[AuthService] ❌ Failed to publish user:updated event to Redis:`,
        error,
      );
    }

    // Publish event to RabbitMQ
    try {
      const eventPayload = {
        userId: user.id,
        email: user.email,
        name: user.name,
        avatar_url: user.avatar_url,
        timestamp: new Date().toISOString(),
      };

      if (isNewUser) {
        await this.rabbitMQService.publishUserCreated(eventPayload);
        console.log(
          `[AuthService] 🐰 Published user.created event to RabbitMQ for user: ${user.id}`,
        );
      } else {
        await this.rabbitMQService.publishUserUpdated(eventPayload);
        console.log(
          `[AuthService] 🐰 Published user.updated event to RabbitMQ for user: ${user.id}`,
        );
      }
    } catch (error) {
      console.error(
        `[AuthService] ❌ Failed to publish event to RabbitMQ:`,
        error,
      );
    }

    return this.generateTokens(user);
  }

  async refreshTokens(refreshToken: string): Promise<TokenResponseDto> {
    try {
      // Check if token is blacklisted
      if (await this.redisService.isTokenBlacklisted(refreshToken)) {
        throw new UnauthorizedException('Token has been revoked');
      }

      const payload = this.jwtService.verify<{ sub: string; email: string }>(
        refreshToken,
        {
          secret: this.configService.get('JWT_REFRESH_SECRET'),
        },
      );

      const user = await this.userRepository.findOne({
        where: { id: payload.sub },
      });

      if (!user) {
        throw new UnauthorizedException('User not found');
      }

      return this.generateTokens(user);
    } catch {
      throw new UnauthorizedException('Invalid refresh token');
    }
  }

  async logout(userId: string, accessToken: string): Promise<void> {
    // Blacklist the access token
    const expiresIn = this.configService.get<string>('JWT_EXPIRES_IN', '15m');
    const seconds = this.parseTimeToSeconds(expiresIn);
    await this.redisService.blacklistToken(accessToken, seconds);

    // Delete session
    await this.redisService.deleteSession(userId);

    // Delete cached tokens
    await this.redisService.del(`token:${userId}`);
  }

  async validateUser(userId: string): Promise<User> {
    const user = await this.userRepository.findOne({
      where: { id: userId },
    });

    if (!user) {
      throw new UnauthorizedException('User not found');
    }

    return user;
  }

  private parseTimeToSeconds(time: string): number {
    const unit = time.slice(-1);
    const value = parseInt(time.slice(0, -1));

    switch (unit) {
      case 's':
        return value;
      case 'm':
        return value * 60;
      case 'h':
        return value * 60 * 60;
      case 'd':
        return value * 24 * 60 * 60;
      default:
        return 900; // default 15 minutes
    }
  }
}

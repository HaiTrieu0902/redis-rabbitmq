import {
  Injectable,
  OnModuleInit,
  OnModuleDestroy,
  Logger,
} from '@nestjs/common';
import { ConfigService } from '@nestjs/config';
import * as amqp from 'amqplib';

export interface UserEventPayload {
  userId: string;
  email: string;
  name: string;
  avatar_url?: string;
  timestamp: string;
  action?: 'created' | 'updated' | 'deleted';
}

@Injectable()
export class RabbitMQService implements OnModuleInit, OnModuleDestroy {
  private readonly logger = new Logger(RabbitMQService.name);
  private connection: any;
  private channel: any;
  private readonly exchangeName = 'user_events';
  private readonly exchangeType = 'topic';

  constructor(private configService: ConfigService) {}

  async onModuleInit() {
    try {
      await this.connect();
      this.logger.log('✅ RabbitMQ connection established');
    } catch (error) {
      this.logger.error('❌ Failed to connect to RabbitMQ:', error);
      throw error;
    }
  }

  async onModuleDestroy() {
    await this.disconnect();
  }

  private async connect(): Promise<void> {
    const host = this.configService.get<string>('RABBITMQ_HOST', 'localhost');
    const port = this.configService.get<number>('RABBITMQ_PORT', 5672);
    const username = this.configService.get<string>('RABBITMQ_USER', 'admin');
    const password = this.configService.get<string>(
      'RABBITMQ_PASSWORD',
      'admin',
    );

    const url = `amqp://${username}:${password}@${host}:${port}`;

    this.connection = await amqp.connect(url);
    this.channel = await this.connection.createChannel();

    // Declare exchange
    await this.channel.assertExchange(this.exchangeName, this.exchangeType, {
      durable: true,
    });

    this.logger.log(`Exchange "${this.exchangeName}" declared`);

    // Handle connection errors
    this.connection.on('error', (err) => {
      this.logger.error('RabbitMQ connection error:', err);
    });

    this.connection.on('close', () => {
      this.logger.warn('RabbitMQ connection closed');
    });
  }

  async disconnect(): Promise<void> {
    try {
      if (this.channel) {
        await this.channel.close();
      }
      if (this.connection) {
        await this.connection.close();
      }
      this.logger.log('RabbitMQ connection closed gracefully');
    } catch (error) {
      this.logger.error('Error closing RabbitMQ connection:', error);
    }
  }

  /**
   * Publish a message to RabbitMQ exchange
   * @param routingKey - The routing key for the message (e.g., 'user.created', 'user.updated')
   * @param message - The message payload
   */
  async publish(routingKey: string, message: any): Promise<boolean> {
    try {
      if (!this.channel) {
        throw new Error('RabbitMQ channel is not initialized');
      }

      const messageBuffer = Buffer.from(JSON.stringify(message));

      const published = this.channel.publish(
        this.exchangeName,
        routingKey,
        messageBuffer,
        {
          persistent: true,
          contentType: 'application/json',
          timestamp: Date.now(),
        },
      );

      if (published) {
        this.logger.log(
          `📤 Published message to "${routingKey}": ${JSON.stringify(message)}`,
        );
      } else {
        this.logger.warn(
          `⚠️ Message to "${routingKey}" was not published (buffer full)`,
        );
      }

      return published;
    } catch (error) {
      this.logger.error(
        `❌ Failed to publish message to "${routingKey}":`,
        error,
      );
      throw error;
    }
  }

  /**
   * Publish user created event
   */
  async publishUserCreated(payload: UserEventPayload): Promise<void> {
    await this.publish('user.created', {
      ...payload,
      action: 'created',
    });
  }

  /**
   * Publish user updated event
   */
  async publishUserUpdated(payload: UserEventPayload): Promise<void> {
    await this.publish('user.updated', {
      ...payload,
      action: 'updated',
    });
  }

  /**
   * Publish user deleted event
   */
  async publishUserDeleted(userId: string): Promise<void> {
    await this.publish('user.deleted', {
      userId,
      action: 'deleted',
      timestamp: new Date().toISOString(),
    });
  }

  /**
   * Setup a consumer for a specific queue and routing key
   * This is typically used by the consumer service (User Service)
   */
  async consume(
    queueName: string,
    routingKeys: string[],
    callback: (message: any) => Promise<void>,
  ): Promise<void> {
    try {
      // Assert queue
      await this.channel.assertQueue(queueName, {
        durable: true,
      });

      // Bind queue to exchange with routing keys
      for (const routingKey of routingKeys) {
        await this.channel.bindQueue(queueName, this.exchangeName, routingKey);
        this.logger.log(
          `Queue "${queueName}" bound to routing key "${routingKey}"`,
        );
      }

      // Start consuming
      await this.channel.consume(
        queueName,
        async (msg) => {
          if (msg) {
            try {
              const content = JSON.parse(msg.content.toString());
              this.logger.log(
                `📥 Received message from "${queueName}": ${JSON.stringify(content)}`,
              );

              await callback(content);

              // Acknowledge message
              this.channel.ack(msg);
            } catch (error) {
              this.logger.error('Error processing message:', error);
              // Reject and requeue the message
              this.channel.nack(msg, false, true);
            }
          }
        },
        {
          noAck: false, // Manual acknowledgment
        },
      );

      this.logger.log(`📨 Started consuming from queue "${queueName}"`);
    } catch (error) {
      this.logger.error(
        `Failed to setup consumer for queue "${queueName}":`,
        error,
      );
      throw error;
    }
  }

  /**
   * Get the channel (for advanced operations)
   */
  getChannel(): amqp.Channel {
    return this.channel;
  }

  /**
   * Check if RabbitMQ is connected
   */
  isConnected(): boolean {
    return !!this.connection && !!this.channel;
  }
}

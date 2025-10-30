package messaging

import (
	"encoding/json"
	"fmt"
	"go-api-streaming/infrastructure/config"
	"log"

	"github.com/streadway/amqp"
)

type RabbitMQClient struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewRabbitMQClient(cfg *config.RabbitMQConfig) (*RabbitMQClient, error) {
	conn, err := amqp.Dial(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	fmt.Println("✓ RabbitMQ connection established")

	return &RabbitMQClient{
		conn:    conn,
		channel: channel,
	}, nil
}

func (r *RabbitMQClient) DeclareQueue(queueName string) error {
	_, err := r.channel.QueueDeclare(
		queueName,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}
	return nil
}

func (r *RabbitMQClient) PublishMessage(queueName string, message interface{}) error {
	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	err = r.channel.Publish(
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	log.Printf("✓ Message published to queue: %s", queueName)
	return nil
}

func (r *RabbitMQClient) Consume(queueName string, handler func([]byte) error) error {
	msgs, err := r.channel.Consume(
		queueName,
		"",    // consumer
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %w", err)
	}

	go func() {
		for msg := range msgs {
			if err := handler(msg.Body); err != nil {
				log.Printf("Error handling message: %v", err)
				msg.Nack(false, true) // requeue on error
			} else {
				msg.Ack(false)
			}
		}
	}()

	log.Printf("✓ Consumer registered for queue: %s", queueName)
	return nil
}

func (r *RabbitMQClient) Close() error {
	if r.channel != nil {
		r.channel.Close()
	}
	if r.conn != nil {
		return r.conn.Close()
	}
	return nil
}

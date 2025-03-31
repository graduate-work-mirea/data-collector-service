package queue

import (
	"context"
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
)

// RabbitMQ represents a RabbitMQ connection
type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   amqp.Queue
}

// NewRabbitMQ creates a new RabbitMQ connection
func NewRabbitMQ(url string) (*RabbitMQ, error) {
	// Connect to RabbitMQ
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	// Create a channel
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	// Declare a queue
	q, err := ch.QueueDeclare(
		"raw_data_queue", // name
		true,             // durable
		false,            // delete when unused
		false,            // exclusive
		false,            // no-wait
		nil,              // arguments
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, err
	}

	return &RabbitMQ{
		conn:    conn,
		channel: ch,
		queue:   q,
	}, nil
}

// Close closes the RabbitMQ connection
func (r *RabbitMQ) Close() error {
	if err := r.channel.Close(); err != nil {
		return err
	}
	return r.conn.Close()
}

// PublishRawData publishes raw data to the queue
func (r *RabbitMQ) PublishRawData(ctx context.Context, data json.RawMessage) error {
	return r.channel.PublishWithContext(
		ctx,
		"",                 // exchange
		r.queue.Name,       // routing key
		false,              // mandatory
		false,              // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        data,
		},
	)
}

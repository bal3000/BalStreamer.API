package helpers

import (
	"fmt"
	"log"

	"github.com/bal3000/BalStreamer.API/configuration"
	"github.com/bal3000/BalStreamer.API/models"
	"github.com/streadway/amqp"
)

// MessageQueue - interface for interacting with message queues
type MessageQueue interface {
	SendMessage(message models.EventMessage) error
}

// RabbitMQConnection - settings to create a connection
type RabbitMQConnection struct {
	configuration *configuration.Configuration
	Channel       *amqp.Channel
}

type rabbitError struct {
	ogErr   error
	message string
}

// NewRabbitMQConnection creates a new rabbit mq connection
func NewRabbitMQConnection(config *configuration.Configuration) RabbitMQConnection {
	conn, err := amqp.Dial(config.RabbitURL)
	failOnError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	failOnError(err, "Failed to create a channel")
	return RabbitMQConnection{configuration: config, Channel: ch}
}

// SendMessage sends the given message
func (mq *RabbitMQConnection) SendMessage(routingKey string, message models.EventMessage) error {
	b, err := message.TransformMessage()
	if err != nil {
		return err
	}

	log.Println("Converted message to JSON and sending")

	return mq.Channel.Publish(
		mq.configuration.ExchangeName, // exchange
		routingKey,                    // routing key
		false,                         // mandatory
		false,                         // immediate
		amqp.Publishing{
			ContentType:  "application/vnd.masstransit+json",
			Body:         []byte(b),
			DeliveryMode: amqp.Persistent,
		})
}

// StartConsumer - starts consuming messages from the given queue
func (mq *RabbitMQConnection) StartConsumer(routingKey string, handler func(d amqp.Delivery) bool, concurrency int) {
	// create the queue if it doesn't already exist
	_, err := mq.Channel.QueueDeclare(mq.configuration.QueueName, true, false, false, false, nil)
	failOnError(err, fmt.Sprintf("Failed to declare a queue: %s", mq.configuration.QueueName))

	// bind the queue to the routing key
	err = mq.Channel.QueueBind(mq.configuration.QueueName, routingKey, mq.configuration.ExchangeName, false, nil)
	failOnError(err, fmt.Sprintf("Failed to bind to queue: %s", mq.configuration.QueueName))

	// prefetch 4x as many messages as we can handle at once
	prefetchCount := concurrency * 4
	err = mq.Channel.Qos(prefetchCount, 0, false)
	failOnError(err, "Failed to setup prefetch")

	msgs, err := mq.Channel.Consume(
		mq.configuration.QueueName, // queue
		"",                         // consumer
		false,                      // auto-ack
		false,                      // exclusive
		false,                      // no-local
		false,                      // no-wait
		nil,                        // args
	)
	failOnError(err, "Failed to get any messages")

	for i := 0; i < concurrency; i++ {
		fmt.Printf("Processing messages on thread %v...\n", i)
		go func() {
			for msg := range msgs {
				// if tha handler returns true then ACK, else NACK
				// the message back into the rabbit queue for
				// another round of processing
				if handler(msg) {
					msg.Ack(false)
				} else {
					msg.Nack(false, true)
				}
			}
			log.Panicln("Rabbit consumer closed - critical Error")
		}()
	}
}

func (err rabbitError) Error() string {

}

func returnErr(err error) error {
	if err != nil {
		return err
	}
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

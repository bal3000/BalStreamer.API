package helpers

import (
	"log"

	"github.com/bal3000/BalStreamer.API/configuration"
	"github.com/bal3000/BalStreamer.API/models"
	"github.com/streadway/amqp"
)

// MessageQueue - interface for interacting with message queues
type MessageQueue interface {
	SendMessage(message models.EventMessage) error
}

// RabbitMQ - settings to create a connection
type RabbitMQ struct {
	configuration *configuration.Configuration
	Connection    *amqp.Connection
	Channel       *amqp.Channel
}

// NewRabbitMQ creates a new rabbit mq connection
func NewRabbitMQ(config *configuration.Configuration) *RabbitMQ {
	conn, err := amqp.Dial(config.RabbitURL)
	failOnError(err, "Failed to connect to RabbitMQ")

	return &RabbitMQ{configuration: config, Connection: conn}
}

// CreateChannel creates a new channel
func (mq *RabbitMQ) CreateChannel() {
	ch, err := mq.Connection.Channel()
	failOnError(err, "Failed to bind a queue")
	mq.Channel = ch
}

// CreateExchange creates an exchange
func (mq *RabbitMQ) CreateExchange() {
	err := mq.Channel.ExchangeDeclare(
		mq.configuration.ExchangeName, // name
		"fanout",                      // type
		mq.configuration.Durable,      // durable
		false,                         // auto-deleted
		false,                         // internal
		false,                         // no-wait
		nil,                           // arguments
	)
	failOnError(err, "Failed to declare an exchange")
}

// DeclareAndBindQueue declares a queue if one does not exist and then binds it to the channel
func (mq *RabbitMQ) DeclareAndBindQueue() {
	q, err := mq.Channel.QueueDeclare(
		mq.configuration.QueueName, // name
		mq.configuration.Durable,   // durable
		false,                      // delete when unused
		false,                      // exclusive
		false,                      // no-wait
		nil,                        // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = mq.Channel.QueueBind(
		q.Name,                        // queue name
		"",                            // routing key
		mq.configuration.ExchangeName, // exchange
		false,
		nil,
	)
	failOnError(err, "Failed to bind a queue")
}

// SendMessage sends the given message
func (mq *RabbitMQ) SendMessage(message models.EventMessage) error {
	b, err := message.TransformMessage()
	if err != nil {
		return err
	}

	log.Println("Converted message to JSON and sending")

	err = mq.Channel.Publish(
		mq.configuration.ExchangeName, // exchange
		"",                            // routing key
		false,                         // mandatory
		false,                         // immediate
		amqp.Publishing{
			ContentType: "application/vnd.masstransit+json",
			Body:        []byte(b),
		})

	if err != nil {
		return err
	}

	log.Println("Message sent")
	return nil
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

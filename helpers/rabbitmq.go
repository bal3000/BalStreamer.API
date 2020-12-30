package helpers

import (
	"log"

	"github.com/bal3000/BalStreamer.API/models"
	"github.com/streadway/amqp"
)

// RabbitMQ - settings to create a connection
type RabbitMQ struct {
	URL          string
	QueueName    string
	ExchangeName string
	Durable      bool
}

// ConnectToRabbitMQ - Connects to RabbitMQ
func (mq *RabbitMQ) ConnectToRabbitMQ() *amqp.Connection {
	conn, err := amqp.Dial(mq.URL)
	failOnError(err, "Failed to connect to RabbitMQ")
	return conn
}

// CreateExchange creates an exchange
func (mq *RabbitMQ) CreateExchange(ch *amqp.Channel) {
	err := ch.ExchangeDeclare(
		mq.ExchangeName, // name
		"fanout",        // type
		mq.Durable,      // durable
		false,           // auto-deleted
		false,           // internal
		false,           // no-wait
		nil,             // arguments
	)
	failOnError(err, "Failed to declare an exchange")
}

// SendMessage sends the given message
func SendMessage(ch *amqp.Channel, message models.EventMessage, exchangeName string) {
	b, err := message.TransformMessage()
	if err != nil {
		log.Fatalln(err)
	}

	err = ch.Publish(
		exchangeName, // exchange
		"",           // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "application/vnd.masstransit+json",
			Body:        []byte(b),
		})

	if err != nil {
		log.Fatalln(err)
	}
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

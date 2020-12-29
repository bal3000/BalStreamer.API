package helpers

import (
	"log"

	"github.com/streadway/amqp"
)

// RabbitSettings - settings to create a connection
type RabbitSettings struct {
	URL       string
	QueueName string
	Durable   bool
}

// ConnectToRabbitMQ - Connects to RabbitMQ
func (settings *RabbitSettings) ConnectToRabbitMQ() *amqp.Connection {
	log.Println("Connecting to RabbitMQ")
	conn, err := amqp.Dial(settings.URL)
	if err != nil {
		panic(err)
	}
	log.Println("Connected to RabbitMQ")
	return conn
}

//SendMessage ensures the queue exists and then sends the given message
func (settings *RabbitSettings) SendMessage(ch *amqp.Channel, message string) {
	q := settings.createQueue(ch)
	err := ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})

	if err != nil {
		panic(err)
	}
}

//CreateQueue creates a queue and assigns the queue name to the settings
func (settings *RabbitSettings) createQueue(ch *amqp.Channel) amqp.Queue {
	q, err := ch.QueueDeclare(
		settings.QueueName, // name
		settings.Durable,   // durable
		false,              // delete when unused
		false,              // exclusive
		false,              // no-wait
		nil,                // arguments
	)

	if err != nil {
		panic(err)
	}

	log.Println("Created a queue called: " + settings.QueueName)
	return q
}

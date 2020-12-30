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

// DeclareAndBindQueue declares a queue if one does not exist and then binds it to the channel
func (mq *RabbitMQ) DeclareAndBindQueue(ch *amqp.Channel) {
	q, err := ch.QueueDeclare(
		mq.QueueName, // name
		mq.Durable,   // durable
		false,        // delete when unused
		true,         // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.QueueBind(
		q.Name,          // queue name
		"",              // routing key
		mq.ExchangeName, // exchange
		false,
		nil,
	)
	failOnError(err, "Failed to bind a queue")
}

// SendMessage sends the given message
func SendMessage(ch *amqp.Channel, message models.EventMessage, exchangeName string) {
	log.Printf("Sending to exchange %s in rabbitMQ", exchangeName)
	b, err := message.TransformMessage()
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Converted message to JSON and sending")

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

	log.Println("Message sent")
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

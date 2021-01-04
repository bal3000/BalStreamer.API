package controllers

import (
	"database/sql"
	"encoding/json"
	"log"
	"time"

	"github.com/bal3000/BalStreamer.API/models"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/streadway/amqp"
)

var (
	upgrader = websocket.Upgrader{}
)

// ChromecastController the controller for the websockets
type ChromecastController struct {
	Database  *sql.DB
	RabbitMQ  *amqp.Channel
	QueueName string
}

// NewChromecastController creates a new ref to chromecast controller
func NewChromecastController(db *sql.DB, ch *amqp.Channel, qn string) *ChromecastController {
	return &ChromecastController{Database: db, RabbitMQ: ch, QueueName: qn}
}

// ChromecastUpdates broadcasts a chromecast to all clients once found
func (controller *ChromecastController) ChromecastUpdates(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	// Prepare insert statement
	insert, insertErr := controller.Database.Prepare("INSERT INTO casterDB VALUES ($1,$2)")
	if insertErr != nil {
		panic(insertErr)
	}
	defer insert.Close()

	msgs, err := controller.RabbitMQ.Consume(
		controller.QueueName, // queue
		"",                   // consumer
		true,                 // auto-ack
		false,                // exclusive
		false,                // no-local
		false,                // no-wait
		nil,                  // args
	)
	if err != nil {
		return err
	}

	forever := make(chan bool)

	go waitAndProcessMsg(msgs, insert, ws, c)

	<-forever

	return nil
}

func waitAndProcessMsg(msgs <-chan amqp.Delivery, insert *sql.Stmt, ws *websocket.Conn, c echo.Context) {
	chromecastEvent := new(models.ChromecastFoundEvent)
	for d := range msgs {
		// need this to test how I can tell the difference between add and remove events
		log.Println("Rabbit message: ", msgs)
		json.Unmarshal(d.Body, chromecastEvent)

		// insert into db - might make this a proc to determine if it already exists and if so update time
		insert.Exec(chromecastEvent.Chromecast, time.Now())

		err := ws.WriteJSON(chromecastEvent)
		if err != nil {
			c.Logger().Error(err)
		}
	}
}

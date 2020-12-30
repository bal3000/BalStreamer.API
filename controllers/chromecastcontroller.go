package controllers

import (
	"database/sql"
	"encoding/json"
	"log"

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
func NewChromecastController(db *sql.DB, ch *amqp.Channel) *ChromecastController {
	return &ChromecastController{Database: db, RabbitMQ: ch}
}

// ChromecastUpdates broadcasts a chromecast to all clients once found
func (controller *ChromecastController) ChromecastUpdates(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

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

	go func() {
		chromecastEvent := new(models.ChromecastFoundEvent)
		for d := range msgs {
			log.Println("Rabbit message: ", msgs)
			json.Unmarshal(d.Body, chromecastEvent)
			err := ws.WriteJSON(chromecastEvent)
			if err != nil {
				c.Logger().Error(err)
			}
		}
	}()

	<-forever

	return nil
}

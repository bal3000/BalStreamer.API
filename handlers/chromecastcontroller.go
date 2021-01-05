package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/bal3000/BalStreamer.API/models"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/streadway/amqp"
)

var (
	upgrader       = websocket.Upgrader{}
	foundEventType = "ChromecastFoundEvent"
	lostEventType  = "ChromecastLostEvent"
)

// ChromecastController the controller for the websockets
type ChromecastHandler struct {
	Database  *sql.DB
	RabbitMQ  *amqp.Channel
	QueueName string
}

// NewChromecastController creates a new ref to chromecast controller
func NewChromecastHandler(db *sql.DB, ch *amqp.Channel, qn string) *ChromecastHandler {
	return &ChromecastHandler{Database: db, RabbitMQ: ch, QueueName: qn}
}

// ChromecastUpdates broadcasts a chromecast to all clients once found
func (controller *ChromecastHandler) ChromecastUpdates(c echo.Context) error {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	// Prepare insert statement
	insert, insertErr := controller.Database.Prepare("INSERT INTO public.\"Chromecasts\" VALUES ($1,$2)")
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
	mtEvent := new(models.MassTransitEvent)
	for d := range msgs {
		// need this to test how I can tell the difference between add and remove events
		log.Println("Rabbit message: ", string(d.Body))

		// convert mass transit message
		chromecastEvent, err := mtEvent.RetrieveMessage(d.Body)
		if err != nil {
			c.Logger().Error(err)
		}

		// insert into db - might make this a proc to determine if it already exists and if so update time
		switch chromecastEvent.EventType {
		case foundEventType:
			insert.Exec(chromecastEvent.Chromecast, time.Now())
		case lostEventType:
			insert.Exec(chromecastEvent.Chromecast, time.Now())
		}

		err = ws.WriteJSON(chromecastEvent)
		if err != nil {
			c.Logger().Error(err)
		}
	}
}

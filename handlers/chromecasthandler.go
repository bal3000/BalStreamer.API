package handlers

import (
	"log"
	"net/http"

	"github.com/bal3000/BalStreamer.API/helpers"
	"github.com/bal3000/BalStreamer.API/models"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/streadway/amqp"
)

var (
	upgrader       = websocket.Upgrader{}
	foundEventType = "ChromecastFoundEvent"
	lostEventType  = "ChromecastLostEvent"
	ws             *websocket.Conn
	chromecasts    = []models.ChromecastEvent{}
)

// ChromecastHandler the controller for the websockets
type ChromecastHandler struct {
	RabbitMQ  *helpers.RabbitMQConnection
	QueueName string
}

// NewChromecastHandler creates a new ref to chromecast controller
func NewChromecastHandler(rabbit *helpers.RabbitMQConnection, qn string) *ChromecastHandler {
	return &ChromecastHandler{RabbitMQ: rabbit, QueueName: qn}
}

// ChromecastUpdates broadcasts a chromecast to all clients once found
func (handler *ChromecastHandler) ChromecastUpdates(c echo.Context) error {
	log.Println("Entered ws")

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	err = handler.RabbitMQ.StartConsumer("chromecast-key", processMsgs, 2)
	if err != nil {
		panic(err)
	}

	forever := make(chan bool)
	<-forever

	return nil
}

func processMsgs(d amqp.Delivery) bool {
	mtEvent := new(models.MassTransitEvent)

	// convert mass transit message
	chromecastEvent, err := mtEvent.RetrieveMessage(d.Body)
	if err != nil {
		log.Println(err)
		return false
	}

	// insert into db - might make this a proc to determine if it already exists and if so update time
	switch chromecastEvent.EventType {
	case foundEventType:
		if !contains(chromecasts, chromecastEvent) {
			chromecasts = append(chromecasts, chromecastEvent)
		}
	case lostEventType:
		if i := find(chromecasts, chromecastEvent); i > -1 {
			chromecasts = remove(chromecasts, i)
		}
	}

	err = ws.WriteJSON(chromecastEvent)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

func contains(a []models.ChromecastEvent, x models.ChromecastEvent) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

func find(a []models.ChromecastEvent, x models.ChromecastEvent) int {
	for i, n := range a {
		if x == n {
			return i
		}
	}
	return len(a)
}

func remove(s []models.ChromecastEvent, i int) []models.ChromecastEvent {
	s[i] = s[len(s)-1]
	// We do not need to put s[i] at the end, as it will be discarded anyway
	return s[:len(s)-1]
}

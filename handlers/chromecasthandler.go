package handlers

import (
	"fmt"
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
	chromecasts    = make(map[string]models.ChromecastEvent)
	handledMsgs    = make(chan models.ChromecastEvent)
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
	log.Println("Entered ws, sending current found chromecasts")

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	// send all chromecasts from last refresh to page
	for _, event := range chromecasts {
		err = ws.WriteJSON(event)
		if err != nil {
			return err
		}
	}

	err = handler.RabbitMQ.StartConsumer("chromecast-key", processMsgs, 2)
	if err != nil {
		panic(err)
	}

	for msg := range handledMsgs {
		err = ws.WriteJSON(msg)
		if err != nil {
			log.Println(err)
			return err
		}
	}
	close(handledMsgs)

	return err
}

func processMsgs(d amqp.Delivery) bool {
	fmt.Printf("processing message: %s", string(d.Body))
	mtEvent := new(models.MassTransitEvent)

	// convert mass transit message
	chromecastEvent, err := mtEvent.RetrieveMessage(d.Body)
	if err != nil {
		log.Println(err)
		return false
	}

	if name, ok := chromecastEvent.Chromecast.(string); ok {
		switch chromecastEvent.EventType {
		case foundEventType:
			chromecasts[name] = chromecastEvent
		case lostEventType:
			delete(chromecasts, name)
		}
	}

	handledMsgs <- chromecastEvent

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
	return s[:len(s)-1]
}

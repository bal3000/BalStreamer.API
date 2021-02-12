package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/bal3000/BalStreamer.API/infrastructure"
	"log"
	"net/http"

	"github.com/bal3000/BalStreamer.API/models"
	"github.com/gorilla/websocket"
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
	RabbitMQ  infrastructure.RabbitMQ
	QueueName string
}

// NewChromecastHandler creates a new ref to chromecast controller
func NewChromecastHandler(rabbit infrastructure.RabbitMQ, qn string) *ChromecastHandler {
	return &ChromecastHandler{RabbitMQ: rabbit, QueueName: qn}
}

// ChromecastUpdates broadcasts a chromecast to all clients once found
func (handler *ChromecastHandler) ChromecastUpdates(res http.ResponseWriter, req *http.Request) {
	log.Println("Entered ws, sending current found chromecasts")

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(res, req, nil)
	if err != nil {
		log.Fatalln(err)
	}
	defer ws.Close()

	// send all chromecasts from last refresh to page
	for _, event := range chromecasts {
		err = ws.WriteJSON(event)
		if err != nil {
			log.Fatalln(err)
		}
	}

	err = handler.RabbitMQ.StartConsumer("chromecast-key", processMsgs, 2)
	if err != nil {
		panic(err)
	}

	for msg := range handledMsgs {
		err = ws.WriteJSON(msg)
		if err != nil {
			log.Fatalln(err)
		}
	}
	close(handledMsgs)
}

func processMsgs(d amqp.Delivery) bool {
	fmt.Printf("processing message: %s, with type: %s", string(d.Body), d.Type)
	event := new(models.ChromecastEvent)

	// convert mass transit message
	err := json.Unmarshal(d.Body, event)
	if err != nil {
		log.Println(err)
		return false
	}

	switch d.Type {
	case foundEventType:
		chromecasts[event.Chromecast] = *event
	case lostEventType:
		delete(chromecasts, event.Chromecast)
	}

	handledMsgs <- *event

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

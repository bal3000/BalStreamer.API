package controllers

import (
	"database/sql"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/streadway/amqp"
)

var (
	upgrader = websocket.Upgrader{}
)

// ChromecastController the controller for the websockets
type ChromecastController struct {
	Database *sql.DB
	RabbitMQ *amqp.Channel
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

	for {
		// Write
		err := ws.WriteJSON([]byte("Hello, Client!"))
		if err != nil {
			c.Logger().Error(err)
		}
	}
}

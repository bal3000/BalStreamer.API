package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/bal3000/BalStreamer.API/helpers"
	"github.com/bal3000/BalStreamer.API/models"
	"github.com/labstack/echo/v4"
)

// CastHandler - controller for casting to chromecast
type CastHandler struct {
	RabbitMQ     *helpers.RabbtMQ
	ExchangeName string
}

// NewCastHandler - constructor to return new controller while passing in dependacies
func NewCastHandler(rabbit *helpers.RabbitMQ, en string) *CastHndler {
	return &CastHandler{RabbitMQ: rabbit, ExchangeName: en}
}

// CastStream - streams given data to given chromecast
func (controller *CastHandler) CastStream(c echo.Context) error {
	castCommand := new(models.StreamToCast)

	if err := c.Bind(castCommand); err != nil {
		log.Println(err)
		return err
	}

	// Send to chromecast
	cast := &models.StreamToChromecastEvent{
		ChromeCastToStream: castCommand.Chromecast,
		Stream:             castCommand.StreamURL,
		StreamDate:         time.Now(),
	}

	go sendMessage(controller.RabbitMQ, cast)

	return c.NoContent(http.StatusNoContent)
}

// StopStream endpoint sends the command to stop the stream on the given chromecast
func (controller *CastHandler) StopStream(c echo.Context) error {
	stopStreamCommand := new(models.StopPlayingStream)

	if err := c.Bind(stopStreamCommand); err != nil {
		log.Println(err)
		return err
	}

	// Send to chromecast
	cast := &models.StopPlayingStreamEvent{
		ChromeCastToStop: stopStreamCommand.ChromeCastToStop,
		StopDateTime:     stopStreamCommand.StopDateTime,
	}

	go sendMessage(controller.RabbitMQ, cast)

	return c.NoContent(http.StatusAccepted)
}

func sendMessage(queue helpers.MessageQueue, event models.EventMessage) {
	queue.SendMessage(event)
}

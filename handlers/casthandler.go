package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/bal3000/BalStreamer.API/helpers"
	"github.com/bal3000/BalStreamer.API/models"
	"github.com/labstack/echo/v4"
	"github.com/streadway/amqp"
)

// CastHandler - controller for casting to chromecast
type CastHandler struct {
	RabbitMQ     *amqp.Channel
	ExchangeName string
}

// NewCastHandler - constructor to return new controller while passing in dependacies
func NewCastHandler(ch *amqp.Channel, en string) *CastHandler {
	return &CastHandler{RabbitMQ: ch, ExchangeName: en}
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

	go helpers.SendMessage(controller.RabbitMQ, cast, controller.ExchangeName)

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

	go helpers.SendMessage(controller.RabbitMQ, cast, controller.ExchangeName)

	return c.NoContent(http.StatusAccepted)
}

package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/bal3000/BalStreamer.API/messaging"
	"github.com/bal3000/BalStreamer.API/models"
	"github.com/labstack/echo/v4"
)

// CastHandler - controller for casting to chromecast
type CastHandler struct {
	RabbitMQ     messaging.RabbitMQ
	ExchangeName string
}

// NewCastHandler - constructor to return new controller while passing in dependacies
func NewCastHandler(rabbit messaging.RabbitMQ, en string) *CastHandler {
	return &CastHandler{RabbitMQ: rabbit, ExchangeName: en}
}

// CastStream - streams given data to given chromecast
func (handler *CastHandler) CastStream(c echo.Context) error {
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

	handler.RabbitMQ.SendMessage("chromecast-key", cast)

	return c.NoContent(http.StatusNoContent)
}

// StopStream endpoint sends the command to stop the stream on the given chromecast
func (handler *CastHandler) StopStream(c echo.Context) error {
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

	go handler.RabbitMQ.SendMessage("chromecast-key", cast)

	return c.NoContent(http.StatusAccepted)
}

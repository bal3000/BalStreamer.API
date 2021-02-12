package handlers

import (
	"encoding/json"
	"github.com/bal3000/BalStreamer.API/infrastructure"
	"log"
	"net/http"
	"time"

	"github.com/bal3000/BalStreamer.API/models"
)

const routingKey string = "chromecast-key"

// CastHandler - controller for casting to chromecast
type CastHandler struct {
	RabbitMQ     infrastructure.RabbitMQ
	ExchangeName string
}

// NewCastHandler - constructor to return new controller while passing in dependencies
func NewCastHandler(rabbit infrastructure.RabbitMQ, en string) *CastHandler {
	return &CastHandler{RabbitMQ: rabbit, ExchangeName: en}
}

// CastStream - streams given data to given chromecast
func (handler *CastHandler) CastStream(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Access-Control-Allow-Origin", "*")
	res.Header().Set("content-type", "application/json")
	castCommand := new(models.StreamToCast)

	if err := json.NewDecoder(req.Body).Decode(castCommand); err != nil {
		log.Println(err)
	}

	// Send to chromecast
	cast := &models.StreamToChromecastEvent{
		ChromeCastToStream: castCommand.Chromecast,
		Stream:             castCommand.StreamURL,
		StreamDate:         time.Now(),
	}

	go handler.RabbitMQ.SendMessage(routingKey, cast)

	res.WriteHeader(http.StatusNoContent)
}

// StopStream endpoint sends the command to stop the stream on the given chromecast
func (handler *CastHandler) StopStream(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Access-Control-Allow-Origin", "*")
	res.Header().Set("content-type", "application/json")
	stopStreamCommand := new(models.StopPlayingStream)

	if err := json.NewDecoder(req.Body).Decode(stopStreamCommand); err != nil {
		log.Println(err)
	}

	// Send to chromecast
	cast := &models.StopPlayingStreamEvent{
		ChromeCastToStop: stopStreamCommand.ChromeCastToStop,
		StopDateTime:     stopStreamCommand.StopDateTime,
	}

	go handler.RabbitMQ.SendMessage(routingKey, cast)

	res.WriteHeader(http.StatusAccepted)
}

package services

import (
	"time"

	"github.com/bal3000/BalStreamer.API/helpers"
	"github.com/bal3000/BalStreamer.API/models"
	"github.com/streadway/amqp"
)

// CastService repersents this service
type CastService struct {
	RabbitMQ *amqp.Channel
}

// CastStream sends the stream message to rabbitmq
func (service *CastService) CastStream(castCommand *models.StreamToCast, exchangeName string) error {
	// Send to chromecast
	cast := &models.StreamToChromecastEvent{
		ChromeCastToStream: castCommand.Chromecast,
		Stream:             castCommand.StreamURL,
		StreamDate:         time.Now(),
	}

	return helpers.SendMessage(service.RabbitMQ, cast, exchangeName)
}

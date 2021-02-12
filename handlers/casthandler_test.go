package handlers

import (
	"github.com/bal3000/BalStreamer.API/infrastructure"
	"github.com/bal3000/BalStreamer.API/models"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/mock"
	"testing"
)

var (
	config   = &infrastructure.Configuration{RabbitURL: "amqp://guest:guest@localhost:5672/", ExchangeName: "bal-streamer-caster", QueueName: "bal-streamer-api", Durable: true}
	castJSON = `{
		"chromecast": "Family room TV",
		"streamURL": "rtmp://cdn.vops.gcp.xeatre.cloud:5222/liveedge-lowlatency-origin-wza-07/src-4506?wUzz3Tsnestarttime=1609777218&wUzz3Tsneendtime=1609781100&wUzz3Tsnehash=PN0KNFTOB-fyV9qdN2wFj5fZ0r74DtGfSdcJNwsh5Oc="
	}`
)

type RabbitChannelMock struct {
	mock.Mock
}

func (m *RabbitChannelMock) SendMessage(routingKey string, message models.EventMessage) error {
	args := m.Called(routingKey, message)
	return args.Error(0)
}

func (m *RabbitChannelMock) StartConsumer(routingKey string, handler func(d amqp.Delivery) bool, concurrency int) error {
	args := m.Called(routingKey, handler, concurrency)
	return args.Error(0)
}

func (m *RabbitChannelMock) CloseChannel() {
}

func TestCastStream(t *testing.T) {
	// Setup
	//rabbitMock := new(RabbitChannelMock)
	//
	//cast := &models.StreamToChromecastEvent{
	//	ChromeCastToStream: "Family room TV",
	//	Stream:             "rtmp://cdn.vops.gcp.xeatre.cloud:5222/liveedge-lowlatency-origin-wza-07/src-4506?wUzz3Tsnestarttime=1609777218&wUzz3Tsneendtime=1609781100&wUzz3Tsnehash=PN0KNFTOB-fyV9qdN2wFj5fZ0r74DtGfSdcJNwsh5Oc=",
	//	StreamDate:         time.Now(),
	//}
	//rabbitMock.On("SendMessage", "chromecast-key", cast).Return(nil)
	//
	//req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(castJSON))
	//req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	//rec := httptest.NewRecorder()
	//castHandle := &CastHandler{RabbitMQ: rabbitMock, ExchangeName: config.ExchangeName}
	//
	//// Assertions
	//if assert.NoError(t, castHandle.CastStream) {
	//	assert.Equal(t, http.StatusNoContent, rec.Code)
	//	assert.Equal(t, "", rec.Body.String())
	//	rabbitMock.AssertExpectations(t)
	//}
}

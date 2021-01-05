package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bal3000/BalStreamer.API/configuration"
	"github.com/bal3000/BalStreamer.API/helpers"
	"github.com/labstack/echo/v4"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	config   = &configuration.Configuration{RabbitURL: "amqp://guest:guest@localhost:5672/", ExchangeName: "bal-streamer-caster", QueueName: "bal-streamer-api", Durable: true}
	castJSON = `{
		"chromecast": "Family room TV",
		"streamURL": "rtmp://cdn.vops.gcp.xeatre.cloud:5222/liveedge-lowlatency-origin-wza-07/src-4506?wUzz3Tsnestarttime=1609777218&wUzz3Tsneendtime=1609781100&wUzz3Tsnehash=PN0KNFTOB-fyV9qdN2wFj5fZ0r74DtGfSdcJNwsh5Oc="
	}`
)

type RabbitChannelMock struct {
	mock.Mock
}

func (m *RabbitChannelMock) Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
	args := m.Mock.Called(exchange, key, mandatory, immediate, msg)
	return args.Error(1)
}

func TestCastStream(t *testing.T) {
	// Setup
	rabbitMock := helpers.NewRabbitMQ(config)
	defer rabbitMock.Connection.Close()

	ch := rabbitMock.CreateChannel()
	defer ch.Close()

	rabbitMock.CreateExchange(ch)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(castJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	castHandle := &CastHandler{RabbitMQ: ch, ExchangeName: config.ExchangeName}

	// Assertions
	if assert.NoError(t, castHandle.CastStream(c)) {
		assert.Equal(t, http.StatusNoContent, rec.Code)
		assert.Equal(t, "", rec.Body.String())
	}
}
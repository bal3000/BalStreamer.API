package app

import (
	"github.com/bal3000/BalStreamer.API/handlers"
	"github.com/labstack/echo/v4"
)

// SetRoutes creates the handlers and routes for those handlers
func (s *Server) SetRoutes() {
	s.Echo.File("/", "public/index.html")
	// Handlers
	cast := handlers.NewCastHandler(s.RabbitMQ, s.Config.ExchangeName)
	chrome := handlers.NewChromecastHandler(s.RabbitMQ, s.Config.QueueName)
	live := handlers.NewLiveStreamHandler(s.Config.LiveStreamURL, s.Config.APIKey)

	CastRoutes(s.Echo, cast)
	ChromecastRoutes(s.Echo, chrome)
	LiveStreamRoutes(s.Echo, live)
}

// CastRoutes sets up the routes for the cast handler
func CastRoutes(e *echo.Echo, cast *handlers.CastHandler) {
	e.POST("/api/cast", cast.CastStream)
	e.DELETE("/api/cast", cast.StopStream)
}

// ChromecastRoutes sets up the routes for the chromecast handler
func ChromecastRoutes(e *echo.Echo, chrome *handlers.ChromecastHandler) {
	e.GET("/chromecasts", chrome.ChromecastUpdates)
}

// LiveStreamRoutes sets up the routes for the live streams handler
func LiveStreamRoutes(e *echo.Echo, live *handlers.LiveStreamHandler) {
	e.GET("/api/livestreams/:sportType/:fromDate/:toDate", live.GetFixtures)
	e.GET("/api/livestreams/:timerId", live.GetStreams)
}

package routes

import (
	"github.com/bal3000/BalStreamer.API/handlers"
	"github.com/labstack/echo/v4"
)

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

package routes

import (
	"github.com/bal3000/BalStreamer.API/handlers"
	"github.com/labstack/echo/v4"
)

// CastRoutes sets up the routes for the cast controller
func CastRoutes(e *echo.Echo, cast *handlers.CastHandler) {
	e.POST("/api/cast", cast.CastStream)
	e.DELETE("/api/cast", cast.StopStream)
}

// ChromecastRoutes sets up the routes for the chromecast controller
func ChromecastRoutes(e *echo.Echo, chrome *handlers.ChromecastHandler) {
	e.GET("/chromecasts", chrome.ChromecastUpdates)
}

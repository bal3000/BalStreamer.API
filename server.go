package main

import (
	"github.com/bal3000/BalStreamer.API/configuration"
	"github.com/bal3000/BalStreamer.API/helpers"
	"github.com/bal3000/BalStreamer.API/routes"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var config configuration.Configuration

func init() {
	config = configuration.ReadConfig()
}

func main() {
	//setup rabbit
	rabbit := helpers.NewRabbitMQConnection(&config)
	defer rabbit.Channel.Close()

	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.DefaultCORSConfig))

	// Routes
	e.File("/", "public/index.html")
	routes.SetRoutes(e, config, &rabbit)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}

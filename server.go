package main

import (
	"context"
	"os"
	"os/signal"
	"time"

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
	go func() {
		if err := e.Start(":8080"); err != nil {
			e.Logger.Info("Shutting down server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt)
	<-shutdown

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}

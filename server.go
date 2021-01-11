package main

import (
	"database/sql"
	"log"

	"github.com/bal3000/BalStreamer.API/configuration"
	"github.com/bal3000/BalStreamer.API/handlers"
	"github.com/bal3000/BalStreamer.API/helpers"
	"github.com/bal3000/BalStreamer.API/routes"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
)

var config configuration.Configuration

func init() {
	config = configuration.ReadConfig()
}

func main() {
	log.Println("Connecting to DB")
	db, err := sql.Open("postgres", config.ConnectionString)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	log.Println("Connected to DB")

	//setup rabbit
	rabbit := helpers.NewRabbitMQ(&config)
	defer rabbit.Connection.Close()

	rabbit.CreateChannel()
	defer rabbit.Channel.Close()

	rabbit.CreateExchange()
	rabbit.DeclareAndBindQueue()

	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Handlers
	cast := handlers.NewCastHandler(rabbit, config.ExchangeName)
	chrome := handlers.NewChromecastHandler(db, rabbit, config.QueueName)
	live := handlers.NewLiveStreamHandler(config.LiveStreamURL, config.APIKey)

	// Routes
	routes.CastRoutes(e, cast)
	routes.ChromecastRoutes(e, chrome)
	routes.LiveStreamRoutes(e, live)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}

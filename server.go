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

	ch := rabbit.CreateChannel()
	defer ch.Close()

	rabbit.CreateExchange(ch)
	rabbit.DeclareAndBindQueue(ch)

	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	cast := handlers.NewCastHandler(ch, config.ExchangeName)
	chrome := handlers.NewChromecastHandler(db, ch, config.QueueName)

	// Routes
	routes.CastRoutes(e, cast)
	routes.ChromecastRoutes(e, chrome)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}

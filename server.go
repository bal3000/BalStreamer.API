package main

import (
	"database/sql"
	"log"

	"github.com/bal3000/BalStreamer.API/configuration"
	"github.com/bal3000/BalStreamer.API/controllers"
	"github.com/bal3000/BalStreamer.API/helpers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	"github.com/streadway/amqp"
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
	rabbit := helpers.RabbitMQ{
		URL:          "amqp://guest:guest@localhost:5672/",
		QueueName:    "caster-q",
		ExchangeName: "bal-streamer-caster",
		Durable:      true,
	}

	conn := rabbit.ConnectToRabbitMQ()
	defer conn.Close()

	ch := createChannel(conn)
	defer ch.Close()

	rabbit.CreateExchange(ch)

	cast := controllers.NewCastController(db, ch, rabbit.ExchangeName)
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	CastRoutes(e, cast)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}

func createChannel(conn *amqp.Connection) *amqp.Channel {
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalln(err)
	}
	return ch
}

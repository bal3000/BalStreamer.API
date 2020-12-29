package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/bal3000/BalStreamer.API/configuration"
	"github.com/bal3000/BalStreamer.API/controllers"
	"github.com/bal3000/BalStreamer.API/helpers"
	"github.com/gorilla/mux"
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
	router := mux.NewRouter()

	// defer func() {
	// 	// incase of error close all db connections
	// 	if r := recover(); r != nil {
	// 		cast.Database.Close()
	// 	}
	// }()

	// configure the router to always run this handler when it couldn't match a request to any other handler
	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(fmt.Sprintf("%s not found\n", r.URL)))
	})

	// create a subrouter just for standard API calls. subrouters are convenient ways to
	// group similar functionality together. this subrouter also verifies that the Content-Type
	// header is correct for a JSON API.
	castRouter := router.Headers("Content-Type", "application/json").Subrouter()
	castRouter.HandleFunc("/api/cast", cast.CastStream).Methods("POST")

	log.Fatalln(http.ListenAndServe(":8080", router))
}

func createChannel(conn *amqp.Connection) *amqp.Channel {
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalln(err)
	}
	return ch
}

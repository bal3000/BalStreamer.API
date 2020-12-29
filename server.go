package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/bal3000/BalStreamer.API/configuration"
	"github.com/bal3000/BalStreamer.API/controllers"
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

	conn := connectToRabbitMQ()
	defer conn.Close()
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	cast := controllers.NewCastController(db, ch)
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

func connectToRabbitMQ() *amqp.Connection {
	log.Println("Connecting to RabbitMQ")
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(err)
	}
	log.Println("Connected to RabbitMQ")
	return conn
}

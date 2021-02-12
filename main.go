package main

import (
	"fmt"
	"github.com/bal3000/BalStreamer.API/app"
	"github.com/bal3000/BalStreamer.API/infrastructure"
	"github.com/gorilla/mux"
	"os"
)

var config infrastructure.Configuration

func init() {
	config = infrastructure.ReadConfig()
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "startup error: %s\n", err)
		os.Exit(1)
	}
}

func run() error {
	//setup rabbit
	rabbit, err := infrastructure.NewRabbitMQConnection(&config)
	if err != nil {
		return err
	}
	defer rabbit.CloseChannel()

	r := mux.NewRouter()

	server := app.NewServer(rabbit, r, config)
	return server.Run()
}

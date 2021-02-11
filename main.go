package main

import (
	"fmt"
	"github.com/bal3000/BalStreamer.API/app"
	"github.com/bal3000/BalStreamer.API/infrastructure"
	"github.com/labstack/echo/v4"
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

	e := echo.New()

	server := app.NewServer(rabbit, e, config)
	return server.Run()
}

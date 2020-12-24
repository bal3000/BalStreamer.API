package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/bal3000/BalStreamer.API/configuration"
	"github.com/bal3000/BalStreamer.API/controllers"
	"github.com/gorilla/mux"
)

var config configuration.Configuration

func init() {
	config = configuration.ReadConfig()
}

func main() {
	cast := controllers.NewCastController(config.ConnectionString)
	router := mux.NewRouter()

	defer func() {
		// incase of error close all db connections
		if r := recover(); r != nil {
			cast.Database.Close()
		}
	}()

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

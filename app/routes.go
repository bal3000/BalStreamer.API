package app

import (
	"github.com/bal3000/BalStreamer.API/handlers"
	"github.com/gorilla/mux"
	"net/http"
)

// SetRoutes creates the handlers and routes for those handlers
func (s *Server) SetRoutes() {
	// Handlers
	cast := handlers.NewCastHandler(s.RabbitMQ, s.Config.ExchangeName)
	chrome := handlers.NewChromecastHandler(s.RabbitMQ, s.Config.QueueName)
	live := handlers.NewLiveStreamHandler(s.Config.LiveStreamURL, s.Config.APIKey)

	CastRoutes(s.Router, cast)
	ChromecastRoutes(s.Router, chrome)
	LiveStreamRoutes(s.Router, live)

	s.Router.PathPrefix("/").Handler(http.FileServer(http.Dir("public/")))
}

// CastRoutes sets up the routes for the cast handler
func CastRoutes(r *mux.Router, cast *handlers.CastHandler) {
	s := r.PathPrefix("/api/cast").Subrouter()
	s.HandleFunc("/", cast.CastStream).Methods(http.MethodPost)
	s.HandleFunc("/", cast.StopStream).Methods(http.MethodDelete)
}

// ChromecastRoutes sets up the routes for the chromecast handler
func ChromecastRoutes(r *mux.Router, chrome *handlers.ChromecastHandler) {
	r.HandleFunc("/chromecasts", chrome.ChromecastUpdates).Methods(http.MethodGet)
}

// LiveStreamRoutes sets up the routes for the live streams handler
func LiveStreamRoutes(r *mux.Router, live *handlers.LiveStreamHandler) {
	s := r.PathPrefix("/api/livestreams").Subrouter()
	s.HandleFunc("/{sportType}/{fromDate}/{toDate}", live.GetFixtures).Methods(http.MethodGet)
	s.HandleFunc("/{timerId}", live.GetStreams).Methods(http.MethodGet)
}

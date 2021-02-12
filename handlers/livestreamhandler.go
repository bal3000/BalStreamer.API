package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"

	"github.com/bal3000/BalStreamer.API/models"
)

// LiveStreamHandler - handler for everything to do with the live stream API
type LiveStreamHandler struct {
	liveStreamURL, apiKey string
}

// NewLiveStreamHandler - Creates a new instance of live stream handler
func NewLiveStreamHandler(liveURL string, key string) *LiveStreamHandler {
	return &LiveStreamHandler{liveStreamURL: liveURL, apiKey: key}
}

// GetFixtures - Gets the fixtures for the given sport and date range
func (handler *LiveStreamHandler) GetFixtures(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	sportType := vars["sportType"]
	fromDate := vars["fromDate"]
	toDate := vars["toDate"]

	url := fmt.Sprintf("%s/%s/%s/%s", handler.liveStreamURL, sportType, fromDate, toDate)
	client := &http.Client{}
	request, err := http.NewRequest(http.MethodGet, url, nil)
	logErrors(err)

	request.Header.Add("APIKey", handler.apiKey)
	response, err := client.Do(request)
	logErrors(err)
	defer response.Body.Close()

	fixtures := &[]models.LiveFixtures{}
	err = json.NewDecoder(response.Body).Decode(fixtures)
	logErrors(err)

	if len(*fixtures) == 0 {
		http.Error(res, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}
	res.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(res).Encode(*fixtures); err != nil {
		log.Fatalln(err)
	}
}

// GetStreams gets the streams for the fixture
func (handler *LiveStreamHandler) GetStreams(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	timerID := vars["timerId"]

	url := fmt.Sprintf("%s/%s", handler.liveStreamURL, timerID)
	client := &http.Client{}
	request, err := http.NewRequest(http.MethodGet, url, nil)
	logErrors(err)

	request.Header.Add("APIKey", handler.apiKey)
	response, err := client.Do(req)
	logErrors(err)
	defer response.Body.Close()

	streams := &models.Streams{}
	err = json.NewDecoder(response.Body).Decode(streams)
	logErrors(err)

	res.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(res).Encode(*streams); err != nil {
		log.Fatalln(err)
	}
}

func logErrors(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/bal3000/BalStreamer.API/models"
	"github.com/labstack/echo/v4"
)

// LiveStreamHandler - handler for everything to do with the live stream API
type LiveStreamHandler struct {
	fixturesURL, streamsURL, apiKey string
}

// NewLiveStreamHandler - Creates a new instance of live stream handler
func NewLiveStreamHandler(fURL string, sURL string, key string) *LiveStreamHandler {
	return &LiveStreamHandler{fixturesURL: fURL, streamsURL: sURL, apiKey: key}
}

// GetFixtures - Gets the fixtures for the given sport and date range
func (handler *LiveStreamHandler) GetFixtures(c echo.Context) error {
	sportType := c.Param("sportType")
	fromDate := c.Param("fromDate")
	toDate := c.Param("toDate")

	url := fmt.Sprintf("%s/%s/%s/%s", handler.fixturesURL, sportType, fromDate, toDate)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	logErrors(c, err)

	req.Header.Add("APIKey", handler.apiKey)
	resp, err := client.Do(req)
	logErrors(c, err)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	logErrors(c, err)

	fixtures := convertFixtureResponse(body)

	if len(fixtures) == 0 {
		return c.String(http.StatusNotFound, "No fixtures found")
	}
	return c.JSON(http.StatusOK, fixtures)
}

func convertFixtureResponse(body []byte) []models.LiveFixtures {
	fixtures := &[]models.LiveFixtures{}
	err := json.Unmarshal(body, fixtures)
	if err != nil {
		panic(err)
	}
	return *fixtures
}

func logErrors(c echo.Context, err error) {
	if err != nil {
		c.Logger().Fatal(err)
	}
}

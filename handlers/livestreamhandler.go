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
	liveStreamURL, apiKey string
}

// NewLiveStreamHandler - Creates a new instance of live stream handler
func NewLiveStreamHandler(liveURL string, key string) *LiveStreamHandler {
	return &LiveStreamHandler{liveStreamURL: liveURL, apiKey: key}
}

// GetFixtures - Gets the fixtures for the given sport and date range
func (handler *LiveStreamHandler) GetFixtures(c echo.Context) error {
	sportType := c.Param("sportType")
	fromDate := c.Param("fromDate")
	toDate := c.Param("toDate")

	url := fmt.Sprintf("%s/%s/%s/%s", handler.liveStreamURL, sportType, fromDate, toDate)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	logErrors(c, err)

	req.Header.Add("APIKey", handler.apiKey)
	resp, err := client.Do(req)
	logErrors(c, err)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	logErrors(c, err)

	fixtures := &[]models.LiveFixtures{}
	convertResponse(body, fixtures)

	if len(*fixtures) == 0 {
		return c.String(http.StatusNotFound, "No fixtures found")
	}
	return c.JSON(http.StatusOK, *fixtures)
}

// GetStreams gets the streams for the fixture
func (handler *LiveStreamHandler) GetStreams(c echo.Context) error {
	timerID := c.Param("timerId")

	url := fmt.Sprintf("%s/%s", handler.liveStreamURL, timerID)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	logErrors(c, err)

	req.Header.Add("APIKey", handler.apiKey)
	res, err := client.Do(req)
	logErrors(c, err)
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	logErrors(c, err)

	streams := &models.Streams{}
	convertResponse(body, streams)
	return c.JSON(http.StatusOK, *streams)
}

func convertResponse(body []byte, toConvertTo interface{}) {
	err := json.Unmarshal(body, toConvertTo)
	if err != nil {
		panic(err)
	}
}

func logErrors(c echo.Context, err error) {
	if err != nil {
		c.Logger().Error(err)
	}
}

package handlers

import (
	"testing"
)

const (
	timerID            = "12345"
	sportType          = "Soccer"
	fromDate           = "2021-01-08"
	toDate             = "2021-01-08"
	baseURL            = "/api/livestreams"
	mockURL            = "http://testurl.com"
	mockKey            = "123234324324"
	mockStreamResponse = `{
		"hls": "https://cdn.cloud/src-2503/playlist.m3u8?wUzz3Tsnestarttime=1610123745&wUzz3Tsneendtime=1610127000&wUzz3T=",
		"hlsDvr": "https://cdn.cloud/src-2503/playlist.m3u8?wUzz3Tsnestarttime=1610123745&wUzz3Tsneendtime=1610127000&wUzz3T=&DVR",
		"dash": "https://cdn.cloud/src-2503/manifest.mpd?wUzz3Tsnestarttime=1610123745&wUzz3Tsneendtime=1610127000&wUzz3Tsne=",
		"dashDvr": "https://cdn.cloud/src-2503/manifest.mpd?wUzz3Tsnestarttime=1610123745&wUzz3Tsneendtime=1610127000&wUzz3Tsne=&DVR",
		"rtmp": "rtmp://cdn.cloud:5222/src-2503?wUzz3Tsnestarttime=1610123745&wUzz3Tsneendtime=1610127000&wUzz3Tsne="
	}`
	mockFixtureResponse = `[
		{
			"stateName": "running",
			"utcStart": "2021-01-08T14:40:00",
			"utcEnd": "2021-01-08T17:20:00",
			"title": "Test 1 vs Test 2",
			"eventId": "dsad",
			"contentTypeName": "Soccer",
			"timerId": "67890",
			"isPrimary": "false",
			"broadcastChannelName": "beIN Sports MENA HD",
			"broadcastNationName": "South Africa",
			"sourceTypeName": "Sat-Receiver"
		},
		{
			"stateName": "running",
			"utcStart": "2021-01-08T14:50:00",
			"utcEnd": "2021-01-08T17:30:00",
			"title": "Test 1 vs Test 2",
			"eventId": "vcxcxv",
			"contentTypeName": "Soccer",
			"timerId": "12345",
			"isPrimary": "true",
			"broadcastChannelName": "Super Sport PSL",
			"broadcastNationName": "South Africa",
			"sourceTypeName": "Sat-Receiver"
		}]`
)

func TestGetStreams(t *testing.T) {
	// TODO: REWORK SO I CAN SEND A MOCK HTTP CLIENT
	// Setup
	//e := echo.New()
	//req := httptest.NewRequest(http.MethodGet, "/", nil)
	//rec := httptest.NewRecorder()
	//c := e.NewContext(req, rec)
	//
	//path := fmt.Sprintf("%s/:timerId", baseURL)
	//c.SetPath(path)
	//c.SetParamNames("timerId")
	//c.SetParamValues(timerID)
	//
	//liveStreamHandler := NewLiveStreamHandler(mockURL, mockKey)
	//
	//// Assertions
	//if assert.NoError(t, liveStreamHandler.GetStreams(c)) {
	//	assert.Equal(t, http.StatusOK, rec.Code)
	//	assert.Equal(t, mockStreamResponse, rec.Body.String())
	//}
}

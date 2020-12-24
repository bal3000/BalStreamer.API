package models

// StreamToCast - the model for the json posted to the cast controller
type StreamToCast struct {
	Chromecast string `json:"chromecast"`
	StreamURL  string `json:"streamURL"`
}

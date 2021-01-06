package models

import "time"

// StreamToCast - the model for the json posted to the cast controller
type StreamToCast struct {
	Chromecast string `json:"chromecast"`
	StreamURL  string `json:"streamURL"`
}

// StopPlayingStream is the model for the json posted to the stop casting endpoint
type StopPlayingStream struct {
	ChromeCastToStop string    `json:"chromeCastToStop"`
	StopDateTime     time.Time `json:"stopDateTime"`
}

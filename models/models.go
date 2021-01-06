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

// LiveFixtures is the model for the json returned from the live stream api
type LiveFixtures struct {
	StateName            string    `json:"stateName"`
	UtcStart             time.Time `json:"utcStart"`
	UtcEnd               time.Time `json:"utcEnd"`
	Title                string    `json:"title"`
	EventID              string    `json:"eventId"`
	ContentTypeName      string    `json:"contentTypeName"`
	TimerID              int       `json:"timerId"`
	IsPrimary            bool      `json:"isPrimary"`
	BroadcastChannelName string    `json:"broadcastChannelName"`
	BroadcastNationName  string    `json:"broadcastNationName"`
	SourceTypeName       string    `json:"sourceTypeName"`
}

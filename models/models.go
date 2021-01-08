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
	StateName            string `json:"stateName"`
	UtcStart             string `json:"utcStart"`
	UtcEnd               string `json:"utcEnd"`
	Title                string `json:"title"`
	EventID              string `json:"eventId"`
	ContentTypeName      string `json:"contentTypeName"`
	TimerID              string `json:"timerId"`
	IsPrimary            string `json:"isPrimary"`
	BroadcastChannelName string `json:"broadcastChannelName"`
	BroadcastNationName  string `json:"broadcastNationName"`
	SourceTypeName       string `json:"sourceTypeName"`
}

// Streams is the model for the json returned from the live stream api
type Streams struct {
	HLS     string `json:"hls"`
	HLSDvr  string `json:"hlsDvr"`
	Dash    string `json:"dash"`
	DashDvr string `json:"dashDvr"`
	RTMP    string `json:"rtmp"`
}

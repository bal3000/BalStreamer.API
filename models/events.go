package models

import (
	"encoding/json"
	"time"
)

// EventMessage interface for transforming messages to masstransit ones
type EventMessage interface {
	TransformMessage() ([]byte, string, error)
}

// StreamToChromecastEvent the send to chromecast event
type StreamToChromecastEvent struct {
	ChromeCastToStream string    `json:"chromeCastToStream"`
	Stream             string    `json:"stream"`
	StreamDate         time.Time `json:"streamDate"`
}

// StopPlayingStreamEvent the stop cast event
type StopPlayingStreamEvent struct {
	ChromeCastToStop string    `json:"chromeCastToStop"`
	StopDateTime     time.Time `json:"stopDateTime"`
}

// ChromecastEvent event when a chromecast is found
type ChromecastEvent struct {
	Chromecast string `json:"chromecast"`
}

// TransformMessage transforms the message to a masstransit one and then turns into JSON
func (message *StreamToChromecastEvent) TransformMessage() ([]byte, string, error) {
	data, err := json.Marshal(message)
	if err != nil {
		return nil, "", err
	}
	return data, "StreamToChromecastEvent", nil
}

// TransformMessage transforms the message to a masstransit one and then turns into JSON
func (message *StopPlayingStreamEvent) TransformMessage() ([]byte, string, error) {
	data, err := json.Marshal(message)
	if err != nil {
		return nil, "", err
	}
	return data, "StopPlayingStreamEvent", nil
}

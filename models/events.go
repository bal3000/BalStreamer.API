package models

import (
	"encoding/json"
	"time"
)

// EventMessage interface for transforming messages to masstransit ones
type EventMessage interface {
	TransformMessage() ([]byte, error)
}

// MassTransitEvent the wrapper so mass transit can accept the events
type massTransitEvent struct {
	Message     interface{} `json:"message"`
	MessageType []string    `json:"messageType"`
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

// ChromecastFoundEvent event when a chromecast is found
type ChromecastFoundEvent struct {
	Chromecast string `json:"chromecast"`
}

// ChromecastLostEvent event when a chromecast is lost
type ChromecastLostEvent struct {
	Chromecast string `json:"chromecast"`
}

// TransformMessage transforms the message to a masstransit one and then turns into JSON
func (message *StreamToChromecastEvent) TransformMessage() ([]byte, error) {
	mtEvent := massTransitEvent{
		Message:     message,
		MessageType: []string{"urn:message:BalStreamer.Shared.EventBus.Events:StreamToChromecastEvent"},
	}

	return json.Marshal(mtEvent)
}

// TransformMessage transforms the message to a masstransit one and then turns into JSON
func (message *StopPlayingStreamEvent) TransformMessage() ([]byte, error) {
	mtEvent := massTransitEvent{
		Message:     message,
		MessageType: []string{"urn:message:BalStreamer.Shared.EventBus.Events:StopPlayingStreamEvent"},
	}

	return json.Marshal(mtEvent)
}

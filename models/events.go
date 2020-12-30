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

// TransformMessage transforms the message to a masstransit one and then turns into JSON
func (message *StreamToChromecastEvent) TransformMessage() ([]byte, error) {
	mtEvent := massTransitEvent{
		Message:     message,
		MessageType: []string{""},
	}

	return json.Marshal(mtEvent)
}

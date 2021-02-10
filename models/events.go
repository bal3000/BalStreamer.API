package models

import (
	"encoding/json"
	"strings"
	"time"
)

// EventMessage interface for transforming messages to masstransit ones
type EventMessage interface {
	TransformMessage() ([]byte, string, error)
}

// MassTransitEvent the wrapper so mass transit can accept the events
type MassTransitEvent struct {
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

// ChromecastEvent event when a chromecast is found
type ChromecastEvent struct {
	EventType  string      `json:"eventType"`
	Chromecast interface{} `json:"chromecast"`
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

// RetrieveMessage converts the incoming message to a struct
func (eve *MassTransitEvent) RetrieveMessage(msg []byte) (ChromecastEvent, error) {
	err := json.Unmarshal(msg, eve)
	if err != nil {
		return ChromecastEvent{}, err
	}

	firstType := strings.Split(eve.MessageType[0], ":")
	eventName := firstType[len(firstType)-1]
	message := eve.Message.(map[string]interface{})

	chromecastEvent := ChromecastEvent{Chromecast: message["chromecast"], EventType: eventName}
	return chromecastEvent, err
}

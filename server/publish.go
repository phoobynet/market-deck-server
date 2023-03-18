package server

import (
	"encoding/json"
	"github.com/phoobynet/market-deck-server/events"
	"github.com/r3labs/sse/v2"
	"github.com/sirupsen/logrus"
)

func Publish[T any](messageType events.EmittableEvents, message T) {
	if sseServer == nil {
		logrus.Panic("SSE server is not ready")
	}

	jsonData, err := json.Marshal(message)

	if err != nil {
		logrus.Errorf("Error marshalling message: %s", err)
	}

	sseServer.Publish(
		string(messageType), &sse.Event{
			Data: jsonData,
		},
	)
}

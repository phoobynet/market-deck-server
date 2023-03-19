package server

import (
	"encoding/json"
	"github.com/r3labs/sse/v2"
	"github.com/sirupsen/logrus"
)

func Publish(message Message) {
	if sseServer == nil {
		logrus.Panic("SSE server is not ready")
	}

	jsonData, err := json.Marshal(message)

	if err != nil {
		logrus.Errorf("Error marshalling message: %s", err)
	}

	sseServer.Publish(
		string(message.Event), &sse.Event{
			Data: jsonData,
		},
	)
}

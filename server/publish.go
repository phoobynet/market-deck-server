package server

import (
	"encoding/json"
	"github.com/phoobynet/market-deck-server/messages"
	"github.com/r3labs/sse/v2"
	"github.com/sirupsen/logrus"
)

func Publish(message messages.Message) {
	if sseServer == nil {
		logrus.Panic("SSE server is not ready")
	}

	jsonData, err := json.Marshal(message)

	if err != nil {
		logrus.Errorf("Error marshalling message: %s\n%v", err, message.Data)
	}

	sseServer.Publish(
		string(message.Event), &sse.Event{
			Data: jsonData,
		},
	)
}

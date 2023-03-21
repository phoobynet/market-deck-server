package server

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/phoobynet/market-deck-server/messages"
	"github.com/r3labs/sse/v2"
	"github.com/sirupsen/logrus"
)

var jsonEncoder = jsoniter.ConfigCompatibleWithStandardLibrary

func Publish(message messages.Message) {
	if sseServer == nil {
		logrus.Panic("SSE server is not ready")
	}

	jsonData, err := jsonEncoder.Marshal(message)

	if err != nil {
		logrus.Panicf("Error marshalling message: %s\n%v", err, message.Data)
	}

	sseServer.Publish(
		string(message.Event), &sse.Event{
			Data: jsonData,
		},
	)
}

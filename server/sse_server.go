package server

import (
	"github.com/julienschmidt/httprouter"
	"github.com/phoobynet/market-deck-server/events"
	"github.com/r3labs/sse/v2"
	"github.com/sirupsen/logrus"
	"net/http"
	"sync"
)

var sseServer *sse.Server
var sseMutex sync.Mutex

func InitSSE() {
	sseMutex.Lock()
	defer sseMutex.Unlock()

	if sseServer == nil {
		sseServer = sse.New()
		sseServer.AutoReplay = false
		sseServer.CreateStream(events.CalendarDayUpdate)
		sseServer.CreateStream(events.RealtimeSymbols)
		sseServer.CreateStream(events.Messages)
		sseServer.CreateStream(events.Errors)
	}
}

func getSSEHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	streamID := r.URL.Query().Get("stream")
	go func() {
		<-r.Context().Done()
		logrus.Infof("Client disconnected from %s\n", streamID)
	}()

	logrus.Infof("Client connected to stream %s...\n", streamID)

	sseServer.ServeHTTP(w, r)
}

package server

import (
	"embed"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/phoobynet/market-deck-server/assets"
	"github.com/phoobynet/market-deck-server/decks"
	"github.com/phoobynet/market-deck-server/snapshots"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"io"
	"io/fs"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Server struct {
	config *Config
	mux    *http.ServeMux
}

func NewServer(
	config *Config,
	dist embed.FS,
	snapshotStream *snapshots.Stream,
) *Server {
	assetRepository := assets.GetRepository()
	deckRepository := decks.GetRepository()

	webServer := &Server{
		config,
		http.NewServeMux(),
	}

	fsys, distFSErr := fs.Sub(dist, "dist")

	if distFSErr != nil {
		panic(distFSErr)
	}

	router := httprouter.New()

	myCors := cors.New(
		cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"*"},
			AllowCredentials: false,
		},
	)

	router.GET("/api/stream", getSSEHandler)

	router.GET(
		"/api/assets", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
			start := time.Now()
			all := assetRepository.GetAll()

			_ = writeJSON(w, http.StatusOK, all)
			end := time.Now()
			logrus.Infof("GetAll took %s", end.Sub(start))
		},
	)

	router.GET(
		"/api/symbols/query", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
			query := r.URL.Query().Get("query")

			if query == "" {
				_ = writeErr(w, http.StatusBadRequest, fmt.Errorf("query parameter is required"))
				return
			}

			results := assetRepository.Search(query)

			_ = writeJSON(w, http.StatusOK, results)
		},
	)

	router.POST(
		"/api/symbols", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
			symbols := r.URL.Query().Get("symbols")

			if symbols == "" {
				snapshotStream.UpdateSymbols([]string{})
			}

			snapshotStream.UpdateSymbols(strings.Split(symbols, ","))

			w.WriteHeader(http.StatusOK)
		},
	)

	router.GET(
		"/api/symbols", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
			deck, err := deckRepository.FindByName("default")

			if err != nil {
				_ = writeErr(w, http.StatusInternalServerError, err)
				return
			}

			var symbols []string

			if len(deck.Symbols) > 0 {
				symbols = strings.Split(deck.Symbols, ",")
				_ = writeJSON(
					w, http.StatusOK, map[string]interface{}{
						"symbols": symbols,
					},
				)
			}
		},
	)

	router.GET(
		"/api/decks", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
			all, err := deckRepository.FindAll()

			if err != nil {
				_ = writeErr(w, http.StatusInternalServerError, err)
				return
			}

			_ = writeJSON(w, http.StatusOK, all)
		},
	)

	router.DELETE(
		"/api/decks/:id", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
			id, err := strconv.Atoi(r.URL.Query().Get("id"))

			if err != nil {
				_ = writeErr(w, http.StatusBadRequest, err)
				return
			}

			err = deckRepository.Delete(uint(id))

			if err != nil {
				_ = writeErr(w, http.StatusInternalServerError, err)
			}

			w.WriteHeader(http.StatusOK)
		},
	)

	router.POST(
		"/api/decks", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
			body, err := r.GetBody()

			if err != nil {
				_ = writeErr(w, http.StatusBadRequest, err)
				return
			}

			b, err := io.ReadAll(body)

			var deck decks.Deck

			err = json.Unmarshal(b, &deck)

			if err != nil {
				_ = writeErr(w, http.StatusBadRequest, err)
				return
			}

			created, err := deckRepository.Create(deck.Name, strings.Split(deck.Symbols, ","))

			if err != nil {
				_ = writeErr(w, http.StatusBadRequest, err)
				return
			}

			_ = writeJSON(w, http.StatusCreated, created)
		},
	)

	router.PUT(
		"/api/decks/:id", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
			body, err := r.GetBody()

			if err != nil {
				_ = writeErr(w, http.StatusBadRequest, err)
				return
			}

			id, err := strconv.Atoi(r.URL.Query().Get("id"))

			if err != nil {
				_ = writeErr(w, http.StatusBadRequest, err)
				return
			}

			b, err := io.ReadAll(body)

			var deck decks.Deck

			err = json.Unmarshal(b, &deck)

			if err != nil {
				_ = writeErr(w, http.StatusBadRequest, err)
				return
			}

			update, err := deckRepository.Update(uint(id), deck.Name, strings.Split(deck.Symbols, ","))

			if err != nil {
				_ = writeErr(w, http.StatusBadRequest, err)
				return
			}

			_ = writeJSON(w, http.StatusOK, update)
		},
	)

	router.GET(
		"/api/bars/intraday", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
			bars := snapshotStream.GetIntradayBars()

			_ = writeJSON(w, http.StatusOK, bars)
		},
	)

	router.GET(
		"/api/bars/ytd", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
			bars := snapshotStream.GetYtdBars()

			_ = writeJSON(w, http.StatusOK, bars)
		},
	)

	webServer.mux.Handle("/", http.FileServer(http.FS(fsys)))
	webServer.mux.Handle("/api/", myCors.Handler(router))

	return webServer
}

func (s *Server) Listen() {
	listenErr := http.ListenAndServe(fmt.Sprintf(":%d", s.config.ServerPort), s.mux)

	if listenErr != nil {
		logrus.Fatalln(listenErr)
	}
}

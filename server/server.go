package server

import (
	"embed"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/phoobynet/market-deck-server/assets"
	"github.com/phoobynet/market-deck-server/decks"
	"github.com/phoobynet/market-deck-server/scrapers/yahoo"
	"github.com/phoobynet/market-deck-server/sec/facts"
	ss "github.com/phoobynet/market-deck-server/snapshots/stream"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"io"
	"io/fs"
	"net/http"
	"strconv"
	"strings"
)

type Server struct {
	config *Config
	mux    *http.ServeMux
}

func NewServer(
	config *Config,
	dist embed.FS,
	snapshotLiteStream *ss.SnapshotStream,
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
			all := assetRepository.GetAll()

			_ = writeJSON(w, http.StatusOK, all)
		},
	)

	router.GET(
		"/api/sec/:ticker/facts",
		func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
			query := facts.FactQuery{
				Ticker:          ps.ByName("ticker"),
				FinancialYear:   0,
				FinancialPeriod: r.URL.Query().Get("fp"),
				Form:            r.URL.Query().Get("form"),
				Concept:         r.URL.Query().Get("concept"),
			}

			fy := r.URL.Query().Get("fy")

			if len(fy) > 0 {
				fyInt, err := strconv.Atoi(fy)

				if err != nil {
					_ = writeErr(w, http.StatusBadRequest, err)
				}

				query.FinancialYear = fyInt
			}

			queryResult := facts.Get(query)

			_ = writeJSON(w, http.StatusOK, queryResult)
		},
	)

	router.POST(
		"/api/symbols", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
			symbols := r.URL.Query().Get("symbols")

			if symbols == "" {
				snapshotLiteStream.UpdateSymbols([]string{})
			} else {
				snapshotLiteStream.UpdateSymbols(strings.Split(symbols, ","))
			}

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
		"/api/scraped/:ticker", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
			summary, err := yahoo.GetTickerSummary(ps.ByName("ticker"))

			if err != nil {
				_ = writeErr(w, http.StatusInternalServerError, err)
			}

			_ = writeJSON(w, http.StatusOK, summary)
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

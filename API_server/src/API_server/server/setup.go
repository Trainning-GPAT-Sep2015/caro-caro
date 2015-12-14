package server

import (
	"net/http"

	"API_server/middlewares"
	"API_server/utils/logs"

	"API_server/handlers"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

var l = logs.New("API_server")

type setupStruct struct {
	Config

	Handler http.Handler
}

func setup(cfg Config) *setupStruct {
	s := &setupStruct{Config: cfg}
	s.setupRoutes()

	return s
}

func commonMiddlewares() func(http.Handler) http.Handler {
	logger := middlewares.NewLogger()
	recovery := middlewares.NewRecovery()

	return func(h http.Handler) http.Handler {
		return recovery(logger(h))
	}
}

func authMiddlewares() func(http.Handler) http.Handler {
	auth := middlewares.NewAuth()

	return func(h http.Handler) http.Handler {
		return auth(h)
	}
}

func (s *setupStruct) setupRoutes() {
	commonMids := commonMiddlewares()
	authMids := authMiddlewares()

	normal := func(h http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			commonMids(h).ServeHTTP(w, r)
		}
	}

	auth := func(h http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			commonMids(authMids(h)).ServeHTTP(w, r)
		}
	}

	router := mux.NewRouter()
	router.HandleFunc("/", normal(handlers.Home)).Methods("GET")
	router.HandleFunc("/{user}", auth(handlers.Restrict)).Methods("GET")

	s.Handler = context.ClearHandler(router)
}
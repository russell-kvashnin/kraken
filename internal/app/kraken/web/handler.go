package web

import (
	"github.com/russell-kvashnin/kraken/internal/pkg/log"
	"net/http"
)

type HttpHandler struct {
	log *log.Logger
	h   Handler
}

type Handler func(w http.ResponseWriter, r *http.Request) error

func (handler HttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := handler.h(w, r); err != nil {
		// handle returned error here.
		w.WriteHeader(503)
		_, err := w.Write([]byte("bad"))

		handler.log.Errorw("Error while writing response", "err", err)
	}
}

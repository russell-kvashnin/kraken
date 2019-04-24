package handler

import (
	kerr "github.com/russell-kvashnin/kraken/internal/pkg/error"
	"net/http"
)

// Handler function
type handlerFun func(w http.ResponseWriter, r *http.Request) error

// Http action handler
type HttpHandler struct {
	errHandler *kerr.Handler
	fun        handlerFun
}

// Http handler constructor
func NewHttpHandler(errHandler *kerr.Handler, fun handlerFun) *HttpHandler {
	handler := new(HttpHandler)
	handler.errHandler = errHandler
	handler.fun = fun

	return handler
}

// ServeHTTP implementation
func (handler *HttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := handler.fun(w, r); err != nil {
		handler.errHandler.Handle(err)

		w.WriteHeader(503)
		_, err = w.Write([]byte("bad"))
	}
}

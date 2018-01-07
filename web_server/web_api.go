package web_server

import (
	"fmt"
	"github.com/celrenheit/lion"
	"github.com/goincremental/negroni-sessions"
	"github.com/goincremental/negroni-sessions/cookiestore"
	"github.com/jmesyan/xingo/logger"
	"github.com/urfave/negroni"
	"net/http"
)

type WebApi struct {
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "this is the custum not found")
}

func NegroniWrap(h negroni.HandlerFunc) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		nw := negroni.NewResponseWriter(w)
		h(nw, r, next)
	}
}

func (this *WebApi) Ready(r *lion.Router) {
	nh := http.HandlerFunc(NotFound)
	lion.WithNotFoundHandler(nh)(r)
	store := cookiestore.New([]byte("secret123"))
	session := sessions.Sessions("my_session", store)
	admin := r.Group("/admin")
	admin.UseNext(NegroniWrap(session))
	admin.Mount("/monitor", Monitor)
	api := r.Group("/api")
	api.Mount("/sess", Api)
	logger.Info("the web api ready")
}

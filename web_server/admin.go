package web_server

import (
	"fmt"
	"github.com/celrenheit/lion"
	"github.com/goincremental/negroni-sessions"
	"net/http"
)

var (
	Monitor *lion.Router
	Api     *lion.Router
)

func init() {
	Monitor = lion.New()
	Monitor.GetFunc("/hello", Hello)
	Monitor.GetFunc("/setSession", setSession)
	Monitor.GetFunc("/showSession", showSession)
	Monitor.GetFunc("/clearSession", clearSession)

	Api = lion.New()
	Api.GetFunc("/showSession", showSession)
}

func Hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "this is the admin hello method")
}

func setSession(w http.ResponseWriter, r *http.Request) {
	session := sessions.GetSession(r)
	session.Set("hello", "world")
	fmt.Fprintf(w, "setOK")
}

func clearSession(w http.ResponseWriter, r *http.Request) {
	session := sessions.GetSession(r)
	session.Clear()
	fmt.Fprintf(w, "clearOK")
}

func showSession(w http.ResponseWriter, r *http.Request) {
	session := sessions.GetSession(r)
	if session == nil {
		fmt.Fprintf(w, "session is not set")
	} else {
		hello := session.Get("hello")
		if hello == nil {
			fmt.Fprintf(w, "hello is not set")
		} else {
			fmt.Fprintf(w, hello.(string))
		}

	}

}

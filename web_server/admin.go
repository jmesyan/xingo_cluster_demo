package web_server

import (
	"fmt"
	"github.com/celrenheit/lion"
	"net/http"
)

var (
	Monitor *lion.Router
)

func init() {
	Monitor = lion.New()
	Monitor.GetFunc("/hello", Hello)
}

func Hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "this is the admin hello method")
}

package web_server

import (
	"github.com/celrenheit/lion"
	"github.com/jmesyan/xingo/logger"
	"net/http"
)

type WebApi struct {
}

func Good(w http.ResponseWriter, r *http.Request) {
	logger.Info("this is good api test")
	w.Write([]byte("good is good"))
}

func (this *WebApi) Ready(r *lion.Router) {
	r.GetFunc("/good", Good)
	logger.Info("this is ready api test")
}

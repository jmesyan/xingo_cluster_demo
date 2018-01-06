package web_server

import (
	"github.com/celrenheit/lion"
	"github.com/jmesyan/xingo/logger"
	// "net/http"
)

type WebApi struct {
}

func (this *WebApi) Ready(r *lion.Router) {
	admin := r.Group("/admin")
	admin.Mount("/monitor", Monitor)
	logger.Info("the web api ready")
}

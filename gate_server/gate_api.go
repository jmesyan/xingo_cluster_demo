package gate_server

import (
	"github.com/jmesyan/xingo/cluster"
	"github.com/jmesyan/xingo/clusterserver"
	"github.com/jmesyan/xingo/logger"
)

type GateRpcApi struct {
}

func (this *GateRpcApi) CreatePlayer(request *cluster.RpcRequest) map[string]interface{} {
	onegame := clusterserver.GlobalClusterServer.ChildsMgr.GetRandomChild("game")
	logger.Info("onegame", onegame)
	if onegame != nil {
		response, err := onegame.CallChildForResult("CreatePlayer")
		if err == nil {
			// logger.Info("gate_api", response)
			return response.Result
		} else {
			logger.Error(err)
		}
	}
	logger.Info("gate_api", "no game server online")
	return map[string]interface{}{
		"pid": 0,
	}
}

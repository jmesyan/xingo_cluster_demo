package gate_server

import (
	"github.com/jmesyan/xingo/cluster"
	"github.com/jmesyan/xingo/clusterserver"
	"github.com/jmesyan/xingo/logger"
	"xingo_cluster_demo/core"
)

type GateRpcApi struct {
}

func (this *GateRpcApi) CreatePlayer(request *cluster.RpcRequest) map[string]interface{} {
	netname := (request.Rpcdata.Args[0]).(string)
	onegame := clusterserver.GlobalClusterServer.ChildsMgr.GetRandomChild("game")
	logger.Info("onegame", onegame)
	if onegame != nil {
		response, err := onegame.CallChildForResult("CreatePlayer", netname)
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

func (this *GateRpcApi) SyncSurrounds(request *cluster.RpcRequest) {
	py, np := request.Rpcdata.Args[0].(core.Player), request.Rpcdata.Args[1].(core.Player)
	netname := py.Net
	net, err := clusterserver.GlobalClusterServer.ChildsMgr.GetChild(netname)
	if err == nil {
		net.CallChildNotForResult("SyncSurrounds", py, np)
	} else {
		logger.Error("can found the net")
	}

}

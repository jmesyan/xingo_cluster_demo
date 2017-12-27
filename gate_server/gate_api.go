package gate_server

import (
	// "github.com/golang/protobuf/proto"
	"github.com/jmesyan/xingo/cluster"
	"github.com/jmesyan/xingo/clusterserver"
	"github.com/jmesyan/xingo/logger"
	"xingo_cluster_demo/core"
	"xingo_cluster_demo/pb"
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
	} else {
		logger.Info("gate_api", "no game server online")
	}
	return map[string]interface{}{
		"pid": 0,
	}
}

func (this *GateRpcApi) SyncSurrounds(request *cluster.RpcRequest) {
	cmd := request.Rpcdata.Args[0].(int32)
	py, np := request.Rpcdata.Args[1].(core.Player), request.Rpcdata.Args[2].(core.Player)
	netname := py.Net
	net, err := clusterserver.GlobalClusterServer.ChildsMgr.GetChild(netname)
	if err == nil {
		net.CallChildNotForResult("SyncSurrounds", cmd, py, np)
	} else {
		logger.Error("can not found the net")
	}

}

func (this *GateRpcApi) BroadCastMsg(request *cluster.RpcRequest) {
	pid := request.Rpcdata.Args[0].(int32)
	content := request.Rpcdata.Args[1].(string)
	childs := clusterserver.GlobalClusterServer.ChildsMgr.GetChildsByPrefix("net")
	for _, net := range childs {
		net.CallChildNotForResult("BroadCastMsg", pid, content)
	}

}

func (this *GateRpcApi) UpdatePos(request *cluster.RpcRequest) {
	pid := request.Rpcdata.Args[0].(int32)
	position := request.Rpcdata.Args[1].(pb.Position)
	logger.Info(pid, position, "gate updatepos")
	onegame := clusterserver.GlobalClusterServer.ChildsMgr.GetRandomChild("game")
	logger.Info("onegame", onegame)
	if onegame != nil {
		onegame.CallChildNotForResult("UpdatePos", pid, position)
	} else {
		logger.Info("gate_api", "no game server online")
	}
}

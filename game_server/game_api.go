package game_server

import (
	"github.com/jmesyan/xingo/cluster"
	"xingo_cluster_demo/core"
)

type GameRpcApi struct {
}

func (this *GameRpcApi) CreatePlayer(request *cluster.RpcRequest) map[string]interface{} {
	p, _ := core.WorldMgrObj.AddPlayer()
	return map[string]interface{}{
		"p": p,
		// "sname":"game1"Pid
	}
}

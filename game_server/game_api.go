package game_server

import (
	"github.com/jmesyan/xingo/cluster"
	"github.com/jmesyan/xingo/clusterserver"
	"github.com/jmesyan/xingo/logger"
	"xingo_cluster_demo/core"
)

type GameRpcApi struct {
}

func syncSurrounds(gate *cluster.Child, py, np core.Player) {
	gate.CallChildNotForResult("SyncSurrounds", py, np)
}

func (this *GameRpcApi) CreatePlayer(request *cluster.RpcRequest) map[string]interface{} {
	var surrounds []core.Player
	netname := request.Rpcdata.Args[0].(string)
	p, _ := core.WorldMgrObj.AddPlayer(netname)
	pids, err := core.WorldMgrObj.AoiObj1.GetSurroundingPids(p)
	if err == nil {
		onegate := clusterserver.GlobalClusterServer.RemoteNodesMgr.GetRandomChild("gate")
		for _, pid := range pids {
			if pid != p.Pid {
				py, err := core.WorldMgrObj.GetPlayer(pid)
				if err == nil {
					surrounds = append(surrounds, *py)
					//给surrounds发送同步消息
					if onegate != nil {
						logger.Info("chose gate: " + onegate.GetName())
						go syncSurrounds(onegate, *py, *p)
					}

				}
			} else {
				surrounds = append(surrounds, *p)
			}
		}
	}

	return map[string]interface{}{
		"p":  p,
		"sr": surrounds,
		// "sname":"game1"Pid
	}
}

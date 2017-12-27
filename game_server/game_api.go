package game_server

import (
	"github.com/jmesyan/xingo/cluster"
	"github.com/jmesyan/xingo/clusterserver"
	"xingo_cluster_demo/core"
	"xingo_cluster_demo/pb"
)

type GameRpcApi struct {
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
						go SyncSurrounds(onegate, 200, *py, *p)
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
	}
}

func (this *GameRpcApi) UpdatePos(request *cluster.RpcRequest) {
	pid := request.Rpcdata.Args[0].(int32)
	ps := request.Rpcdata.Args[1].(pb.Position)
	p, _ := core.WorldMgrObj.GetPlayer(pid)
	p.UpdatePos(ps.X, ps.Y, ps.Z, ps.V)
}

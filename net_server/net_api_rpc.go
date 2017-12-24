package net_server

import (
	"github.com/jmesyan/xingo/cluster"
	// "github.com/jmesyan/xingo/clusterserver"
	"github.com/jmesyan/xingo/logger"
	// "github.com/jmesyan/xingo/utils"
	"xingo_cluster_demo/core"
	"xingo_cluster_demo/pb"
)

type NetRpcApi struct {
}

func (this *NetRpcApi) SyncSurrounds(request *cluster.RpcRequest) {
	py, np := request.Rpcdata.Args[0].(core.Player), request.Rpcdata.Args[1].(core.Player)
	if p, ok := NetPlayers[py.Pid]; ok {
		position := &pb.Position{
			X: np.X,
			Y: np.Y,
			Z: np.Z,
			V: np.V,
		}
		//出现在自己的视野中
		data := &pb.BroadCast{
			Pid: np.Pid,
			Tp:  2,
			Data: &pb.BroadCast_P{
				P: position,
			},
		}

		SendMsg(p.Fconn, 200, data)
	} else {
		// netname := utils.GlobalObj.Name
		logger.Info("no player find in net:")
	}
}

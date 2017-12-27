package net_server

import (
	"github.com/jmesyan/xingo/cluster"
	// "github.com/jmesyan/xingo/clusterserver"
	"github.com/jmesyan/xingo/logger"
	// "github.com/jmesyan/xingo/utils"
	// "github.com/golang/protobuf/proto"
	"xingo_cluster_demo/core"
	"xingo_cluster_demo/pb"
)

type NetRpcApi struct {
}

func (this *NetRpcApi) SyncSurrounds(request *cluster.RpcRequest) {
	cmd := request.Rpcdata.Args[0].(int32)
	py, np := request.Rpcdata.Args[1].(core.Player), request.Rpcdata.Args[2].(core.Player)
	logger.Info("SyncSurrounds", cmd, py, np, NetPlayers)
	SyncPosition(py)
	SyncPosition(np)
	if cmd == 200 {
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

	if cmd == 201 {
		if p, ok := NetPlayers[py.Pid]; ok {
			data := &pb.SyncPid{
				Pid: np.Pid,
			}
			SendMsg(p.Fconn, 201, data)
		} else {
			// netname := utils.GlobalObj.Name
			logger.Info("no player find in net:")
		}
	}

	if cmd == 211 {
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
				Tp:  4,
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
}

func (this *NetRpcApi) BroadCastMsg(request *cluster.RpcRequest) {
	pid := request.Rpcdata.Args[0].(int32)
	content := request.Rpcdata.Args[1].(string)

	data := &pb.BroadCast{
		Pid: pid,
		Tp:  1,
		Data: &pb.BroadCast_Content{
			Content: content,
		},
	}
	for _, py := range NetPlayers {
		if py.Fconn != nil {
			SendBuffMsg(py.Fconn, 200, data)
		}
	}
	logger.Info("BroadCastMsg to client")
}

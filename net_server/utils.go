package net_server

import (
	"github.com/golang/protobuf/proto"
	"github.com/jmesyan/xingo/cluster"
	"github.com/jmesyan/xingo/clusterserver"
	"github.com/jmesyan/xingo/iface"
	"github.com/jmesyan/xingo/logger"
	"github.com/jmesyan/xingo/utils"
	"xingo_cluster_demo/core"
)

var NetPlayers map[int32]core.Player

func init() {
	NetPlayers = make(map[int32]core.Player)
}

func SendMsg(fconn iface.Iconnection, msgId uint32, data proto.Message) {
	if fconn != nil {
		packdata, err := utils.GlobalObject.Protoc.GetDataPack().Pack(msgId, data)
		if err == nil {
			fconn.Send(packdata)
		} else {
			logger.Error("pack data error")
		}
	}
}

func SendBuffMsg(fconn iface.Iconnection, msgId uint32, data proto.Message) {
	if fconn != nil {
		packdata, err := utils.GlobalObject.Protoc.GetDataPack().Pack(msgId, data)
		if err == nil {
			fconn.SendBuff(packdata)
		} else {
			logger.Error("pack data error")
		}
	}
}

func GetRandomGate() *cluster.Child {
	return clusterserver.GlobalClusterServer.RemoteNodesMgr.GetRandomChild("gate")
}

func SyncPosition(player core.Player) {
	pid := player.Pid
	p, ok := NetPlayers[pid]
	if ok {
		player.Fconn = p.Fconn
		NetPlayers[pid] = player
	}
}

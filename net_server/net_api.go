package net_server

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/jmesyan/xingo/clusterserver"
	"github.com/jmesyan/xingo/fnet"
	"github.com/jmesyan/xingo/iface"
	"github.com/jmesyan/xingo/logger"
	"github.com/jmesyan/xingo/utils"
	"time"
	"xingo_cluster_demo/core"
	"xingo_cluster_demo/pb"
)

func DoConnectioned(fconn iface.Iconnection) {
	st := time.Now()
	logger.Info("connection connect , I get it")
	//请求pid
	onegate := clusterserver.GlobalClusterServer.RemoteNodesMgr.GetRandomChild("gate")

	if onegate != nil {
		logger.Info("chose gate: " + onegate.GetName())
		response, err := onegate.CallChildForResult("CreatePlayer")
		if err == nil {
			pid, _ := response.Result["pid"].(int32)
			if pid > 0 {
				logger.Info("get pid success, pid:", pid)
				fconn.SetProperty("pid", pid)
				//同步Pid
				msg := &pb.SyncPid{
					Pid: pid,
				}
				packdata, err := utils.GlobalObject.Protoc.GetDataPack().Pack(1, msg)
				if err == nil {
					fconn.Send(packdata)
					diff := time.Now().Sub(st).Nanoseconds()
					logger.Info("get pid total consume:", (diff / 1e6), "ms")
				} else {
					logger.Error("pack data error")
				}

			} else {
				logger.Info("no game server serve")
				fconn.LostConnection()
			}
		} else {
			logger.Error(err)
		}
	}

}

func init() {
	utils.GlobalObject.OnConnectioned = DoConnectioned
}

type NetApiRouter struct {
}

func (this *NetApiRouter) Api_0(request *fnet.PkgAll) {
	logger.Debug("call Api_0")
	// request.Fconn.SendBuff(0, nil)
	packdata, err := utils.GlobalObject.Protoc.GetDataPack().Pack(0, nil)
	if err == nil {
		request.Fconn.Send(packdata)
	} else {
		logger.Error("pack data error")
	}
}

func (this *NetApiRouter) Api_2(request *fnet.PkgAll) {
	msg := &pb.Talk{}
	err := proto.Unmarshal(request.Pdata.Data, msg)
	if err == nil {
		logger.Debug(fmt.Sprintf("user talk: content: %s.", msg.Content))
		// pid, err1 := request.Fconn.GetProperty("pid")
		pid, err1 := request.Fconn.GetProperty("pid")
		if err1 == nil {
			p, _ := core.WorldMgrObj.GetPlayer(pid.(int32))
			p.Talk(msg.Content)
		} else {
			logger.Error(err1)
			request.Fconn.LostConnection()
		}

	} else {
		logger.Error(err)
		request.Fconn.LostConnection()
	}
}

func (this *NetApiRouter) Api_3(request *fnet.PkgAll) {
	msg := &pb.Position{}
	err := proto.Unmarshal(request.Pdata.Data, msg)
	if err == nil {
		logger.Debug(fmt.Sprintf("user move: (%f, %f, %f, %f)", msg.X, msg.Y, msg.Z, msg.V))
		pid, err1 := request.Fconn.GetProperty("pid")
		if err1 == nil {
			p, _ := core.WorldMgrObj.GetPlayer(pid.(int32))
			p.UpdatePos(msg.X, msg.Y, msg.Z, msg.V)
		} else {
			logger.Error(err1)
			request.Fconn.LostConnection()
		}

	} else {
		logger.Error(err)
		request.Fconn.LostConnection()
	}
}

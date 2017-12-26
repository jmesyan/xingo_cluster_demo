package net_server

import (
	"fmt"
	"github.com/golang/protobuf/proto"
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
	onegate := GetRandomGate()

	if onegate != nil {
		logger.Info("chose gate: " + onegate.GetName())
		response, err := onegate.CallChildForResult("CreatePlayer", utils.GlobalObject.Name)
		if err == nil {
			self, _ := response.Result["p"].(core.Player)
			pid := self.Pid
			self.Fconn = fconn
			NetPlayers[pid] = self
			if pid > 0 {
				logger.Info("get pid success, pid:", pid)
				fconn.SetProperty("pid", pid)
				//同步Pid
				msg := &pb.SyncPid{
					Pid: pid,
				}
				SendMsg(fconn, 1, msg)
				position := &pb.Position{
					X: self.X,
					Y: self.Y,
					Z: self.Z,
					V: self.V,
				}

				//出现在自己的视野中
				data := &pb.BroadCast{
					Pid: pid,
					Tp:  2,
					Data: &pb.BroadCast_P{
						P: position,
					},
				}

				SendMsg(fconn, 200, data)

				//同步周围玩家
				sr, _ := response.Result["sr"].([]core.Player)
				for _, spy := range sr {
					msg2 := &pb.SyncPlayers{}
					p := &pb.Player{
						Pid: spy.Pid,
						P: &pb.Position{
							X: spy.X,
							Y: spy.Y,
							Z: spy.Z,
							V: spy.V,
						},
					}

					msg2.Ps = append(msg2.Ps, p)
					SendMsg(fconn, 202, msg2)
				}

				diff := time.Now().Sub(st).Nanoseconds()
				logger.Info("get pid total consume:", (diff / 1e6), "ms")

			} else {
				logger.Info("no game server serve")
				fconn.LostConnection()
			}
		} else {
			logger.Error(err)
		}
	}

}

type NetApiRouter struct {
}

func (this *NetApiRouter) Api_0(request *fnet.PkgAll) {
	logger.Debug("call Api_0")
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
		pid, err1 := request.Fconn.GetProperty("pid")
		if err1 == nil {
			onegate := GetRandomGate()

			if onegate != nil {
				logger.Info("chose gate: " + onegate.GetName())
				onegate.CallChildNotForResult("BroadCastMsg", pid, msg.Content)
			}
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
			onegate := GetRandomGate()
			if onegate != nil {
				logger.Info("chose gate: " + onegate.GetName())
				onegate.CallChildNotForResult("UpdatePos", pid, msg)
			}
		} else {
			logger.Error(err1)
			request.Fconn.LostConnection()
		}

	} else {
		logger.Error(err)
		request.Fconn.LostConnection()
	}
}

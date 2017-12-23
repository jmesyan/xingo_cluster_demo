package net_server

import (
	"github.com/jmesyan/xingo/clusterserver"
	"github.com/jmesyan/xingo/iface"
	"github.com/jmesyan/xingo/logger"
	"github.com/jmesyan/xingo/utils"
)

func DoConnectioned(fconn iface.Iconnection) {
	logger.Info("connection connect , I get it")
	//请求pid
	onegate := clusterserver.GlobalClusterServer.RemoteNodesMgr.GetRandomChild("gate")

	if onegate != nil {
		logger.Info("chose gate: " + onegate.GetName())
		response, err := onegate.CallChildForResult("CreatePlayer")
		if err == nil {
			pid, _ := response.Result["pid"].(int32)
			if pid > 0 {
				logger.Info("get pid success")
				fconn.SetProperty("pid", pid)
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

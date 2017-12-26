package main

import (
	"github.com/jmesyan/xingo"
	"github.com/jmesyan/xingo/utils"
	_ "net/http"
	_ "net/http/pprof"
	"os"
	"path/filepath"
	"strings"
	"xingo_cluster_demo/core"
	"xingo_cluster_demo/game_server"
	"xingo_cluster_demo/gate_server"
	"xingo_cluster_demo/net_server"
)

func main() {

	//server code
	args := os.Args
	dir, err := filepath.Abs(filepath.Dir("."))
	if err == nil {
		if true {
			sname := args[1]
			s := xingo.NewXingoCluterServer(sname, filepath.Join(dir, "conf", "clusterconf.json"))
			/*
				注册分布式服务器
			*/
			// //net server
			s.AddModule("net", &net_server.NetApiRouter{}, nil, &net_server.NetRpcApi{})
			// //gate server
			s.AddModule("gate", nil, nil, &gate_server.GateRpcApi{})
			// //admin server
			s.AddModule("game", nil, nil, &game_server.GameRpcApi{})

			if strings.HasPrefix(sname, "game") {
				core.WorldMgrObjInit()
			}

			if strings.HasPrefix(sname, "net") {
				utils.GlobalObject.OnConnectioned = net_server.DoConnectioned
			}
			s.StartClusterServer()
		} else {
			s := xingo.NewXingoCluterServer(args[1], filepath.Join(dir, "conf", "clusterconf_测试网关有root和http.json"))
			/*
				注册分布式服务器
			*/
			// //net server
			// s.AddModule("net", &net_server.TestNetApi2{}, &net_server.TestNetHttp{}, &net_server.TestNetRpc{})
			// //game server
			// s.AddModule("game", nil, nil, &game_server.TestGameRpc{})

			s.StartClusterServer()
		}

	}
}

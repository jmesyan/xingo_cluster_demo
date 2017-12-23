package main

import (
	"github.com/jmesyan/xingo"
	_ "net/http"
	_ "net/http/pprof"
	"os"
	"path/filepath"
	// "xingo_cluster_demo/admin_server"
	"xingo_cluster_demo/game_server"
	"xingo_cluster_demo/gate_server"
	_ "xingo_cluster_demo/net_server"
)

func main() {
	//pprof
	//go func() {
	//	println(http.ListenAndServe("localhost:6060", nil))
	//}()

	//server code
	args := os.Args
	dir, err := filepath.Abs(filepath.Dir("."))
	if err == nil {
		if true {
			s := xingo.NewXingoCluterServer(args[1], filepath.Join(dir, "conf", "clusterconf.json"))
			/*
				注册分布式服务器
			*/
			// //net server
			// s.AddModule("net", &net_server.TestNetApi{}, nil, &net_server.TestNetRpc{})
			// //gate server
			s.AddModule("gate", nil, nil, &gate_server.GateRpcApi{})
			// //admin server
			s.AddModule("game", nil, nil, &game_server.GameRpcApi{})

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

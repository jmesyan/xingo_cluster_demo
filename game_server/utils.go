package game_server

import (
	"github.com/jmesyan/xingo/cluster"
	"xingo_cluster_demo/core"
)

func SyncSurrounds(gate *cluster.Child, cmd int32, py, np core.Player) {
	gate.CallChildNotForResult("SyncSurrounds", cmd, py, np)
}

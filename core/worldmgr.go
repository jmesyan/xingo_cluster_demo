package core

import (
	"errors"
	"github.com/golang/protobuf/proto"
	"github.com/jmesyan/xingo/logger"
	// "github.com/jmesyan/xingo/utils"
	"sync"
	// "xingo_demo/pb"
	"github.com/jmesyan/xingo/clusterserver"
)

type WorldMgr struct {
	PlayerNumGen int32
	Players      map[int32]*Player
	AoiObj1      *AOIMgr //地图1
	sync.RWMutex
}

var WorldMgrObj *WorldMgr

func WorldMgrObjInit() {
	WorldMgrObj = &WorldMgr{
		PlayerNumGen: 0,
		Players:      make(map[int32]*Player),
		AoiObj1:      NewAOIMgr(85, 410, 75, 400, 10, 20),
	}
}

func (this *WorldMgr) AddPlayer(netname string) (*Player, error) {
	this.Lock()
	this.PlayerNumGen += 1
	p := NewPlayer(netname, this.PlayerNumGen)
	this.Players[p.Pid] = p
	this.Unlock()
	this.AoiObj1.Add2AOI(p)
	// //同步Pid
	// msg := &pb.SyncPid{
	// 	Pid: p.Pid,
	// }
	// p.SendMsg(1, msg)
	// //加到aoi
	//
	// //周围的人
	// p.SyncSurrouding()
	return p, nil
}

func (this *WorldMgr) RemovePlayer(pid int32) {
	this.Lock()
	defer this.Unlock()
	//从aoi移除
	this.AoiObj1.LeaveAOI(this.Players[pid])
	delete(this.Players, pid)
}

func (this *WorldMgr) Move(p *Player) {
	onegate := clusterserver.GlobalClusterServer.RemoteNodesMgr.GetRandomChild("gate")
	/*aoi*/
	pids, err := this.AoiObj1.GetSurroundingPids(p)
	if err == nil {
		for _, pid := range pids {
			player, err1 := this.GetPlayer(pid)
			if err1 == nil {
				go SyncSurrounds(onegate, 211, *player, *p)
			}
		}
	}
}

func (this *WorldMgr) SendMsgByPid(pid int32, msgId uint32, data proto.Message) {
	p, err := this.GetPlayer(pid)
	if err == nil {
		p.SendMsg(msgId, data)
	}
}

func (this *WorldMgr) GetPlayer(pid int32) (*Player, error) {
	this.RLock()
	defer this.RUnlock()
	p, ok := this.Players[pid]
	if ok {
		return p, nil
	} else {
		return nil, errors.New("no player in the world!!!")
	}
}

func (this *WorldMgr) Broadcast(msgId uint32, data proto.Message) {
	this.RLock()
	defer this.RUnlock()
	for _, p := range this.Players {
		p.SendMsg(msgId, data)
	}
}

func (this *WorldMgr) BroadcastBuff(msgId uint32, data proto.Message) {
	this.RLock()
	defer this.RUnlock()
	for _, p := range this.Players {
		p.SendBuffMsg(msgId, data)
	}
}

func (this *WorldMgr) AOIBroadcast(p *Player, msgId uint32, data proto.Message) {
	/*aoi*/
	pids, err := WorldMgrObj.AoiObj1.GetSurroundingPids(p)
	if err == nil {
		for _, pid := range pids {
			player, err1 := WorldMgrObj.GetPlayer(pid)
			if err1 == nil {
				player.SendMsg(msgId, data)
			}
		}
	} else {
		logger.Error(err)
	}
}

package core

import (
	"encoding/gob"
	"errors"
	"github.com/golang/protobuf/proto"
	"github.com/jmesyan/xingo/cluster"
	"github.com/jmesyan/xingo/clusterserver"
	"github.com/jmesyan/xingo/iface"
	"github.com/jmesyan/xingo/logger"
	"math/rand"
	"xingo_demo/pb"
)

func init() {
	gob.Register(Player{})
	gob.Register([]Player{})
}

type Player struct {
	Net   string
	Fconn iface.Iconnection
	Pid   int32
	X     float32 //平面x
	Y     float32 //高度
	Z     float32 //平面y!!!!!注意不是Y
	V     float32 //旋转0-360度
}

func NewPlayer(net string, pid int32) *Player {
	p := &Player{
		Net:   net,
		Pid:   pid,
		X:     float32(rand.Intn(10) + 160),
		Y:     0,
		Z:     float32(rand.Intn(17) + 134),
		V:     0,
		Fconn: nil,
	}

	return p
}

func SyncSurrounds(gate *cluster.Child, cmd int32, py, np Player) {
	gate.CallChildNotForResult("SyncSurrounds", cmd, py, np)
}

/*
同步周围玩家
*/
func (this *Player) SyncSurrouding() {
	pids, err := WorldMgrObj.AoiObj1.GetSurroundingPids(this)

	if err == nil {
		msg := &pb.SyncPlayers{}
		for _, pid := range pids {
			player, err1 := WorldMgrObj.GetPlayer(pid)
			if err1 == nil {
				p := &pb.Player{
					Pid: pid,
					P: &pb.Position{
						X: player.X,
						Y: player.Y,
						Z: player.Z,
						V: player.V,
					},
				}
				msg.Ps = append(msg.Ps, p)
				//出现在周围人的视野
				data := &pb.BroadCast{
					Pid: this.Pid,
					Tp:  2,
					Data: &pb.BroadCast_P{
						P: &pb.Position{
							X: this.X,
							Y: this.Y,
							Z: this.Z,
							V: this.V,
						},
					},
				}
				player.SendMsg(200, data)
			}
		}
		//分包发送
		per := 20
		ps := msg.Ps
		for i := 0; ; i++ {
			if i*per > len(ps)-1 {
				break
			}
			if i*per+per > len(ps)-1 {
				msg.Ps = ps[i*per:]
			} else {
				msg.Ps = ps[i*per : i*per+per]
			}
			this.SendMsg(202, msg)
		}
		//this.SendMsg(202, msg)
	} else {
		logger.Error(err)
	}

}

func (this *Player) UpdatePos(x float32, y float32, z float32, v float32) {
	oldGridId := WorldMgrObj.AoiObj1.GetGridIDByPos(this.X, this.Z)
	//更新位置的时候判断是否需要更新gridID
	newGridId := WorldMgrObj.AoiObj1.GetGridIDByPos(x, z)

	if newGridId < 0 || newGridId >= WorldMgrObj.AoiObj1.lenX*WorldMgrObj.AoiObj1.lenY {
		//更新的坐标有误直接返回
		return
	}
	//更新
	this.X = x
	this.Y = y
	this.Z = z
	this.V = v

	if oldGridId != newGridId {
		WorldMgrObj.AoiObj1.LeaveAOIFromGrid(this, oldGridId)
		WorldMgrObj.AoiObj1.Add2AOI(this)
		//需要处理老的aoi消失和新的aoi出生
		this.OnExchangeAoiGrid(oldGridId, newGridId)
	}
	WorldMgrObj.Move(this)
}

func (this *Player) OnExchangeAoiGrid(oldGridId int32, newGridId int32) error {
	oldAoiGrids, err1 := WorldMgrObj.AoiObj1.GetSurroundingByGridId(oldGridId)
	newAoiGrids, err2 := WorldMgrObj.AoiObj1.GetSurroundingByGridId(newGridId)
	if err1 != nil || err2 != nil {
		logger.Error(err1, err2)
		return errors.New("OnExchangeAoiGrid")
	}
	alls := make([]*Grid, 0)
	alls = append(alls, oldAoiGrids...)
	alls = append(alls, newAoiGrids...)
	//并集
	union := make(map[int32]*Grid, 0)
	for _, v := range alls {
		if _, ok := union[v.ID]; ok != true {
			union[v.ID] = v
		}
	}
	oldAoiGridsMap := make(map[int32]bool, 0)
	for _, oldGrid := range oldAoiGrids {
		if _, ok := oldAoiGridsMap[oldGrid.ID]; ok != true {
			oldAoiGridsMap[oldGrid.ID] = true
		}
	}

	newAoiGridsMap := make(map[int32]bool, 0)
	for _, newGrid := range newAoiGrids {
		if _, ok := newAoiGridsMap[newGrid.ID]; ok != true {
			newAoiGridsMap[newGrid.ID] = true
		}
	}

	onegate := clusterserver.GlobalClusterServer.RemoteNodesMgr.GetRandomChild("gate")
	if onegate == nil {
		return errors.New("can not found gate")
	}

	for gid, grid := range union {
		//出生
		if _, ok := oldAoiGridsMap[gid]; ok != true {
			for _, pid := range grid.GetPids() {
				if pid != this.Pid {
					p, err := WorldMgrObj.GetPlayer(pid)
					if err == nil {
						go SyncSurrounds(onegate, 200, *p, *this)
						go SyncSurrounds(onegate, 200, *this, *p)
					}
				}

			}
		}
		if _, ok := newAoiGridsMap[gid]; ok != true {
			//消失
			for _, pid := range grid.GetPids() {
				if pid != this.Pid {
					p, err := WorldMgrObj.GetPlayer(pid)
					if err == nil {
						go SyncSurrounds(onegate, 201, *p, *this)
						go SyncSurrounds(onegate, 201, *this, *p)
					}
				}
			}
		}
	}
	return nil
}

func (this *Player) Talk(content string) {
	data := &pb.BroadCast{
		Pid: this.Pid,
		Tp:  1,
		Data: &pb.BroadCast_Content{
			Content: content,
		},
	}

	WorldMgrObj.BroadcastBuff(200, data)
	//WorldMgrObj.Broadcast(200, data)
}

func (this *Player) LostConnection() {
	msg := &pb.SyncPid{
		Pid: this.Pid,
	}
	WorldMgrObj.Broadcast(201, msg)
}

func (this *Player) SendMsg(msgId uint32, data proto.Message) {
	// if this.Fconn != nil {
	// 	packdata, err := utils.GlobalObject.Protoc.GetDataPack().Pack(msgId, data)
	// 	if err == nil {
	// 		this.Fconn.Send(packdata)
	// 	} else {
	// 		logger.Error("pack data error")
	// 	}
	// }
}

func (this *Player) SendBuffMsg(msgId uint32, data proto.Message) {
	// if this.Fconn != nil {
	// 	packdata, err := utils.GlobalObject.Protoc.GetDataPack().Pack(msgId, data)
	// 	if err == nil {
	// 		this.Fconn.SendBuff(packdata)
	// 	} else {
	// 		logger.Error("pack data error")
	// 	}
	// }
}

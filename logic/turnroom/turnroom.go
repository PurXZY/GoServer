package turnroom

import (
	"GoServer/i"
	"GoServer/log"
	"GoServer/proto/usercmd"
	"github.com/golang/protobuf/proto"
)

type TurnRoom struct {
	TurnLogic
	uniqId uint32
	owner  i.IRoomPlayer
}

func NewTurnRoom(uniqId uint32, owner i.IRoomPlayer) *TurnRoom {
	room := &TurnRoom{
		uniqId: uniqId,
		owner:  owner,
	}
	room.init()
	return room
}

func (tr *TurnRoom) init() {
	log.Info.Printf("new TurnRoom id:%v, owner:%v", tr.uniqId, tr.owner.GetName())
	tr.notifyPlayerIntoRoom()
	tr.TurnLogic.Init(tr)
}

func (tr *TurnRoom) notifyPlayerIntoRoom() {
	msg := usercmd.IntoRoomS2CMsg{
		RoomId: 7777,
	}
	tr.BroadcastMsg(usercmd.CmdType_IntoRoomRes, &msg)
}

func (tr *TurnRoom) BroadcastMsg(cmdType usercmd.CmdType, msg proto.Message) {
	tr.owner.SendMsg(cmdType, msg)
}
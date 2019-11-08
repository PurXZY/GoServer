package logic

import (
	"GoServer/i"
	"GoServer/log"
	"GoServer/logic/turnroom"
	"GoServer/proto/usercmd"
	"GoServer/util"
	"github.com/golang/protobuf/proto"
)

type Avatar struct {
	name string
	task i.ITcpTask
}

func NewAvatar() *Avatar {
	return &Avatar{}
}

func (a *Avatar) SetTcpTask(t i.ITcpTask) {
	a.task = t
}

func (a *Avatar) GetName() string {
	return a.name
}

func (a *Avatar) SendMsg(msgType usercmd.CmdType, msg proto.Message) {
	data, err := util.EncodeCmd(uint16(msgType), msg)
	if err != nil {
		return
	}
	a.task.SendData(data)
}

func (a *Avatar) ParseMsg(data []byte) {
	cmdType := usercmd.CmdType(util.GetCmdType(data))
	log.Info.Printf("ParseMsg %T %+v", cmdType, cmdType)
	switch cmdType {
	case usercmd.CmdType_LoginReq:
		recvCmd, ok := util.DecodeCmd(data, &usercmd.LoginC2SMsg{}).(*usercmd.LoginC2SMsg)
		if !ok {
			log.Error.Println("decode cmd fail cmdType:", cmdType)
			return
		}
		a.name = recvCmd.GetName()
		log.Debug.Println("recv msg UserCmd_LoginReq name:", recvCmd.GetName())
		msg := usercmd.LoginS2CMsg{
			PlayerId: 6666,
		}
		a.SendMsg(usercmd.CmdType_LoginRes, &msg)
	case usercmd.CmdType_IntoRoomReq:
		a.ReqIntoRoom()
	default:
		log.Error.Println("unknown cmdType:", cmdType)
		return
	}
}

func (a *Avatar) ReqIntoRoom() {
	turnroom.NewTurnRoom(7777, a)
}
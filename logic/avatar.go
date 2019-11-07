package logic

import (
	"GoServer/i"
	"GoServer/log"
	"GoServer/proto/usercmd"
	"GoServer/util"
	"github.com/golang/protobuf/proto"
)

type Avatar struct {
	task i.ITcpTask
}

func NewAvatar() *Avatar {
	return &Avatar{}
}

func (a *Avatar) SetTcpTask(t i.ITcpTask) {
	a.task = t
}

func (a *Avatar) SendMsg(msgType usercmd.UserCmd, msg proto.Message) {
	data, err := util.EncodeCmd(uint16(msgType), msg)
	if err != nil {
		return
	}
	a.task.SendData(data)
}

func (a *Avatar) ParseMsg(data []byte) {
	log.Info.Println("ParseMsg")
	cmdType := usercmd.UserCmd(util.GetCmdType(data))
	switch cmdType {
	case usercmd.UserCmd_LoginReq:
		recvCmd, ok := util.DecodeCmd(data, &usercmd.LoginC2SMsg{}).(*usercmd.LoginC2SMsg)
		if !ok {
			log.Error.Println("decode cmd fail cmdType:", cmdType)
			return
		}
		log.Debug.Println("recv msg UserCmd_LoginReq name:", recvCmd.GetName())
		msg := usercmd.LoginS2CMsg{
			PlayerId: 6666,
		}
		a.SendMsg(usercmd.UserCmd_LoginRes, &msg)
	case usercmd.UserCmd_IntoRoomReq:
		a.ReqIntoRoom()
	default:
		log.Error.Println("unknown cmdType:", cmdType)
		return
	}
}

func (a *Avatar) ReqIntoRoom() {
	msg := usercmd.IntoRoomS2CMsg{
		RoomId: 7777,
	}
	a.SendMsg(usercmd.UserCmd_IntoRoomRes, &msg)
}
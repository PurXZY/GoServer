package i

import (
	"GoServer/proto/usercmd"
	"github.com/golang/protobuf/proto"
)

type ITaskOwner interface {
	SetTcpTask(ITcpTask)
	ParseMsg([]byte)
	SendMsg(usercmd.CmdType, proto.Message)
}

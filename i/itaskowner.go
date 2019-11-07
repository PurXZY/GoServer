package i

type ITaskOwner interface {
	SetTcpTask(ITcpTask)
	ParseMsg([]byte)
}

package i

type ITaskOwner interface {
	ParseMsg([]byte)
}

package gamenet

import (
	"GoServer/log"
	"net"
)

const (
	osNetWriteBufferSize = 64 * 1024
	osNetReadBufferDataSize = 64 * 1024
)

type TcpTask struct {
	conn *net.TCPConn
}

func NewTcpTask(conn *net.TCPConn) *TcpTask {
	log.Debug.Println("new connection: ", conn.RemoteAddr())
	_ = conn.SetKeepAlive(true)
	_ = conn.SetNoDelay(true)
	_ = conn.SetWriteBuffer(osNetWriteBufferSize)
	_ = conn.SetReadBuffer(osNetReadBufferDataSize)
	t := &TcpTask{
		conn: conn,
	}
	t.init()
	return t
}

func (t *TcpTask) init() {
	log.Debug.Println("start TcpTask addr: ", t.conn.RemoteAddr())
	go t.sendLoop()
	go t.recvLoop()
}

func (t *TcpTask) sendLoop() {

}

func (t *TcpTask) recvLoop() {

}
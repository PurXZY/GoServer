package gamenet

import (
	"GoServer/i"
	"GoServer/log"
	"GoServer/logic"
	"GoServer/util"
	"io"
	"net"
	"runtime/debug"
	"sync"
	"sync/atomic"
)

const (
	osNetWriteBufferSize    = 64 * 1024
	osNetReadBufferDataSize = 64 * 1024
	sendDataMaxSize         = 64 * 1024
	cmdDataHeadSize         = 4
)

type TcpTask struct {
	owner      i.ITaskOwner
	conn       *net.TCPConn
	sendBuff   *ByteBuffer
	sendMutex  sync.Mutex
	sendSignal chan bool
	isClosed   int32
}

func NewTcpTask(conn *net.TCPConn) *TcpTask {
	log.Debug.Println("new connection: ", conn.RemoteAddr())
	_ = conn.SetKeepAlive(true)
	_ = conn.SetNoDelay(true)
	_ = conn.SetWriteBuffer(osNetWriteBufferSize)
	_ = conn.SetReadBuffer(osNetReadBufferDataSize)
	t := &TcpTask{
		conn:       conn,
		sendBuff:   NewByteBuffer(),
		sendSignal: make(chan bool, 1),
		isClosed:   0,
	}
	t.owner = logic.NewAvatar()
	t.init()
	t.owner.SetTcpTask(t)
	return t
}

func (t *TcpTask) IsClosed() bool {
	return atomic.LoadInt32(&t.isClosed) != 0
}

func (t *TcpTask) Close() {
	if !atomic.CompareAndSwapInt32(&t.isClosed, 0, 1) {
		return
	}
	log.Debug.Println("close TcpTask addr:", t.conn.RemoteAddr())
	_ = t.conn.Close()
	t.sendSignal <- true
	close(t.sendSignal)
}

func (t *TcpTask) init() {
	log.Debug.Println("start TcpTask addr: ", t.conn.RemoteAddr())
	go t.sendLoop()
	go t.recvLoop()
}

func (t *TcpTask) SendData(data []byte) {
	if t.IsClosed() {
		log.Error.Println("TcpTask isClosed")
		return
	}
	dataSize := len(data) + cmdDataHeadSize
	if dataSize == 0 {
		log.Warn.Println("send usercmd size:", dataSize)
		return
	}
	t.sendMutex.Lock()
	curBuffSize := t.sendBuff.RdSize()
	if curBuffSize+dataSize > sendDataMaxSize {
		log.Error.Printf("send buff over limit cur:%d new:%d", curBuffSize, dataSize)
		t.Close()
	}
	t.sendBuff.Append(util.Uint32ToBytes(uint32(dataSize)))
	t.sendBuff.Append(data)
	t.sendMutex.Unlock()
	t.sendSignal <- false
}

func (t *TcpTask) sendLoop() {
	defer func() {
		if err := recover(); err != nil {
			log.Error.Println("err: ", err, "\n", string(debug.Stack()))
		}
	}()

	defer t.Close()

	var (
		tmpByte = NewByteBuffer()
		sendNum int
		err     error
	)

	for needClose := range t.sendSignal {
		if needClose {
			return
		} else {
			t.sendMutex.Lock()
			if t.sendBuff.RdReady() {
				// 发送数据由发送缓冲区移动到发送协程中
				tmpByte.Append(t.sendBuff.RdBuf()[:t.sendBuff.RdSize()])
				t.sendBuff.Reset()
			}
			t.sendMutex.Unlock()
			if !tmpByte.RdReady() {
				continue
			}
			for tmpByte.RdReady() {
				// 发送完整
				sendNum, err = t.conn.Write(tmpByte.RdBuf()[:tmpByte.RdSize()])
				if err != nil {
					log.Error.Println("send loop addr:", t.conn.RemoteAddr(), ", err:", err)
					return
				}
				tmpByte.RdFlip(sendNum)
			}
		}
	}
}

func (t *TcpTask) recvLoop() {
	defer func() {
		if err := recover(); err != nil {
			log.Error.Println("err:", err, "\n", string(debug.Stack()))
		}
	}()

	defer t.Close()

	var (
		recvBuff  = NewByteBuffer()
		totalSize int
		dataSize  int
		needNum   int
		readNum   int
		err       error
		msgBuff   []byte
	)

	for {
		totalSize = recvBuff.RdSize()
		for totalSize < cmdDataHeadSize {
			needNum = cmdDataHeadSize - totalSize
			readNum, err = t.conn.Read(recvBuff.WrBuf())
			if err == io.EOF {
				log.Debug.Println("remote close io eof addr:", t.conn.RemoteAddr())
				return
			} else if err != nil {
				log.Error.Printf("recv loop addr:%s, err:%T %+v", t.conn.RemoteAddr(), err, err)
				return
			}
			recvBuff.WrFlip(readNum)
			totalSize = recvBuff.RdSize()
		}
		msgBuff = recvBuff.RdBuf()
		dataSize = int(util.BytesToUint32(msgBuff[:cmdDataHeadSize]))
		if dataSize > sendDataMaxSize {
			log.Error.Println("recv too big usercmd over limit size:", dataSize)
			return
		}
		for totalSize < dataSize {
			needNum = dataSize - totalSize
			if recvBuff.WrSize() < needNum {
				recvBuff.WrGrow(needNum)
			}
			readNum, err = t.conn.Read(recvBuff.WrBuf())
			if err != nil {
				log.Error.Println("recv loop addr:", t.conn.RemoteAddr(), ", err:", err)
				return
			}
			recvBuff.WrFlip(readNum)
			totalSize = recvBuff.RdSize()
		}
		msgBuff = recvBuff.RdBuf()
		t.owner.ParseMsg(msgBuff[cmdDataHeadSize:dataSize])
		recvBuff.RdFlip(dataSize)
	}
}
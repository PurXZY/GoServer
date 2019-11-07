package gamenet

import (
	"GoServer/log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

type Server struct {
	listener *net.TCPListener
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Start() {
	defer s.final()

	s.dealWithSignal()
	var address string = ":8888"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", address)
	if err != nil {
		log.Error.Println("bind fail ", address)
		return
	}
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Error.Println("listen fail ", address)
		return
	}
	log.Info.Println("bind success ", address)
	s.listener = listener
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			if _, ok := err.(*net.OpError); !ok {
				log.Error.Printf("err:%v", err)
			}
			return
		}
		NewTcpTask(conn)
	}
}

func (s *Server) dealWithSignal() {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGABRT, syscall.SIGTERM, syscall.SIGPIPE, syscall.SIGHUP)
	go func() {
		for sig := range ch {
			log.Info.Println("receive signal: ", sig)
			switch sig {
			case syscall.SIGPIPE:
				log.Error.Println("SIGPIPE")
			default:
				s.Terminate()
			}
		}
	}()
}

func (s *Server) Terminate() {
	_ = s.listener.Close()
}

func (s *Server) final() {
	log.Info.Println("server final")
}
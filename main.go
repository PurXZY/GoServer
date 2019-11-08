package main

import (
	"GoServer/gamenet"
	"GoServer/log"
	"GoServer/logic/datamgr"
	"os"
)

func main() {
	log.InitLog(os.Stdout, os.Stdout, os.Stdout)

	datamgr.GetMe().Init()

	s := gamenet.NewServer()
	s.Start()
}

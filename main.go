package main

import (
	"GoServer/gamenet"
	"GoServer/log"
	"os"
)

func main() {
	log.InitLog(os.Stdout, os.Stdout, os.Stdout)
	s := gamenet.NewServer()
	s.Start()
}

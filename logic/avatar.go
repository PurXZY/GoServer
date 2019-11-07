package logic

import "GoServer/log"

type Avatar struct {

}

func NewAvatar() *Avatar {
	return &Avatar{}
}

func (a *Avatar) ParseMsg(data []byte) {
	log.Info.Println("ParseMsg")
}

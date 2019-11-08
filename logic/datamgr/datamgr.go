package datamgr

import (
	"GoServer/log"
	"encoding/json"
	"io/ioutil"
	"sync"
)

type DataMgr struct {
	propData map[uint32]map[string]uint32
}


var (
	mgr     *DataMgr
	mgrOnce sync.Once
)

func GetMe() *DataMgr {
	if mgr == nil {
		mgrOnce.Do(func() {
			mgr = &DataMgr{}
		})
	}
	return mgr
}

func (dm *DataMgr) Init() {
	dm.initPropData()
}

func (dm *DataMgr) initPropData() {
	dm.propData = make(map[uint32]map[string]uint32)

	fileData, err := ioutil.ReadFile("D:\\GoPath\\src\\jsondata\\prop.json")
	if err != nil {
		log.Error.Println("initPropData err:", err)
		return
	}

	err = json.Unmarshal(fileData, &dm.propData)
	if err != nil {
		log.Error.Println("initPropData err:", err)
		return
	}
	log.Info.Println("initPropData success")
}

func (dm *DataMgr) GetPropData(u uint32) map[string]uint32 {
	if data, ok := dm.propData[u]; ok {
		return data
	}
	log.Error.Println("GetPropData wrong u:", u)
	return nil
}
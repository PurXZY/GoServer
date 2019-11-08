package battleentity

import (
	"GoServer/proto/usercmd"
	"sort"
)

type TurnEntityMgr struct {
	battleEntities map[uint32]*BattleEntity
}

func NewTurnEntityMgr() *TurnEntityMgr {
	tem := &TurnEntityMgr{
		battleEntities: make(map[uint32]*BattleEntity),
	}
	return tem
}

func (tem *TurnEntityMgr) CreateAllBattleEntities() {
	enemyData := map[usercmd.PosIndex]uint32{
		usercmd.PosIndex_PosELeft:   1,
		usercmd.PosIndex_PosECenter: 1,
		usercmd.PosIndex_PosERight:  2,
	}
	for posIndex, entityType := range enemyData {
		tem.addNewBattleEntity(uint32(posIndex), entityType)
	}

	myData := map[usercmd.PosIndex]uint32{
		usercmd.PosIndex_PosBCenter: 3,
	}
	for posIndex, entityType := range myData {
		tem.addNewBattleEntity(uint32(posIndex), entityType)
	}
}

func (tem *TurnEntityMgr) addNewBattleEntity(posIndex uint32, entityType uint32) {
	e := NewBattleEntity(posIndex, entityType)
	tem.battleEntities[posIndex] = e
}

func (tem *TurnEntityMgr) GetAllBattleEntities() []*BattleEntity {
	var v []*BattleEntity
	for _, value := range tem.battleEntities {
		v = append(v, value)
	}
	return v
}

func (tem *TurnEntityMgr) GetEntityById(posIndex uint32) *BattleEntity {
	return tem.battleEntities[posIndex]
}

func (tem *TurnEntityMgr) GetSortedEntitiesByMoveSpeed() []uint32 {
	type EntitySpeed struct {
		PosIndex  uint32
		MoveSpeed uint32
	}
	var sortList []EntitySpeed
	for posIndex, entity := range tem.battleEntities {
		sortList = append(sortList, EntitySpeed{posIndex, entity.MoveSpeed})
	}
	sort.Slice(sortList, func(i, j int) bool {
		return sortList[i].MoveSpeed > sortList[j].MoveSpeed
	})
	var ret []uint32
	for _, v := range sortList {
		ret = append(ret, v.PosIndex)
	}
	return ret
}

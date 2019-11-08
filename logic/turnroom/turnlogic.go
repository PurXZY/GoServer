package turnroom

import (
	"GoServer/log"
	"GoServer/logic/turnroom/battleentity"
	"GoServer/proto/usercmd"
)

type TurnLogic struct {
	tem            *battleentity.TurnEntityMgr
	turnRoom       *TurnRoom
	curBigTurn     uint32
	curSmallTurn   uint32
	entityTurnInfo []uint32
}

func (tl *TurnLogic) Init(room *TurnRoom) {
	tl.turnRoom = room
	tl.tem = battleentity.NewTurnEntityMgr()
	tl.tem.CreateAllBattleEntities()
	tl.SyncAllBattleEntities()
	tl.beginFirstTurn()
}

func (tl *TurnLogic) beginFirstTurn() {
	tl.curBigTurn = 1
	tl.curSmallTurn = 0
	tl.sortBattleEntityMoveSpeed()
	tl.SyncTurnInfo()
}

func (tl *TurnLogic) curTurnEntityPosIndex() uint32 {
	return tl.entityTurnInfo[tl.curSmallTurn]
}

func (tl *TurnLogic) curTurnEntity() *battleentity.BattleEntity {
	return tl.tem.GetEntityById(tl.curTurnEntityPosIndex())
}

func (tl *TurnLogic) sortBattleEntityMoveSpeed() {
	tl.entityTurnInfo = tl.entityTurnInfo[:0]
	tl.entityTurnInfo = tl.tem.GetSortedEntitiesByMoveSpeed()
}

// ------------------- sync -------------------

func (tl *TurnLogic) SyncAllBattleEntities() {
	msg := usercmd.SyncAllBattleEntitiesS2CMsg{}
	for _, entity := range tl.tem.GetAllBattleEntities() {
		msg.Entities = append(msg.Entities, entity.GetPropData())
	}
	tl.turnRoom.BroadcastMsg(usercmd.CmdType_SyncAllBattleEntities, &msg)
}

func (tl *TurnLogic) SyncTurnInfo() {
	log.Info.Printf("New Turn Big:%v Small:%v Pos:%+v", tl.curBigTurn, tl.curSmallTurn, tl.curTurnEntityPosIndex())
	msg := usercmd.TurnInfoS2CMsg{
		BigTurnIndex:      tl.curBigTurn,
		SmallTurnIndex:    tl.curSmallTurn,
		CurEntityPosIndex: tl.curTurnEntityPosIndex(),
		SkillSet:          tl.curTurnEntity().Skills,
	}
	tl.turnRoom.BroadcastMsg(usercmd.CmdType_TurnInfo, &msg)
}

package battleentity

import (
	"GoServer/logic/datamgr"
	"GoServer/proto/usercmd"
)

type EntityProp struct {
	entity         *BattleEntity
	Health         uint32
	PhysicalAttack uint32
	MagicAttack    uint32
	PhysicalDefend uint32
	MagicDefend    uint32
	MoveSpeed      uint32
}

func (ep *EntityProp) InitProp(e *BattleEntity, propId uint32) {
	ep.entity = e
	data := datamgr.GetMe().GetPropData(propId)
	if data == nil {
		return
	}
	ep.Health = data["Health"]
	ep.PhysicalAttack = data["PhysicalAttack"]
	ep.MagicAttack = data["MagicAttack"]
	ep.PhysicalAttack = data["PhysicalAttack"]
	ep.MagicDefend = data["MagicDefend"]
	ep.MoveSpeed = data["MoveSpeed"]
}

func (ep *EntityProp) GetPropData() *usercmd.BattleEntity {
	e := &usercmd.BattleEntity{
		PosIndex:       ep.entity.PosIndex,
		EntityType:     ep.entity.EntityType,
		Health:         ep.Health,
		PhysicalAttack: ep.PhysicalAttack,
		MagicAttack:    ep.MagicAttack,
		PhysicalDefend: ep.PhysicalDefend,
		MagicDefend:    ep.MagicDefend,
		MoveSpeed:      ep.MoveSpeed,
	}
	return e
}

package battleentity

import "GoServer/log"

type BattleEntity struct {
	EntityProp
	EntitySkill
	PosIndex   uint32
	EntityType uint32
}

func NewBattleEntity(PosIndex uint32, entityType uint32) *BattleEntity {
	entity := &BattleEntity{
		PosIndex:   PosIndex,
		EntityType: entityType,
	}
	entity.InitProp(entity, entityType)
	entity.InitSkill(entity)
	log.Debug.Printf("new BattleEntity %+v %+v %+v", entity.PosIndex, entity.EntityType, entity.EntityProp)
	return entity
}

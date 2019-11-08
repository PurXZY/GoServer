package battleentity

type EntitySkill struct {
	entity *BattleEntity
	Skills []uint32
}

func (es *EntitySkill) InitSkill(e *BattleEntity) {
	es.Skills = []uint32 {1, 2, 3}
}

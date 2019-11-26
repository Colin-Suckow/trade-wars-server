package ecs

import "math/rand"

type BasicEntity struct {
	id            uint64
	parentEntity  *BasicEntity
	childEntities []BasicEntity
}

func makeBasicEntity() BasicEntity {
	return BasicEntity{rand.Uint64(), nil, []BasicEntity{}}
}

func (e BasicEntity) addChild(child BasicEntity) {
	e.childEntities = append(e.childEntities, child)
	child.parentEntity = &e
}

package ecs

import "math/rand"

type BasicEntity struct {
	id            uint64
	parentEntity  *BasicEntity
	childEntities []BasicEntity
}

func makeBasicEntity() BasicEntity {
	return BasicEntity{rand.Uint64()}
}

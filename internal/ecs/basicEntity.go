package ecs

import "math/rand"

type BasicEntity struct {
	id uint64
}

func makeBasicEntity() BasicEntity {
	return BasicEntity{rand.Uint64()}
}

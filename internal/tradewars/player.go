package tradewars

import "github.com/EngoEngine/ecs"

type Player struct {
	ecs.BasicEntity
	PositionComponent
	OwnershipComponent
}

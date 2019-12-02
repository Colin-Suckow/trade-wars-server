package tradewars

import "github.com/EngoEngine/ecs"

import tradewars "suckow.dev/trade-wars-server/internal/tradewars/Components"

type Player struct {
	ecs.BasicEntity
	tradewars.PositionComponent
	tradewars.OwnershipComponent
}

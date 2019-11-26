package tradewars

import (
	"github.com/EngoEngine/ecs"
	tradewars "suckow.dev/trade-wars-server/internal/tradewars/Components"
)

//Individual Cargo Item

type CargoItem struct {
	ecs.BasicEntity
	tradewars.NameComponent
	tradewars.ValueComponent
}

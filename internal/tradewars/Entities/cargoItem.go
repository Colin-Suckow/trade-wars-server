package tradewars

import (
	"suckow.dev/trade-wars-server/internal/ecs"
	tradewars "suckow.dev/trade-wars-server/internal/tradewars/Components"
)

//Individual Cargo Item

type CargoItem struct {
	ecs.BasicEntity
	tradewars.NameComponent
	tradewars.ValueComponent
}

package tradewars

import "suckow.dev/trade-wars-server/internal/ecs"

import tradewars "suckow.dev/trade-wars-server/internal/tradewars/Components"

type Spacecraft struct {
	ecs.BasicEntity
	tradewars.NameComponent
	tradewars.HealthComponent
}

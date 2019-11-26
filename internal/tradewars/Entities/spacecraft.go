package tradewars

import "github.com/EngoEngine/ecs"

import tradewars "suckow.dev/trade-wars-server/internal/tradewars/Components"

type Spacecraft struct {
	ecs.BasicEntity
	tradewars.NameComponent
	tradewars.HealthComponent
}

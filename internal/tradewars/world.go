package tradewars

import (
	"github.com/EngoEngine/ecs"
	evbus "github.com/asaskevich/EventBus"
	tradewars "suckow.dev/trade-wars-server/internal/tradewars/Systems"
)

func InitializeWorld(bus evbus.EventBus) {
	world := ecs.World{}

	//Setup bus system so every other system can have access to main server bus
	world.AddSystem(&tradewars.BusSystem{Bus: &bus})
}

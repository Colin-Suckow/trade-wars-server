package tradewars

import (
	"github.com/EngoEngine/ecs"
	evbus "github.com/asaskevich/EventBus"
	components "suckow.dev/trade-wars-server/internal/tradewars/Components"
	entities "suckow.dev/trade-wars-server/internal/tradewars/Entities"
	tradewars "suckow.dev/trade-wars-server/internal/tradewars/Systems"
)

var MainWorld ecs.World

func InitializeWorld(bus evbus.Bus) {
	world := ecs.World{}

	//Setup bus system so every other system can have access to main server bus
	world.AddSystem(&tradewars.BusSystem{Bus: bus})

	world.AddSystem(&tradewars.MapSystem{Bus: bus})

	MainWorld = world

}

func NewPlayer(callsign string) ecs.BasicEntity {
	entity := ecs.NewBasic()
	player := entities.Player{entity, components.PositionComponent{0, 0, 0}, components.OwnershipComponent{callsign, entity.ID()}}

	for _, system := range MainWorld.Systems() {
		switch sys := system.(type) {

		case *tradewars.MapSystem:
			sys.Add(&player.BasicEntity, &player.PositionComponent)

		default:
			//Do nothing
		}

	}

	return entity

}

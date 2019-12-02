package tradewars

import (
	"github.com/EngoEngine/ecs"
	evbus "github.com/asaskevich/EventBus"
	components "suckow.dev/trade-wars-server/internal/tradewars/Components"
	entities "suckow.dev/trade-wars-server/internal/tradewars/Entities"
	systems "suckow.dev/trade-wars-server/internal/tradewars/Systems"
)

var MainWorld ecs.World
var MainBus evbus.Bus

func InitializeWorld(bus *evbus.Bus) {
	world := ecs.World{}

	//Setup bus system so every other system can have access to main server bus
	world.AddSystem(&systems.BusSystem{Bus: bus})

	mapSystem := &systems.MapSystem{Bus: bus}
	world.AddSystem(mapSystem)
	
	tradewars.MainBus.Subscribe("tradewars:position", mapSystem.BroadcastIndividualPosition)

	MainWorld = world
	MainBus = mBus

}

func NewPlayer(callsign string) ecs.BasicEntity {
	entity := ecs.NewBasic()
	player := entities.Player{entity, components.PositionComponent{0, 0, 0}, components.OwnershipComponent{callsign, entity.ID()}}

	for _, system := range MainWorld.Systems() {
		switch sys := system.(type) {

		case *systems.MapSystem:
			sys.Add(&player.BasicEntity, &player.PositionComponent)

		default:
			//Do nothing
		}

	}

	return entity

}

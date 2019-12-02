package tradewars

import (
	"github.com/EngoEngine/ecs"
	evbus "github.com/asaskevich/EventBus"
)

var MainWorld ecs.World
var MainBus evbus.Bus

func InitializeWorld(bus *evbus.Bus) {
	world := ecs.World{}

	//Setup bus system so every other system can have access to main server bus
	world.AddSystem(&BusSystem{Bus: bus})

	mapSystem := &MapSystem{Bus: bus}
	world.AddSystem(mapSystem)

	MainBus = *bus

	MainBus.Subscribe("tradewars:position", mapSystem.BroadcastIndividualPosition)

	MainWorld = world

}

func NewPlayer(callsign string) ecs.BasicEntity {
	entity := ecs.NewBasic()
	player := Player{entity, PositionComponent{0, 0, 0}, OwnershipComponent{callsign, entity.ID()}}

	for _, system := range MainWorld.Systems() {
		switch sys := system.(type) {

		case *MapSystem:
			sys.Add(&player.BasicEntity, &player.PositionComponent)

		default:
			//Do nothing
		}

	}

	return entity

}

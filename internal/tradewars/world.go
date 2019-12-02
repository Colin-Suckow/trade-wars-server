package tradewars

import (
	"github.com/EngoEngine/ecs"
)

var MainWorld ecs.World

func InitializeWorld() {
	world := ecs.World{}

	//Setup bus system so every other system can have access to main server bus
	world.AddSystem(&BusSystem{})

	mapSystem := &MapSystem{}
	world.AddSystem(mapSystem)

	WebsocketBus.Subscribe("tradewars:position", mapSystem.BroadcastIndividualPosition)
	WebsocketBus.Subscribe("tradewars:movePosition", mapSystem.moveIndividualPosition)

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

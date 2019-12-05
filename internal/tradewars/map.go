package tradewars

import (
	"encoding/json"

	"github.com/EngoEngine/ecs"
	evbus "github.com/asaskevich/EventBus"
)

//Handles location of all players and objects

type positionalEntity struct {
	*ecs.BasicEntity
	*PositionComponent
}

type MapSystem struct {
	entities []positionalEntity
	Bus      *evbus.Bus
}

func (m *MapSystem) Add(basic *ecs.BasicEntity, position *PositionComponent) {
	m.entities = append(m.entities, positionalEntity{basic, position})
}

func (m *MapSystem) Remove(basic ecs.BasicEntity) {
	var delete int = -1
	for index, entity := range m.entities {
		if entity.ID() == basic.ID() {
			delete = index
			break
		}
	}
	if delete >= 0 {
		m.entities = append(m.entities[:delete], m.entities[delete+1:]...)
	}
}

func (m *MapSystem) New(world *ecs.World) {

}

func (m *MapSystem) Update(dt float32) {
	//Do nothing.

}

func (m *MapSystem) BroadcastIndividualPosition(targetClient *client) {
	for _, entity := range m.entities {
		if entity.ID() == targetClient.entity.ID() {
			jsonData, err := json.Marshal(entity.PositionComponent)
			if err != nil {
				return
			}
			client := *targetClient
			println("Callsign: ")
			BroadcastJson(AddTargetToJson(string(jsonData), client.callsign))
			return
		}
	}

}

func (m *MapSystem) moveIndividualPosition(targetEntity ecs.BasicEntity, dx int, dy int) {
	for _, entity := range m.entities {
		if entity.ID() == targetEntity.ID() {

			entity.X += dx
			entity.Y += dy

			//broadcast new position
			jsonData, err := json.Marshal(entity.PositionComponent)
			if err != nil {
				break //Break from loop and precede to error handler
			}

			BroadcastJson(string(jsonData))
			return
		}
	}
	BroadcastJson(string("Test"))

}

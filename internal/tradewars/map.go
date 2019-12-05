package tradewars

import (
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

			broadcastEvent(buildEvent("positionUpdate", *targetClient, map[string]interface{}{
				"x": entity.PositionComponent.X,
				"y": entity.PositionComponent.Y,
			}))
			return
		}
	}

}

func (m *MapSystem) moveIndividualPosition(targetClient *client, x int, y int) {
	for _, entity := range m.entities {
		if entity.ID() == targetClient.entity.ID() {

			entity.X = x
			entity.Y = y

			broadcastEvent(buildEvent("positionUpdate", *targetClient, map[string]interface{}{
				"x": entity.PositionComponent.X,
				"y": entity.PositionComponent.Y,
			}))
			return
		}
	}
	BroadcastJson(string("Test"))

}

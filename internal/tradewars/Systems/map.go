package tradewars

import (
	"encoding/json"
	"github.com/EngoEngine/ecs"
	evbus "github.com/asaskevich/EventBus"
	components "suckow.dev/trade-wars-server/internal/tradewars/Components"
	tradewars "suckow.dev/trade-wars-server/internal/tradewars/Components"
	types "suckow.dev/trade-wars-server/internal/tradewars/Types"
)

//Handles location of all players and objects

type positionalEntity struct {
	*ecs.BasicEntity
	*tradewars.PositionComponent
}

type MapSystem struct {
	entities []positionalEntity
	Bus      evbus.Bus
}

func (m *MapSystem) Add(basic *ecs.BasicEntity, position *components.PositionComponent) {
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
	//Get bus from world
	m.Bus = FindBusSystem(world.Systems()).Bus

}

func (m MapSystem) broadcastIndividualPosition(message types.Message) {

	targetEntity := message.Target

	for _, entity := range m.entities {
		if entity.ID() == targetEntity.ID() {
			json, err := json.Marshal(entity.PositionComponent)
			if err != nil {
				break //Break from loop and precede to error handler
			}
			m.Bus.Publish("network:broadcast:individualposition", json)
			return
		}
	}
	m.Bus.Publish("network:error:internal")
}

package tradewars

import (
	"github.com/EngoEngine/ecs"
	evbus "github.com/asaskevich/EventBus"
	tradewars "suckow.dev/trade-wars-server/internal/tradewars/Components"
)

//Handles location of all players and objects

type positionalEntity struct {
	*ecs.BasicEntity
	*tradewars.PositionComponent
}

type MapSystem struct {
	entities []positionalEntity
	bus      *evbus.Bus
}

func (m *MapSystem) Add(basic *ecs.BasicEntity, position *tradewars.PositionComponent) {
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
	m.bus = FindBusSystem(world.Systems()).Bus
}

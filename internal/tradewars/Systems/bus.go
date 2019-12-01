package tradewars

import (
	"github.com/EngoEngine/ecs"
	evbus "github.com/asaskevich/EventBus"
)

//This system holds the entire applications message bus. Ideally it should not be a system, but I can't think of a better way to make it avilable to other systems
//Doesn't actually do anything

type BusSystem struct {
	entities []ecs.BasicEntity
	Bus      evbus.Bus
}

func (m *BusSystem) Add(basic *ecs.BasicEntity) {
	m.entities = append(m.entities)
}

func (m *BusSystem) Remove(basic ecs.BasicEntity) {
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

func (m *BusSystem) Update(dt float32) {
	//Do nothing. BusSystem should never do anything

}

//FindBusSystem Find self in list of other systems. Used by other systems to get a reference to the main bus
func FindBusSystem(systems []ecs.System) *BusSystem {
	for _, system := range systems {
		switch system.(type) {
		case *BusSystem:
			return system.(*BusSystem)
		default:
			//Do nothing
		}
	}
	panic("FindBusSystem: Could not find BusSystem in world!")
}

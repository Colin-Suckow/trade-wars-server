package tradewars

import (
	"encoding/json"
	"github.com/EngoEngine/ecs"
	evbus "github.com/asaskevich/EventBus"
	"log"
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
	//Get bus from world
	//log.Println("Assigning bus to mapsystem")
	m.Bus = FindBusSystem(world.Systems()).Bus

}

func (m *MapSystem) Update(dt float32) {
	//Do nothing.

}

func (m *MapSystem) BroadcastIndividualPosition(targetEntity ecs.BasicEntity) {
	log.Println("Started function")
	for i, entity := range m.entities {
		log.Println("loop " + string(i))
		if entity.ID() == targetEntity.ID() {
			_, err := json.Marshal(entity.PositionComponent)
			if err != nil {
				log.Println("err")
				break //Break from loop and precede to error handler
			}
			log.Println("Found position")

			if MainBus.HasCallback("network:broadcast:json") == true {
				log.Println("Has callback")
			} else {
				log.Println("Does not have callback")
			}
			MainBus.Publish("network:broadcast:json", "test")
			log.Println("Published")
			return
		}
	}
	log.Println("Found nothing")
}

func (m MapSystem) moveIndividualPosition(id uint64, dx int, dy int) {
	mBus := *m.Bus
	for _, entity := range m.entities {
		if entity.ID() == id {

			entity.X += dx
			entity.Y += dy

			//broadcast new position
			json, err := json.Marshal(entity.PositionComponent)
			if err != nil {
				break //Break from loop and precede to error handler
			}

			mBus.Publish("network:broadcast:json", json)
			return
		}
	}
	mBus.Publish("network:error:internal")
}

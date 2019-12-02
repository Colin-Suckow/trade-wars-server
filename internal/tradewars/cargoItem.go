package tradewars

import (
	"github.com/EngoEngine/ecs"
	
)

//Individual Cargo Item

type CargoItem struct {
	ecs.BasicEntity
	NameComponent
	ValueComponent
}

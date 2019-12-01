package tradewars

import "github.com/EngoEngine/ecs"

type Message struct {
	Key    string
	Target ecs.BasicEntity
	data   string
}

package networking

import (
	"log"
	"net/http"

	"github.com/EngoEngine/ecs"
	"github.com/asaskevich/EventBus"
	"github.com/gorilla/websocket"
	"suckow.dev/trade-wars-server/internal/tradewars"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type client struct {
	conn   *websocket.Conn
	entity ecs.BasicEntity
}

var WebsocketBus EventBus.Bus

//Store socket connects so we can write to them
var Connections []client

func initializeWebsocketConnection(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		BadRequstError(w)
		return
	}
	log.Println("Client connected")
	ws.WriteMessage(1, []byte("Hi new client!"))

	//Create new player for new connection
	newPlayer := tradewars.NewPlayer("")

	Connections = append(Connections, client{ws, newPlayer})
	reader(ws)
}

func reader(conn *websocket.Conn) {
	for {
		//Read in new message
		_, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		log.Println("Recieved message: " + string(p))

	}
}

package tradewars

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/EngoEngine/ecs"
	"github.com/asaskevich/EventBus"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type client struct {
	conn     *websocket.Conn
	entity   ecs.BasicEntity
	callsign string
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
	newPlayer := NewPlayer("")

	Connections = append(Connections, client{ws, newPlayer, ""})
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
		decodeCommand(p, conn)
	}
}

func decodeCommand(jsonData []byte, conn *websocket.Conn) {
	var objmap map[string]interface{}
	if err := json.Unmarshal(jsonData, &objmap); err != nil {
		respondInvalid(conn)
		return
	}

	command := objmap["command"]
	if command == nil {
		respondInvalid(conn)
		return
	}

	switch command {
	case "ping":
		conn.WriteMessage(websocket.TextMessage, []byte("pong"))
		return
	case "getOwnPosition":
		WebsocketBus.Publish("tradewars:position", getClientFromConnection(conn).entity)
		return
	case "changeOwnPosition":
		changePosition(conn, jsonData)
	case "setCallsign":
		setCallsign(getClientFromConnection(conn), jsonData)
	default:
		respondInvalid(conn)
	}

}

func respondInvalid(conn *websocket.Conn) {
	conn.WriteMessage(websocket.TextMessage, []byte("Invalid message"))
}

func BroadcastJson(jsonData string) {
	log.Println("Sending data")
	for _, client := range Connections {
		log.Println("Sent a line")
		client.conn.WriteMessage(websocket.TextMessage, []byte(jsonData))
	}
}

func getClientFromConnection(conn *websocket.Conn) client {
	for _, client := range Connections {
		if client.conn == conn {
			return client
		}
	}
	log.Println("Could not find connection")
	panic("Could not find connection!")
}

func setCallsign(cli client, jsonData []byte) {
	objmap := readJson(jsonData)
	cli.callsign = objmap["callsign"].(string)
	BroadcastJson(cli.callsign)
}

func changePosition(conn *websocket.Conn, jsonData []byte) {
	objmap := readJson(jsonData)
	WebsocketBus.Publish("tradewars:movePosition", getClientFromConnection(conn).entity, int(objmap["x"].(float64)), int(objmap["y"].(float64)))
}

func readJson(jsonData []byte) map[string]interface{} {
	var objmap map[string]interface{}
	if err := json.Unmarshal(jsonData, &objmap); err != nil {
		panic("could not read json")
	}
	return objmap
}

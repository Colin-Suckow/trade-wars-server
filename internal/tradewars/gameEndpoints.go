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
var Connections []*client

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
	Connections = append(Connections, &client{ws, newPlayer, "NULL"})
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

	client := getClientFromConnection(conn)

	switch command {
	case "ping":
		conn.WriteMessage(websocket.TextMessage, []byte("pong"))
		return
	case "getOwnPosition":
		WebsocketBus.Publish("tradewars:position", client)
		return
	case "changeOwnPosition":
		changePosition(client, jsonData)
	case "setCallsign":
		setCallsign(client, jsonData)
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

func getClientFromConnection(conn *websocket.Conn) *client {
	for _, client := range Connections {
		if client.conn == conn {
			log.Println("Found client: " + client.callsign)
			return client
		}
	}
	panic("Could not find connection!")
}

func getClientFromEntity(entity ecs.BasicEntity) *client {
	for _, client := range Connections {
		if client.entity.ID() == entity.ID() {
			return client
		}
	}
	panic("Could not find connection!")
}

func setCallsign(cli *client, jsonData []byte) {
	objmap := readJson(jsonData)
	if objmap["callsign"] == nil {
		respondInvalid(cli.conn)
		return
	}
	client := *cli
	client.callsign = objmap["callsign"].(string)
	*cli = client
	BroadcastJson("{\"event\":\"changedCallsign\"}")
}

func changePosition(client *client, jsonData []byte) {
	objmap := readJson(jsonData)
	WebsocketBus.Publish("tradewars:movePosition", client.entity, int(objmap["x"].(float64)), int(objmap["y"].(float64)))
}

func readJson(jsonData []byte) map[string]interface{} {
	var objmap map[string]interface{}
	if err := json.Unmarshal(jsonData, &objmap); err != nil {
		panic("could not read json")
	}
	return objmap
}

func AddTargetToJson(jsonData string, target string) string {
	//Quick dirty implementation
	//Assume the input data is correct and just add the field manually after the first character, which is assumed to be a {
	jsonToInsert := "\"target\":\"" + target + "\","
	return jsonData[:1] + jsonToInsert + jsonData[1:]
}

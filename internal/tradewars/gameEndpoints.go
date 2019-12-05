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

type Event struct {
	EventType   string
	Target      string
	EventParams map[string]interface{}
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

	client := getClientFromConnection(conn)

	var objmap map[string]interface{}
	if err := json.Unmarshal(jsonData, &objmap); err != nil {
		respondInvalid(client)
		return
	}

	command := objmap["command"]
	if command == nil {
		respondInvalid(client)
		return
	}

	//Commands that are valid without a callsign set
	switch command {
	case "ping":
		conn.WriteMessage(websocket.TextMessage, []byte("pong"))
		return
	case "setCallsign":
		setCallsign(client, jsonData)
		return
	default:
		//Let the other switch handle it
	}

	if client.callsign == "NULL" {
		respondInvalid(client)
		return
	}

	//Commands that require a callsign set
	switch command {
	case "changeOwnPosition":
		changePosition(client, jsonData)
	case "getOwnPosition":
		WebsocketBus.Publish("tradewars:position", client)

	default:
		respondInvalid(client)
	}

}

func respondInvalid(cli *client) {
	respondEvent(cli, buildEmptyEvent("invalidRequest", *cli))
}

func BroadcastJson(jsonData string) {
	log.Println("Sending data")
	for _, client := range Connections {
		log.Println("Sent a line")
		client.conn.WriteMessage(websocket.TextMessage, []byte(jsonData))
	}
}

func broadcastEvent(event Event) {
	jsonData, err := json.Marshal(event)

	if err != nil {
		return
	}

	BroadcastJson(string(jsonData))
}

func respondEvent(cli *client, event Event) {
	jsonData, err := json.Marshal(event)

	if err != nil {
		return
	}

	cli.conn.WriteMessage(websocket.TextMessage, jsonData)
}

func buildEvent(eventName string, target client, arguments map[string]interface{}) Event {
	return Event{eventName, target.callsign, arguments}
}

func buildEmptyEvent(eventName string, target client) Event {
	return buildEvent(eventName, target, map[string]interface{}{})
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
		respondInvalid(cli)
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
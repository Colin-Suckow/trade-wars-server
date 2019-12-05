package tradewars

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

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
	lastPing time.Time
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
	Connections = append(Connections, &client{ws, newPlayer, "NULL", time.Now()})
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

	log.Println("Recieved " + command.(string) + " from " + client.callsign)

	//Commands that are valid without a callsign set
	switch command {
	case "ping":
		handlePing(client)
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
	case "setOwnPosition":
		changePosition(client, jsonData)
	case "getOwnPosition":
		WebsocketBus.Publish("tradewars:position", client)

	case "getAllPosition":
		for _, cli := range Connections {
			WebsocketBus.Publish("tradewars:positionRespondAll", client, cli)
		}

	case "chatMessage":
		recieveChat(client, jsonData)

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
	oldCallsign := cli.callsign
	cli.callsign = objmap["callsign"].(string)
	broadcastEvent(buildEvent("callsignChange", *cli, map[string]interface{}{"old": oldCallsign}))
}

func changePosition(client *client, jsonData []byte) {
	objmap := readJson(jsonData)
	WebsocketBus.Publish("tradewars:movePosition", client, int(objmap["x"].(float64)), int(objmap["y"].(float64)))
}

func readJson(jsonData []byte) map[string]interface{} {
	var objmap map[string]interface{}
	if err := json.Unmarshal(jsonData, &objmap); err != nil {
		panic("could not read json")
	}
	return objmap
}

func recieveChat(cli *client, jsonData []byte) {
	objmap := readJson(jsonData)
	if objmap["message"] == nil {
		respondInvalid(cli)
		return
	}
	broadcastEvent(buildEvent("chatMessage", *cli, map[string]interface{}{"message": objmap["message"]}))
}

func handlePing(cli *client) {

	cli.lastPing = time.Now()

	//Directly write output for simplicity
	cli.conn.WriteMessage(websocket.TextMessage, []byte("pong"))
}

func checkConnections() {
	pongTimeout := 20 * time.Second
	currentTime := time.Now()
	for _, client := range Connections {
		if currentTime.Sub(client.lastPing) > pongTimeout {
			disconnectClient(client)
		}
	}
}

func disconnectClient(cli *client) {
	log.Println(cli.callsign + ":" + string(cli.entity.GetBasicEntity().ID()) + " disconnected. Reason: Timeout")
	for i, client := range Connections {
		if client.entity.ID() == cli.entity.ID() {
			respondEvent(cli, buildEvent("disconnect", *cli, map[string]interface{}{"reason": "timeout"}))
			client.conn.Close()
			Connections = remove(Connections, i)
		}
	}
}

func AddTargetToJson(jsonData string, target string) string {
	//Quick dirty implementation
	//Assume the input data is correct and just add the field manually after the first character, which is assumed to be a {
	jsonToInsert := "\"target\":\"" + target + "\","
	return jsonData[:1] + jsonToInsert + jsonData[1:]
}

func remove(s []*client, i int) []*client {
	s[i] = s[len(s)-1]
	// We do not need to put s[i] at the end, as it will be discarded anyway
	return s[:len(s)-1]
}

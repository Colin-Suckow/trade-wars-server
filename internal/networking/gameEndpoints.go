package networking

import "github.com/gorilla/websocket"

import "net/http"

import "log"

import "github.com/asaskevich/EventBus"

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var WebsocketBus EventBus.Bus

func initializeWebsocketConnection(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		BadRequstError(w)
		return
	}
	log.Println("Client connected")
	ws.WriteMessage(1, []byte("Hi new client!"))
	reader(ws)
}

func reader(conn *websocket.Conn) {

}

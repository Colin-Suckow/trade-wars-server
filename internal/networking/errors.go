package networking

import (
	"github.com/gorilla/websocket"
	"net/http"
)

//Generic bad request error
func BadRequstError(w http.ResponseWriter) {
	http.Error(w, "Bad Request", 400)
}

func wsInternalError(conn *websocket.Conn) {
	conn.WriteMessage(websocket.TextMessage, []byte("Internal Error"))
}

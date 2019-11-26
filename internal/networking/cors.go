package networking

import "net/http"

//Allow any client to access api. Begone cors
func EnableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

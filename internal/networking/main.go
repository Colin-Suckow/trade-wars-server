package networking

import (
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func ServeMux(mux *http.ServeMux) {
	err := godotenv.Load()
	portVariable := os.Getenv("PORT")

	log.Println("Starting server on port " + portVariable)
	err = http.ListenAndServe(":"+portVariable, mux)
	log.Fatal(err)
}

func SetupRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", ClientPlayersScreen)
	mux.HandleFunc("/players", ClientPlayersScreen)
	mux.HandleFunc("/map", ClientGameScreen)
	mux.HandleFunc("/startSession", IssueCookie)
	mux.HandleFunc("/getCallsign", testReturnCallsign)
}

func testReturnCallsign(w http.ResponseWriter, r *http.Request) {
	EnableCors(w)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	var cookie, err = r.Cookie("callsign")
	if err != nil {
		http.Error(w, "Couldn't read callsign. :(", 500)
		return
	}
	callsign := cookie.Value
	w.Write([]byte("Hello " + callsign))
}

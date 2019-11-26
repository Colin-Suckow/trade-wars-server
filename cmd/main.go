package main

import (
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	networking "suckow.dev/trade-wars-server/internal/networking"
)

func testReturnCallsign(w http.ResponseWriter, r *http.Request) {
	//networking.EnableCors(&w)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	var cookie, err = r.Cookie("callsign")
	if err != nil {
		http.Error(w, "Couldn't read callsign. :(", 500)
		return
	}
	callsign := cookie.Value
	w.Write([]byte("Hello " + callsign))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", networking.ClientPlayersScreen)
	mux.HandleFunc("/players", networking.ClientPlayersScreen)
	mux.HandleFunc("/map", networking.ClientGameScreen)
	mux.HandleFunc("/startSession", networking.IssueCookie)
	mux.HandleFunc("/getCallsign", testReturnCallsign)

	err := godotenv.Load()
	portVariable := os.Getenv("PORT")

	log.Println("Starting server on port " + portVariable)
	err = http.ListenAndServe(":"+portVariable, mux)
	log.Fatal(err)
}

package main

import (
	"log"
	"net/http"
	"os"
	networking "suckow.dev/trade-wars-server/internal/networking"
)

func handleHome(w http.ResponseWriter, r *http.Request) {
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
	mux.HandleFunc("/", handleHome)
	mux.HandleFunc("/players", networking.ClientPlayersScreen)
	mux.HandleFunc("/map", networking.ClientGameScreen)
	mux.HandleFunc("/startSession", networking.IssueCookie)

	portVariable := os.Getenv("PORT")
	if portVariable == "" {
		portVariable = "4000"
	}

	log.Println("Starting server on port " + portVariable)
	err := http.ListenAndServe(":"+portVariable, mux)
	log.Fatal(err)
}

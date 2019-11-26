package main

import (
	"log"
	"net/http"
	"os"
	networking "suckow.dev/trade-wars-server/internal/networking"
)

func handleHome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world"))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleHome)
	mux.HandleFunc("/players", networking.ClientPlayersScreen)
	mux.HandleFunc("/map", networking.ClientGameScreen)

	portVariable := os.Getenv("PORT")
	if portVariable == "" {
		portVariable = "4000"
	}

	log.Println("Starting server on port " + portVariable)
	err := http.ListenAndServe(":"+portVariable, mux)
	log.Fatal(err)
}

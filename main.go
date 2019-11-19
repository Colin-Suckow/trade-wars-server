package main

import (
	"log"
	"net/http"
)

type Player struct {
	Name      string
	ID        int
	Inventory []Cargo
	Sector    string
	System    string
	X         int
	Y         int
}

type Cargo struct {
	Name        string
	ID          int
	Description string
	Value       int
}

type Sector struct {
	Name        string
	Description string
	ID          int
	Systems     []System
	X           int
	Y           int
}

type System struct {
	Name        string
	Description string
	ID          int
	Planets     []Planet
	X           int
	Y           int
}

type Planet struct {
	Name        string
	Description string
	ID          int
	X           int
	Y           int
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world"))
}

func clientPlayersScreen(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://carson-key.github.io/trade-wars-static/#/", 308)
}

func clientGameScreen(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://carson-key.github.io/trade-wars-static/#/game", 308)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", clientPlayersScreen)
	mux.HandleFunc("/players", clientPlayersScreen)
	mux.HandleFunc("/map", clientGameScreen)

	log.Println("Starting server on port 4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}

package main

import (
	"log"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world"))
}

func redirectToClient(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://carson-key.github.io/trade-wars-static/#/", 308)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/redirect", redirectToClient)

	log.Println("Starting server on port 4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}

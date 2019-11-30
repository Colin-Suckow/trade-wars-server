package main

import (
	"net/http"
	"suckow.dev/trade-wars-server/internal/networking"
)

func main() {
	mux := http.NewServeMux()
	networking.SetupRoutes(mux)
	networking.ServeMux(mux)
}

package main

import (
	"github.com/asaskevich/EventBus"
	"net/http"
	"suckow.dev/trade-wars-server/internal/networking"
)

func main() {

	bus := EventBus.New()

	mux := http.NewServeMux()
	networking.SetupRoutes(mux)
	networking.ServeMux(mux)
}

package main

import (
	"github.com/asaskevich/EventBus"
	"net/http"
	"suckow.dev/trade-wars-server/internal/networking"
	"suckow.dev/trade-wars-server/internal/tradewars"
)

func main() {

	bus := EventBus.New()
	tradewars.InitializeWorld(bus)

	mux := http.NewServeMux()
	networking.SetupRoutes(mux)
	networking.InitializeBridge(bus)
	networking.ServeMux(mux)
}

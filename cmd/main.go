package main

import (
	"net/http"

	"github.com/asaskevich/EventBus"
	networking "suckow.dev/trade-wars-server/internal/networking"
	"suckow.dev/trade-wars-server/internal/tradewars"
)

func main() {

	bus := EventBus.New()
	tradewars.InitializeWorld(&bus)

	mux := http.NewServeMux()
	networking.SetupRoutes(mux)
	networking.InitializeBridge(&bus)
	networking.ServeMux(mux)
}

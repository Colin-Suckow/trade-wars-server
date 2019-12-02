package main

import (
	"net/http"

	"github.com/asaskevich/EventBus"
	"suckow.dev/trade-wars-server/internal/tradewars"
)

func main() {

	bus := EventBus.New()

	mux := http.NewServeMux()
	tradewars.SetupRoutes(mux)
	tradewars.InitializeBridge(&bus)

	tradewars.InitializeWorld()

	tradewars.ServeMux(mux)
}

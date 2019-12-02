package tradewars

import "net/http"

func ClientPlayersScreen(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://carson-key.github.io/trade-wars-static/#/", 301)
}

func ClientGameScreen(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://carson-key.github.io/trade-wars-static/#/game", 301)
}

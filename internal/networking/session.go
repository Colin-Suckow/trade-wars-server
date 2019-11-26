package networking

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

type callsignPostBody struct {
	Callsign string `json:"callsign"`
}

func IssueCookie(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			//TODO respond with error
			panic(err)
		}
		var callsignData callsignPostBody
		err = json.Unmarshal(body, &callsignData)
		if err != nil {
			//Todo respond with error again
			panic(err)
		}
		cookie := http.Cookie{
			Name:    "callsign",
			Value:   callsignData.Callsign,
			Expires: time.Now().AddDate(0, 0, 1),
			Path:    "/",
		}
		http.SetCookie(w, &cookie)
		w.Write([]byte("Sent cookie: " + callsignData.Callsign))

	}
}

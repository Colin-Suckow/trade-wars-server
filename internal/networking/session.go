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
		EnableCors(&w)
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			BadRequstError(w)
			return
		}
		var callsignData callsignPostBody
		err = json.Unmarshal(body, &callsignData)
		if err != nil || callsignData.Callsign == "" {
			BadRequstError(w)
			return
		}
		cookie := http.Cookie{
			Name:    "callsign",
			Value:   callsignData.Callsign,
			Expires: time.Now().AddDate(0, 0, 1),
			Path:    "/",
		}
		http.SetCookie(w, &cookie)
		w.Write([]byte("{\n\"success\":true,\n\"callsign\":\"" + callsignData.Callsign + "\"\n}")) //TODO make this less bad

	}
}

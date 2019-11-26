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
	//EnableCors(&w)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	if r.Method == http.MethodOptions {
		return
	} else if r.Method == http.MethodPost {

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
		
	}
		w.Write([]byte("{\n\"success\":true,\n\"callsign\":\"" + callsignData.Callsign + "\"\n}")) //TODO make this less bad

	}
}

package networking

import "net/http"

//Generic bad request error
func BadRequstError(w http.ResponseWriter) {
	http.Error(w, "Bad Request", 400)
}

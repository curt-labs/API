package access_control

import (
	"net/http"
)

func Tokenize(w http.ResponseWriter, r *http.Request) {

	params := r.URL.Query()
	if len(params.Get("key")) == 0 {
		http.Redirect(w, r, "/unauthorized", http.StatusFound)
	}
	return
}

var FilterUser = func(w http.ResponseWriter, r *http.Request) {
	if r.URL.User == nil || r.URL.User.Username() != "admin" {
		http.Error(w, "", http.StatusUnauthorized)
	}
}

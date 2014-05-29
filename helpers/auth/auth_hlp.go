package auth

import (
	"github.com/curt-labs/GoAPI/helpers/database"
	"github.com/martini-contrib/sessions"
	"net/http"

	//"time"
)

var AuthHandler = func(w http.ResponseWriter, r *http.Request, session sessions.Session) {

	params := r.URL.Query()
	key := params.Get("key")

	if r.Method != "GET" && key == "" {
		key = r.FormValue("key")
	}

	if len(key) == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if r.Method == "GET" {
		// First we'll try checking the session to see if we have already authenticated the key
		val := session.Get(key)
		if val != nil {
			return
		}

		if checkKey(key) {
			session.Set(key, key)
			return
		}

		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	} else {
		if checkPrivateKey(key) {
			return
		}
	}

	http.Error(w, "Unauthorized", http.StatusUnauthorized)
	return
}

func checkKey(key string) bool {

	qry, err := database.GetStatement("AuthStmt")
	if err != nil {
		return false
	}

	params := struct {
		Key string
	}{key}

	rows, _, err := qry.Exec(params)
	if database.MysqlError(err) {
		return false
	}
	if len(rows) == 0 {
		return false
	}

	return true
}

func checkPrivateKey(key string) bool {

	qry, err := database.GetStatement("PrivateAuthStmt")
	if err != nil {
		return false
	}

	params := struct {
		Key string
	}{key}
	rows, _, err := qry.Exec(params)
	if database.MysqlError(err) || len(rows) == 0 {
		return false
	}
	return true
}

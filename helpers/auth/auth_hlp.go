package auth

import (
	"../../plate"
	"../database"
	"net/http"

	//"time"
)

var (

	//  Prepared statements would go here
	authStmt = `select id from ApiKey where api_key = '%s'`
)

var AuthHandler = func(w http.ResponseWriter, r *http.Request) {

	params := r.URL.Query()
	key := params.Get("key")

	if len(key) == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// First we'll try checking the session to see if we have already authenticated the key
	session := plate.Session.Get(r)
	if session[key] != nil {
		return
	}

	if !checkKey(key) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	return
}

func checkKey(key string) bool {

	rows, _, err := database.Db.Query(authStmt, key)
	if database.MysqlError(err) {
		return false
	}
	if len(rows) == 0 {
		return false
	}

	return true
}

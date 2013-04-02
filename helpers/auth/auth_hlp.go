package auth

import (
	"../../helpers/plate"
	"../database"
	"net/http"

	//"time"
)

var (

	//  Prepared statements would go here
	authStmt = `select id from ApiKey where api_key = ?`

	privateAuthStmt = `select ak.id from ApiKey as ak
				join ApiKeyType as akt on ak.type_id = akt.id
				where akt.type = 'PRIVATE'
				&& api_key = ?`
)

var AuthHandler = func(w http.ResponseWriter, r *http.Request) {

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
		session := plate.Session.Get(r)
		if session[key] != nil {
			return
		}

		if checkKey(key) {
			return
		}

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

	qry, err := database.Db.Prepare(authStmt)
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

	qry, err := database.Db.Prepare(privateAuthStmt)
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

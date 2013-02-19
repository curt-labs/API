package auth

import (
	"../../plate"
	"../database"
	"../mymysql/autorc"
	// "github.com/ziutek/mymysql/mysql"
	_ "../mymysql/thrsafe"
	"net/http"

	//"time"
)

const (
	db_proto = "tcp"
	db_addr  = "curtsql.cloudapp.net:3306"
	db_user  = "root"
	db_pass  = "eC0mm3rc3"
	db_name  = "CurtDev"
)

var (
	// MySQL Connection Handler
	db = autorc.New(db_proto, "", db_addr, db_user, db_pass, db_name)

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

	rows, _, err := db.Query(authStmt, key)
	if database.MysqlError(err) {
		return false
	}
	if len(rows) == 0 {
		return false
	}

	return true
}

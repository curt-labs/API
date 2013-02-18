package auth

import (
	"../../plate"
	"github.com/ziutek/mymysql/autorc"
	// "github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/thrsafe"
	//"log"
	"net/http"
	//"time"
)

const (
	db_proto = "tcp"
	db_addr  = "curtsql.cloudapp.net"
	db_user  = "root"
	db_pass  = "eC0mm3rc3"
	db_name  = "CurtDev"
)

var (
	// MySQL Connection Handler
	db = autorc.New(db_proto, "", db_addr, db_user, db_pass, db_name)

	//  Prepared statements would go here
	//  stmt *autorc.Stmt
)

var AuthHandler = func(w http.ResponseWriter, r *http.Request) bool {

	params := r.URL.Query()
	key := params.Get("key")

	if len(key) == 0 {
		return false
	}

	// First we'll try checking the session to see if we have already authenticated the key
	session := plate.Session.Get(r)
	if session[key] != nil {
		return true
	}

	// check the database
	//hd, err := hood.Open("mymysql", "dataSourceName")

	return true
}

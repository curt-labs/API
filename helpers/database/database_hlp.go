package database

import (
	"../mymysql/autorc"
	// "github.com/ziutek/mymysql/mysql"
	_ "../mymysql/thrsafe"
	"log"
	"os"
)

var (
	db_proto = "tcp"
	db_addr  = os.Getenv("DB_HOST")
	db_user  = os.Getenv("API_DB_USER")
	db_pass  = os.Getenv("API_DB_PASS")
	db_name  = os.Getenv("API_DB")

	// MySQL Connection Handler
	Db = autorc.New(db_proto, "", db_addr, db_user, db_pass, db_name)

	//  Prepared statements would go here
	//  stmt *autorc.Stmt
)

func MysqlError(err error) (ret bool) {
	ret = (err != nil)
	if ret {
		log.Println("MySQL error: ", err)
	}
	return
}

func MysqlErrExit(err error) {
	if MysqlError(err) {
		os.Exit(1)
	}
}

package database

import (
	//"../mymysql/autorc"
	"../mymysql/thrsafe"
	"log"
	"os"
)

var (

	// MySQL Connection Handler
	Db = thrsafe.New(db_proto, "", db_addr, db_user, db_pass, db_name)

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
		log.Println(err)
		os.Exit(1)
	}
}

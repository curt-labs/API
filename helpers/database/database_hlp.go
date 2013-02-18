package database

import (
	"../mymysql/autorc"
	// "github.com/ziutek/mymysql/mysql"
	_ "../mymysql/thrsafe"
	"log"
	"os"
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

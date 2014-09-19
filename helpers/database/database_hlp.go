package database

import (
	"flag"
	"github.com/ziutek/mymysql/autorc"
	_ "github.com/ziutek/mymysql/thrsafe"
	"log"
	"os"
)

var (

	// MySQL Connection Handler
	Db      = autorc.New("tcp", "", "127.0.0.1:3306", "root", "", "CurtDev")
	EmptyDb = flag.String("clean", "", "bind empty database with structure defined")
)

func BindDatabase() {
	if addr := os.Getenv("DATABASE_HOST"); addr != "" {
		proto := os.Getenv("DATABASE_PROTOCOL")
		user := os.Getenv("DATABASE_USERNAME")
		pass := os.Getenv("DATABASE_PASSWORD")
		name := os.Getenv("CURT_DEV_NAME")
		Db = autorc.New(proto, "", addr, user, pass, name)
		return
	}

	// flag.Parse()
	if EmptyDb != nil && *EmptyDb != "" {
		Db = autorc.New("tcp", "", "127.0.0.1:3306", "root", "", "CurtDev_Empty")
	}
}

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

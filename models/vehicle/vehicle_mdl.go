package vehicle

import (
	"../../helpers/database"
	"github.com/ziutek/mymysql/mysql"
	// _ "github.com/ziutek/mymysql/thrsafe"
)

type ConfigOption struct {
	Type    string
	Options []mysql.Row
}

type Vehicle struct {
	Year                  float64
	Make, Model, Submodel string
	Configuration         []string
	Parts                 []interface{}
	Groups                []interface{}
}

func (vehicle *Vehicle) GetYears() (opt ConfigOption) {
	db := database.Db

	opt.Type = "Years"

	smt, err := db.Prepare("select YearID from vcdb_Year order by YearID desc")
	database.MysqlErrExit(err)

	rows, _, err := smt.Exec()
	if database.MysqlError(err) {
		return opt
	}

	opt.Options = rows
	return
}

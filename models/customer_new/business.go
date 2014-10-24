package customer_new

import (
	"database/sql"

	"github.com/curt-labs/GoAPI/helpers/database"
	"github.com/curt-labs/GoAPI/helpers/sortutil"
	_ "github.com/go-sql-driver/mysql"
)

var (
	getBusinessClassesStmt = `select BusinessClassID, name, sort, showOnWebsite from BusinessClass`
)

type BusinessClasses []BusinessClass
type BusinessClass struct {
	ID            int    `json:"id" xml:"id"`
	Name          string `json:"name" xml:"name"`
	Sort          int    `json:"sort" xml:"sort"`
	ShowOnWebsite bool   `json:"show" xml:"show"`
}

func GetAllBusinessClasses() (classes BusinessClasses, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return
	}
	defer db.Close()

	stmt, err := db.Prepare(getBusinessClassesStmt)
	if err != nil {
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return
	}
	var bc BusinessClass
	for rows.Next() {
		bc = BusinessClass{}
		err = rows.Scan(
			&bc.ID,
			&bc.Name,
			&bc.Sort,
			&bc.ShowOnWebsite,
		)
		if err != nil {
			return
		}
		classes = append(classes, bc)
	}
	defer rows.Close()

	sortutil.AscByField(classes, "Sort")

	return
}

package part

import (
	"database/sql"
	"github.com/curt-labs/GoAPI/helpers/database"
	_ "github.com/go-sql-driver/mysql"
)

type Attribute struct {
	Key, Value string
}

var (
	partAttrStmt = `select field, value from PartAttribute where partID = ?`
)

func (p *Part) GetAttributes() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return
	}
	defer db.Close()

	qry, err := db.Prepare(partAttrStmt)
	if err != nil {
		return
	}
	defer qry.Close()

	rows, err := qry.Query(p.PartId)
	if err != nil || rows == nil {
		return err
	}

	var attrs []Attribute
	for rows.Next() {
		var attr Attribute
		if err := rows.Scan(&attr.Key, &attr.Value); err == nil {
			attrs = append(attrs, attr)
		}
	}
	p.Attributes = attrs

	return
}

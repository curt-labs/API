package products

import (
	"database/sql"
	"github.com/curt-labs/GoAPI/helpers/database"
	_ "github.com/go-sql-driver/mysql"
)

var (
	getYearsStmt = `
		select distinct y.YearID from vcdb_Year as y
		join BaseVehicle as bv on y.YearID = bv.YearID
		join vcdb_Vehicle as v on bv.ID = v.BaseVehicleID
		join vcdb_VehiclePart as vp on v.ID = vp.VehicleID
		join Part as p on vp.PartNumber = p.partID
		where p.status = 800 || p.status = 900
		order by y.YearID desc`
)

func (l *Lookup) GetYears() error {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(getYearsStmt)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Query()
	if err != nil {
		return err
	}

	l.Years = make([]int, 0)
	for res.Next() {
		var year int
		err = res.Scan(&year)
		if err == nil {
			l.Years = append(l.Years, year)
		}
	}
	defer res.Close()

	l.Pagination = Pagination{
		TotalItems:    len(l.Years),
		ReturnedCount: len(l.Years),
		Page:          1,
		PerPage:       len(l.Years),
		TotalPages:    1,
	}

	return nil
}

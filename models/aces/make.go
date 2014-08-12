package aces

import (
	"database/sql"
	"github.com/curt-labs/GoAPI/helpers/database"
	_ "github.com/go-sql-driver/mysql"
)

var (
	getMakeStmt = `
		select distinct m.MakeName from vcdb_Make as m
		join BaseVehicle as bv on m.ID = bv.MakeID
		join vcdb_Vehicle as v on bv.ID = v.BaseVehicleID
		join vcdb_VehiclePart as vp on v.ID = vp.VehicleID
		join Part as p on vp.PartNumber = p.partID
		where (p.status = 800 || p.status = 900) && bv.YearID = ?
		order by m.MakeName`
)

func (l *Lookup) GetMakes() error {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(getMakeStmt)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Query(l.Vehicle.Base.Year)
	if err != nil {
		return err
	}

	l.Makes = make([]string, 0)
	for res.Next() {
		var ma string
		err = res.Scan(&ma)
		if err == nil {
			l.Makes = append(l.Makes, ma)
		}
	}

	l.Pagination = Pagination{
		TotalItems:    len(l.Makes),
		ReturnedCount: len(l.Makes),
		Page:          1,
		PerPage:       len(l.Makes),
		TotalPages:    1,
	}

	return nil
}

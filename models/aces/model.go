package aces

import (
	"database/sql"
	"github.com/curt-labs/GoAPI/helpers/database"
	_ "github.com/go-sql-driver/mysql"
)

var (
	getModelStmt = `select distinct mo.ModelName from vcdb_Model as mo
		join BaseVehicle as bv on mo.ID = bv.ModelID
		join vcdb_Make as ma on bv.MakeID = ma.ID
		join vcdb_Vehicle as v on bv.ID = v.BaseVehicleID
		join vcdb_VehiclePart as vp on v.ID = vp.VehicleID
		join Part as p on vp.PartNumber = p.partID
		where (p.status = 800 || p.status = 900) && bv.YearID = ? && ma.MakeName = ?
		order by mo.ModelName`
)

func (l *Lookup) GetModels() error {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(getModelStmt)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Query(l.Vehicle.Base.Year, l.Vehicle.Base.Make)
	if err != nil {
		return err
	}

	l.Models = make([]string, 0)
	for res.Next() {
		var m string
		err = res.Scan(&m)
		if err == nil {
			l.Models = append(l.Models, m)
		}
	}

	l.Pagination = Pagination{
		TotalItems:    len(l.Models),
		ReturnedCount: len(l.Models),
		Page:          1,
		PerPage:       len(l.Models),
		TotalPages:    1,
	}

	return nil
}

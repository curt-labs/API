package products

import (
	"database/sql"
	"github.com/curt-labs/GoAPI/helpers/database"
	_ "github.com/go-sql-driver/mysql"
)

var (
	getSubmodelStmt = `
		select distinct s.SubmodelName from vcdb_Vehicle as v
		join Submodel as s on v.SubModelID = s.ID
		join BaseVehicle as bv on v.BaseVehicleID = bv.ID
		join vcdb_Model as mo on bv.ModelID = mo.ID
		join vcdb_Make as ma on bv.MakeID = ma.ID
		join vcdb_VehiclePart as vp on v.ID = vp.VehicleID
		join Part as p on vp.PartNumber = p.partID
		where (p.status = 800 || p.status = 900) && bv.YearID = ? && ma.MakeName = ? && mo.ModelName = ?
		&& p.brandID in (?)
		order by s.SubmodelName`
)

func (l *Lookup) GetSubmodels(brandIds string) error {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(getSubmodelStmt)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Query(l.Vehicle.Base.Year, l.Vehicle.Base.Make, l.Vehicle.Base.Model, brandIds)
	if err != nil {
		return err
	}

	l.Submodels = make([]string, 0)
	for res.Next() {
		var m string
		err = res.Scan(&m)
		if err == nil {
			l.Submodels = append(l.Submodels, m)
		}
	}
	defer res.Close()

	l.Pagination = Pagination{
		TotalItems:    len(l.Submodels),
		ReturnedCount: len(l.Submodels),
		Page:          1,
		PerPage:       len(l.Submodels),
		TotalPages:    1,
	}

	return nil
}

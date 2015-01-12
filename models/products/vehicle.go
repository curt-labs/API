package products

import (
	"database/sql"
	"github.com/curt-labs/GoAPI/helpers/database"
	"github.com/curt-labs/GoAPI/helpers/sortutil"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"strings"
)

var (
	getVcdbVehicleIDWithSubmodel = `
		select v.VehicleID from Vehicle as v
		join BaseVehicle as bv on v.BaseVehicleID = bv.BaseVehicleID
		join Make as ma on bv.MakeID = ma.MakeID
		join Model as mo on bv.ModelID = mo.ModelID
		join Submodel as s on v.SubmodelID = s.SubmodelID
		where bv.YearID = ? && ma.MakeName = ? && mo.ModelName = ? && s.SubmodelName = ?
		limit 1`
	getVcdbVehicleID = `
		select v.VehicleID from Vehicle as v
		join BaseVehicle as bv on v.BaseVehicleID = bv.BaseVehicleID
		join Make as ma on bv.MakeID = ma.MakeID
		join Model as mo on bv.ModelID = mo.ModelID
		where bv.YearID = ? && ma.MakeName = ? && mo.ModelName = ? && (v.SubmodelID = 0 || v.SubmodelID is null)
		limit 1`
	getVehicleParts = `
		select distinct vp.PartNumber
		from vcdb_Vehicle as v
		join Submodel as s on v.SubModelID = s.ID
		join BaseVehicle as bv on v.BaseVehicleID = bv.ID
		join vcdb_Model as mo on bv.ModelID = mo.ID
		join vcdb_Make as ma on bv.MakeID = ma.ID
		join vcdb_VehiclePart as vp on v.ID = vp.VehicleID
		join Part as p on vp.PartNumber = p.partID
		where (p.status = 800 || p.status = 900) && 
		bv.YearID = ? && ma.MakeName = ? &&
		mo.ModelName = ? && s.SubmodelName = ? &&
		(v.ConfigID = 0 || v.ConfigID is null)
		&& p.brandID in(?)
		order by vp.PartNumber`
	getBaseVehicleParts = `
		select distinct vp.PartNumber
		from vcdb_Vehicle as v
		join BaseVehicle as bv on v.BaseVehicleID = bv.ID
		join vcdb_Model as mo on bv.ModelID = mo.ID
		join vcdb_Make as ma on bv.MakeID = ma.ID
		join vcdb_VehiclePart as vp on v.ID = vp.VehicleID
		join Part as p on vp.PartNumber = p.partID
		where (p.status = 800 || p.status = 900) &&
		bv.YearID = ? && ma.MakeName = ? &&
		mo.ModelName = ? && (v.SubmodelID = 0 || v.SubmodelID is null) &&
		(v.ConfigID = 0 || v.ConfigID is null)
		&& p.brandID in (?)
		order by vp.PartNumber`
)

type Vehicle struct {
	Base           BaseVehicle     `json:"base" xml:"base"`
	Submodel       string          `json:"submodel" xml:"submodel"`
	Configurations []Configuration `json:"configurations" xml:"configurations"`
}

type Lookup struct {
	Years          []int                 `json:"available_years,omitempty" xml:"available_years,omitempty"`
	Makes          []string              `json:"available_makes,omitempty" xml:"available_makes,omitempty"`
	Models         []string              `json:"available_models,omitempty" xml:"available_models,omitempty"`
	Submodels      []string              `json:"available_submodels,omitempty" xml:"available_submodels,omitempty"`
	Configurations []ConfigurationOption `json:"available_configurations,omitempty" xml:"available_configurations,omitempty"`
	Vehicle        Vehicle               `json:"vehicle" xml:"vehicle"`
	Parts          []Part                `json:"parts" xml:"parts"`
	Filter         interface{}           `json:"filter" xml:"filter"`
	Pagination     Pagination            `json:"pagination" xml:"pagination"`
	CustomerKey    string                `json:"-" xml:"-"`
	Brands         []int                 `json:"-" xml:"-"`
}

type Pagination struct {
	TotalItems    int `json:"total_items" xml:"total_items"`
	ReturnedCount int `json:"returned_count" xml:"returned_count"`
	Page          int `json:"page" xml:"page"`
	PerPage       int `json:"per_page" xml:"per_page"`
	TotalPages    int `json:"total_pages" xml:"total_pages"`
}

func (l *Lookup) LoadParts(ch chan []Part) {
	parts := make([]Part, 0)

	vehicleChan := make(chan error)
	baseVehicleChan := make(chan error)
	go l.loadVehicleParts(vehicleChan)
	go l.loadBaseVehicleParts(baseVehicleChan)

	if len(l.Vehicle.Configurations) > 0 {
		configs, err := l.Vehicle.getDefinedConfigurations(l.CustomerKey)
		if err != nil || configs == nil {
			ch <- parts
			return
		}

		chosenValArr := make(map[string]string, 0)
		for _, config := range l.Vehicle.Configurations {
			chosenValArr[strings.ToLower(config.Value)] = strings.TrimSpace(strings.ToLower(config.Value))
		}

		for _, config := range *configs {
			// configValArr := make(map[string]string, 0)
			matches := true
			for _, val := range config {

				v := strings.TrimSpace(strings.ToLower(val.Value))

				if _, ok := chosenValArr[v]; !ok {
					matches = false
				}
			}
			if matches {
				for _, partID := range config[0].Parts {
					p := Part{ID: partID}
					l.Parts = append(l.Parts, p)
				}
			}
		}
	}

	<-vehicleChan
	<-baseVehicleChan
	removeDuplicates(&l.Parts)

	parts = make([]Part, 0)
	for i, p := range l.Parts {
		if err := p.Get(l.CustomerKey); err == nil && p.ShortDesc != "" {
			parts = append(parts, p)
		} else if len(parts) > 0 {
			parts = append(parts[:i], parts[i+1:]...)
		}
	}
	l.Parts = parts

	sortutil.AscByField(l.Parts, "ID")

	l.Pagination = Pagination{
		TotalItems:    len(l.Parts),
		ReturnedCount: len(l.Parts),
		Page:          1,
		PerPage:       len(l.Parts),
		TotalPages:    1,
	}

	ch <- parts
}

func (l *Lookup) loadVehicleParts(ch chan error) {
	stmtBeginning := `select distinct vp.PartNumber
		from vcdb_Vehicle as v
		join Submodel as s on v.SubModelID = s.ID
		join BaseVehicle as bv on v.BaseVehicleID = bv.ID
		join vcdb_Model as mo on bv.ModelID = mo.ID
		join vcdb_Make as ma on bv.MakeID = ma.ID
		join vcdb_VehiclePart as vp on v.ID = vp.VehicleID
		join Part as p on vp.PartNumber = p.partID
		where (p.status = 800 || p.status = 900) && 
		bv.YearID = ? && ma.MakeName = ? &&
		mo.ModelName = ? && s.SubmodelName = ? &&
		(v.ConfigID = 0 || v.ConfigID is null)`
	stmtEnd := `order by vp.PartNumber`
	brandStmt := " && p.brandID in ("
	for _, b := range l.Brands {
		brandStmt += strconv.Itoa(b) + ","
	}
	brandStmt = strings.TrimRight(brandStmt, ",") + ")"
	wholeStmt := stmtBeginning + brandStmt + stmtEnd

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		ch <- err
		return
	}
	defer db.Close()

	stmt, err := db.Prepare(wholeStmt)
	if err != nil {
		ch <- err
		return
	}

	rows, err := stmt.Query(l.Vehicle.Base.Year, l.Vehicle.Base.Make, l.Vehicle.Base.Model, l.Vehicle.Submodel)
	if err != nil || rows == nil {
		ch <- err
		return
	}

	for rows.Next() {
		var p Part
		if err = rows.Scan(&p.ID); err == nil {
			l.Parts = append(l.Parts, p)
		}
	}
	defer rows.Close()
	ch <- nil
	return
}

func (l *Lookup) loadBaseVehicleParts(ch chan error) {
	stmtBeginning := `select distinct vp.PartNumber
		from vcdb_Vehicle as v
		join BaseVehicle as bv on v.BaseVehicleID = bv.ID
		join vcdb_Model as mo on bv.ModelID = mo.ID
		join vcdb_Make as ma on bv.MakeID = ma.ID
		join vcdb_VehiclePart as vp on v.ID = vp.VehicleID
		join Part as p on vp.PartNumber = p.partID
		where (p.status = 800 || p.status = 900) &&
		bv.YearID = ? && ma.MakeName = ? &&
		mo.ModelName = ? && (v.SubmodelID = 0 || v.SubmodelID is null) &&
		(v.ConfigID = 0 || v.ConfigID is null)`
	stmtEnd := `order by vp.PartNumber`
	brandStmt := " && p.brandID in ("
	for _, b := range l.Brands {
		brandStmt += strconv.Itoa(b) + ","
	}
	brandStmt = strings.TrimRight(brandStmt, ",") + ")"
	wholeStmt := stmtBeginning + brandStmt + stmtEnd

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		ch <- err
		return
	}
	defer db.Close()

	stmt, err := db.Prepare(wholeStmt)
	if err != nil {
		ch <- err
		return
	}

	rows, err := stmt.Query(l.Vehicle.Base.Year, l.Vehicle.Base.Make, l.Vehicle.Base.Model)
	if err != nil || rows == nil {
		ch <- err
		return
	}

	for rows.Next() {
		var p Part
		if err = rows.Scan(&p.ID); err == nil {
			l.Parts = append(l.Parts, p)
		}
	}
	defer rows.Close()

	ch <- nil
	return
}

func removeDuplicates(xs *[]Part) {
	found := make(map[int]bool)
	j := 0
	for i, x := range *xs {
		if !found[x.ID] {
			found[x.ID] = true
			(*xs)[j] = (*xs)[i]
			j++
		}
	}
	*xs = (*xs)[:j]
}

func (v *Vehicle) GetVcdbID() (int, error) {
	db, err := sql.Open("mysql", database.VcdbConnectionString())
	if err != nil {
		return 0, err
	}
	defer db.Close()

	var row *sql.Row
	if v.Submodel != "" {
		stmt, err := db.Prepare(getVcdbVehicleIDWithSubmodel)
		if err != nil {
			return 0, err
		}
		defer stmt.Close()

		row = stmt.QueryRow(v.Base.Year, v.Base.Make, v.Base.Model, v.Submodel)
		if row == nil {
			return 0, err
		}
	} else {
		stmt, err := db.Prepare(getVcdbVehicleID)
		if err != nil {
			return 0, err
		}
		defer stmt.Close()

		row = stmt.QueryRow(v.Base.Year, v.Base.Make, v.Base.Model)
		if row == nil {
			return 0, err
		}
	}

	var id int
	err = row.Scan(&id)

	return id, err
}

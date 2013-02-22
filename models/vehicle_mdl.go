package models

import (
	"../helpers/database"
	"fmt"
	"log"
	"strings"
)

type ConfigResponse struct {
	ConfigOption ConfigOption
	Matched      *ProductMatch
}

type ConfigOption struct {
	Type    string
	Options []string
}

type ProductMatch struct {
	Parts  []Part
	Groups []int
}

type Vehicle struct {
	Year                  float64
	Make, Model, Submodel string
	Configuration         []string
	Parts                 []Part
	Groups                []interface{}
}

var (
	db = database.Db

	yearStmt = `select YearID from vcdb_Year order by YearID desc`

	makeStmt = `select distinct ma.MakeName as make from BaseVehicle bv
					join vcdb_Make ma on bv.MakeID = ma.ID
					join vcdb_Vehicle v on bv.ID = v.BaseVehicleID
					join vcdb_VehiclePart vp on v.ID = vp.VehicleID
					where bv.YearID = %f
					order by ma.MakeName`

	modelStmt = `select distinct mo.ModelName as model from BaseVehicle bv
					join vcdb_Model mo on bv.ModelID = mo.ID
					join vcdb_Make ma on bv.MakeID = ma.ID
					join vcdb_Vehicle v on bv.ID = v.BaseVehicleID
					join vcdb_VehiclePart vp on v.ID = vp.VehicleID
					where bv.YearID = %f and ma.MakeName = '%s'
					order by mo.ModelName`

	submodelStmt = `select distinct sm.SubmodelName as submodel from BaseVehicle bv
					join vcdb_Model mo on bv.ModelID = mo.ID
					join vcdb_Make ma on bv.MakeID = ma.ID
					join vcdb_Vehicle v on bv.ID = v.BaseVehicleID
					join Submodel sm on v.SubmodelID = sm.ID
					join vcdb_VehiclsePart vp on v.ID = vp.VehicleID
					where bv.YearID = %f and ma.MakeName = '%s' 
					and mo.ModelName = '%s'`

	configStmt = `select cat.name, ca.value from ConfigAttributeType cat
					join ConfigAttribute ca on cat.ID = ca.ConfigAttributeTypeID
					join VehicleConfigAttribute vca on ca.ID = vca.AttributeID
					join vcdb_Vehicle v on vca.VehicleConfigID = v.ConfigID
					join BaseVehicle bv on v.BaseVehicleID = bv.ID
					join Submodel sm on v.SubModelID = sm.ID
					join vcdb_Make ma on bv.MakeID = ma.ID
					join vcdb_Model mo on bv.ModelID = mo.ID
					where bv.YearID = %f and ma.MakeName = '%s'
					and mo.ModelName = '%s' and sm.SubmodelName = '%s' %s order by cat.sort`

	nestedConfigStmt = `and cat.name not in (
					select cat.name from ConfigAttributeType cat
					join ConfigAttribute ca on cat.ID = ca.ConfigAttributeTypeID
					where ca.value in ('%s'))`

	vehiclePartsStmt = `select distinct vp.PartNumber as part from vcdb_VehiclePart vp
					join vcdb_Vehicle v on vp.VehicleID = v.ID
					left join VehicleConfigAttribute vca on v.ConfigID = vca.VehicleConfigID
					left join ConfigAttribute ca on vca.AttributeID = ca.ID
					left join ConfigAttributeType cat on ca.ConfigAttributeTypeID = cat.ID
					join BaseVehicle bv on v.BaseVehicleID = bv.ID
					left join Submodel sm on v.SubModelID = sm.ID
					join vcdb_Make ma on bv.MakeID = ma.ID
					join vcdb_Model mo on bv.ModelID = mo.ID
					where 
					(bv.YearID = %f and ma.MakeName = '%s'
					and mo.ModelName = '%s')
					or
					(bv.YearID = %f and ma.MakeName = '%s'
					and mo.ModelName = '%s' and sm.SubmodelName = '%s')
					or
					(bv.YearID = %f and ma.MakeName = '%s'
					and mo.ModelName = '%s' and sm.SubmodelName = '%s'
					and ca.value in ('%s'))
					order by part;`

	reverseLookupStmt = `select bv.YearID, ma.MakeName, mo.ModelName, sm.SubmodelName
				from BaseVehicle bv
				join vcdb_Vehicle v on bv.ID = v.BaseVehicleID
				join vcdb_VehiclePart vp on v.ID = vp.VehicleID
				left join Submodel sm on v.SubModelID = sm.ID
				left join vcdb_Make ma on bv.MakeID = ma.ID
				left join vcdb_Model mo on bv.ModelID = mo.ID
				where vp.PartNumber = %d
				order by bv.YearID desc, ma.MakeName, mo.ModelName`

	vehicleNotesStmt = `select distinct n.note from vcdb_VehiclePart vp
					left join Note n on vp.ID = n.vehiclePartID
					join vcdb_Vehicle v on vp.VehicleID = v.ID
					left join VehicleConfigAttribute vca on v.ConfigID = vca.VehicleConfigID
					left join ConfigAttribute ca on vca.AttributeID = ca.ID
					join BaseVehicle bv on v.BaseVehicleID = bv.ID
					left join Submodel sm on v.SubModelID = sm.ID
					join vcdb_Make ma on bv.MakeID = ma.ID
					join vcdb_Model mo on bv.ModelID = mo.ID
					where bv.YearID = %f and ma.MakeName = '%s'
					and mo.ModelName = '%s' and (sm.SubmodelName = '%s' or sm.SubmodelName is null)
					and (ca.value in ('%s') or ca.value is null) and vp.PartNumber = %d;`
)

func (vehicle *Vehicle) GetYears() (opt ConfigOption) {
	db := database.Db

	opt.Type = "Years"

	rows, _, err := db.Query(yearStmt)
	if database.MysqlError(err) {
		return
	}

	years := make([]string, 0)

	for _, row := range rows {
		years = append(years, row.Str(0))
	}
	opt.Options = years
	return
}

func (vehicle *Vehicle) GetMakes() (opt ConfigOption) {
	db := database.Db

	opt.Type = "Makes"

	rows, _, err := db.Query(makeStmt, vehicle.Year)
	if database.MysqlError(err) {
		return
	}

	makes := make([]string, 0)
	for _, row := range rows {
		makes = append(makes, row.Str(0))
	}
	opt.Options = makes
	return

}

func (vehicle *Vehicle) GetModels() (opt ConfigOption) {
	db := database.Db

	opt.Type = "Models"

	rows, _, err := db.Query(modelStmt, vehicle.Year, vehicle.Make)
	if database.MysqlError(err) {
		return
	}

	models := make([]string, 0)
	for _, row := range rows {
		models = append(models, row.Str(0))
	}

	opt.Options = models
	return
}

func (vehicle *Vehicle) GetSubmodels() (opt ConfigOption) {
	db := database.Db

	rows, _, err := db.Query(submodelStmt, vehicle.Year, vehicle.Make, vehicle.Model)
	if database.MysqlError(err) {
		return
	}

	subs := make([]string, 0)
	for _, row := range rows {
		subs = append(subs, row.Str(0))
	}

	opt.Type = "Submodels"
	opt.Options = subs
	return
}

func (vehicle *Vehicle) GetConfiguration() (opt ConfigOption) {
	db := database.Db

	var nested string
	if len(vehicle.Configuration) > 0 {
		nested = fmt.Sprintf(nestedConfigStmt, strings.Join(vehicle.Configuration, "/"))
	}

	rows, _, err := db.Query(configStmt,
		vehicle.Year,
		vehicle.Make,
		vehicle.Model,
		vehicle.Submodel,
		nested)
	if database.MysqlError(err) {
		return
	}

	if len(rows) > 0 {
		config_type := rows[0].Str(0)

		config_vals := make([]string, 0)
		for _, row := range rows {
			if row.Str(0) == config_type {
				config_vals = append(config_vals, row.Str(1))
			} else {
				break
			}
		}

		opt.Type = config_type
		opt.Options = config_vals
	}

	return
}

func (vehicle *Vehicle) GetProductMatch() (match *ProductMatch) {

	match = new(ProductMatch)

	parts, err := vehicle.GetParts()
	if err != nil {
		return
	}

	c := make(chan int)
	var part_objs []Part
	for _, id := range parts {
		go func(pId int) {
			p := Part{
				PartId: pId,
			}
			log.Println(p.PartId)
			p.GetWithVehicle(vehicle)
			part_objs = append(part_objs, p)
			c <- 1
		}(id)
	}

	for _, _ = range parts {
		<-c
	}

	match.Parts = part_objs
	match.Groups = make([]int, 0)

	return
}

func (v *Vehicle) GetParts() (parts []int, err error) {
	db := database.Db

	log.Printf(vehiclePartsStmt,
		v.Year,
		v.Make,
		v.Model,
		v.Year,
		v.Make,
		v.Model,
		v.Submodel,
		v.Year,
		v.Make,
		v.Model,
		v.Submodel,
		strings.Join(v.Configuration, ","))

	rows, _, err := db.Query(vehiclePartsStmt,
		v.Year,
		v.Make,
		v.Model,
		v.Year,
		v.Make,
		v.Model,
		v.Submodel,
		v.Year,
		v.Make,
		v.Model,
		v.Submodel,
		strings.Join(v.Configuration, ","))

	if err != nil {
		return
	}

	for _, row := range rows {
		parts = append(parts, row.Int(0))
	}

	return
}

func (vehicle *Vehicle) GetGroupsByBase() (groups []int) {
	return
}

func (vehicle *Vehicle) GetGroupsBySubmodel() (groups []int) {
	return
}

func (vehicle *Vehicle) GetGroupsByConfig() (groups []int) {
	return
}

func (v *Vehicle) GetNotes(partId int) (notes []string, err error) {
	db := database.Db

	rows, _, err := db.Query(vehicleNotesStmt, v.Year, v.Make, v.Model, v.Submodel, strings.Join(v.Configuration, ","), partId)
	for _, row := range rows {
		notes = append(notes, row.Str(0))
	}
	return
}

func ReverseLookup(partId int) (vehicles []Vehicle, err error) {
	db := database.Db

	rows, res, err := db.Query(reverseLookupStmt, partId)
	if database.MysqlError(err) {
		return
	}

	year := res.Map("YearID")
	make := res.Map("MakeName")
	model := res.Map("ModelName")
	submodel := res.Map("SubmodelName")

	for _, row := range rows {

		if database.MysqlError(err) {
			break
		}
		if row == nil {
			break // end of result
		}

		v := Vehicle{
			Year:     row.Float(year),
			Make:     row.Str(make),
			Model:    row.Str(model),
			Submodel: row.Str(submodel),
		}

		vehicles = append(vehicles, v)

	}
	return
}

func AppendIfMissing(existing []int, slice []int) []int {
	for i, s := range slice {
		for _, ex := range existing {
			if s == ex {
				slice = append(slice[:i], slice[i+1:]...)
				return AppendIfMissing(existing, slice)
			}
		}
	}
	return append(existing, slice...)
}

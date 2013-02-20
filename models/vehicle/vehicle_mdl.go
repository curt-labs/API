package vehicle

import (
	"../../helpers/database"
	"fmt"
	"sort"
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
	Parts  []int
	Groups []int
}

type Vehicle struct {
	Year                  float64
	Make, Model, Submodel string
	Configuration         []string
	Parts                 []interface{}
	Groups                []interface{}
}

type Attribute struct {
	Key, Value string
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
					join vcdb_VehiclePart vp on v.ID = vp.VehicleID
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

	partsBaseStmt = `select distinct vp.PartNumber as part from vcdb_VehiclePart vp
				join vcdb_Vehicle v on vp.VehicleID = v.ID
				join BaseVehicle bv on v.BaseVehicleID = bv.ID
				join vcdb_Model mo on bv.ModelID = mo.ID
				join vcdb_Make ma on bv.MakeID = ma.ID
				where bv.YearID = %f and ma.MakeName = '%s'
				and mo.ModelName = '%s'
				order by part`

	partsSubStmt = `select distinct vp.PartNumber as part from vcdb_VehiclePart vp
					join vcdb_Vehicle v on vp.VehicleID = v.ID
					join Submodel sm on v.SubmodelID = sm.ID
					join BaseVehicle bv on v.BaseVehicleID = bv.ID
					join vcdb_Model mo on bv.ModelID = mo.ID
					join vcdb_Make ma on bv.MakeID = ma.ID
					where bv.YearID = %f and ma.MakeName = '%s'
					and mo.ModelName = '%s' and sm.SubmodelName = '%s'
					order by part`

	partsConfigStmt = `select distinct vp.PartNumber as part from vcdb_VehiclePart vp
					join vcdb_Vehicle v on vp.VehicleID = v.ID
					join VehicleConfigAttribute vca on v.ConfigID = vca.VehicleConfigID
					join ConfigAttribute ca on vca.AttributeID = ca.ID
					join ConfigAttributeType cat on ca.ConfigAttributeTypeID = cat.ID
					join BaseVehicle bv on v.BaseVehicleID = bv.ID
					join Submodel sm on v.SubModelID = sm.ID
					join vcdb_Make ma on bv.MakeID = ma.ID
					join vcdb_Model mo on bv.ModelID = mo.ID
					where bv.YearID = %f and ma.MakeName = '%s'
					and mo.ModelName = '%s' and sm.SubmodelName = '%s'
					and ca.value not in ('%s')
					order by part`
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

	base_parts := make([]int, 0)
	sub_parts := make([]int, 0)
	config_parts := make([]int, 0)

	subChan := make(chan int)
	configChan := make(chan int)

	base_parts = vehicle.GetPartsByBase()

	go func() {
		sub_parts = vehicle.GetPartsBySubmodel()
		subChan <- 1
	}()

	go func() {
		config_parts = vehicle.GetPartsByConfig()
		configChan <- 1
	}()

	<-subChan
	<-configChan

	parts := AppendIfMissing(AppendIfMissing(base_parts, sub_parts), config_parts)
	sort.Ints(parts)

	match.Parts = parts
	match.Groups = make([]int, 0)

	return
}

func (vehicle *Vehicle) GetPartsByBase() (parts []int) {

	db := database.Db

	rows, _, err := db.Query(partsBaseStmt, vehicle.Year, vehicle.Make, vehicle.Model)
	if database.MysqlError(err) {
		return
	}

	parts = make([]int, 0)
	for _, row := range rows {
		parts = append(parts, row.Int(0))
	}

	return
}

func (vehicle *Vehicle) GetPartsBySubmodel() (parts []int) {

	db := database.Db

	rows, _, err := db.Query(partsSubStmt, vehicle.Year, vehicle.Make, vehicle.Model, vehicle.Submodel)
	if database.MysqlError(err) {
		return
	}

	parts = make([]int, 0)
	for _, row := range rows {
		parts = append(parts, row.Int(0))
	}
	return
}

func (vehicle *Vehicle) GetPartsByConfig() (parts []int) {
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

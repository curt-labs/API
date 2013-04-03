package models

import (
	"../helpers/database"
	"../helpers/redis"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
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
	yearStmt = `select YearID from vcdb_Year order by YearID desc`

	makeStmt = `select distinct ma.MakeName as make from BaseVehicle bv
					join vcdb_Make ma on bv.MakeID = ma.ID
					join vcdb_Vehicle v on bv.ID = v.BaseVehicleID
					join vcdb_VehiclePart vp on v.ID = vp.VehicleID
					where bv.YearID = ?
					order by ma.MakeName`

	modelStmt = `select distinct mo.ModelName as model
				from BaseVehicle as bv 
				join vcdb_Make as ma on bv.MakeID = ma.ID 
				join vcdb_Model as mo on bv.ModelID = mo.ID
				join vcdb_Vehicle as v on bv.ID = v.BaseVehicleID
				join vcdb_VehiclePart as vp on v.ID = vp.VehicleID
				where bv.YearID = ? and ma.MakeName = ?
				order by mo.ModelName`

	submodelStmt = `select distinct sm.SubmodelName as submodel from BaseVehicle bv
					join vcdb_Model mo on bv.ModelID = mo.ID
					join vcdb_Make ma on bv.MakeID = ma.ID
					join vcdb_Vehicle v on bv.ID = v.BaseVehicleID
					join Submodel sm on v.SubmodelID = sm.ID
					join vcdb_VehiclePart vp on v.ID = vp.VehicleID
					where bv.YearID = ? and ma.MakeName = ? 
					and mo.ModelName = ?`

	configStmt = `select cat.name, ca.value from ConfigAttributeType cat
					join ConfigAttribute ca on cat.ID = ca.ConfigAttributeTypeID
					join VehicleConfigAttribute vca on ca.ID = vca.AttributeID
					join vcdb_Vehicle v on vca.VehicleConfigID = v.ConfigID
					join BaseVehicle bv on v.BaseVehicleID = bv.ID
					join Submodel sm on v.SubModelID = sm.ID
					join vcdb_Make ma on bv.MakeID = ma.ID
					join vcdb_Model mo on bv.ModelID = mo.ID
					where bv.YearID = ? and ma.MakeName = ?
					and mo.ModelName = ? and sm.SubmodelName = ? %s
					group by ca.value
					order by cat.sort`

	nestedConfigBegin = `and cat.name not in (
					select cat.name from ConfigAttributeType cat
					join ConfigAttribute ca on cat.ID = ca.ConfigAttributeTypeID
					where ca.value in (`
	nestedConfigEnd = `))`

	vehiclePartsStmt = `select distinct vp.PartNumber as part from vcdb_VehiclePart vp
					join Part as p on vp.PartNumber = p.partID
					join vcdb_Vehicle v on vp.VehicleID = v.ID
					left join VehicleConfigAttribute vca on v.ConfigID = vca.VehicleConfigID
					left join ConfigAttribute ca on vca.AttributeID = ca.ID
					left join ConfigAttributeType cat on ca.ConfigAttributeTypeID = cat.ID
					join BaseVehicle bv on v.BaseVehicleID = bv.ID
					left join Submodel sm on v.SubModelID = sm.ID
					join vcdb_Make ma on bv.MakeID = ma.ID
					join vcdb_Model mo on bv.ModelID = mo.ID
					where p.status in (800,900)
					and (
						(bv.YearID = ? and ma.MakeName = ?
						and mo.ModelName = ?)
						or
						(bv.YearID = ? and ma.MakeName = ?
						and mo.ModelName = ? and sm.SubmodelName = ?)
						or
						(bv.YearID = ? and ma.MakeName = ?
						and mo.ModelName = ? and sm.SubmodelName = ?
						and ca.value in (?))
					)
					order by part;`

	reverseLookupStmt = `select bv.YearID, ma.MakeName, mo.ModelName, sm.SubmodelName
				from BaseVehicle bv
				join vcdb_Vehicle v on bv.ID = v.BaseVehicleID
				join vcdb_VehiclePart vp on v.ID = vp.VehicleID
				left join Submodel sm on v.SubModelID = sm.ID
				left join vcdb_Make ma on bv.MakeID = ma.ID
				left join vcdb_Model mo on bv.ModelID = mo.ID
				where vp.PartNumber = ?
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
					where bv.YearID = ? and ma.MakeName = ?
					and mo.ModelName = ? and (sm.SubmodelName = ? or sm.SubmodelName is null)
					and (ca.value in (?) or ca.value is null) and vp.PartNumber = ?;`
)

func (vehicle *Vehicle) GetYears() (opt ConfigOption) {

	years := make([]string, 0)

	year_bytes, err := redis.RedisClient.Get("vehicle_years")
	if err != nil || len(year_bytes) == 0 {
		rows, _, err := database.Db.Query(yearStmt)
		if database.MysqlError(err) {
			return
		}

		for _, row := range rows {
			years = append(years, row.Str(0))
		}

		if year_bytes, err = json.Marshal(years); err == nil {
			redis.RedisClient.Set("vehicle_years", year_bytes)
			redis.RedisClient.Expire("vehicle_years", int64(time.Duration.Hours(24)))
		}
	} else {
		_ = json.Unmarshal(year_bytes, &years)
	}

	opt.Type = "Years"
	opt.Options = years
	return
}

func (vehicle *Vehicle) GetMakes() (opt ConfigOption) {

	qry, err := database.Db.Prepare(makeStmt)
	if err != nil {
		return
	}

	rows, _, err := qry.Exec(vehicle.Year)
	if database.MysqlError(err) {
		return
	}

	makes := make([]string, 0)
	for _, row := range rows {
		makes = append(makes, row.Str(0))
	}

	opt.Type = "Makes"
	opt.Options = makes
	return
}

func (vehicle *Vehicle) GetModels() (opt ConfigOption) {
	qry, err := database.Db.Prepare(modelStmt)
	if err != nil {
		return
	}

	params := struct {
		Year float64
		Make string
	}{
		vehicle.Year,
		vehicle.Make,
	}

	rows, _, err := qry.Exec(params)
	if database.MysqlError(err) {
		return
	}

	models := make([]string, 0)
	for _, row := range rows {
		models = append(models, row.Str(0))
	}

	opt.Type = "Models"
	opt.Options = models
	return
}

func (vehicle *Vehicle) GetSubmodels() (opt ConfigOption) {
	qry, err := database.Db.Prepare(submodelStmt)
	if err != nil {
		return
	}

	params := struct {
		Year  float64
		Make  string
		Model string
	}{
		vehicle.Year,
		vehicle.Make,
		vehicle.Model,
	}

	rows, _, err := qry.Exec(params)
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

	var nested string
	if len(vehicle.Configuration) > 0 {
		nested = nestedConfigBegin + " ? " + nestedConfigEnd
	}

	stmt := fmt.Sprintf(configStmt, nested)

	qry, err := database.Db.Prepare(stmt)
	if err != nil {
		return
	}

	params := struct {
		Year          float64
		Make          string
		Model         string
		Submodel      string
		Configuration string
	}{
		vehicle.Year,
		vehicle.Make,
		vehicle.Model,
		vehicle.Submodel,
		strings.Join(vehicle.Configuration, ","),
	}

	rows, _, err := qry.Exec(params)

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

func (vehicle *Vehicle) GetProductMatch(key string) (match *ProductMatch) {

	log.Println(time.Now())

	match = new(ProductMatch)

	parts, err := vehicle.GetParts()
	if err != nil {
		return
	}

	populated, err := GetWithVehicleByGroup(parts, vehicle, key)

	var ps []Part
	for _, v := range populated {
		ps = append(ps, v)
	}

	match.Parts = ps
	match.Groups = make([]int, 0)

	return
}

func (v *Vehicle) GetParts() (parts map[int]Part, err error) {
	qry, err := database.Db.Prepare(vehiclePartsStmt)
	if err != nil {
		return
	}

	params := struct {
		Year1     float64
		Make1     string
		Model1    string
		Year2     float64
		Make2     string
		Model2    string
		Submodel1 string
		Year3     float64
		Make3     string
		Model3    string
		Submodel2 string
		Config    string
	}{
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
		strings.Join(v.Configuration, ","),
	}

	rows, _, err := qry.Exec(params)
	if err != nil {
		return
	}

	parts = make(map[int]Part, len(rows))

	for _, row := range rows {
		parts[row.Int(0)] = Part{PartId: row.Int(0)}
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
	qry, err := database.Db.Prepare(reverseLookupStmt)
	if err != nil {
		return
	}

	rows, res, err := qry.Exec(partId)
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

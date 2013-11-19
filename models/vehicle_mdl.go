package models

import (
	"encoding/json"
	"fmt"
	"github.com/curt-labs/GoAPI/helpers/database"
	"github.com/curt-labs/GoAPI/helpers/redis"
	"strconv"
	"strings"
	"time"
)

type Lookup struct {
	Parts  []*Part
	Groups []int

	Vehicle Vehicle
}

type ConfigResponse struct {
	ConfigOption ConfigOption
	Matched      *ProductMatch
}

type ConfigOption struct {
	Type    string
	Options []string
}

type ProductMatch struct {
	Parts  []*Part
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
						and mo.ModelName = ? and vca.ID is null)
						or
						(bv.YearID = ? and ma.MakeName = ?
						and mo.ModelName = ? and sm.SubmodelName = ? and vca.ID is null)
						or
						(bv.YearID = ? and ma.MakeName = ?
						and mo.ModelName = ? and sm.SubmodelName = ?
						and ca.value in (`
	vehiclePartsStmtEnd = `))) order by part;`

	vehicleConnectorStmt = `select distinct vp.PartNumber as part from vcdb_VehiclePart vp
					join Part as p on vp.PartNumber = p.partID
					join Class as pc on p.classID = pc.classID
					join vcdb_Vehicle v on vp.VehicleID = v.ID
					left join VehicleConfigAttribute vca on v.ConfigID = vca.VehicleConfigID
					left join ConfigAttribute ca on vca.AttributeID = ca.ID
					left join ConfigAttributeType cat on ca.ConfigAttributeTypeID = cat.ID
					join BaseVehicle bv on v.BaseVehicleID = bv.ID
					left join Submodel sm on v.SubModelID = sm.ID
					join vcdb_Make ma on bv.MakeID = ma.ID
					join vcdb_Model mo on bv.ModelID = mo.ID
					where p.status in (800,900) and UPPER(pc.class) = 'WIRING' `
	vehicleConnectorStmtWithConfig = `and (
						(bv.YearID = ? and ma.MakeName = ?
						and mo.ModelName = ? and vca.ID is null)
						or
						(bv.YearID = ? and ma.MakeName = ?
						and mo.ModelName = ? and sm.SubmodelName = ? and vca.ID is null)
						or
						(bv.YearID = ? and ma.MakeName = ?
						and mo.ModelName = ? and sm.SubmodelName = ?
						and ca.value in (`
	vehicleConnectorStmtWithConfigEnd = `))) order by part`

	vehicleConnectorStmtWithoutConfig = `and (
						(bv.YearID = ? and ma.MakeName = ?
						and mo.ModelName = ? and vca.ID is null)
						or
						(bv.YearID = ? and ma.MakeName = ?
						and mo.ModelName = ? and sm.SubmodelName = ? and vca.ID is null)) order by part`

	reverseLookupStmt = `select v.ID,bv.YearID, ma.MakeName, mo.ModelName, sm.SubmodelName, ca.value
				from BaseVehicle bv
				join vcdb_Vehicle v on bv.ID = v.BaseVehicleID
				join vcdb_VehiclePart vp on v.ID = vp.VehicleID
				left join Submodel sm on v.SubModelID = sm.ID
				left join vcdb_Make ma on bv.MakeID = ma.ID
				left join vcdb_Model mo on bv.ModelID = mo.ID
				left join VehicleConfigAttribute vca on v.ConfigID = vca.VehicleConfigID
				left join ConfigAttribute ca on vca.AttributeID = ca.ID
				where vp.PartNumber = ?
				order by bv.YearID desc, ma.MakeName, mo.ModelName, sm.SubmodelName`

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
					and (ca.value in (`
	vehicleNotesStmtEnd = `) or ca.value is null) and vp.PartNumber = ?;`

	vehicleNotesStmt_Grouped = `select distinct n.note, vp.PartNumber from vcdb_VehiclePart vp
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
					and (ca.value in (`
	vehicleNotesStmtMiddle_Grouped = `) or ca.value is null) and vp.PartNumber IN (`
	vehicleNotesStmtEnd_Grouped    = `)`
)

func (lookup *Lookup) GetYears() (opt ConfigOption) {

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

func (lookup *Lookup) GetMakes() (opt ConfigOption) {

	qry, err := database.Db.Prepare(makeStmt)
	if err != nil {
		return
	}

	rows, _, err := qry.Exec(lookup.Vehicle.Year)
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

func (lookup *Lookup) GetModels() (opt ConfigOption) {
	qry, err := database.Db.Prepare(modelStmt)
	if err != nil {
		return
	}

	params := struct {
		Year float64
		Make string
	}{
		lookup.Vehicle.Year,
		lookup.Vehicle.Make,
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

func (lookup *Lookup) GetSubmodels() (opt ConfigOption) {
	qry, err := database.Db.Prepare(submodelStmt)
	if err != nil {
		return
	}

	params := struct {
		Year  float64
		Make  string
		Model string
	}{
		lookup.Vehicle.Year,
		lookup.Vehicle.Make,
		lookup.Vehicle.Model,
	}

	rows, _, err := qry.Exec(params)
	if database.MysqlError(err) {
		return
	}

	subs := make([]string, 0)

	for _, row := range rows {
		subs = append(subs, strings.TrimSpace(row.Str(0)))
	}

	opt.Type = "Submodels"
	opt.Options = subs
	return
}

func (lookup *Lookup) GetConfiguration() (opt ConfigOption) {

	var nested string
	if len(lookup.Vehicle.Configuration) > 0 {
		nested = nestedConfigBegin
		for i, c := range lookup.Vehicle.Configuration {
			if len(c) > 0 {
				nested = nested + "'" + database.Db.Escape(c) + "'"
				if i < len(lookup.Vehicle.Configuration)-1 {
					nested = nested + ","
				}
			}

		}
		nested = nested + nestedConfigEnd
	}

	stmt := fmt.Sprintf(configStmt, nested)

	qry, err := database.Db.Prepare(stmt)
	if err != nil {
		return
	}

	params := struct {
		Year     float64
		Make     string
		Model    string
		Submodel string
	}{
		lookup.Vehicle.Year,
		lookup.Vehicle.Make,
		lookup.Vehicle.Model,
		lookup.Vehicle.Submodel,
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
				val := strings.TrimSpace(row.Str(1))
				config_vals = append(config_vals, val)
			} else {
				break
			}
		}

		opt.Type = config_type
		opt.Options = config_vals
	}

	return
}

func (lookup *Lookup) GetProductMatch(key string) (match *ProductMatch) {

	match = new(ProductMatch)

	err := lookup.GetParts()
	if err != nil {
		return
	}

	err = lookup.GetWithVehicle(key)

	var ps []*Part
	for _, v := range lookup.Parts {
		ps = append(ps, v)
	}

	match.Parts = ps
	match.Groups = make([]int, 0)

	return
}

func (lookup *Lookup) GetParts() error {

	stmt := vehiclePartsStmt
	if len(lookup.Vehicle.Configuration) > 0 {
		for i, c := range lookup.Vehicle.Configuration {
			stmt = stmt + "'" + database.Db.Escape(c) + "'"
			if i < len(lookup.Vehicle.Configuration)-1 {
				stmt = stmt + ","
			}
		}
	} else {
		stmt = stmt + "''"
	}
	stmt = stmt + vehiclePartsStmtEnd

	qry, err := database.Db.Prepare(stmt)
	if err != nil {
		return err
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
	}{
		lookup.Vehicle.Year,
		lookup.Vehicle.Make,
		lookup.Vehicle.Model,
		lookup.Vehicle.Year,
		lookup.Vehicle.Make,
		lookup.Vehicle.Model,
		lookup.Vehicle.Submodel,
		lookup.Vehicle.Year,
		lookup.Vehicle.Make,
		lookup.Vehicle.Model,
		lookup.Vehicle.Submodel,
	}

	rows, _, err := qry.Exec(params)
	if err != nil {
		return err
	}

	for _, row := range rows {
		lookup.Parts = append(lookup.Parts, &Part{PartId: row.Int(0)})
	}

	return nil
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

	stmt := vehicleNotesStmt
	if len(v.Configuration) > 0 {
		for i, c := range v.Configuration {
			stmt = stmt + "'" + database.Db.Escape(c) + "'"
			if i < len(v.Configuration)-1 {
				stmt = stmt + ","
			}
		}
		stmt = stmt + vehicleNotesStmtEnd
	}

	qry, err := database.Db.Prepare(stmt)
	if err != nil {
		return
	}

	params := struct {
		Year     float64
		Make     string
		Model    string
		Submodel string
		Part     int
	}{
		v.Year,
		v.Make,
		v.Model,
		v.Submodel,
		partId,
	}

	rows, _, err := qry.Exec(params)
	if database.MysqlError(err) || len(rows) == 0 {
		return
	}

	for _, row := range rows {
		notes = append(notes, row.Str(0))
	}
	return
}

func (lookup *Lookup) GetNotes() error {

	var ids []string
	for _, p := range lookup.Parts {
		ids = append(ids, strconv.Itoa(p.PartId))
	}
	if len(ids) == 0 {
		return nil
	}

	stmt := vehicleNotesStmt_Grouped
	if len(lookup.Vehicle.Configuration) > 0 {
		for i, c := range lookup.Vehicle.Configuration {
			stmt = stmt + "'" + database.Db.Escape(c) + "'"
			if i < len(lookup.Vehicle.Configuration)-1 {
				stmt = stmt + ","
			}
		}
		stmt = stmt + vehicleNotesStmtMiddle_Grouped
	}

	if len(lookup.Parts) > 0 {
		for i, p := range lookup.Parts {
			stmt = stmt + "'" + strconv.Itoa(p.PartId) + "'"
			if i < len(lookup.Parts)-1 {
				stmt = stmt + ","
			}
		}
		stmt = stmt + vehicleNotesStmtEnd_Grouped
	}

	qry, err := database.Db.Prepare(stmt)
	if database.MysqlError(err) {
		return err
	}

	params := struct {
		Year     float64
		Make     string
		Model    string
		Submodel string
	}{
		lookup.Vehicle.Year,
		lookup.Vehicle.Make,
		lookup.Vehicle.Model,
		lookup.Vehicle.Submodel,
	}

	rows, _, err := qry.Exec(params)
	if database.MysqlError(err) || len(rows) == 0 {
		return err
	}

	notes := make(map[int][]string, len(lookup.Parts))

	for _, row := range rows {
		pId := row.Int(1)
		notes[pId] = append(notes[pId], row.Str(0))
	}

	for _, p := range lookup.Parts {
		p.VehicleAttributes = notes[p.PartId]
	}

	return nil
}

func (lookup *Lookup) GetConnector(key string) error {

	stmt := vehicleConnectorStmt
	var params interface{}
	if len(lookup.Vehicle.Configuration) > 0 {
		stmt = stmt + vehicleConnectorStmtWithConfig
		for i, c := range lookup.Vehicle.Configuration {
			stmt = stmt + "'" + database.Db.Escape(c) + "'"
			if i < len(lookup.Vehicle.Configuration)-1 {
				stmt = stmt + ","
			}
		}
		stmt = stmt + vehicleConnectorStmtWithConfigEnd
		params = struct {
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
		}{
			lookup.Vehicle.Year,
			lookup.Vehicle.Make,
			lookup.Vehicle.Model,
			lookup.Vehicle.Year,
			lookup.Vehicle.Make,
			lookup.Vehicle.Model,
			lookup.Vehicle.Submodel,
			lookup.Vehicle.Year,
			lookup.Vehicle.Make,
			lookup.Vehicle.Model,
			lookup.Vehicle.Submodel,
		}
	} else {
		stmt = stmt + vehicleConnectorStmtWithoutConfig
		params = struct {
			Year1     float64
			Make1     string
			Model1    string
			Year2     float64
			Make2     string
			Model2    string
			Submodel1 string
		}{
			lookup.Vehicle.Year,
			lookup.Vehicle.Make,
			lookup.Vehicle.Model,
			lookup.Vehicle.Year,
			lookup.Vehicle.Make,
			lookup.Vehicle.Model,
			lookup.Vehicle.Submodel,
		}
	}

	qry, err := database.Db.Prepare(stmt)
	if err != nil {
		return err
	}

	rows, _, err := qry.Exec(params)
	if err != nil {
		return err
	}

	for _, row := range rows {
		lookup.Parts = append(lookup.Parts, &Part{PartId: row.Int(0)})
	}

	err = lookup.GetWithVehicle(key)

	var ps []*Part
	for _, v := range lookup.Parts {
		ps = append(ps, v)
	}

	lookup.Parts = ps
	return nil
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

	id := res.Map("ID")
	year := res.Map("YearID")
	vMake := res.Map("MakeName")
	model := res.Map("ModelName")
	submodel := res.Map("SubmodelName")
	config_val := res.Map("value")

	vehicleArray := make(map[int]Vehicle, 0)

	for _, row := range rows {

		if database.MysqlError(err) || row == nil {
			break
		}

		v, ok := vehicleArray[row.Int(id)]
		if ok {
			// Vehicle Record exists for this ID
			// so we'll simply append this configuration variable
			v.Configuration = append(v.Configuration, row.Str(config_val))
		} else {
			// New Vehicle record
			v = Vehicle{
				Year:          row.Float(year),
				Make:          row.Str(vMake),
				Model:         row.Str(model),
				Submodel:      row.Str(submodel),
				Configuration: []string{row.Str(config_val)},
			}
		}
		vehicleArray[row.Int(id)] = v
	}

	for _, v := range vehicleArray {
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

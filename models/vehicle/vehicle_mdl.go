package vehicle

import (
	"../../helpers/database"
	// "github.com/ziutek/mymysql/mysql"
	// _ "github.com/ziutek/mymysql/thrsafe"
	// "log"
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
	Parts  []interface{}
	Groups []interface{}
}

type Vehicle struct {
	Year                  float64
	Make, Model, Submodel string
	Configuration         []string
	Parts                 []interface{}
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

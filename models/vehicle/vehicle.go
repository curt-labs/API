package vehicle

import (
	"github.com/curt-labs/GoAPI/helpers/database"
)

type ConfigOption struct {
	Type    string
	Options []string
}

type Vehicle struct {
	Year                  float64
	Make, Model, Submodel string
	Configuration         []string
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

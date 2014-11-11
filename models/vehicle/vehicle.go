package vehicle

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/curt-labs/GoAPI/helpers/redis"

	"github.com/curt-labs/GoAPI/helpers/api"
	"github.com/curt-labs/GoAPI/helpers/database"
	_ "github.com/go-sql-driver/mysql"
)

type ConfigOption struct {
	Type    string
	Options []string
}

type Vehicle struct {
	ID                    int
	Year                  int
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

	reverseLookupStmt = `select
												v.ID,bv.YearID, ma.MakeName, mo.ModelName,
												IFNULL(sm.SubmodelName, "") as SubmodelName,
												IFNULL(ca.value, "") as value
												from BaseVehicle bv
												join vcdb_Vehicle v on bv.ID = v.BaseVehicleID
												join vcdb_VehiclePart vp on v.ID = vp.VehicleID
												left join Submodel sm on v.SubModelID = sm.ID
												left join vcdb_Make ma on bv.MakeID = ma.ID
												left join vcdb_Model mo on bv.ModelID = mo.ID
												left join VehicleConfigAttribute vca on v.ConfigID = vca.VehicleConfigID
												left join ConfigAttribute ca on vca.AttributeID = ca.ID
												where vp.PartNumber = ?
												group by bv.YearID, ma.MakeName, mo.ModelName, sm.SubmodelName, ca.value
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
	qrystmt := vehicleNotesStmt
	if len(v.Configuration) > 0 {
		for i, c := range v.Configuration {
			qrystmt = qrystmt + "'" + api_helpers.Escape(c) + "'"
			if i < len(v.Configuration)-1 {
				qrystmt = qrystmt + ","
			}
		}
		qrystmt = qrystmt + vehicleNotesStmtEnd
	}

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return
	}
	defer db.Close()

	stmt, err := db.Prepare(qrystmt)
	if err != nil {
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(v.Year, v.Make, v.Model, v.Submodel, partId)
	if err != nil {
		return
	}

	for rows.Next() {
		var note string
		if err = rows.Scan(&note); err == nil {
			notes = append(notes, note)
		}
	}
	defer rows.Close()

	return
}

func ReverseLookup(partId int) (vehicles []Vehicle, err error) {

	redis_key := fmt.Sprintf("part:%d:vehicles", partId)

	data, err := redis.Get(redis_key)
	if err == nil && len(data) > 0 {
		if err = json.Unmarshal(data, vehicles); err != nil {
			return
		}
	}

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return
	}
	defer db.Close()

	stmt, err := db.Prepare(reverseLookupStmt)
	if err != nil {
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(partId)
	if err != nil {
		return
	}

	vehicleArray := make(map[int]Vehicle, 0)

	for rows.Next() {
		var id int
		var year int
		var ma string
		var mo string
		var sm string
		var configVal string

		if err = rows.Scan(&id, &year, &ma, &mo, &sm, &configVal); err != nil {
			break
		}

		v, ok := vehicleArray[id]
		if ok {
			// Vehicle Record exists for this ID
			// so we'll simply append this configuration variable
			v.Configuration = append(v.Configuration, configVal)
		} else {
			v = Vehicle{
				ID:            id,
				Year:          year,
				Make:          ma,
				Model:         mo,
				Submodel:      sm,
				Configuration: []string{configVal},
			}
		}
		vehicleArray[v.ID] = v
	}
	defer rows.Close()

	for _, v := range vehicleArray {
		vehicles = append(vehicles, v)
	}

	go redis.Setex(redis_key, vehicles, redis.CacheTimeout)

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

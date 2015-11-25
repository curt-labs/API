package vehicle

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/curt-labs/API/helpers/api"
	"github.com/curt-labs/API/helpers/database"
	"github.com/curt-labs/API/helpers/redis"
	_ "github.com/go-sql-driver/mysql"
)

type ConfigOption struct {
	Type    string   `json:"type,omitempty" xml:"type,omitempty"`
	Options []string `json:"options,omitempty" xml:"options,omitempty"`
}

type Vehicle struct {
	ID            int      `json:"id,omitempty" xml:"id,omitempty"`
	Year          int      `json:"year,omitempty" xml:"year,omitempty"`
	Make          string   `json:"make,omitempty" xml:"make,omitempty"`
	Model         string   `json:"model,omitempty" xml:"model,omitempty"`
	Submodel      string   `json:"submodel,omitempty" xml:"submodel,omitempty"`
	Configuration []Config `json:"configuration,omitempty" xml:"configuration,omitempty"`
}

type Config struct {
	Type  string `json:"type,omitempty" xml:"type,omitempty"`
	Value string `json:"value,omitempty" xml:"value,omitempty"`
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
												IFNULL(ca.value, "") as value,
												IFNULL(cat.name,"") as type
												from BaseVehicle bv
												join vcdb_Vehicle v on bv.ID = v.BaseVehicleID
												join vcdb_VehiclePart vp on v.ID = vp.VehicleID
												left join Submodel sm on v.SubModelID = sm.ID
												left join vcdb_Make ma on bv.MakeID = ma.ID
												left join vcdb_Model mo on bv.ModelID = mo.ID
												left join VehicleConfigAttribute vca on v.ConfigID = vca.VehicleConfigID
												left join ConfigAttribute ca on vca.AttributeID = ca.ID
												left join ConfigAttributeType cat on cat.ID =ca.ConfigAttributeTypeID
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
					and mo.ModelName = ? and (sm.SubmodelName = ? or sm.SubmodelName is null) and (
					`
	vehicleNotesStmtEnd = `  ca.value is null) and vp.PartNumber = ?;`

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

	getVehicleNewStmt = `select v.ID, bv.AAIABaseVehicleID, bv.YearID, ma.ID, ma.MakeName, mo.ID, mo.ModelName, sm.AAIASubmodelID, sm.SubmodelName,
	 	 group_concat(ca.ID)
	 	from vcdb_Vehicle as v
	 	join BaseVehicle as bv on bv.ID = v.BaseVehicleID
	 	left join vcdb_Model as mo on mo.ID = bv.ModelID
	 	left join vcdb_Make as ma on ma.ID = bv.MakeID
	 	left join Submodel as sm on sm.ID = v.SubmodelID
	 	left join VehicleConfigAttribute as vca on vca.VehicleConfigID = v.ConfigID
	 	left join ConfigAttribute as ca on ca.ID = vca.AttributeID
	 	left join ConfigAttributeType as cat on cat.ID = ca.ConfigAttributeTypeID
	 	where bv.AAIABaseVehicleID = ?
	 	and sm.AAIASubmodelID = ?
	 	group by v.ID`

	getConfigsStmt = `select cat.ID, cat.name,  ca.ID, ca.value
			from ConfigAttribute as ca
			join ConfigAttributeType as cat on cat.ID = ca.ConfigAttributeTypeID
			where ca.ID = ?`

	getVehicleConfigs = `select cat.name , ca.value
		from VehicleConfigAttribute as vca 
		left join ConfigAttribute as ca on ca.ID = vca.AttributeID
		left join ConfigAttributeType as cat on cat.ID = ca.ConfigAttributeTypeID 
		left join vcdb_Vehicle as v on v.ConfigID = vca.VehicleConfigID
		left join BaseVehicle as bv on bv.ID = v.BaseVehicleID
		left join Submodel as sm on sm.ID = v.SubmodelID
		where bv.AAIABaseVehicleID = ?  
		and sm.AAIASubmodelID = ?`
	getVehicleByBaseStmt = `select v.ID, bv.AAIABaseVehicleID, bv.YearID, ma.ID, ma.MakeName, mo.ID, mo.ModelName 
	 	from vcdb_Vehicle as v
	 	join BaseVehicle as bv on bv.ID = v.BaseVehicleID
	 	left join vcdb_Model as mo on mo.ID = bv.ModelID
	 	left join vcdb_Make as ma on ma.ID = bv.MakeID
	 	left join Submodel as sm on sm.ID = v.SubmodelID
	 	where bv.AAIABaseVehicleID = ?
	 	and (sm.AAIASubmodelID = 20 or v.SubmodelID is null or v.SubmodelID = 0)`
	getVehicleBySubmodelStmt = `select v.ID, bv.AAIABaseVehicleID, bv.YearID, ma.ID, ma.MakeName, mo.ID, mo.ModelName, sm.AAIASubmodelID, sm.SubmodelName
	 	from vcdb_Vehicle as v
	 	join BaseVehicle as bv on bv.ID = v.BaseVehicleID
	 	left join vcdb_Model as mo on mo.ID = bv.ModelID
	 	left join vcdb_Make as ma on ma.ID = bv.MakeID
	 	left join Submodel as sm on sm.ID = v.SubmodelID
	 	where bv.AAIABaseVehicleID = ?
	 	and sm.AAIASubmodelID = ?`
	getConfigAttributes = `select ID from ConfigAttribute where value = trim(lower(?))`
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
		qrystmt += "  ca.value in (`"
		for i, c := range v.Configuration {
			qrystmt = qrystmt + "'" + api_helpers.Escape(c.Value) + "'"
			if i < len(v.Configuration)-1 {
				qrystmt = qrystmt + ","
			}
		}
		qrystmt += ") or"
	}
	qrystmt = qrystmt + vehicleNotesStmtEnd

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
		if err = json.Unmarshal(data, &vehicles); err == nil {
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
		var configType string

		if err = rows.Scan(&id, &year, &ma, &mo, &sm, &configVal, &configType); err != nil {
			break
		}

		v, ok := vehicleArray[id]
		if ok && configType != "" && configVal != "" {
			// Vehicle Record exists for this ID
			// so we'll simply append this configuration variable
			config := Config{Type: configType, Value: configVal}
			v.Configuration = append(v.Configuration, config)
		} else {
			var config Config
			if configType != "" && configVal != "" {
				config = Config{
					Type:  configType,
					Value: configVal,
				}
			}
			v = Vehicle{
				ID:       id,
				Year:     year,
				Make:     ma,
				Model:    mo,
				Submodel: sm,
			}
			if config.Type != "" && config.Value != "" {
				v.Configuration = append(v.Configuration, config)
			}
		}
		vehicleArray[v.ID] = v
	}
	defer rows.Close()

	vehicles = make([]Vehicle, 0)
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

//For TrucksPlus Aces XML Lookup
//takes aaia base and sub Ids
func GetVehicle(baseId, subId int, configs []string) (Vehicle, error) {
	var err error
	var v Vehicle

	//get config Attribute IDs from configs
	configIds, err := getConfigAttributeIDs(configs)
	if err != nil {
		return v, err
	}

	if subId == 20 || subId == 0 {
		//basevehicle
		return GetVehicleByBase(baseId)
	}
	for _, c := range configIds {
		if c > 0 {
			//config
			return GetVehicleByConfig(baseId, subId, configIds)
		}
	}
	//submodel
	return GetVehicleBySubmodel(baseId, subId)
}

func GetVehicleByBase(baseId int) (Vehicle, error) {
	var err error
	var v Vehicle
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return v, err
	}
	defer db.Close()

	stmt, err := db.Prepare(getVehicleByBaseStmt)
	if err != nil {
		return v, err
	}
	defer stmt.Close()

	var makeId, modelId *int

	err = stmt.QueryRow(baseId).Scan(
		&v.ID,
		&baseId,
		&v.Year,
		&makeId,
		&v.Make,
		&modelId,
		&v.Model,
	)
	return v, err
}

func GetVehicleBySubmodel(baseId, subId int) (Vehicle, error) {
	var err error
	var v Vehicle
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return v, err
	}
	defer db.Close()

	stmt, err := db.Prepare(getVehicleBySubmodelStmt)
	if err != nil {
		return v, err
	}
	defer stmt.Close()

	var makeId, modelId *int
	var submodel *string

	err = stmt.QueryRow(baseId, subId).Scan(
		&v.ID,
		&baseId,
		&v.Year,
		&makeId,
		&v.Make,
		&modelId,
		&v.Model,
		&subId,
		&submodel,
	)
	if err != nil {
		return v, err
	}
	if submodel != nil {
		v.Submodel = *submodel
	}
	return v, err
}

func GetVehicleByConfig(baseId, subId int, configs []int) (Vehicle, error) {
	var err error
	var v Vehicle
	var outputVehicle Vehicle

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return v, err
	}
	defer db.Close()

	stmt, err := db.Prepare(getVehicleNewStmt)
	if err != nil {
		return v, err
	}
	defer stmt.Close()

	//get Base+Submodel
	var makeId, modelId *int
	var submodel, configIDConcat *string

	res, err := stmt.Query(baseId, subId)
	if err != nil {
		return v, err
	}
	for res.Next() {
		err = res.Scan(
			&v.ID,
			&baseId,
			&v.Year,
			&makeId,
			&v.Make,
			&modelId,
			&v.Model,
			&subId,
			&submodel,
			&configIDConcat,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				err = nil
			}
			return outputVehicle, err
		}
		if submodel != nil {
			v.Submodel = *submodel
		}

		//get int map of configs associated with this DB vehicle (maps configId to itself)
		var configsArray []string
		var configsIntArray []int
		configsMap := make(map[int]int)
		if configIDConcat != nil {
			configsArray = strings.Split(*configIDConcat, ",")
			if err != nil {
				return outputVehicle, err
			}
			for _, configStr := range configsArray {
				thisInt, err := strconv.Atoi(configStr)
				if err != nil {
					return outputVehicle, err
				}
				configsIntArray = append(configsIntArray, thisInt) //array of these configIds
				configsMap[thisInt] = thisInt                      //map of these configIds

			}
		}
		notHere := false

		for _, idFromParams := range configs {
			if idFromParams > 0 {
				if _, ok := configsMap[idFromParams]; !ok {
					notHere = true
				}
			}
		}

		if notHere == false {
			//actually get the configurations, assign to the vehicle being returned
			v.Configuration, err = getConfigurations(configsIntArray)
			if err != nil {
				return outputVehicle, err
			}
			outputVehicle = v
			return outputVehicle, nil
		}
	}
	return outputVehicle, err
}

//converts a string array of attr values to an int array of their attribute ids
func getConfigAttributeIDs(configs []string) ([]int, error) {
	var err error
	var conIds []int

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return conIds, err
	}
	defer db.Close()
	for _, configStr := range configs {
		stmt, err := db.Prepare(getConfigAttributes)
		if err != nil {
			return conIds, err
		}
		defer stmt.Close()
		var id int
		configStr = strings.ToLower(strings.TrimSpace(configStr))
		err = stmt.QueryRow(configStr).Scan(&id)
		if err == sql.ErrNoRows {
			err = nil
		}
		if err != nil {
			return conIds, err
		}
		conIds = append(conIds, id)
	}
	return conIds, err
}

//gets an array of Configs from an array of attribute Ids
func getConfigurations(configIds []int) ([]Config, error) {
	var err error
	var configArray []Config
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return configArray, err
	}
	defer db.Close()
	for _, id := range configIds {
		stmt, err := db.Prepare(getConfigsStmt)
		if err != nil {
			return configArray, err
		}
		defer stmt.Close()
		var c Config
		var catId, caId *int
		err = stmt.QueryRow(id).Scan(&catId, &c.Type, &caId, &c.Value)
		if err != nil {
			return configArray, err
		}
		configArray = append(configArray, c)
	}
	return configArray, err
}

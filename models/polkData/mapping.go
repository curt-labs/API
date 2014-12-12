package polkData

import (
	"database/sql"
	"github.com/curt-labs/GoAPI/helpers/database"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
)

var (
	partNumberMap = `select partID, oldPartNumber from Part where oldPartNumber is not null`
	baseMapStmt   = `select bv.ID, bv.AAIABaseVehicleID from BaseVehicle as bv`
	subMapStmt    = `select sm.ID, sm.AAIASubmodelID from Submodel as sm`
	superConfigs  = `select ca.ConfigAttributeTypeID, cat.AcesTypeID, ca.vcdbID, ca.ID 
			from CurtDev.ConfigAttribute as ca 
			join CurtDev.ConfigAttributeType as cat on cat.ID = ca.ConfigAttributeTypeID`
)

//maps old part number to current part number
func GetPartNumberMap() (map[string]int, error) {
	partMap := make(map[string]int)
	var err error
	var tempID *int
	var tempOldID *string
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return partMap, err
	}
	defer db.Close()

	stmt, err := db.Prepare(partNumberMap)
	if err != nil {
		return partMap, err
	}
	defer stmt.Close()
	res, err := stmt.Query()
	for res.Next() {
		err = res.Scan(&tempID, &tempOldID)
		if err != nil {
			return partMap, err
		}
		if tempID != nil && tempOldID != nil {
			partMap[*tempOldID] = *tempID
		}
	}
	return partMap, err
}

//maps aaia baseID to curt baseID
func GetBaseVehicleMap() (map[int]int, error) {
	var err error
	baseMap := make(map[int]int)
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return baseMap, err
	}
	defer db.Close()

	stmt, err := db.Prepare(baseMapStmt)
	if err != nil {
		return baseMap, err
	}
	defer stmt.Close()
	var id, aaia int
	res, err := stmt.Query()
	for res.Next() {
		err = res.Scan(&id, &aaia)
		if err != nil {
			return baseMap, err
		}
		baseMap[aaia] = id
	}
	return baseMap, err
}

//maps aaia subID to curt subID
func GetSubmodelMap() (map[int]int, error) {
	var err error
	subMap := make(map[int]int)
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return subMap, err
	}
	defer db.Close()

	stmt, err := db.Prepare(subMapStmt)
	if err != nil {
		return subMap, err
	}
	defer stmt.Close()

	var id, aaia int
	res, err := stmt.Query()
	for res.Next() {
		err = res.Scan(&id, &aaia)
		if err != nil {
			return subMap, err
		}
		subMap[aaia] = id
	}
	return subMap, err
}

func GetConfigMap() (map[string]string, error) {
	var err error
	configMap := make(map[string]string)
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return configMap, err
	}
	defer db.Close()

	stmt, err := db.Prepare(superConfigs)
	if err != nil {
		return configMap, err
	}
	defer stmt.Close()
	var typeID, acesTypeID, acesValID, valID *int
	var k, v string
	res, err := stmt.Query()
	for res.Next() {
		err = res.Scan(&typeID, &acesTypeID, &acesValID, &valID)
		if err != nil {
			return configMap, err
		}
		if *acesTypeID > 0 && *acesValID > 0 {
			k = strconv.Itoa(*acesTypeID) + "," + strconv.Itoa(*acesValID)
			v = strconv.Itoa(*typeID) + "," + strconv.Itoa(*valID)
			configMap[k] = v
		}

	}
	return configMap, err
}

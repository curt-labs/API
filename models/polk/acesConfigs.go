package polk

import (
	"database/sql"
	"github.com/curt-labs/GoAPI/helpers/database"
	_ "github.com/go-sql-driver/mysql"
)

func GetConfigMaps() map[string](map[int]string) {
	output := make(map[string](map[int]string))
	configTables := map[string]string{
		"Aspiration":              `select AspirationID as ID, AspirationName as value from Aspiration order by AspirationName`,
		"BedLength":               `select BedLengthID as ID, BedLength as value from BedLength order by BedLength`,
		"BedType":                 `select BedTypeID as ID, BedTypeName as value from BedType order by BedTypeName`,
		"BodyType":                `select BodyTypeID as ID, BodyTypeName as value from BodyType order by BodyTypeName`,
		"BodyNumDoors":            `select BodyNumDoorsID as ID, BodyNumDoors as value from BodyNumDoors order by BodyNumDoors`,
		"BrakeABS":                `select BrakeABSID as ID, BrakeABSName as value from BrakeABS order by BrakeABSName`,
		"BrakeSystem":             `select BrakeSystemID as ID, BrakeSystemName as value from BrakeSystem order by BrakeSystemName`,
		"BrakeType":               `select BrakeTypeID as ID, BrakeTypeName as value from BrakeType order by BrakeTypeName`,
		"CylinderHeadType":        `select CylinderHeadTypeID as ID, CylinderHeadTypeName as value from CylinderHeadType order by CylinderHeadTypeName`,
		"DriveType":               `select DriveTypeID as ID, DriveTypeName as value from DriveType order by DriveTypeName`,
		"ElecControlled":          `select ElecControlledID as ID, ElecControlled as value from ElecControlled order by ElecControlled`,
		"Engine":                  `select eb.EngineBaseID as ID, CONCAT(eb.Liter, ' Liter ', eb.BlockType,'-',eb.Cylinders) as value from EngineBase as eb`,
		"EngineDesignation":       `select EngineDesignationID as ID, EngineDesignationName as value from EngineDesignation order by EngineDesignationName`,
		"EngineVersion":           `select EngineVersionID as ID, EngineVersion as value from EngineVersion order by EngineVersion`,
		"EngineVIN":               `select EngineVINID as ID, EngineVINName as value from EngineVIN order by EngineVINName`,
		"EnglishPhrase":           `select EnglishPhraseID as ID, EnglishPhrase as value from EnglishPhrase order by EnglishPhrase`,
		"FuelDeliverySubType":     `select FuelDeliverySubTypeID as ID, FuelDeliverySubTypeName as value from FuelDeliverySubType order by FuelDeliverySubTypeName`,
		"FuelDeliveryType":        `select FuelDeliveryTypeID as ID, FuelDeliveryTypeName as value from FuelDeliveryType order by FuelDeliveryTypeName`,
		"FuelSystemControlType":   `select FuelSystemControlTypeID as ID, FuelSystemControlTypeName as value from FuelSystemControlType order by FuelSystemControlTypeName`,
		"FuelSystemDesign":        `select FuelSystemDesignID as ID, FuelSystemDesignName as value from FuelSystemDesign order by FuelSystemDesignName`,
		"FuelType":                `select FuelTypeID as ID, FuelTypeName as value from FuelType order by FuelTypeName`,
		"IgnitionSystem":          `select IgnitionSystemTypeID as ID, IgnitionSystemTypeName as value from IgnitionSystemType order by IgnitionSystemTypeName`,
		"Mfr":                     `select MfrID as ID, MfrName as value from Mfr order by MfrName`,
		"MfrBodyCode":             `select MfrBodyCodeID as ID, MfrBodyCodeName as value from MfrBodyCode order by MfrBodyCodeName`,
		"SpringType":              `select SpringTypeID as ID, SpringTypeName as value from SpringType order by SpringTypeName`,
		"SteeringSystem":          `select SteeringSystemID as ID, SteeringSystemName as value from SteeringSystem order by SteeringSystemName`,
		"SteeringType":            `select SteeringTypeID as ID, SteeringTypeName as value from SteeringType order by SteeringTypeName`,
		"TransmissionControlType": `select TransmissionControlTypeID as ID, TransmissionControlTypeName as value from TransmissionControlType order by TransmissionControlTypeName`,
		"TransmissionMfrCode":     `select TransmissionMfrCodeID as ID, TransmissionMfrCode as value from TransmissionMfrCode order by TransmissionMfrCode`,
		"TransmissionNumSpeeds":   `select TransmissionNumSpeedsID as ID, TransmissionNumSpeeds as value from TransmissionNumSpeeds order by TransmissionNumSpeeds`,
		"TransmissionType":        `select TransmissionTypeID as ID, TransmissionTypeName as value from TransmissionType order by TransmissionTypeName`,
		"Valves":                  `select ValvesID as ID, ValvesPerEngine as value from Valves order by ValvesPerEngine`,
		"WheelBase":               `select WheelBaseID as ID, WheelBase as value from WheelBase order by WheelBase`,
	}
	for name, query := range configTables {
		output[name] = ConfigMap(query)
	}

	return output
}

//Make Config Map
func ConfigMap(query string) map[int]string {
	configMap := make(map[int]string)
	db, err := sql.Open("mysql", database.VcdbConnectionString())
	if err != nil {
		return configMap
	}
	defer db.Close()

	stmt, err := db.Prepare(query)
	if err != nil {
		return configMap
	}
	defer stmt.Close()
	res, err := stmt.Query()
	var id int
	var val string
	for res.Next() {
		res.Scan(
			&id,
			&val,
		)
		configMap[id] = val
	}
	return configMap
}

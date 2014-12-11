package polk

import (
	"database/sql"
	"encoding/csv"
	"github.com/curt-labs/GoAPI/helpers/database"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
	"strconv"
)

//Comments Contain CurtDev.AcesType.ID and vcdb.table
type CsvDatum struct {
	Make                       string `json:"make omitempty" xml:"make omitempty"`
	Model                      string `json:"model omitempty" xml:"model omitempty"`
	SubModel                   string `json:"submodel omitempty" xml:"submodel omitempty"`
	Year                       string `json:"year omitempty" xml:"year omitempty"`
	GVW                        int
	VehicleID                  int
	BaseVehicleID              int
	YearID                     int
	MakeID                     int
	ModelID                    int
	SubmodelID                 int
	VehicleTypeID              int
	FuelTypeID                 int     // 6 FuelType
	FuelDeliveryID             int     //20 FuelDeliveryType
	AcesLiter                  float64 //EngineBase.Liter
	AcesCC                     float64 //EngineBase.CC
	AcesCID                    int     //EngineBase.CID
	AcesCyl                    int     //EngineBase.Cylinders
	AcesBlockType              string  //EngineBase.BlockType
	AspirationID               int     // 8 Aspiration
	DriveTypeID                int     // 3 DriveType
	BodyTypeID                 int     // 2 BodyType
	BodyNumDoorsID             int     // 4 BodyNumDoors
	EngineVinID                int     // 16 EngineVIN
	RegionID                   int     //Region
	PowerOutputID              int     // 25 PowerOutput
	FuelDelConfigID            int     //FuelDeliveryConfig
	BodyStyleConfigID          int     //BodyStyleConfig
	ValvesID                   int     // 40 Valves
	CylHeadTypeID              int     // 12 CylinderHeadType
	BlockType                  string  //EngineBase.BlockType
	EngineBaseID               int     // 7 EngineBase
	EngineConfigID             int     //EngineConfig
	PCDBPartTerminologyName    string
	Position                   []byte
	PartNumber                 string
	PartDesc                   string
	VehicleCount               int
	DistributedPartOpportunity int
	MaximumPartOpportunity     int
	PartID                     int //actual CURT Id
}
type CsvData []CsvDatum

type CurtVehicleConfig struct {
	ID                int
	AAIAVehicleID     int
	AAIABaseVehicleID int
	AAIASubModelID    int
	CurtBaseID        int
	CurtSubmodelID    int
	AcesConfigTypeID  int
	AcesConfigValueID int
	PartID            int
	ConfigType        int
	ConfigTypeName    string
	ConfigValue       int
	ConfigValueName   string
}

var (
	getAcesBaseVehicleByPartNumber = `select bv.AAIABaseVehicleID, sm.AAIASubmodelID, at.ID, ca.vcdbID, ca.value, at.name from vcdb_VehiclePart as vvp 
join vcdb_Vehicle as v on v.ID = vvp.VehicleID
join BaseVehicle as bv on bv.ID = v.BaseVehicleID
left join Submodel as sm on sm.ID = v.SubmodelID
left join VehicleConfigAttribute as vca on vca.VehicleConfigID = v.ConfigID
left join ConfigAttribute as ca on ca.ID = vca.AttributeID
left join ConfigAttributeType as cat on cat.ID = ca.ConfigAttributeTypeID
left join AcesType as at on at.ID = cat.AcesTypeID
where vvp.PartNumber = ?`

	// 	getAcesConfigsByAcesVehicleID = `select v.VehicleID, bed.BedConfigID, body.BodyStyleConfigID, brake.BrakeConfigID, drive.DriveTypeID, eng.EngineConfigID, mfr.MfrBodyCodeID, spring.SpringTypeConfigID, steer.SteeringConfigID, trans.TransmissionID, wheel.WheelBaseID
	// from Vehicle as v

	// join BaseVehicle as bv on bv.BaseVehicleID =v.BaseVehicleID
	// join Submodel as s on s.SubmodelID = v.SubmodelID

	// join VehicleToBedConfig as bed on bed.VehicleID = v.VehicleID
	// join VehicleToBodyStyleConfig as body on body.VehicleID = v.VehicleID
	// join VehicleToBrakeConfig as brake on brake.VehicleID = v.VehicleID
	// join VehicleToDriveType as drive on drive.VehicleID = v.VehicleID
	// join VehicleToEngineConfig as eng on eng.VehicleID = v.VehicleID
	// join VehicleToMfrBodyCode as mfr on mfr.VehicleID = v.VehicleID
	// join VehicleToSpringTypeConfig as spring on spring.VehicleID = v.VehicleID
	// join VehicleToSteeringConfig as steer on steer.VehicleID = v.VehicleID
	// join VehicleToTransmission as trans on trans.VehicleID = v.VehicleID
	// join VehicleToWheelbase as wheel on wheel.VehicleID = v.VehicleID

	// where v.VehicleID = ? `

	getCurtVehicleAndPartFromAcesVehicleID = `select CDvv.ID, CDvv.ConfigID, CDvv.AppID, CDat.ID as "acesType", CDat.name, CDca.vcdbID as "acesValue", CDca.value, CDvp.PartNumber, Vv.VehicleID

from vcdb.Vehicle as Vv

join CurtDev.BaseVehicle as CDbv on CDbv.AAIABaseVehicleID = Vv.BaseVehicleID
join CurtDev.Submodel as CDs on CDs.AAIASubmodelID = Vv.SubmodelID
join CurtDev.vcdb_Vehicle as CDvv on CDvv.BaseVehicleID = CDbv.ID and CDvv.SubmodelID = Cds.ID
left join CurtDev.VehicleConfigAttribute as CDvca on CDvca.VehicleConfigID = CDvv.ConfigID
left join CurtDev.ConfigAttribute as CDca on CDca.ID = CDvca.AttributeID
left join CurtDev.ConfigAttributeType as CDcat on CDcat.ID = CDca.ConfigAttributeTypeID
left join CurtDev.AcesType as CDat on CDat.ID = CDcat.AcesTypeID
left join CurtDev.vcdb_VehiclePart as CDvp on CDvp.VehicleID = CDvv.ID

where Vv.VehicleID = ? `

	getBaseVehicleSubmodelFromAcesID = `select CDvv.ID 
		from vcdb.Vehicle as Vv
		join CurtDev.BaseVehicle as CDbv on CDbv.AAIABaseVehicleID = Vv.BaseVehicleID
		left join CurtDev.Submodel as CDs on CDs.AAIASubmodelID = Vv.SubmodelID
		join CurtDev.vcdb_Vehicle as CDvv on CDvv.BaseVehicleID = CDbv.ID and CDvv.SubmodelID = Cds.ID
		where Vv.VehicleID = ? `

	getCurtConfigValueFromAcesConfig = `select ca.ID, ca.value
	from CurtDev.ConfigAttributeType as cat
	join CurtDev.ConfigAttribute as ca on cat.ID = ca.ConfigAttributeTypeID
	where cat.AcesTypeID = ?
	and ca.vcdbID = ?`

	getCurtConfigTypeFromAcesConfig = `select cat.ID, cat.Name
		from CurtDev.ConfigAttributeType as cat
		where cat.AcesTypeID = ? `

	getCurtBaseFromAcesBase = `select bv.ID from BaseVehicle as bv
		where bv.AAIABaseVehicleID = ?`

	getCurtSubFromAcesSub = `select sm.ID from Submodel as sm
		where sm.AAIASubmodelID = ?`

	getVehicleWithConfigs = `select CDvv.ID from CurtDev.vcdb_Vehicle as CDvv
		left join CurtDev.VehicleConfigAttribute as CDvca on CDvca.VehicleConfigID = CDvv.ConfigID
		left join CurtDev.ConfigAttribute as CDca on CDca.ID = CDvca.AttributeID
		left join CurtDev.ConfigAttributeType as CDcat on CDcat.ID = CDca.ConfigAttributeTypeID
		where CDvv.BaseVehicleID = ?
		and CDvv.SubModelID = ? 
		and CDcat.ID = ?
		and CDcat.ID = ?`

	getCurtVehicleIdFromAAIABaseVehicle = `select v.ID from vcdb_Vehicle as v
		left join BaseVehicle as bv on bv.ID = v.BaseVehicleID
		where v.SubmodelID = 0
		and bv.AAIABaseVehicleID = ?`
	getCurtVehicleIdFromAAIASubVehicle = `select v.ID from vcdb_Vehicle as v
		left join Submodel as s on s.ID = v.SubModelID
		left join BaseVehicle as bv on bv.ID = v.BaseVehicleID 
		where v.ConfigID = 0
		and s.AAIASubModelID = ?
		and bv.AAIABaseVehicleID = ?`
	addPartToVehicle           = `insert into vcdb_VehiclePart (VehicleID, PartNumber) values (?,?)`
	createCurtVehicle          = `insert into vcdb_Vehicle (BaseVehicleID, SubmodelID, ConfigID, AppID, RegionID) values (?,?,?,0,0)`
	partNumberMap              = `select partID, oldPartNumber from Part where oldPartNumber is not null`
	insertConfigAttributeValue = `insert into ConfigAttribute (ConfigAttributeTypeID, vcdbID, value) values (?,?,?)`
	checkVehiclePart           = `select vv.ID from vcdb_VehiclePart as vv where vv.VehicleID = ? and PartNUmber = ?`
)

func RunDiff(filename string, headerLines int, useOldPartNumbers bool) error {
	var err error
	var cs CsvData

	outputFile, err := os.Create("PartNumbersNeeded")
	defer outputFile.Close()

	//csv into memory
	cs, partsNeeded, err := CaptureCsv(filename, headerLines, useOldPartNumbers)
	if err != nil {
		return err
	}

	//write missing parts to PartNumbersNeeded file
	if len(partsNeeded) > 0 {
		for i, str := range partsNeeded {
			for _, vehicle := range str {
				outputFile.WriteString("part: " + i + "  AAIAvehicleID: " + strconv.Itoa(vehicle.VehicleID) + ", AAIABaseID: " + strconv.Itoa(vehicle.BaseVehicleID) + ", AAIASubmodel: " + strconv.Itoa(vehicle.SubmodelID) + "\n")
			}
		}
	}
	outputFile.Sync()

	//create basevehicle  map
	baseMap := make(map[int][]CsvDatum)
	for _, c := range cs {
		baseMap[c.BaseVehicleID] = append(baseMap[c.BaseVehicleID], c)
	}

	err = AuditBaseVehicle(baseMap)
	if err == nil {
		return nil
	}

	return err
}

//Csv to array of structs
func CaptureCsv(filename string, headerLines int, useOldPartNumbers bool) ([]CsvDatum, map[string][]CsvDatum, error) {
	var err error
	var cs []CsvDatum
	partsNeeded := make(map[string][]CsvDatum)

	csvfile, err := os.Open(filename)
	if err != nil {
		return cs, partsNeeded, err
	}

	defer csvfile.Close()
	reader := csv.NewReader(csvfile)
	reader.FieldsPerRecord = -1 //flexible number of fields

	lines, err := reader.ReadAll()
	if err != nil {
		return cs, partsNeeded, err
	}

	lines = lines[headerLines:] //axe header

	for _, line := range lines {
		//get values
		Make := line[0]
		Model := line[1]
		SubModel := line[2]
		Year := line[3]
		GVW, err := strconv.Atoi(line[4])
		VehicleID, err := strconv.Atoi(line[5])
		BaseVehicleID, err := strconv.Atoi(line[6])
		YearID, err := strconv.Atoi(line[7])
		MakeID, err := strconv.Atoi(line[8])
		ModelID, err := strconv.Atoi(line[9])
		SubmodelID, err := strconv.Atoi(line[10])
		VehicleTypeID, err := strconv.Atoi(line[11])
		FuelTypeID, err := strconv.Atoi(line[12])
		FuelDeliveryID, err := strconv.Atoi(line[13])
		AcesLiter, err := strconv.ParseFloat(line[14], 64)
		AcesCC, err := strconv.ParseFloat(line[15], 64)
		AcesCID, err := strconv.Atoi(line[16])
		AcesCyl, err := strconv.Atoi(line[17])
		AcesBlockType := line[18]
		AspirationID, err := strconv.Atoi(line[19])
		DriveTypeID, err := strconv.Atoi(line[20])
		BodyTypeID, err := strconv.Atoi(line[21])
		BodyNumDoorsID, err := strconv.Atoi(line[22])
		EngineVinID, err := strconv.Atoi(line[23])
		RegionID, err := strconv.Atoi(line[24])
		PowerOutputID, err := strconv.Atoi(line[25])
		FuelDelConfigID, err := strconv.Atoi(line[26])
		BodyStyleConfigID, err := strconv.Atoi(line[27])
		ValvesID, err := strconv.Atoi(line[28])
		CylHeadTypeID, err := strconv.Atoi(line[29])
		BlockType := line[30]
		EngineBaseID, err := strconv.Atoi(line[31])
		EngineConfigID, err := strconv.Atoi(line[32])
		PCDBPartTerminologyName := line[33]
		Position := []byte(line[34])
		PartNumber := line[35]
		PartDesc := line[36]
		VehicleCount, err := strconv.Atoi(line[37])
		DistributedPartOpportunity, err := strconv.Atoi(line[38])
		MaximumPartOpportunity, err := strconv.Atoi(line[39])
		if err != nil {
			return cs, partsNeeded, err
		}
		//assign to struct
		c := CsvDatum{
			Make:                       Make,
			Model:                      Model,
			SubModel:                   SubModel,
			Year:                       Year,
			GVW:                        GVW,
			VehicleID:                  VehicleID,
			BaseVehicleID:              BaseVehicleID,
			YearID:                     YearID,
			MakeID:                     MakeID,
			ModelID:                    ModelID,
			SubmodelID:                 SubmodelID,
			VehicleTypeID:              VehicleTypeID,
			FuelTypeID:                 FuelTypeID,
			FuelDeliveryID:             FuelDeliveryID,
			AcesLiter:                  AcesLiter,
			AcesCC:                     AcesCC,
			AcesCID:                    AcesCID,
			AcesCyl:                    AcesCyl,
			AcesBlockType:              AcesBlockType,
			AspirationID:               AspirationID,
			DriveTypeID:                DriveTypeID,
			BodyTypeID:                 BodyTypeID,
			BodyNumDoorsID:             BodyNumDoorsID,
			EngineVinID:                EngineVinID,
			RegionID:                   RegionID,
			PowerOutputID:              PowerOutputID,
			FuelDelConfigID:            FuelDelConfigID,
			BodyStyleConfigID:          BodyStyleConfigID,
			ValvesID:                   ValvesID,
			CylHeadTypeID:              CylHeadTypeID,
			BlockType:                  BlockType,
			EngineBaseID:               EngineBaseID,
			EngineConfigID:             EngineConfigID,
			PCDBPartTerminologyName:    PCDBPartTerminologyName,
			Position:                   Position,
			PartNumber:                 PartNumber,
			PartDesc:                   PartDesc,
			VehicleCount:               VehicleCount,
			DistributedPartOpportunity: DistributedPartOpportunity,
			MaximumPartOpportunity:     MaximumPartOpportunity,
		}
		if useOldPartNumbers == true {
			partMap, err := GetPartNumberMap()
			if err != nil {
				return cs, partsNeeded, err
			}
			//get new part id, if there is one
			if newPartNum, ok := partMap[c.PartNumber]; ok {
				c.PartID = newPartNum
			} else {
				//no new part number -> append to partsNeeded for output file write
				partsNeeded[c.PartNumber] = append(partsNeeded[c.PartNumber], c)
			}
		} else {
			c.PartID, err = strconv.Atoi(c.PartNumber)
		}

		cs = append(cs, c)
	}
	return cs, partsNeeded, err
}

func AuditBaseVehicle(baseMap map[int][]CsvDatum) error {
	var err error
	submodelMap := make(map[int][]CsvDatum)

	for _, baseVehicle := range baseMap {
		baseFlag := false
		for i, base := range baseVehicle {
			if i > 0 {
				if base.PartNumber != baseVehicle[i-1].PartNumber {
					baseFlag = true
					break
				}
			}
		}

		if baseFlag == false {
			log.Print("All the same part ", baseVehicle[0].PartNumber, " for this CsvBasevehicle: ", baseVehicle[0].BaseVehicleID)
			//check for vehiclePartExistence
			vehiclePartID, err := baseVehicle[0].CheckVehiclePartTableByAAIABaseVehicle()
			if vehiclePartID == 0 {
				log.Print("need to add vehicle part")
				//add part to base vehicle
				err = AddPartToBaseVehicle(baseVehicle[0])
				if err != nil {
					//TODO need to add base v  & try again
					log.Print("Error adding to baseVehicle (no baseVehicle) ", err)
				}
			} else {
				log.Print("Vehicle part EXISTS FOR BASE VEHICLE")
			}
		} else {
			log.Print("Diff parts for base vehicle ", baseVehicle[0].BaseVehicleID, ", try submodel")
			//There are different parts for this base vehicle, try submodel
			//build map of AAIAsubmodelID's to CsvData
			for _, base := range baseVehicle {
				submodelMap[base.SubmodelID] = append(submodelMap[base.SubmodelID], base)
			}
		}
	}
	//audit submodels
	if len(submodelMap) > 0 {
		// log.Print("SUBMODEL MAP ", submodelMap)
		err = AuditSubmodel(submodelMap)
	}
	return err
}

func (c *CsvDatum) CheckVehiclePartTableByAAIABaseVehicle() (int, error) {
	var curtVehicleID int
	var vehiclePartID int
	var err error

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return vehiclePartID, err
	}
	defer db.Close()
	stmt, err := db.Prepare(getCurtVehicleIdFromAAIABaseVehicle)
	if err != nil {
		return vehiclePartID, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(c.BaseVehicleID).Scan(&curtVehicleID)
	if err != nil {
		return vehiclePartID, err
	}

	stmt, err = db.Prepare(checkVehiclePart)
	if err != nil {
		return vehiclePartID, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(curtVehicleID, c.PartID).Scan(&vehiclePartID)
	if err != nil {
		return vehiclePartID, err
	}
	return vehiclePartID, err
}

func AuditSubmodel(subMap map[int][]CsvDatum) error {
	var err error
	vIDmap := make(map[int][]CsvDatum)

	for _, subVehicle := range subMap {
		subFlag := false
		for i, sub := range subVehicle {
			if i > 0 {
				if sub.PartNumber != subVehicle[i-1].PartNumber {
					subFlag = true
					break
				}
			}
		}
		if subFlag == false {
			//add part to sub vehicle
			log.Print("All the same part ", subVehicle[0].PartNumber, " for this CsvSubmodel: ", subVehicle[0].SubmodelID, ". CsvBaseID: ", subVehicle[0].BaseVehicleID)
			//check table for vehiclePart Existence
			vehiclePartID, err := subVehicle[0].CheckVehiclePartTableByAAIASubmodel()
			if vehiclePartID == 0 {
				err = AddPartToSubVehicle(subVehicle[0])
				if err != nil {
					//TODO need to add sub vehicle & try again
					log.Print("Error adding to submodel (no submodel) ", err)
				}
			} else {
				log.Print("VEHICLE PART EXISTS FOR SUBMODEL")
			}
		} else {
			log.Print("Diff parts for submodel ", subVehicle[0].SubmodelID, ", try configs")
			//config breakdown
			//make map of un-added AAIAVehicleID to CsvData
			for _, sub := range subVehicle {
				vIDmap[sub.VehicleID] = append(vIDmap[sub.VehicleID], sub)
			}
		}
	}
	//audit configs by AAIAVehicleID
	if len(vIDmap) > 0 {
		//TODO audit configs
		// log.Print("VEHICLE ID MAP ", vIDmap)
		err = AuditConfigs(vIDmap)
	}
	return err
}

func (c *CsvDatum) CheckVehiclePartTableByAAIASubmodel() (int, error) {
	var curtVehicleID int
	var vehiclePartID int
	var err error

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return vehiclePartID, err
	}
	defer db.Close()
	stmt, err := db.Prepare(getCurtVehicleIdFromAAIASubVehicle)
	if err != nil {
		return vehiclePartID, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(c.BaseVehicleID, c.SubmodelID).Scan(&curtVehicleID)
	if err != nil {
		return vehiclePartID, err
	}

	stmt, err = db.Prepare(checkVehiclePart)
	if err != nil {
		return vehiclePartID, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(curtVehicleID, c.PartID).Scan(&vehiclePartID)
	if err != nil {
		return vehiclePartID, err
	}
	return vehiclePartID, err
}

func AuditConfigs(vIDmap map[int][]CsvDatum) error {
	var err error
	//get Maps
	ConfigMaps := GetConfigMaps()
	for _, vehicle := range vIDmap {
		log.Print("--NEW VEHICLE --")
		for _, v := range vehicle {
			// log.Print("V - ", v)
			err = v.AddPartToVehicle(ConfigMaps)
			if err != nil {
				return err
			}
		}
	}
	return err
}

func AddPartToBaseVehicle(c CsvDatum) error {
	var err error
	var curt CurtVehicleConfig
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(getCurtVehicleIdFromAAIABaseVehicle)
	if err != nil {
		return err
	}
	defer stmt.Close()
	errBaseVehicle := stmt.QueryRow(c.BaseVehicleID).Scan(&curt.ID)
	if errBaseVehicle != nil {
		//create vcdb_Vehicle that is 0 submodel, 0 config
		curt.AAIABaseVehicleID = c.BaseVehicleID
		err = curt.GetCurtBaseVehicleFromAcesBaseVehicle()
		if err != nil {
			return err
		}
		curt.CurtSubmodelID = 0 //no submodel differentiation
		log.Print("We'd create vcdb here: ", curt.CurtBaseID, "  ", curt.CurtSubmodelID)
		//TODO - uncomment
		// err = curt.CreateVcdb_Vehicle()
		// if err != nil {
		// 	return err
		// }
	}
	log.Print("Adding part ", c.PartID, " to Curt base vehicle ", curt.ID)
	curt.PartID = c.PartID
	//TODO - uncomment
	// err = curt.AddVehicleToVehiclePart()
	// if err != nil {
	// 	return err
	// }
	return err
}

func AddPartToSubVehicle(c CsvDatum) error {
	var err error
	var curt CurtVehicleConfig
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(getCurtVehicleIdFromAAIASubVehicle)
	if err != nil {
		return err
	}
	defer stmt.Close()
	errSubmodel := stmt.QueryRow(c.SubmodelID, c.BaseVehicleID).Scan(&curt.ID)
	if errSubmodel != nil {
		//create vcdb_Vehicle that is 0 config
		curt.AAIABaseVehicleID = c.BaseVehicleID
		err = curt.GetCurtBaseVehicleFromAcesBaseVehicle()
		if err != nil {
			return err
		}
		curt.AAIASubModelID = c.SubmodelID
		err = curt.GetCurtSubmodelFromAcesSubmodel()
		if err != nil {
			return err
		}
		log.Print("We'd create vcdb here: ", curt.CurtBaseID, "  ", curt.CurtSubmodelID)
		//TODO - uncomment
		// err = curt.CreateVcdb_Vehicle()
		// if err != nil {
		// 	return err
		// }
	}
	log.Print("Adding part ", c.PartID, " to Curt sub vehicle ", curt.ID)
	curt.PartID = c.PartID
	//TODO - uncomment
	// err = curt.AddVehicleToVehiclePart()
	// if err != nil {
	// 	return err
	// }
	return err
}

func (curt *CurtVehicleConfig) CreateVcdb_Vehicle() error {
	var err error
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(createCurtVehicle)
	if err != nil {
		return err
	}
	defer stmt.Close()
	res, err := stmt.Exec(curt.CurtBaseID, curt.CurtSubmodelID, 0) //no config
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	curt.ID = int(id)
	return err
}

func (c *CsvDatum) AddPartToVehicle(configMaps map[string](map[int]string)) error {
	var err error
	var curt CurtVehicleConfig

	//assign AAIA base/submodel aces IDs to new CURT vehicle
	curt.AAIAVehicleID = c.VehicleID
	curt.AAIABaseVehicleID = c.BaseVehicleID
	curt.AAIASubModelID = c.SubmodelID
	curt.PartID = c.PartID

	//get Curt Base Vehicle and Submodel IDs
	err = curt.GetCurtBaseVehicleFromAcesBaseVehicle()
	if err != nil {
		return err
	}
	err = curt.GetCurtSubmodelFromAcesSubmodel()
	if err != nil {
		return err
	}

	//ugly way to range over configs supplied by Polk
	acesConfigTypeArray := [...]int{6, 20, 8, 3, 2, 4, 16, 25, 40, 12, 7} //There's got to be a better way

	for _, acesConfigType := range acesConfigTypeArray {
		//assign aces type from loop of aces types
		curt.AcesConfigTypeID = acesConfigType
		//get Aces value from CsvDatum struct
		curt.AcesConfigValueID = c.getConfigType(curt.AcesConfigTypeID)
		//get Curt type from aces type
		curt.ConfigType, curt.ConfigTypeName, _ = GetCurtConfigTypeAcesConfig(curt.AcesConfigTypeID)
		//get Curt value from aces value and type, if there is one
		curt.ConfigValue, curt.ConfigValueName, _ = GetCurtConfigValueAcesConfig(curt.AcesConfigTypeID, curt.AcesConfigValueID)

		//see if there is a CURT config value for this config type
		if curt.ConfigValue > 0 {
			err = curt.GetCurtVechicleWithConfig()
			if err == nil {
				//add part to this curt vehicle
				log.Print("We'd add a part to vehicle config ", curt.ConfigTypeName, "  -  ", curt.ConfigValueName)
				// //TODO UNCOMMENT
				// err = curt.AddVehicleToVehiclePart()
				// if err != nil {
				// 	return err
				// }
			}
			if err == sql.ErrNoRows {
				//add vehicle with this config before adding part
				log.Print("We'd add a vehicle before part. Aces Base/Sub: ", c.BaseVehicleID, " ", c.SubmodelID, "| Curt Base/Sub: ", curt.CurtBaseID, " ", curt.CurtSubmodelID, " ", curt.ConfigTypeName, "  -  ", curt.ConfigValueName)
				err = nil
				// //TODO UNCOMMENT
				err = curt.InsertConfigurationValue()
				if err != nil {
					return err
				}
				// err = curt.AddVehicleToVehiclePart()
				// if err != nil {
				// 	return err
				// }
			}
			if err != nil {
				log.Print("Another err ", err)
				return err
			}
			// err = curt.Print()
		} else {
			//exclude configs without aces data
			if curt.AcesConfigValueID > 0 {
				//what am I missing - TODO - should we add all these config values that we don't have values for
				//need to add config value, then vehicle with that config, then vehiclePart
				log.Print("ACES ", curt.AcesConfigTypeID, curt.AcesConfigValueID, "   CURT ", curt.ConfigType, curt.ConfigValue)
			}
		}
	}
	return err
}

func (c *CsvDatum) getConfigType(n int) int {
	staticMap := make(map[int]int)
	staticMap[6] = c.FuelTypeID
	staticMap[20] = c.FuelDeliveryID
	staticMap[8] = c.AspirationID
	staticMap[3] = c.DriveTypeID
	staticMap[2] = c.BodyTypeID
	staticMap[4] = c.BodyNumDoorsID
	staticMap[16] = c.EngineVinID
	staticMap[25] = c.PowerOutputID
	staticMap[40] = c.ValvesID
	staticMap[12] = c.CylHeadTypeID
	staticMap[7] = c.EngineBaseID
	return staticMap[n]
}

func (curt *CurtVehicleConfig) AddVehicleToVehiclePart() error {
	var err error
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(addPartToVehicle)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(curt.ID, curt.PartID)
	if err != nil {
		return err
	}
	return err
}

func (curt *CurtVehicleConfig) InsertConfigurationValue() error {
	var err error
	log.Print("INSERT CONFIG VAL")

	// db, err := sql.Open("mysql", database.ConnectionString())
	// if err != nil {
	// 	return
	// }
	// defer db.Close()

	// stmt, err := db.Prepare(insertConfigAttributeValue)
	// if err != nil {
	// 	return
	// }
	// defer stmt.Close()
	// res, err := stmt.Exec(curt.ConfigType, curt.AcesConfigValueID, value)
	// if err != nil {
	// 	return err
	// }
	// id, err := res.LastInsertId()
	// curt.ConfigValue = int(id)
	return err
}

func GetCurtConfigTypeAcesConfig(acesType int) (int, string, error) {
	var err error
	var configType int
	var configTypeName string
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return 0, "", err
	}
	defer db.Close()

	stmt, err := db.Prepare(getCurtConfigTypeFromAcesConfig)
	if err != nil {
		return 0, "", err
	}
	defer stmt.Close()
	var configTypeByte []byte
	err = stmt.QueryRow(acesType).Scan(&configType, &configTypeByte)
	if err != nil {
		return 0, "", err
	}
	if configTypeByte != nil {
		configTypeName = string(configTypeByte[:])
	}
	return configType, configTypeName, err
}

func GetCurtConfigValueAcesConfig(acesType, acesValue int) (int, string, error) {
	var err error
	var configValue int
	var configValueName string
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return 0, "", err
	}
	defer db.Close()

	stmt, err := db.Prepare(getCurtConfigValueFromAcesConfig)
	if err != nil {
		return 0, "", err
	}
	defer stmt.Close()
	var configValByte []byte
	err = stmt.QueryRow(acesType, acesValue).Scan(&configValue, &configValByte)
	if err != nil {
		return 0, "", err
	}
	if configValByte != nil {
		configValueName = string(configValByte[:])
	}
	return configValue, configValueName, err
}

func (c *CurtVehicleConfig) GetCurtBaseVehicleFromAcesBaseVehicle() error {
	var err error
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(getCurtBaseFromAcesBase)
	if err != nil {
		return err
	}
	defer stmt.Close()
	err = stmt.QueryRow(c.AAIABaseVehicleID).Scan(&c.CurtBaseID)
	if err != nil {
		return err
	}
	return err
}

func (c *CurtVehicleConfig) GetCurtSubmodelFromAcesSubmodel() error {
	var err error
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(getCurtSubFromAcesSub)
	if err != nil {
		return err
	}
	defer stmt.Close()
	err = stmt.QueryRow(c.AAIASubModelID).Scan(&c.CurtSubmodelID)
	if err != nil {
		return err
	}
	return err
}

func (c *CurtVehicleConfig) GetCurtVechicleWithConfig() error {
	var err error
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(getVehicleWithConfigs)
	if err != nil {
		return err
	}
	defer stmt.Close()
	err = stmt.QueryRow(c.CurtBaseID, c.CurtSubmodelID, c.ConfigValue, c.ConfigType).Scan(&c.ID)
	return err
}

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

func (curt *CurtVehicleConfig) Print() error {
	log.Printf(`Curt Vehicle ID: %d
		AAIA VehicleID: %d
		AAIA Base ID: %d
		AAIA Submodel ID: %d
		Curt Base ID: %d
		Curt Submodel ID: %d
		AcesConfigTypeID: %d
		AcesConfigValueID: %d
		PartID: %d
		ConfigType: %d
		ConfigTypeName: %v
		ConfigValue: %d
		ConfigValueName: %v
		`, curt.ID, curt.AAIAVehicleID, curt.AAIABaseVehicleID, curt.AAIASubModelID, curt.CurtBaseID, curt.CurtSubmodelID, curt.AcesConfigTypeID, curt.AcesConfigValueID, curt.PartID, curt.ConfigType, curt.ConfigTypeName, curt.ConfigValue, curt.ConfigValueName)
	return nil
}

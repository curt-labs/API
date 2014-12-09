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
	PartID            string
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
	addPartToVehicle  = `insert into vcdb_VehiclePart (VehicleID, PartNumber) values (?,?)`
	createCurtVehicle = `insert into vcdb_Vehicle (BaseVehicleID, SubmodelID, ConfigID, AppID, RegionID) values (?,?,?,0,0)`
)

func RunDiff(filename string, headerLines int) error {
	var err error
	var cs CsvData

	//csv into memory
	cs, err = CaptureCsv(filename, headerLines)
	if err != nil {
		return err
	}

	baseMap := make(map[int][]CsvDatum)
	subMap := make(map[int][]CsvDatum)

	//create basevehicle and submodel maps
	for _, c := range cs {
		//create basevehicle map
		baseMap[c.BaseVehicleID] = append(baseMap[c.BaseVehicleID], c)

		//create submodel map
		subMap[c.SubmodelID] = append(subMap[c.SubmodelID], c)

		//check for curtVehicle from AcesVehicle + configs
		// err = c.InsertData()
	}

	err = AuditBaseVehicle(baseMap, subMap)
	if err == nil {
		return nil
	}

	// log.Print("NEED TO CONFIG", err)
	err = nil
	//now, go config by config
	return err
}

//Csv to array of structs
func CaptureCsv(filename string, headerLines int) ([]CsvDatum, error) {
	var err error
	var cs []CsvDatum

	csvfile, err := os.Open(filename)
	if err != nil {
		return cs, err
	}

	defer csvfile.Close()
	reader := csv.NewReader(csvfile)
	reader.FieldsPerRecord = -1 //flexible number of fields

	lines, err := reader.ReadAll()
	if err != nil {
		return cs, err
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
			return cs, err
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
		cs = append(cs, c)
	}
	return cs, err
}

func AuditBaseVehicle(baseMap map[int][]CsvDatum, subMap map[int][]CsvDatum) error {
	var err error
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
			//add part to base vehicle
			log.Print("All the same part ", baseVehicle[0].PartNumber, " for this CsvBasevehicle: ", baseVehicle[0].BaseVehicleID)
			err = AddPartToBaseVehicle(baseVehicle[0])
			if err != nil {
				//TODO need to add base v  & try again
				log.Print("Error adding to baseVehicle (no baseVehicle) ", err)
			}
		} else {
			log.Print("Diff parts for base vehicle ", baseVehicle[0].BaseVehicleID, ", try submodel")
			//There are different parts for this base vehicle, try submodel
			err = AuditSubmodel(subMap)
			if err != nil {
				return err
			}
		}
	}
	return err
}

func AuditSubmodel(subMap map[int][]CsvDatum) error {
	var err error
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
			err = AddPartToSubVehicle(subVehicle[0])
			if err != nil {
				//TODO need to add sub vehicle & try again
				log.Print("Error adding to submodel (no submodel) ", err)
			}
		} else {
			log.Print("Diff parts for submodel ", subVehicle[0].SubmodelID, ", try configs")
			//TODO - config breakdown
			err = AuditConfigs(subMap)
			if err != nil {
				return err
			}
		}
	}
	return err
}

func AuditConfigs(subMap map[int][]CsvDatum) error {
	var err error

	//TODO - finish him
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
	log.Print("Adding part ", c.PartNumber, " to Curt base vehicle ", curt.ID)
	//TODO - uncomment
	// stmt, err = db.Prepare(addPartToVehicle)
	// if err != nil {
	// 	return err
	// }
	// defer stmt.Close()
	// _, err = stmt.Exec(curt.ID, c.PartNumber)
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
	log.Print("Adding part ", c.PartNumber, " to Curt sub vehicle ", curt.ID)
	//TODO - uncomment
	// stmt, err = db.Prepare(addPartToVehicle)
	// if err != nil {
	// 	return err
	// }
	// defer stmt.Close()
	// _, err = stmt.Exec(curt.ID, c.PartNumber)
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

func (c *CsvDatum) InsertData() error {
	var err error
	var curt CurtVehicleConfig

	vehiclesToNotAdd := 0

	//assign base/submodel aces IDs
	curt.AAIAVehicleID = c.VehicleID
	curt.AAIABaseVehicleID = c.BaseVehicleID
	curt.AAIASubModelID = c.SubmodelID

	//get Curt Base Vehicle and Submodel IDs
	err = curt.GetCurtBaseVehicleFromAcesBaseVehicle()
	if err != nil {
		return err
	}
	err = curt.GetCurtSubmodelFromAcesSubmodel()
	if err != nil {
		return err
	}

	//which part are we talking about?
	curt.PartID = c.PartNumber

	//ugly way to range over configs supplied by Polk
	acesConfigTypeArray := [...]int{6, 20, 8, 3, 2, 4, 16, 25, 40, 12, 7} //There's got to be a better way

	for _, acesConfigType := range acesConfigTypeArray {

		curt.AcesConfigValueID = c.getConfigType(curt.AcesConfigTypeID)

		curt.AcesConfigTypeID = acesConfigType
		curt.ConfigType, curt.ConfigTypeName, _ = GetCurtConfigTypeAcesConfig(curt.AcesConfigTypeID)
		curt.ConfigValue, curt.ConfigValueName, _ = GetCurtConfigValueAcesConfig(curt.AcesConfigTypeID, curt.AcesConfigValueID) //TODO

		err = curt.GetCurtVechicleWithConfig()
		if err == nil {
			//add part
			vehiclesToNotAdd++
		}
		if err == sql.ErrNoRows {
			//add vehicle config
			err = nil
		}

		err = curt.Print()

	}
	log.Print(vehiclesToNotAdd)
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

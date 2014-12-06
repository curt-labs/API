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
	PartNumber                 int
	PartDesc                   string
	VehicleCount               int
	DistributedPartOpportunity int
	MaximumPartOpportunity     int
}
type CsvData []CsvDatum

type CurtVehicleConfig struct {
	ID                int
	AAIABaseVehicleID int
	AAIASubModelID    int
	AcesConfigTypeID  int
	AcesConfigValueID int
	PartID            int
	ConfigType        int
	ConfigValue       int
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

	getAcesConfigsByAcesVehicleID = `select v.VehicleID, bed.BedConfigID, body.BodyStyleConfigID, brake.BrakeConfigID, drive.DriveTypeID, eng.EngineConfigID, mfr.MfrBodyCodeID, spring.SpringTypeConfigID, steer.SteeringConfigID, trans.TransmissionID, wheel.WheelBaseID
from Vehicle as v

join BaseVehicle as bv on bv.BaseVehicleID =v.BaseVehicleID
join Submodel as s on s.SubmodelID = v.SubmodelID

join VehicleToBedConfig as bed on bed.VehicleID = v.VehicleID
join VehicleToBodyStyleConfig as body on body.VehicleID = v.VehicleID
join VehicleToBrakeConfig as brake on brake.VehicleID = v.VehicleID
join VehicleToDriveType as drive on drive.VehicleID = v.VehicleID
join VehicleToEngineConfig as eng on eng.VehicleID = v.VehicleID
join VehicleToMfrBodyCode as mfr on mfr.VehicleID = v.VehicleID
join VehicleToSpringTypeConfig as spring on spring.VehicleID = v.VehicleID
join VehicleToSteeringConfig as steer on steer.VehicleID = v.VehicleID
join VehicleToTransmission as trans on trans.VehicleID = v.VehicleID
join VehicleToWheelbase as wheel on wheel.VehicleID = v.VehicleID

where v.VehicleID = ? `

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

	getCurtConfigValueFromAcesConfig = `select  ca.ID
	from CurtDev.ConfigAttributeType as cat
	join CurtDev.ConfigAttribute as ca on cat.ID = ca.ConfigAttributeTypeID
	where cat.AcesTypeID = ?
	and ca.vcdbID = ?`

	getCurtConfigTypeFromAcesConfig = `select cat.ID
		from CurtDev.ConfigAttributeType as cat
		where cat.AcesTypeID = ? `
)

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
		PartNumber, err := strconv.Atoi(line[35])
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
	// log.Print("CS", cs)
	return cs, err
}

func (c *CsvDatum) InsertData() error {
	var err error
	var curt CurtVehicleConfig

	//assign base/submodel aces IDs
	curt.AAIABaseVehicleID = c.BaseVehicleID
	curt.AAIASubModelID = c.SubmodelID

	acesConfigTypeArray := [...]int{6, 20, 8, 3, 2, 4, 16, 25, 40, 12, 7} //There's got to be a better way

	for _, acesConfigType := range acesConfigTypeArray {
		curt.AcesConfigTypeID = acesConfigType
		curt.ConfigType, _ = GetCurtConfigTypeAcesConfig(curt.AcesConfigTypeID)
		curt.ConfigValue, _ = GetCurtConfigValueAcesConfig(curt.AcesConfigTypeID, c.FuelDeliveryID)
		log.Print(curt)
		//TODO - errors mean we need those configs

		//TODO insert vehiclePart when there is a match

	}

	return err
}
func GetCurtConfigTypeAcesConfig(acesType int) (int, error) {
	var err error
	var configType int
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return 0, err
	}
	defer db.Close()

	stmt, err := db.Prepare(getCurtConfigTypeFromAcesConfig)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	err = stmt.QueryRow(acesType).Scan(&configType)
	if err != nil {
		return 0, err
	}
	return configType, err
}

func GetCurtConfigValueAcesConfig(acesType, acesValue int) (int, error) {
	var err error
	var configValue int
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return 0, err
	}
	defer db.Close()

	stmt, err := db.Prepare(getCurtConfigValueFromAcesConfig)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	err = stmt.QueryRow(acesType, acesValue).Scan(&configValue)
	if err != nil {
		return 0, err
	}
	// log.Print(acesType, acesValue, configValue)
	return configValue, err
}

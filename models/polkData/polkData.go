package polkData

import (
	"database/sql"
	"encoding/csv"
	"github.com/curt-labs/GoAPI/helpers/database"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
	"strconv"
)

var (
	findPartID = `select vp.PartNumber from vcdb_VehiclePart as vp
		where vp.VehicleID = ?
		and vp.PartNumber = ?`
	findVehicle = `select v.ID from vcdb_Vehicle as v 
		where v.BaseVehicleID = ?
		and v.SubmodelID = ?
		and v.ConfigID = ?`
	getCurtConfigValueFromAcesConfig = `select ca.ID, ca.value
		from CurtDev.ConfigAttributeType as cat
		join CurtDev.ConfigAttribute as ca on cat.ID = ca.ConfigAttributeTypeID
		where cat.AcesTypeID = ?
		and ca.vcdbID = ?`

	getCurtConfigTypeFromAcesConfig = `select cat.ID, cat.Name
		from CurtDev.ConfigAttributeType as cat
		where cat.AcesTypeID = ? `

	getVehicleWithConfig = `select CDvv.ID from CurtDev.vcdb_Vehicle as CDvv
		left join CurtDev.VehicleConfigAttribute as CDvca on CDvca.VehicleConfigID = CDvv.ConfigID
		left join CurtDev.ConfigAttribute as CDca on CDca.ID = CDvca.AttributeID
		left join CurtDev.ConfigAttributeType as CDcat on CDcat.ID = CDca.ConfigAttributeTypeID
		where CDvv.BaseVehicleID = ?
		and CDvv.SubModelID = ? 
		and CDcat.ID = ?
		and CDcat.ID = ?`
)

//Comments Contain CurtDev.AcesType.ID and vcdb.table
type CsvDatum struct {
	CsvVehicle CsvVehicle
	Part       Part
	PartData
	CurtVehicle CurtVehicle
}
type CsvData []CsvDatum

type CsvVehicle struct {
	Make              string
	Model             string
	SubModel          string
	Year              string
	GVW               int
	VehicleID         int
	BaseVehicleID     int
	YearID            int
	MakeID            int
	ModelID           int
	SubmodelID        int
	VehicleTypeID     int
	FuelTypeID        int     // 6 FuelType
	FuelDeliveryID    int     //20 FuelDeliveryType
	AcesLiter         float64 //EngineBase.Liter
	AcesCC            float64 //EngineBase.CC
	AcesCID           int     //EngineBase.CID
	AcesCyl           int     //EngineBase.Cylinders
	AcesBlockType     string  //EngineBase.BlockType
	AspirationID      int     // 8 Aspiration
	DriveTypeID       int     // 3 DriveType
	BodyTypeID        int     // 2 BodyType
	BodyNumDoorsID    int     // 4 BodyNumDoors
	EngineVinID       int     // 16 EngineVIN
	RegionID          int     //Region
	PowerOutputID     int     // 25 PowerOutput
	FuelDelConfigID   int     //FuelDeliveryConfig
	BodyStyleConfigID int     //BodyStyleConfig
	ValvesID          int     // 40 Valves
	CylHeadTypeID     int     // 12 CylinderHeadType
	BlockType         string  //EngineBase.BlockType
	EngineBaseID      int     // 7 EngineBase
	EngineConfigID    int     //EngineConfig
}

type Part struct {
	ID                      int
	OldID                   string
	PCDBPartTerminologyName string
	Position                []byte
	PartDesc                string
}

type PartData struct {
	VehicleCount               int
	DistributedPartOpportunity int
	MaximumPartOpportunity     int
}

type CurtVehicle struct {
	ID                int
	CurtBaseID        int
	CurtSubmodelID    int
	AAIAVehicleID     int
	AAIABaseVehicleID int
	AAIASubModelID    int
	CurtConfigs       CurtConfigs
}

type CurtConfig struct {
	AcesConfigTypeID  int
	AcesConfigValueID int
	ConfigTypeID      int
	ConfigTypeName    string
	ConfigValueID     int
	ConfigValueName   string
}
type CurtConfigs []CurtConfig

func RunDiff(filename string, headerLines int, useOldPartNumbers bool, insertMissingData bool) error {
	var err error
	var cs CsvData

	//csv into memory
	cs, err = CaptureCsv(filename, headerLines, useOldPartNumbers, insertMissingData)
	if err != nil {
		return err
	}

	// //write missing parts to PartsNeeded file
	// partsNeededFile, err := os.Create("PartNumbersNeeded")
	// defer partsNeededFile.Close()
	// if len(*partsNeeded) > 0 {
	// 	for _, c := range *partsNeeded {
	// 		partsNeededFile.WriteString("Old Part ID: " + c.Part.OldID + "  AAIAvehicleID: " + strconv.Itoa(c.CsvVehicle.VehicleID) + ", AAIABaseID: " + strconv.Itoa(c.CsvVehicle.BaseVehicleID) + ", AAIASubmodel: " + strconv.Itoa(c.CsvVehicle.SubmodelID) + "\n")
	// 	}
	// }
	// partsNeededFile.Sync()

	// //write missing baseVehicles to missingBaseVehicles file
	// baseVehiclesNeededFile, err := os.Create("BaseVehiclesNeeded")
	// defer baseVehiclesNeededFile.Close()
	// if len(*missingBaseVehicles) > 0 {
	// 	for _, c := range *missingBaseVehicles {
	// 		baseVehiclesNeededFile.WriteString("Old Part ID: " + c.Part.OldID + "  AAIAvehicleID: " + strconv.Itoa(c.CsvVehicle.VehicleID) + ", AAIABaseID: " + strconv.Itoa(c.CsvVehicle.BaseVehicleID) + ", AAIASubmodel: " + strconv.Itoa(c.CsvVehicle.SubmodelID) + "\n")
	// 	}
	// }
	// baseVehiclesNeededFile.Sync()

	// //write missing submodels to missingSubmodesls file
	// submodelsNeededFile, err := os.Create("SubmodelsNeeded")
	// defer submodelsNeededFile.Close()
	// if len(*missingSubmodels) > 0 {
	// 	for _, c := range *missingSubmodels {
	// 		submodelsNeededFile.WriteString("Old Part ID: " + c.Part.OldID + "  AAIAvehicleID: " + strconv.Itoa(c.CsvVehicle.VehicleID) + ", AAIABaseID: " + strconv.Itoa(c.CsvVehicle.BaseVehicleID) + ", AAIASubmodel: " + strconv.Itoa(c.CsvVehicle.SubmodelID) + "\n")
	// 	}
	// }
	// submodelsNeededFile.Sync()

	//create basevehicle  map
	baseMap := make(map[int][]CsvDatum)
	for _, c := range cs {
		baseMap[c.CurtVehicle.CurtBaseID] = append(baseMap[c.CurtVehicle.CurtBaseID], c)
	}
	//begin audits
	err = Audits(baseMap)
	if err == nil {
		return nil
	}

	return err
}

//Csv to array of structs
func CaptureCsv(filename string, headerLines int, useOldPartNumbers bool, insertMissingData bool) ([]CsvDatum, error) {
	var err error
	var cs []CsvDatum

	//base and submodel maps
	baseMap, err := GetBaseVehicleMap()
	subMap, err := GetSubmodelMap()
	if err != nil {
		return cs, err
	}

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

	//setup files
	//write missing parts to PartsNeeded file
	partsNeededFile, err := os.Create("PartNumbersNeeded")
	defer partsNeededFile.Close()
	partOffset, err := WriteVehicleHeader(partsNeededFile)

	//write missing baseVehicles to missingBaseVehicles file
	baseVehiclesNeededFile, err := os.Create("BaseVehiclesNeeded")
	defer baseVehiclesNeededFile.Close()
	baseOffset, err := WriteVehicleHeader(baseVehiclesNeededFile)

	//write missing submodels to missingSubmodesls file
	submodelsNeededFile, err := os.Create("SubmodelsNeeded")
	defer submodelsNeededFile.Close()
	submodelOffset, err := WriteVehicleHeader(submodelsNeededFile)

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
			CsvVehicle: CsvVehicle{
				Make:              Make,
				Model:             Model,
				SubModel:          SubModel,
				Year:              Year,
				GVW:               GVW,
				VehicleID:         VehicleID,
				BaseVehicleID:     BaseVehicleID,
				YearID:            YearID,
				MakeID:            MakeID,
				ModelID:           ModelID,
				SubmodelID:        SubmodelID,
				VehicleTypeID:     VehicleTypeID,
				FuelTypeID:        FuelTypeID,
				FuelDeliveryID:    FuelDeliveryID,
				AcesLiter:         AcesLiter,
				AcesCC:            AcesCC,
				AcesCID:           AcesCID,
				AcesCyl:           AcesCyl,
				AcesBlockType:     AcesBlockType,
				AspirationID:      AspirationID,
				DriveTypeID:       DriveTypeID,
				BodyTypeID:        BodyTypeID,
				BodyNumDoorsID:    BodyNumDoorsID,
				EngineVinID:       EngineVinID,
				RegionID:          RegionID,
				PowerOutputID:     PowerOutputID,
				FuelDelConfigID:   FuelDelConfigID,
				BodyStyleConfigID: BodyStyleConfigID,
				ValvesID:          ValvesID,
				CylHeadTypeID:     CylHeadTypeID,
				BlockType:         BlockType,
				EngineBaseID:      EngineBaseID,
				EngineConfigID:    EngineConfigID,
			},
			Part: Part{
				PCDBPartTerminologyName: PCDBPartTerminologyName,
				Position:                Position,
				OldID:                   PartNumber,
				PartDesc:                PartDesc,
			},
			PartData: PartData{
				VehicleCount:               VehicleCount,
				DistributedPartOpportunity: DistributedPartOpportunity,
				MaximumPartOpportunity:     MaximumPartOpportunity,
			},
		}

		//link new part numbers is boolean is set true
		if useOldPartNumbers == true {
			partMap, err := GetPartNumberMap()
			if err != nil {
				return cs, err
			}
			//get new part id, if there is one
			if newPartNum, ok := partMap[c.Part.OldID]; ok {
				c.Part.ID = newPartNum
			} else {
				//no new part number -> append to partsNeeded for output file write
				// partsNeeded = append(partsNeeded, c)
				partOffset, err = WriteVehicle(partsNeededFile, partOffset, c, 0, 0)
			}
		} else {
			//Curt (new) Part number is the only part number
			c.Part.ID, err = strconv.Atoi(c.Part.OldID)
		}

		//link curt base vehicle
		if curtBaseID, ok := baseMap[c.CsvVehicle.BaseVehicleID]; ok {
			c.CurtVehicle.CurtBaseID = curtBaseID
		} else {
			//missing base vehicle
			if insertMissingData == true {
				//TODO create base vehicle  in BaseVehicle table and link
			} else {
				// missingBaseVehicles = append(missingBaseVehicles, c)
				baseOffset, err = WriteVehicle(baseVehiclesNeededFile, baseOffset, c, 0, 0)
			}
		}

		//link curt submodel
		if curtSubID, ok := subMap[c.CsvVehicle.SubmodelID]; ok {
			c.CurtVehicle.CurtSubmodelID = curtSubID
		} else {
			//missing submodel
			if insertMissingData == true {
				//TODO create submodel in Submodel table and link
			} else {
				// missingSubmodels = append(missingSubmodels, c)
				submodelOffset, err = WriteVehicle(submodelsNeededFile, submodelOffset, c, 0, 0)
			}
		}

		cs = append(cs, c)
	}
	//return array of csv values, partsNeeded doc, err
	return cs, err
}

//Audits
func Audits(baseMap map[int][]CsvDatum) error {
	var err error
	subMap := make(map[int][]CsvDatum)
	var vehicleArray []CsvDatum

	subMap, err = AuditBaseVehicle(baseMap)
	if len(subMap) > 0 {
		vehicleArray, err = AuditSubmodel(subMap)
		if err != nil {
			return err
		}
	}
	if len(vehicleArray) > 0 {
		err = HandleVehicles(vehicleArray)
		if err != nil {
			return err
		}
	}
	return err
}

func AuditBaseVehicle(baseMap map[int][]CsvDatum) (map[int][]CsvDatum, error) {
	var err error
	submodelMap := make(map[int][]CsvDatum)

	//file prep
	baseNeeded, err := os.Create("NeedBaseVehicleInVcdbVehicleTable")
	defer baseNeeded.Close()
	off, err := WriteVehicleHeader(baseNeeded)

	for _, baseVehicle := range baseMap {
		baseFlag := false
		for i, base := range baseVehicle {
			if i > 0 {
				if base.Part.ID != baseVehicle[i-1].Part.ID {
					baseFlag = true //not all of these basevehicles have the same part number
					break
				}
			}
		}
		if baseFlag == false {
			log.Print("All the same part ", baseVehicle[0].Part.ID, " for this CsvBasevehicle: ", baseVehicle[0].CurtVehicle.CurtBaseID)
			//check for vehiclePartExistence
			vehicle, err := FindVehicle(baseVehicle[0].CurtVehicle.CurtBaseID, 0, 0)
			baseVehicle[0].CurtVehicle.ID = vehicle.ID
			if err != nil {
				if err != sql.ErrNoRows {
					return submodelMap, err
				} else {
					err = nil
					//log needed base vehicles in vcdbVehicle table
					off, err = WriteVehicle(baseNeeded, off, baseVehicle[0], 0, 0)
					//TODO - uncomment to insert; need to add base vehicle, assign ID to vehicle
					// err = baseVehicle[0].InsertBaseVehicleIntoVcdbVehicles()
					// if err != nil {
					// 	return submodelMap, err
					// }
				}
			}

			//check if vehiclepart exists
			err = baseVehicle[0].FindPartID()
			if err != nil {
				if err != sql.ErrNoRows {
					return submodelMap, err
				} else {
					//TODO -uncomment to add vehiclePart; need to add vehiclePart
					// err = baseVehicle[0].InsertPartIntoVehiclePart()
					// if err != nil {
					// 	return submodelMap, err
					// }
				}
			}
		} else {
			log.Print("Diff parts for base vehicle ", baseVehicle[0].CurtVehicle.CurtBaseID, ", try submodel")
			//There are different parts for this base vehicle, try submodel
			//build map of AAIAsubmodelID's to CsvData
			for _, base := range baseVehicle {
				submodelMap[base.CurtVehicle.CurtSubmodelID] = append(submodelMap[base.CurtVehicle.CurtSubmodelID], base)
			}
		}
	}
	return submodelMap, err
}

func AuditSubmodel(subMap map[int][]CsvDatum) ([]CsvDatum, error) {
	var err error
	// vIDmap := make(map[int][]CsvDatum)
	var vehicleArray []CsvDatum

	subNeeded, err := os.Create("NeedSubmodelInVcdbVehicleTable")
	if err != nil {
		return vehicleArray, err
	}
	defer subNeeded.Close()
	off, err := WriteVehicleHeader(subNeeded)

	for _, subVehicle := range subMap {
		subFlag := false
		for i, sub := range subVehicle {
			if i > 0 {
				if sub.Part.ID != subVehicle[i-1].Part.ID {
					subFlag = true
					break
				}
			}
		}
		if subFlag == false {
			//add part to sub vehicle
			log.Print("All the same part ", subVehicle[0].CurtVehicle.CurtSubmodelID, " for this CsvSubmodel: ", subVehicle[0].CurtVehicle.CurtSubmodelID, ". CsvBaseID: ", subVehicle[0].CurtVehicle.CurtBaseID)
			//check for vehiclePartExistence
			vehicle, err := FindVehicle(subVehicle[0].CurtVehicle.CurtBaseID, subVehicle[0].CurtVehicle.CurtSubmodelID, 0)
			subVehicle[0].CurtVehicle.ID = vehicle.ID
			if err != nil {
				if err != sql.ErrNoRows {
					return vehicleArray, err
				} else {
					err = nil
					//log submodel needed in vcdbVehicle table
					off, err = WriteVehicle(subNeeded, off, subVehicle[0], 0, 0)
					//TODO - uncomment to insert; need to add submodel, assign ID to vehicle
					// err = subVehicle[0].InsertSubmodelIntoVcdbVehicles()
					// if err != nil {
					// 	return vIDmap, err
					// }
				}
			}

			//check if vehiclepart exists
			err = subVehicle[0].FindPartID()
			if err != nil {
				if err != sql.ErrNoRows {
					return vehicleArray, err
				} else {
					// TODO -uncomment to add vehiclePart; add vehiclePart
					// err = subVehicle[0].InsertPartIntoVehiclePart()
					// if err != nil {
					// 	return vIDmap, err
					// }
				}
			}
		} else {
			log.Print("Diff parts for submodel ", subVehicle[0].CurtVehicle.CurtSubmodelID, ", try configs")
			//config breakdown
			//make map of un-added AAIAVehicleID to CsvData
			// for _, sub := range subVehicle {
			// vIDmap[sub.CsvVehicle.VehicleID] = append(vIDmap[sub.CsvVehicle.VehicleID], sub)
			// }
			for _, sub := range subVehicle {
				vehicleArray = append(vehicleArray, sub)
			}
		}
	}
	return vehicleArray, err
}

func HandleVehicles(vehicleArray []CsvDatum) error {
	var err error
	//File prep
	configsNeededFile, err := os.Create("ConfigsNeeded")
	if err != nil {
		return err
	}
	defer configsNeededFile.Close()
	off, err := WriteVehicleHeader(configsNeededFile)

	acesConfigTypeArray := [...]int{6, 20, 8, 3, 2, 4, 16, 25, 40, 12, 7} //There's got to be a better way
	//assign configs to each vehicles' config array
	log.Print("ARRAY ", vehicleArray)
	for _, v := range vehicleArray {
		for _, acesConfigTypeID := range acesConfigTypeArray {
			var config CurtConfig
			config.AcesConfigTypeID = acesConfigTypeID
			config.AcesConfigValueID = v.getAcesConfigValue(acesConfigTypeID)

			//get curt config type from aces type
			config.ConfigTypeID, config.ConfigTypeName, err = GetCurtConfigTypeAcesConfig(config.AcesConfigTypeID)
			if err != nil {
				log.Print("MISSING configtype, ", err)
				return err
			}
			//get Curt value from aces value and type, if there is one
			config.ConfigValueID, config.ConfigValueName, err = GetCurtConfigValueAcesConfig(config.AcesConfigTypeID, config.AcesConfigValueID)
			if err != nil {
				if err != sql.ErrNoRows {
					return err
				} else {
					err = nil
					//log missing config val
					off, err = WriteVehicle(configsNeededFile, off, v, config.AcesConfigTypeID, config.AcesConfigValueID)

					if err != nil {
						return err
					}
					//TODO add missing configs?
				}
			}

			v.CurtVehicle.CurtConfigs = append(v.CurtVehicle.CurtConfigs, config)
		}
	}
	return err
}

func WriteVehicle(f *os.File, off int64, c CsvDatum, acesConfigTypeID int, acesConfigValueID int) (int64, error) {
	var err error
	csvVehicleString := c.CsvVehicle.Make + "," +
		c.CsvVehicle.Model + "," +
		c.CsvVehicle.SubModel + "," +
		c.CsvVehicle.Year + "," +
		strconv.Itoa(c.CsvVehicle.GVW) + "," +
		strconv.Itoa(c.CsvVehicle.VehicleID) + "," +
		strconv.Itoa(c.CsvVehicle.BaseVehicleID) + "," +
		strconv.Itoa(c.CsvVehicle.YearID) + "," +
		strconv.Itoa(c.CsvVehicle.MakeID) + "," +
		strconv.Itoa(c.CsvVehicle.ModelID) + "," +
		strconv.Itoa(c.CsvVehicle.SubmodelID) + "," +
		strconv.Itoa(c.Part.ID) + "," +
		c.Part.PartDesc + "," +
		strconv.Itoa(c.CurtVehicle.ID) + "," +
		strconv.Itoa(c.CurtVehicle.CurtBaseID) + "," +
		strconv.Itoa(c.CurtVehicle.CurtSubmodelID) + "," +
		strconv.Itoa(acesConfigTypeID) + "," +
		strconv.Itoa(acesConfigValueID) + "," + "\n"
	b := []byte(csvVehicleString)
	n, err := f.WriteAt(b, off)
	if err != nil {
		return off, err
	}
	off += int64(n)
	return off, err
}

func WriteVehicleHeader(f *os.File) (int64, error) {
	off := int64(0)
	b := []byte(`CsvVehicleMake,CsvVehicleModel,CsvVehicleSubmodel,CsvVehicleYear,CsvVehicleGVW,CsvVehicleVehicleID,
		CsvVehicleBaseVehicleID,CsvVehicleYearID,CsvVehicleMakeID,CsvVehicleModelID,CsvVehicleSubmodelID,
		PartID,PartDesc,CurtVehicleID,CurtVehicleBaseID,CurtVehicleSubmodelID,AAIAConfigTypeID,AAIAConfigValueID` + "\n")
	n, err := f.WriteAt(b, off)
	off += int64(n)
	return off, err
}

func (c *CsvDatum) getAcesConfigValue(n int) int {
	staticMap := make(map[int]int)
	staticMap[6] = c.CsvVehicle.FuelTypeID
	staticMap[20] = c.CsvVehicle.FuelDeliveryID
	staticMap[8] = c.CsvVehicle.AspirationID
	staticMap[3] = c.CsvVehicle.DriveTypeID
	staticMap[2] = c.CsvVehicle.BodyTypeID
	staticMap[4] = c.CsvVehicle.BodyNumDoorsID
	staticMap[16] = c.CsvVehicle.EngineVinID
	staticMap[25] = c.CsvVehicle.PowerOutputID
	staticMap[40] = c.CsvVehicle.ValvesID
	staticMap[12] = c.CsvVehicle.CylHeadTypeID
	staticMap[7] = c.CsvVehicle.EngineBaseID
	return staticMap[n]
}

//See if part exists in vehiclePart table
func (c *CsvDatum) FindPartID() error {
	var err error
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(findPartID)
	if err != nil {
		return err
	}
	defer stmt.Close()
	var id int
	err = stmt.QueryRow(c.CurtVehicle.ID, c.Part.ID).Scan(&id)
	if err != nil {
		return err
	}
	return err
}

//get curt vehicleID
func FindVehicle(curtBaseID, curtSubmodelID, configID int) (CurtVehicle, error) {
	var err error
	var curt CurtVehicle
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return curt, err
	}
	defer db.Close()

	stmt, err := db.Prepare(findVehicle)
	if err != nil {
		return curt, err
	}
	defer stmt.Close()
	err = stmt.QueryRow(curtBaseID, curtSubmodelID, configID).Scan(&curt.ID)
	if err != nil {
		return curt, err
	}
	return curt, err
}

func (c *CsvDatum) GetCurtVechicleWithConfig(curtConfigValueID, curtConfigTypeID int) error {
	var err error
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(getVehicleWithConfig)
	if err != nil {
		return err
	}
	defer stmt.Close()
	err = stmt.QueryRow(c.CurtVehicle.CurtBaseID, c.CurtVehicle.CurtSubmodelID, curtConfigValueID, curtConfigTypeID).Scan(&c.CurtVehicle.ID)
	return err
}

//curt configs
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

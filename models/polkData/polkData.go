package polkData

import (
	"database/sql"
	"encoding/csv"
	"github.com/curt-labs/GoAPI/helpers/database"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"strings"
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
	GVW               uint8
	VehicleID         int
	BaseVehicleID     int
	YearID            int
	MakeID            int
	ModelID           int
	SubmodelID        int
	VehicleTypeID     uint8
	FuelTypeID        uint8   // 6 FuelType
	FuelDeliveryID    uint8   //20 FuelDeliveryType
	AcesLiter         float64 //EngineBase.Liter
	AcesCC            float64 //EngineBase.CC
	AcesCID           uint16  //EngineBase.CID
	AcesCyl           uint8   //EngineBase.Cylinders
	AcesBlockType     string  //EngineBase.BlockType
	AspirationID      uint8   // 8 Aspiration
	DriveTypeID       uint8   // 3 DriveType
	BodyTypeID        uint8   // 2 BodyType
	BodyNumDoorsID    uint8   // 4 BodyNumDoors
	EngineVinID       uint8   // 16 EngineVIN
	RegionID          uint8   //Region
	PowerOutputID     uint16  // 25 PowerOutput
	FuelDelConfigID   uint8   //FuelDeliveryConfig
	BodyStyleConfigID uint8   //BodyStyleConfig
	ValvesID          uint8   // 40 Valves
	CylHeadTypeID     uint8   // 12 CylinderHeadType
	BlockType         string  //EngineBase.BlockType
	EngineBaseID      uint16  // 7 EngineBase
	EngineConfigID    uint16  //EngineConfig
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

func Run(filename string, headerLines int, useOldPartNumbers bool, insertMissingData bool) error {
	var err error
	var cs CsvData
	var m runtime.MemStats

	//csv into memory
	cs, err = CaptureCsv(filename, headerLines, useOldPartNumbers, insertMissingData)
	if err != nil {
		return err
	}
	runtime.ReadMemStats(&m)
	log.Print("Allocated: ", m.Alloc, " NextGC: ", m.NextGC, " LastCG: ", m.LastGC, " Frees: ", m.Frees)

	ch := make(chan map[int][]CsvDatum)
	arrayChan := make(chan []CsvDatum)
	go func() {
		baseMap := make(map[int][]CsvDatum)
		for _, c := range cs {
			baseMap[c.CurtVehicle.CurtBaseID] = append(baseMap[c.CurtVehicle.CurtBaseID], c)
		}
		ch <- baseMap
	}()
	baseMap := <-ch
	runtime.GC()
	log.Print("Allocated: ", m.Alloc, " NextGC: ", m.NextGC, " LastCG: ", m.LastGC, " Frees: ", m.Frees)

	go func() {
		subMap, err := AuditBaseVehicle(baseMap, insertMissingData)
		if err != nil {
			return
		}
		ch <- subMap
	}()
	subMap := <-ch
	runtime.GC()
	log.Print("Allocated: ", m.Alloc, " NextGC: ", m.NextGC, " LastCG: ", m.LastGC, " Frees: ", m.Frees)

	go func() {
		vehicleArray, err := AuditSubmodel(subMap, insertMissingData)
		if err != nil {
			return
		}
		arrayChan <- vehicleArray
	}()
	vehicleArray := <-arrayChan
	runtime.GC()
	log.Print("Allocated: ", m.Alloc, " NextGC: ", m.NextGC, " LastCG: ", m.LastGC, " Frees: ", m.Frees)

	if len(vehicleArray) > 0 {
		err = HandleVehicles(vehicleArray, insertMissingData)
		if err != nil {
			return err
		}
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
				GVW:               uint8(GVW),
				VehicleID:         VehicleID,
				BaseVehicleID:     BaseVehicleID,
				YearID:            YearID,
				MakeID:            MakeID,
				ModelID:           ModelID,
				SubmodelID:        SubmodelID,
				VehicleTypeID:     uint8(VehicleTypeID),
				FuelTypeID:        uint8(FuelTypeID),
				FuelDeliveryID:    uint8(FuelDeliveryID),
				AcesLiter:         AcesLiter,
				AcesCC:            AcesCC,
				AcesCID:           uint16(AcesCID),
				AcesCyl:           uint8(AcesCyl),
				AcesBlockType:     AcesBlockType,
				AspirationID:      uint8(AspirationID),
				DriveTypeID:       uint8(DriveTypeID),
				BodyTypeID:        uint8(BodyTypeID),
				BodyNumDoorsID:    uint8(BodyNumDoorsID),
				EngineVinID:       uint8(EngineVinID),
				RegionID:          uint8(RegionID),
				PowerOutputID:     uint16(PowerOutputID),
				FuelDelConfigID:   uint8(FuelDelConfigID),
				BodyStyleConfigID: uint8(BodyStyleConfigID),
				ValvesID:          uint8(ValvesID),
				CylHeadTypeID:     uint8(CylHeadTypeID),
				BlockType:         BlockType,
				EngineBaseID:      uint16(EngineBaseID),
				EngineConfigID:    uint16(EngineConfigID),
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
				//create base vehicle  in BaseVehicle table and link
				err = c.InsertBaseVehicle()
				if err != nil {
					return cs, err
				}
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
				//create submodel in Submodel table and link
				err = c.InsertSubmodel()
				if err != nil {
					return cs, err
				}
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

func AuditBaseVehicle(baseMap map[int][]CsvDatum, insertMissingData bool) (map[int][]CsvDatum, error) {
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
					if insertMissingData == false {
						//log needed base vehicles in vcdbVehicle table
						off, err = WriteVehicle(baseNeeded, off, baseVehicle[0], 0, 0)
					} else {
						// add base vehicle, assign ID to vehicle
						err = baseVehicle[0].InsertBaseVehicleIntoVcdbVehicles()
						if err != nil {
							return submodelMap, err
						}
					}
				}
			}

			//check if vehiclepart exists
			err = baseVehicle[0].FindPartID()
			if err != nil {
				if err != sql.ErrNoRows {
					return submodelMap, err
				} else {
					if insertMissingData == true {
						//add vehiclePart
						err = baseVehicle[0].InsertPartIntoVehiclePart()
						if err != nil {
							return submodelMap, err
						}
					}
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
	baseMap = nil
	return submodelMap, err
}

func AuditSubmodel(subMap map[int][]CsvDatum, insertMissingData bool) ([]CsvDatum, error) {
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
					if insertMissingData == false {
						//log submodel needed in vcdbVehicle table
						off, err = WriteVehicle(subNeeded, off, subVehicle[0], 0, 0)
					} else {
						// add submodel, assign ID to vehicle
						err = subVehicle[0].InsertSubmodelIntoVcdbVehicles()
						if err != nil {
							return vehicleArray, err
						}
					}
				}
			}

			//check if vehiclepart exists
			err = subVehicle[0].FindPartID()
			if err != nil {
				if err != sql.ErrNoRows {
					return vehicleArray, err
				} else {
					if insertMissingData == true {
						//  add vehiclePart
						err = subVehicle[0].InsertPartIntoVehiclePart()
						if err != nil {
							return vehicleArray, err
						}
					}
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
	subMap = nil
	return vehicleArray, err
}

func HandleVehicles(vehicleArray []CsvDatum, insertMissingData bool) error {
	var err error
	vehicleIDmap := make(map[int][]CsvDatum)
	//get configmap
	configMap, err := GetConfigMap()
	//File prep
	// configsNeededFile, err := os.Create("ConfigsNeeded")
	// if err != nil {
	// 	return err
	// }
	// defer configsNeededFile.Close()
	// off, err := WriteVehicleHeader(configsNeededFile)

	acesConfigTypeArray := [...]int{6, 20, 8, 3, 2, 4, 16, 25, 40, 12, 7} //There's got to be a better way
	//assign configs to each vehicles' config array
	// log.Print("ARRAY ", vehicleArray)
	for _, v := range vehicleArray {
		for _, acesConfigTypeID := range acesConfigTypeArray {
			var config CurtConfig
			config.AcesConfigTypeID = acesConfigTypeID
			config.AcesConfigValueID = v.getAcesConfigValue(acesConfigTypeID)

			//get Curt value from aces value and type, if there is one
			acesTypeAndValue := strconv.Itoa(config.AcesConfigTypeID) + "," + strconv.Itoa(config.AcesConfigValueID)
			//check against map
			if acesTV, ok := configMap[acesTypeAndValue]; ok {
				curtTVArray := strings.Split(acesTV, ",")
				config.ConfigTypeID, err = strconv.Atoi(curtTVArray[0])
				config.ConfigValueID, err = strconv.Atoi(curtTVArray[1])
				// log.Print("FIND PARTS ? ", v.CurtVehicle)
				//find part matches?

				// } else {
				// 	off, err = WriteVehicle(configsNeededFile, off, v, config.AcesConfigTypeID, config.AcesConfigValueID)
			}

			v.CurtVehicle.CurtConfigs = append(v.CurtVehicle.CurtConfigs, config)
		}
		vehicleIDmap[v.CsvVehicle.VehicleID] = append(vehicleIDmap[v.CsvVehicle.VehicleID], v)
	}

	vehicleArray = nil
	err = diffVehicleConfigs(vehicleIDmap, insertMissingData)
	return err
}

func diffVehicleConfigs(vehicleIDmap map[int][]CsvDatum, insertMissingData bool) error {
	var err error
	configsDiff, err := os.Create("ConfigsDiff")
	if err != nil {
		return err
	}
	defer configsDiff.Close()
	off, err := WriteVehicleHeader(configsDiff)
	for _, vehicles := range vehicleIDmap {
		allConfigsEqualFlag := true
		for _, vehicle := range vehicles {
			for i, config := range vehicle.CurtVehicle.CurtConfigs {
				if i > 0 {
					allConfigsEqualFlag = reflect.DeepEqual(config, vehicle.CurtVehicle.CurtConfigs[i-1])
				}
			}
		}
		if allConfigsEqualFlag == false {
			for _, vehicle := range vehicles {
				off, err = WriteVehicle(configsDiff, off, vehicle, 0, 0)
				//check and add part to every config
			}
		} else {
			//Check and add part to submodel
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
		strconv.Itoa(int(c.CsvVehicle.GVW)) + "," +
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
	staticMap[6] = int(c.CsvVehicle.FuelTypeID)
	staticMap[20] = int(c.CsvVehicle.FuelDeliveryID)
	staticMap[8] = int(c.CsvVehicle.AspirationID)
	staticMap[3] = int(c.CsvVehicle.DriveTypeID)
	staticMap[2] = int(c.CsvVehicle.BodyTypeID)
	staticMap[4] = int(c.CsvVehicle.BodyNumDoorsID)
	staticMap[16] = int(c.CsvVehicle.EngineVinID)
	staticMap[25] = int(c.CsvVehicle.PowerOutputID)
	staticMap[40] = int(c.CsvVehicle.ValvesID)
	staticMap[12] = int(c.CsvVehicle.CylHeadTypeID)
	staticMap[7] = int(c.CsvVehicle.EngineBaseID)
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

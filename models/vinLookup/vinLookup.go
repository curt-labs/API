package vinLookup

import (
	"bytes"
	"database/sql"
	"encoding/base64"
	"encoding/xml"
	"errors"
	"github.com/curt-labs/GoAPI/helpers/apicontext"
	"github.com/curt-labs/GoAPI/helpers/database"
	"io/ioutil"
	"net/http"
	"strconv"
)

type AcesVehicle struct {
	AcesID            int
	AAIABaseVehicleID int
	AAIAMakeID        int
	AAIAModelID       int
	AAIAYearID        int
	AAIASubmodelID    int
	AAIARegionID      int
}
type CurtVehicle struct {
	ID            int
	BaseVehicle   BaseVehicle
	Submodel      Submodel
	Configuration VehicleConfiguration
	// Parts         []products.Part
}

type BaseVehicle struct {
	ID        int
	ModelID   int
	MakeID    int
	YearID    int
	ModelName string
	MakeName  string
}

type Submodel struct {
	ID   int
	Name string
}

type VehicleConfiguration struct {
	TypeID      int //aka key id
	ValueID     int
	Type        string //aka key
	Value       string
	AcesValueID int
}

type ConfigurationBits struct {
	WheelBase                        interface{} //WHL_BAS_SHRST_INCHS
	BodyType                         interface{} //ACES_BODY_TYPE
	DriveType                        interface{} //ACES_DRIVE_ID
	NumberOfDoors                    interface{} //DOOR_CNT
	FuelType                         interface{}
	Engine                           interface{} //ACES_LITERS + ACES_CYLINDERS--not quite
	Aspiration                       interface{} //ACES_ASP_ID
	BedLength                        interface{} //TRK_BED_LEN_CD
	BedType                          interface{}
	BrakeABS                         interface{}
	BrakeSystem                      interface{}
	CylinderHeadType                 interface{}
	EngineDesignation                interface{}
	EngineManufacturer               interface{}
	EngineVersion                    interface{}
	EngineVin                        interface{} //ACES_ENG_VIN_ID
	FrontBrakeType                   interface{}
	FrontSpringType                  interface{}
	FuelDeliverySubType              interface{}
	FuelDeliveryType                 interface{} //ACES_FUEL
	FuelSystemControlType            interface{}
	FuelSystemDesign                 interface{} //ACES_FUEL
	IgnitionSystemDesign             interface{}
	ManufacturerBodyCode             interface{}
	PowerOutput                      interface{}
	RearBrakeType                    interface{}
	RearSpringType                   interface{}
	SteeringSystem                   interface{}
	SteeringType                     interface{}
	TransmissionElectronicControlled interface{}
	Transmission                     interface{} //TRANS_CD
	TransmissionControlType          interface{}
	TransmissionManufacturerCode     interface{}
	TransmissionNumberOfSpeeds       interface{} //TRANS_OPT1_SPEED_CD
	TransmissionType                 interface{}
	ValvesPerEngine                  interface{}
	Region                           interface{} //ACES_REGION_ID
}

//reponse
type XMLResponse struct {
	XMLName xml.Name
	Body    Body
}
type Body struct {
	XMLName           xml.Name
	DecodeVinResponse DecodeVinResponse `xml:"decodeVinResponse"`
}

type DecodeVinResponse struct {
	XMLName     xml.Name    `xml:"decodeVinResponse"`
	VinResponse VinResponse `xml:"VinResponse"`
}

type VinResponse struct {
	Vin          string  `xml:"vin"`
	ReturnCode   string  `xml:"returnCode"`
	CorrectedVin string  `xml:"correctedVin"`
	ErrorBytes   string  `xml:"errorBytes"`
	Fields       []Field `xml:"fields"`
}
type Field struct {
	Key   string `xml:"name,attr"`
	Value string `xml:",chardata"`
}

//request
type Envelope struct {
	XMLName xml.Name       `xml:"soapenv:Envelope"`
	SoapEnv string         `xml:"xmlns:soapenv,attr"`
	Web     string         `xml:"xmlns:web,attr"`
	Header  EnvelopeHeader `xml:"soapenv:Header"`
	Body    EnvelopeBody   `xml:"soapenv:Body"`
}

type EnvelopeHeader struct {
}
type EnvelopeBody struct {
	DecodeVin DecodeVin `xml:"web:decodeVin"`
}
type DecodeVin struct {
	Vin             string `xml:"VinRequest>vin"`
	RequestedFields string `xml:"RequestedFields"`
}

var (
	getCurtVehiclesPreConfig = `SELECT vv.ID, vmd.ID,vmd.ModelName, vmk.ID, vmk.MakeName, vyr.YearID, sm.ID, sm.SubmodelName, cat.name, cat.ID, ca.value, ca.ID, ca.vcdbID
								FROM vcdb_Vehicle AS vv
								LEFT JOIN BaseVehicle AS bv ON bv.ID = vv.BaseVehicleID
								LEFT JOIN vcdb_Model AS vmd ON vmd.ID = bv.ModelID
								LEFT JOIN vcdb_Make AS vmk ON vmk.ID = bv.MakeID
								LEFT JOIN vcdb_Year AS vyr ON vyr.YearID = bv.YearID
								LEFT JOIN Submodel AS sm ON sm.ID = vv.SubmodelID
								LEFT JOIN VehicleConfigAttribute AS vca ON vca.VehicleConfigID = vv.ConfigID
								LEFT JOIN ConfigAttribute AS ca ON ca.ID = vca.AttributeID
								LEFT JOIN ConfigAttributeType AS cat ON cat.ID = ca.ConfigAttributeTypeID
								WHERE bv.AAIABaseVehicleID = ?
								AND (sm.AAIASubmodelID = ?  OR sm.AAIASubmodelID IS NULL) `

	getPartID = `SELECT PartNumber FROM vcdb_VehiclePart WHERE VehicleID = ?`
)

const (
	soapRequestedFields = `ACES_BASE_VEHICLE,ACES_MAKE_ID,ACES_MDL_ID,ACES_SUB_MDL_ID,ACES_YEAR_ID,ACES_REGION_ID,ACES_VEHICLE_ID,
		ACES_FUEL,ACES_FUEL_DELIVERY,ACES_ENG_VIN_ID,ACES_ASP_ID,ACES_DRIVE_ID,ACES_BODY_TYPE,ACES_REGION_ID,ACES_LITERS,ACES_CC_DISPLACEMENT,ACES_CI_DISPLACEMENT,
		ACES_CYLINDERS,ACES_RESERVED,DOOR_CNT,BODY_STYLE_DESC,WHL_BAS_SHRST_INCHS,TRK_BED_LEN_DESC,TRANS_CD,TRK_BED_LEN_CD,ENG_FUEL_DESC`
)

func VinPartLookup(vin string, dtx *apicontext.DataContext) (l products.Lookup, err error) {
	//get ACES vehicles
	av, configMap, err := getAcesVehicle(vin)
	if err != nil {
		return l, err
	} else if av.AAIABaseVehicleID == 0 {
		return l, errors.New("failed to decode VIN")
	}

	//get CURT vehicle
	l, err = av.getCurtVehicles(configMap)
	if err != nil {
		return l, err
	}

	//get parts
	var ps []products.Part
	ch := make(chan []products.Part)
	go l.LoadParts(ch, dtx)
	ps = <-ch

	l.Parts = ps
	if len(l.Parts) == 0 {
		err = sql.ErrNoRows
	}
	return l, err
}

func GetVehicleConfigs(vin string) (l products.Lookup, err error) {
	//get ACES vehicles
	av, configMap, err := getAcesVehicle(vin)
	if err != nil {
		return l, err
	} else if av.AAIABaseVehicleID == 0 {
		return l, errors.New("failed to decode VIN")
	}

	//get CURT vehicle
	l, err = av.getCurtVehicles(configMap)
	if len(l.Parts) == 0 {
		err = sql.ErrNoRows
	}
	return l, err
}

//already have vehicleID (vcdb_vehicle.ID)? get parts
func (v *CurtVehicle) GetPartsFromVehicleConfig(dtx *apicontext.DataContext) (ps []products.Part, err error) {
	//get parts
	var p products.Part
	//get part id
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return ps, err
	}
	defer db.Close()
	stmt, err := db.Prepare(getPartID)
	if err != nil {
		return ps, err
	}
	defer stmt.Close()
	res, err := stmt.Query(v.ID)
	for res.Next() {
		err = res.Scan(&p.ID)
		if err != nil {
			return ps, err
		}
		//get part -- adds some weight
		err = p.FromDatabase(dtx)
		if err != nil {
			return ps, err
		}

		ps = append(ps, p)
	}
	defer res.Close()
	return ps, err
}

func query(vin string) (output []byte, err error) {
	var e Envelope
	e.SoapEnv = "http://schemas.xmlsoap.org/soap/envelope/"
	e.Web = "http://webservice.vindecoder.polk.com/"
	e.Body.DecodeVin.Vin = vin
	e.Body.DecodeVin.RequestedFields = soapRequestedFields

	output, err = xml.MarshalIndent(e, " ", "\t")
	if err != nil {
		return output, err
	}
	return output, err
}

func getAcesVehicle(vin string) (av AcesVehicle, configMap map[int]interface{}, err error) {
	data := []byte(database.VintelligencePass())
	password := base64.StdEncoding.EncodeToString(data)

	b, err := query(vin)
	if err != nil {
		return av, configMap, err
	}
	buffer := bytes.NewReader(b)
	client := http.Client{}
	req, err := http.NewRequest("POST", "https://vintelligence3.polk.com/vindecoder/VinDecoderService", buffer)
	if err != nil {
		return av, configMap, err
	}
	req.Header.Add("Authorization", "Basic "+password)
	req.Header.Add("Content-Type", "text/xml;charset=utf-8")
	req.Header.Add("Host", "\"api.curtmfg.com\"")

	resp, err := client.Do(req)
	if err != nil {
		return av, configMap, err
	}

	if resp.StatusCode != 200 {
		err = errors.New(resp.Status)
		return av, configMap, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return av, configMap, err
	}

	var x XMLResponse
	err = xml.Unmarshal(body, &x)
	if err != nil {
		return av, configMap, err
	}

	for _, field := range x.Body.DecodeVinResponse.VinResponse.Fields {
		switch field.Key {
		case "ACES_BASE_VEHICLE":
			if field.Value != "" {
				av.AAIABaseVehicleID, err = strconv.Atoi(field.Value)
			}
		case "ACES_MAKE_ID":
			if field.Value != "" {
				av.AAIAMakeID, err = strconv.Atoi(field.Value)
			}
		case "ACES_MDL_ID":
			if field.Value != "" {
				av.AAIAModelID, err = strconv.Atoi(field.Value)
			}
		case "ACES_SUB_MDL_ID":
			if field.Value != "" {
				av.AAIASubmodelID, err = strconv.Atoi(field.Value)
			}
		case "ACES_YEAR_ID":
			if field.Value != "" {
				av.AAIAYearID, err = strconv.Atoi(field.Value)
			}
		case "ACES_REGION_ID":
			if field.Value != "" {
				av.AAIARegionID, err = strconv.Atoi(field.Value)
			}
		case "ACES_VEHICLE_ID":
			if field.Value != "" {
				av.AcesID, err = strconv.Atoi(field.Value)
			}
		}
	}
	//return code error?
	rc := x.Body.DecodeVinResponse.VinResponse
	returnCode, err := strconv.Atoi(rc.ReturnCode)

	if returnCode > 3 {
		switch returnCode {
		case 4:
			err = errors.New("Could not decode. Check digit calculates properly. Return Code: " + rc.ReturnCode)
			return av, configMap, err
		case 5:
			err = errors.New("Could not decode. Check digit does not calculate properly. Return Code: " + rc.ReturnCode)
			return av, configMap, err
		case 6:
			err = errors.New("Customer is not licensed to receive data. Return Code: " + rc.ReturnCode)
			return av, configMap, err
		default:
			err = errors.New("Error decoding VIN. Return Code: " + rc.ReturnCode)
			return av, configMap, err
		}
	}
	//check out them configs
	configMap, err = av.checkConfigs(x.Body.DecodeVinResponse.VinResponse.Fields)

	return av, configMap, err
}

//creates a map of config options from the SOAP request to check against curt vehicles
func (av *AcesVehicle) checkConfigs(responseFields []Field) (configMap map[int]interface{}, err error) {
	//map of configAttributeType AcesID to configAttribute Aces ID
	configMap = make(map[int]interface{})
	for _, field := range responseFields {
		switch field.Key {
		case "WHL_BAS_SHRST_INCHS":
			if field.Value != "" {
				configMap[1], err = strconv.Atoi(field.Value)
			}
		case "ACES_BODY_TYPE":
			if field.Value != "" {
				configMap[2], err = strconv.Atoi(field.Value)
			}
		case "ACES_DRIVE_ID":
			if field.Value != "" {
				configMap[3], err = strconv.Atoi(field.Value)
			}
		case "DOOR_CNT":
			if field.Value != "" {
				configMap[4], err = strconv.Atoi(field.Value)
			}
		case "ACES_ASP_ID":
			if field.Value != "" {
				configMap[8], err = strconv.Atoi(field.Value)
			}
		case "ACES_ENG_VIN_ID":
			if field.Value != "" {
				configMap[16], err = strconv.Atoi(field.Value)
			}
		case "ACES_FUEL":
			if field.Value != "" {
				configMap[20], err = strconv.Atoi(field.Value)
			}
		case "TRANS_CD":
			if field.Value != "" {
				configMap[34] = field.Value
			}
		case "TRANS_OPT1_SPEED_CD":
			if field.Value != "" {
				configMap[38], err = strconv.Atoi(field.Value)
			}
			if err != nil {
				return configMap, err
			}
		}

	}

	return configMap, err
}

//sierra 3500 vin 1GTJK34131E957990

func (av *AcesVehicle) getCurtVehicles(configMap map[int]interface{}) (products.Lookup, error) { //get CURT vehicles
	var l products.Lookup
	var err error
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return l, err
	}
	defer db.Close()
	stmt, err := db.Prepare(getCurtVehiclesPreConfig)
	if err != nil {
		return l, err
	}
	defer stmt.Close()
	res, err := stmt.Query(av.AAIABaseVehicleID, av.AAIASubmodelID)
	if err != nil {
		return l, err
	}

	var sub, configKey, configValue *string
	var subID, configKeyID, configValueID, acesConfigValID *int
	var cv CurtVehicle

	var pco products.ConfigurationOption
	var vehicleConfig products.Configuration

	pcoMap := make(map[string][]string)

	for res.Next() {

		err = res.Scan(
			&cv.ID,
			&cv.BaseVehicle.ModelID,
			&cv.BaseVehicle.ModelName,
			&cv.BaseVehicle.MakeID,
			&cv.BaseVehicle.MakeName,
			&cv.BaseVehicle.YearID,
			&subID,
			&sub,
			&configKey,
			&configKeyID,
			&configValue,
			&configValueID,
			&acesConfigValID,
		)
		if subID != nil {
			cv.Submodel.ID = *subID
		}
		if sub != nil {
			cv.Submodel.Name = *sub
		}
		if configKey != nil {
			cv.Configuration.Type = *configKey
		}
		if configValue != nil {
			cv.Configuration.Value = *configValue
		}
		if configKeyID != nil {
			cv.Configuration.TypeID = *configKeyID
		}
		if configValueID != nil {
			cv.Configuration.ValueID = *configValueID
		}
		if acesConfigValID != nil {
			cv.Configuration.AcesValueID = *acesConfigValID
		}

		//configs - assign to map, flag
		configValFlag := true
		if vs, ok := pcoMap[cv.Configuration.Type]; ok {
			for _, v := range vs {
				if v == cv.Configuration.Value {
					configValFlag = false
				}
			}
		}
		if configValFlag == true {
			if name, ok := configMap[cv.Configuration.TypeID]; ok {
				//configMap contains this config type

				if cv.Configuration.AcesValueID == name {
					pcoMap[cv.Configuration.Type] = append(pcoMap[cv.Configuration.Type], cv.Configuration.Value)

					//vehicleConfigs (not l.ConfugurationOption)
					vehicleConfig.Key = cv.Configuration.Type
					vehicleConfig.Value = cv.Configuration.Value
					l.Vehicle.Configurations = append(l.Vehicle.Configurations, vehicleConfig)
				}
			} else {
				pcoMap[cv.Configuration.Type] = append(pcoMap[cv.Configuration.Type], cv.Configuration.Value)

				//vehicleConfigs (not l.ConfugurationOption)
				vehicleConfig.Key = cv.Configuration.Type
				vehicleConfig.Value = cv.Configuration.Value
				l.Vehicle.Configurations = append(l.Vehicle.Configurations, vehicleConfig)
			}

		}

		l.Vehicle.Base.Make = cv.BaseVehicle.MakeName
		l.Vehicle.Base.Model = cv.BaseVehicle.ModelName
		l.Vehicle.Base.Year = cv.BaseVehicle.YearID
		l.Vehicle.Submodel = cv.Submodel.Name

	} //end scan loop
	defer res.Close()

	//assign configs
	for key, val := range pcoMap {
		pco.Type = key
		pco.Options = val
		l.Configurations = append(l.Configurations, pco)
	}

	l.Makes = append(l.Makes, l.Vehicle.Base.Make)
	l.Models = append(l.Models, l.Vehicle.Base.Model)
	l.Years = append(l.Years, l.Vehicle.Base.Year)
	l.Submodels = append(l.Submodels, l.Vehicle.Submodel)

	return l, err
}
